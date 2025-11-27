package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/block"
	"github.com/LatticeBCLab/go-lattice/wallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog/log"
)

// emptyChainId is the empty chain ID. Represents the mainnet chain id.
const emptyChainId = ""

// JsonRpcBody is the request structure for JSON-RPC requests.
type JsonRpcBody struct {
	Id      int           `json:"id,omitempty"`
	JsonRpc string        `json:"jsonrpc,omitempty"`
	Method  string        `json:"method,omitempty"` // 方法名
	Params  []interface{} `json:"params,omitempty"` // 方法参数
}

// JsonRpcResponse is the response structure for JSON-RPC requests.
type JsonRpcResponse[T any] struct {
	Id      int           `json:"id,omitempty"`
	JsonRpc string        `json:"jsonrpc,omitempty"`
	Result  T             `json:"result,omitempty"`
	Error   *JsonRpcError `json:"error,omitempty"`
}

type JsonRpcError struct {
	Code    int16  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Error returns the error message.
func (e *JsonRpcError) Error() error {
	return fmt.Errorf("%d:%s", e.Code, e.Message)
}

// NewJsonRpcBody creates a new JSON-RPC body.
func NewJsonRpcBody(method string, params ...interface{}) *JsonRpcBody {
	return &JsonRpcBody{
		Id:      1,
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

// NewJwt creates a new JWT instance.
func NewJwt(secret string, expirationDuration time.Duration) Jwt {
	if secret == "" {
		return nil
	}
	return &jwtImpl{
		Secret:             secret,
		Algorithm:          jwt.SigningMethodHS256,
		ExpirationDuration: expirationDuration,
		TokenCache:         new(JwtTokenCache),
	}
}

// JwtTokenCache is the cache for the JWT token.
type JwtTokenCache struct {
	Token    string
	ExpireAt time.Time
}

// IsValid checks if the token is valid.
//
// Returns:
//   - error
func (cache *JwtTokenCache) IsValid() error {
	if cache.Token == "" {
		return errors.New("token is empty")
	}

	if time.Now().After(cache.ExpireAt) {
		return errors.New("token is expired")
	}

	return nil
}

// Jwt is the interface for the JWT operations.
type Jwt interface {
	GenerateToken() (string, error)
	ParseToken(token string) (*jwt.Token, error)
	GetToken() (string, error)
}

// jwtImpl is the implementation of the Jwt interface.
type jwtImpl struct {
	Secret             string            // jwt的secret
	Algorithm          jwt.SigningMethod // jwt.SigningMethodHS256
	ExpirationDuration time.Duration     // token过期时长
	TokenCache         *JwtTokenCache    // token缓存
}

// GenerateToken generates a JWT token.
func (j *jwtImpl) GenerateToken() (string, error) {
	now := time.Now()
	expiresAt := now.Add(j.ExpirationDuration).Add(-3 * time.Minute) // 提前3分钟过期
	t := jwt.NewWithClaims(j.Algorithm, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt), // 该项验证
		IssuedAt:  jwt.NewNumericDate(now),       // 该项验证
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "lattice_go",
		Subject:   "jwt",
		ID:        "1",
		Audience:  []string{"somebody_else"},
	})
	token, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}

	j.TokenCache.Token = token
	j.TokenCache.ExpireAt = expiresAt

	return token, nil
}

// ParseToken parses a JWT token.
func (j *jwtImpl) ParseToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("AllYourBase"), nil
	})
	if err != nil {
		return nil, err
	}

	switch {
	case t.Valid:
		return t, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, errors.New("that's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, errors.New("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, errors.New("token is either expired or not active yet")
	default:
		return nil, errors.New("couldn't handle this token")
	}
}

// GetToken returns the cached JWT token or generates a new one if it's expired.
func (j *jwtImpl) GetToken() (string, error) {
	if err := j.TokenCache.IsValid(); err != nil {
		token, err := j.GenerateToken()
		if err != nil {
			return "", err
		}
		return token, nil
	}
	return j.TokenCache.Token, nil
}

// HttpApiInitParam 初始化HTTP API的参数
type HttpApiInitParam struct {
	NodeAddress                string            // 节点的IP:Port
	HttpUrl                    string            // 节点的URL
	GinServerUrl               string            // 节点gin服务路径
	Transport                  http.RoundTripper // tr
	JwtSecret                  string            // jwt的secret信息
	JwtTokenExpirationDuration time.Duration     // jwt token的过期时间
}

// NewHttpApi creates a new HTTP API for the Lattice node.
func NewHttpApi(args *HttpApiInitParam) HttpApi {
	if args.Transport == nil {
		args.Transport = http.DefaultTransport
	}
	return &httpApi{
		HttpAddress:  args.NodeAddress,
		NodeUrl:      args.HttpUrl,
		GinServerUrl: args.GinServerUrl,
		transport:    args.Transport,
		jwtApi:       NewJwt(args.JwtSecret, args.JwtTokenExpirationDuration),
	}
}

// HttpApi This is the interface for the HTTP API of the Lattice node.
type HttpApi interface {
	NewHeaders(chainId string) map[string]string
	GetTransport() http.RoundTripper
	// CanDial 测试是否可以连接到节点
	CanDial(timeout time.Duration) error
	// Forward 转发原始的http请求
	Forward(w http.ResponseWriter, r *http.Request)
	// GetLatestBlock 获取当前账户的最新的区块信息，不包括pending中的交易
	GetLatestBlock(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error)
	// GetLatestBlockWithPending 获取当前账户的最新的区块信息，包括pending中的交易
	GetLatestBlockWithPending(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error)
	// SendSignedTransaction 发送已签名的交易
	SendSignedTransaction(ctx context.Context, chainId string, signedTX *block.Transaction) (*common.Hash, error)
	// SendSignedTransactions batch send transactions
	SendSignedTransactions(ctx context.Context, chainId string, signedTXs []*block.Transaction) ([]*common.Hash, error)
	// PreCallContract 预执行合约
	PreCallContract(ctx context.Context, chainId string, unsignedTX *block.Transaction) (*types.Receipt, error)
	// GetReceipt 获取交易回执
	GetReceipt(ctx context.Context, chainId, hash string) (*types.Receipt, error)
	// GetReceipts 批量查询回执
	GetReceipts(ctx context.Context, chainId string, hashes []string) ([]*types.Receipt, error)
	// GetContractLifecycleProposal 获取合约生命周期提案
	GetContractLifecycleProposal(ctx context.Context, chainId, contractAddress string, state types.ProposalState, startDate, endDate string) ([]types.Proposal[types.ContractLifecycleProposal], error)
	// UploadFile 上传文件到链上
	UploadFile(ctx context.Context, chainId, filePath string) (*types.UploadFileResponse, error)
	// DownloadFile 从链上下载文件
	DownloadFile(ctx context.Context, cid, filePath string) error
	// GetNodeInfo 获取节点信息
	GetNodeInfo(ctx context.Context) (*types.NodeInfo, error)
	// JoinSubchain 让当前节点加入子链（通道）
	//
	// Parameters:
	//   - ctx context.Context
	//   - subchainId uint64：要加入的链
	//   - networkId uint64: 网络ID
	//   - inode string: 已在该链中的节点的INode信息
	JoinSubchain(ctx context.Context, subchainId, networkId uint64, inode string) error
	// StartSubchain 启动子链（通道）
	StartSubchain(ctx context.Context, subchainId string) error
	// StopSubchain 停止子链（通道）
	StopSubchain(ctx context.Context, subchainId string) error
	// DeleteSubchain 删除子链（通道）
	DeleteSubchain(ctx context.Context, subchainId string) error
	// GetSubchain 获取子链的配置信息
	GetSubchain(ctx context.Context, subchainId string) (*types.Subchain, error)
	// GetCreatedSubchain 获取所有通道
	GetCreatedSubchain(ctx context.Context) ([]uint64, error)
	// GetJoinedSubchain 获取已加入通道
	GetJoinedSubchain(ctx context.Context) ([]uint64, error)
	// GetSubchainRunningStatus 获取子链的运行状态
	GetSubchainRunningStatus(ctx context.Context, subchainID string) (*types.SubchainRunningStatus, error)
	// GetSubchainBriefInfo get subchain brief info
	// If subchainID is empty, return all subchains brief info, otherwise return the subchain brief info of the specified subchainID
	GetSubchainBriefInfo(ctx context.Context, subchainID string) ([]*types.SubchainBriefInfo, error)
	// GetConsensusNodesStatus 查询共识节点的状态
	GetConsensusNodesStatus(ctx context.Context, chainID string) ([]*types.ConsensusNodeStatus, error)
	// GetGenesisNodeAddress 获取共识节点地址
	GetGenesisNodeAddress(ctx context.Context, chainID string) (string, error)
	// GetLatestDaemonBlock 获取最新的守护区块信息
	GetLatestDaemonBlock(ctx context.Context, chainID string) (*types.DaemonBlock, error)
	// GetNodePeers 获取在主链上的节点的对等节点信息
	GetNodePeers(ctx context.Context) ([]*types.NodePeer, error)
	// GetSubchainPeers 获取子链（通道）的节点对等节点信息
	GetSubchainPeers(ctx context.Context, subchainId string) (map[string]*types.SubchainPeer, error)
	// GetLatcPeers 查询peer节点
	GetLatcPeers(ctx context.Context, chainId string) (map[string]*types.SubchainPeer, error)
	// GetNodeConfig 查询节点的配置信息
	GetNodeConfig(ctx context.Context, chainID string) (*types.NodeConfiguration, error)
	// GetContractInformation 获取合约的信息
	GetContractInformation(ctx context.Context, chainID, contractAddress string) (*types.ContractInformation, error)
	// GetContractManagement 获取合约管理信息
	GetContractManagement(ctx context.Context, chainID, contractAddress string, daemonBlockHeight *big.Int) (*types.ContractManagement, error)
	// GetDaemonBlockByHash 根据守护区块哈希查询守护区块信息
	GetDaemonBlockByHash(ctx context.Context, chainId, hash string) (*types.DaemonBlock, error)
	// GetDaemonBlockByHeight 根据守护区块高度查询守护区块信息
	GetDaemonBlockByHeight(ctx context.Context, chainId string, height uint64) (*types.DaemonBlock, error)
	// ExistsBusinessContractAddress 检查存证的业务合约地址是否存在
	ExistsBusinessContractAddress(ctx context.Context, chainId, address string) (bool, error)
	// GetTransactionBlockByHash 根据哈希查询交易区块的信息
	GetTransactionBlockByHash(ctx context.Context, chainId, hash string) (*types.TransactionBlock, error)
	// GetNodeProtocol 获取节点的网络协议信息
	GetNodeProtocol(ctx context.Context, chainId string) (*types.NodeProtocol, error)
	// GetVoteById 查询投票详情
	GetVoteById(ctx context.Context, chainId, voteId string) (*types.VoteDetails, error)
	// GetProposal 查询提案
	GetProposal(ctx context.Context, chainId, proposalId string, ty types.ProposalType, state types.ProposalState, proposalAddress, contractAddress, startDate, endDate string, result interface{}) error
	// GetRawProposal 查询提案
	GetRawProposal(ctx context.Context, chainId, proposalId string, ty types.ProposalType, state types.ProposalState, proposalAddress, contractAddress, startDate, endDate string) (json.RawMessage, error)
	// GetTransactionsPagination 根据守护区块高度分页查询交易
	GetTransactionsPagination(ctx context.Context, chainId string, startDaemonBlockHeight uint64, pageSize uint16) (*types.TransactionsPagination, error)
	// GetEvidences 获取留痕信息，date string: 格式20240511
	GetEvidences(ctx context.Context, chainId, date string, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error)
	// GetErrorEvidences 获取错误留痕
	//
	// Parameters:
	//   - chainId string
	//   - date string: 格式20240511
	//   - level string: 日志级别, "none":执行日志, "error":error级别的错误日志, "crit":crit级别的错误日志, 不填则默认为执行日志
	//   - eviType string: 日志类型, "vote":投票, "tblock":账户交易, "dblock":守护区块, "sign":签名, "pre":预执行合约, "onChain":发布合约交易, "execute":执行合约交易, "update":合约升级, "upgrade":升级合约的账户交易, "deploy":合约部署, "call":合约调用, "revoke":合约吊销, "freeze":合约冻结, "release":合约解冻, "error":error错误, "crit":crit错误, "add":增加账户, "del":删除账户, "lock":锁定账户, "unlock":解锁账户, "oracle":预言机
	GetErrorEvidences(ctx context.Context, chainId, date string, evidenceLevel types.EvidenceLevel, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error)
	// GetNodeConfirmedConfiguration 获取节点确认的配置信息
	GetNodeConfirmedConfiguration(ctx context.Context, chainId string) (*types.NodeConfirmedConfiguration, error)
	// GetNodeVersion 获取节点程序的版本信息
	GetNodeVersion(ctx context.Context) (*types.NodeVersion, error)
	// GetNodeSaintKey 获取节点的SaintKey信息
	GetNodeSaintKey(ctx context.Context) (*wallet.FileKey, error)
	// GetNodeConfiguration 获取节点的配置
	GetNodeConfiguration(ctx context.Context) (*types.NodeConfiguration, error)
	LoadNodeConfiguration(ctx context.Context, chainId string) (*types.NodeConfiguration, error)
	// GetNodeWorkingDirectory 获取节点的工作目录（绝对路径）
	GetNodeWorkingDirectory(ctx context.Context) (string, error)
	// GetSnapshot 获取快照
	GetSnapshot(ctx context.Context, chainId string, daemonBlockHeight *big.Int) (*types.NodeProtocolConfig, error)
	// GetLatcInfo 获取协议配置信息
	GetLatcInfo(ctx context.Context, chainId string) (*types.NodeProtocolConfig, error)
	// GetProposalById 根据提案ID查询提案
	GetProposalById(ctx context.Context, chainId, proposalId string, result interface{}) error
	// GetSubchainIdByProposalId get subchain id by proposal id
	GetSubchainIdByProposalId(ctx context.Context, chainId, proposalId string) (uint32, error)
	// Freeze 立即压缩区块
	Freeze(ctx context.Context, chainId string, dblockNumber *big.Int) (uint64, error)
	// GetFreezeDBlockByHash 通过区块哈希查询压缩的区块
	GetFreezeDBlockByHash(ctx context.Context, chainId string, hash string) (*types.DaemonBlock, error)
	// GetFreezeDBlockByNumber 通过区块高度查询压缩的区块
	GetFreezeDBlockByNumber(ctx context.Context, chainId string, dblockNumber *big.Int) (*types.DaemonBlock, error)
	// GetFreezeTBlockByHash 通过区块哈希查询压缩的交易
	GetFreezeTBlockByHash(ctx context.Context, chainId string, hash string) (*types.TransactionBlock, error)
	// GetFreezeTBlockByNumber 通过账户地址和区块高度查询压缩的交易
	GetFreezeTBlockByNumber(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.TransactionBlock, error)
	// GetFreezeReceipt 通过账户地址和区块高度查询压缩的收据
	GetFreezeReceipt(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.Receipt, error)
	// GetFreezeSaveSpace 获取当前区块压缩已节省的空间
	GetFreezeSaveSpace(ctx context.Context, chainId string) (*types.FreezeSaveSpace, error)
	// GetFreezeInterval 获取区块压缩的时间间隔
	GetFreezeInterval(ctx context.Context, chainId string) (*types.FreezeInterval, error)
	// ImportFileKey interface to the wallet module, import the file key to the blockchain node.
	// string, the address of the account
	ImportFileKey(ctx context.Context, fileKey string) (string, error)
	// ImportRawKey interface to the wallet module, import the account to the blockchain node.
	ImportRawKey(ctx context.Context, privateKey, password string) (bool, error)
	// GetAccounts get accounts from the node
	GetAccounts(ctx context.Context, chainId string) ([]string, error)
	// GetTBlockState 获取账户区块状态
	GetTBlockState(ctx context.Context, chainId, hash string) (types.TBlockState, error)
	// GetElapsed 获取节点各关键操作耗时统计
	GetElapsed(ctx context.Context) (map[string]int64, error)
	// GetNodeCertificate 查看节点数字证书
	GetNodeCertificate(ctx context.Context) (*types.NodeCertificate, error)
	// GetPeerNodeCertificate 根据SerialNumber查询节点数字证书
	GetPeerNodeCertificate(ctx context.Context, serialNumber string) (*types.NodeCertificate, error)
	// GetPeerNodeCertificateByAddress 根据节点地址查询节点数字证书
	GetPeerNodeCertificateByAddress(ctx context.Context, nodeAddress string) (*types.NodeCertificate, error)
	// GetConsensus 获取节点前置共识和后置共识
	GetConsensus(ctx context.Context) (*types.LatcConsensus, error)
	// GetSyncStatus 获取节点同步状态
	GetSyncStatus(ctx context.Context) (*types.SyncStatus, error)
	// GetLastBatchDBlockNumber 获取最近一次并行处理的守护区块高度
	GetLastBatchDBlockNumber(ctx context.Context, chainId string) (*big.Int, error)
	// ConnectNodeAsync 连接新的节点，异步接口，是否连接成功需要借助node_peers判断
	ConnectNodeAsync(ctx context.Context, inode string) error
	// ConnectPeerAsync 连接新的链的对等节点，异步接口，是否连接成功需要借助latc_peers 判断
	ConnectPeerAsync(ctx context.Context, nodeHash string) error
	// DisconnectPeerAsync 断开对等节点，异步接口，是否断开连接成功需要借助latc_peers 判断
	DisconnectPeerAsync(ctx context.Context, nodeHash string) error
	// GetDBlockProof 获取守护区块签名
	GetDBlockProof(ctx context.Context, chainId string, dblockNumber *big.Int) (*types.WitnessProof, error)
	// GetTBlockProof 获取交易区块签名
	GetTBlockProof(ctx context.Context, chainId string, accountAddress string, tblockNumber *big.Int) (*types.WitnessProof, error)
	// GetCurrentTBlock 获取当前账户最新的区块信息
	GetCurrentTBlock(ctx context.Context, chainId string, accountAddress string) (*types.TransactionBlock, error)
	// GetTBlockByHeight 查询当前账户指定高度的区块信息
	GetTBlockByHeight(ctx context.Context, chainId, accountAddress string, height uint64) (*types.TransactionBlock, error)
	// GetBalanceWithPending 获取账户余额
	GetBalanceWithPending(ctx context.Context, chainId, accountAddress string) (*types.AccountBalance, error)
	// GetGenesisBlock 获取创世区块信息
	GetGenesisBlock(ctx context.Context, chainId string) (*types.TransactionBlock, error)
	// GetTBlocksByHeights 根据账户高度区间查询交易
	GetTBlocksByHeights(ctx context.Context, chainId string, accountAddress string, heights []uint64) ([]*types.TransactionBlock, error)
	// GetDBlocksByHeights 根据守护区块高度区间查询守护区块
	GetDBlocksByHeights(ctx context.Context, chainId string, heights []uint64) ([]*types.DaemonBlock, error)
	ProxyReEncryption(ctx context.Context, chainId string, ciphertext, businessAddress, initiator, whitelist string) (string, error)
	GetRecentDBlocks(ctx context.Context, chainId string, limit uint32) ([]*types.DaemonBlock, error)
	GetTBlockCount(ctx context.Context, chainId string) (*types.TBlockCount, error)
	ImportCertificate(ctx context.Context, chainId string, pemCertificate string) error
	// PublishCertificates 根据节点的公钥发布证书，返回证书序列号
	PublishCertificates(ctx context.Context, chainId string, publicKeys []string) ([]string, error)
	GetCertificate(ctx context.Context, chainId string, serialNumber string) (*types.NodeCertificate, error)
}

type httpApi struct {
	HttpAddress  string            // 节点的IP:Port
	NodeUrl      string            // 节点的Http请求路径
	GinServerUrl string            // 节点的Gin服务请求路径
	transport    http.RoundTripper // http transport
	jwtApi       Jwt               // jwt api
}

func (api *httpApi) forwardErrorHandler(rw http.ResponseWriter, _ *http.Request, err error) {
	rw.WriteHeader(http.StatusBadGateway)
	rw.Write([]byte(err.Error()))
}

// Forward : 透传原始http请求到链上
func (api *httpApi) Forward(rw http.ResponseWriter, r *http.Request) {
	nodeUrl, err := url.Parse(api.NodeUrl)
	if err != nil {
		api.forwardErrorHandler(rw, r, fmt.Errorf("failed to parse node url %s, err=%v", api.NodeUrl, err))
		return
	}
	log.Debug().Msgf("正在转发http请求到链上，转发%s到%s", r.URL, nodeUrl)
	headers := api.newHeaders("")
	for k, v := range headers {
		if r.Header.Get(k) == "" && v != "" {
			r.Header.Set(k, v)
		}
	}
	proxy := httputil.NewSingleHostReverseProxy(nodeUrl)
	proxy.Transport = api.transport
	proxy.ErrorHandler = api.forwardErrorHandler
	proxy.ServeHTTP(rw, r)
}

const (
	headerContentType = "Content-Type"
	headerChainID     = "ChainId"
	headerAuthorize   = "Authorization"
	headerConnection  = "Connection"
)

// 设置http的请求头
//
// Parameters:
//   - chainId string
//
// Returns:
//   - map[string]string
func (api *httpApi) newHeaders(chainId string) map[string]string {
	headers := map[string]string{
		headerContentType: "application/json",
		headerChainID:     chainId,
	}
	if api.jwtApi != nil {
		token, _ := api.jwtApi.GetToken()
		headers[headerAuthorize] = fmt.Sprintf("Bearer %s", token)
	}
	return headers
}

func (api *httpApi) NewHeaders(chainId string) map[string]string {
	return api.newHeaders(chainId)
}

func (api *httpApi) GetTransport() http.RoundTripper {
	return api.transport
}

func (api *httpApi) CanDial(timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", api.HttpAddress, timeout)
	if err != nil {
		return err
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func (api *httpApi) SendSignedTransaction(ctx context.Context, chainId string, signedTX *block.Transaction) (*common.Hash, error) {
	response, err := Post[common.Hash](ctx, api.NodeUrl, NewJsonRpcBody("wallet_sendRawTBlock", signedTX), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) SendSignedTransactions(ctx context.Context, chainId string, signedTXs []*block.Transaction) ([]*common.Hash, error) {
	response, err := Post[[]*common.Hash](ctx, api.NodeUrl, NewJsonRpcBody("wallet_sendRawBatchTBlock", signedTXs), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) PreCallContract(ctx context.Context, chainId string, unsignedTX *block.Transaction) (*types.Receipt, error) {
	response, err := Post[types.Receipt](ctx, api.NodeUrl, NewJsonRpcBody("wallet_preExecuteContract", unsignedTX), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetReceipt(ctx context.Context, chainId, hash string) (*types.Receipt, error) {
	response, err := Post[types.Receipt](ctx, api.NodeUrl, NewJsonRpcBody("latc_getReceipt", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetReceipts(ctx context.Context, chainId string, hashes []string) ([]*types.Receipt, error) {
	response, err := Post[[]*types.Receipt](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockReceipts", hashes), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) ExistsBusinessContractAddress(ctx context.Context, chainId, address string) (bool, error) {
	response, err := Post[bool](ctx, api.NodeUrl, NewJsonRpcBody("wallet_confirmTaggedContract", address), api.newHeaders(chainId), api.transport)
	if err != nil {
		return false, nil
	}
	if response.Error != nil {
		return false, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetEvidences(ctx context.Context, chainId, date string, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error) {
	response, err := Post[types.Evidences](ctx, api.NodeUrl, NewJsonRpcBody("latc_getEvidences", date, evidenceType, page, pageSize), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetErrorEvidences(ctx context.Context, chainId, date string, evidenceLevel types.EvidenceLevel, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error) {
	response, err := Post[types.Evidences](ctx, api.NodeUrl, NewJsonRpcBody("latc_getErrorEvidences", date, evidenceLevel, evidenceType, page, pageSize), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetElapsed(ctx context.Context) (map[string]int64, error) {
	response, err := Post[map[string]int64](ctx, api.NodeUrl, NewJsonRpcBody("latc_getElapsed"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

// GetConsensus implements HttpApi.
func (api *httpApi) GetConsensus(ctx context.Context) (*types.LatcConsensus, error) {
	response, err := Post[types.LatcConsensus](ctx, api.NodeUrl, NewJsonRpcBody("latc_getConsensus"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// GetSyncStatus implements HttpApi.
func (api *httpApi) GetSyncStatus(ctx context.Context) (*types.SyncStatus, error) {
	response, err := Post[types.SyncStatus](ctx, api.NodeUrl, NewJsonRpcBody("sync_status"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// Post send http request use post method
//
// Parameters:
//   - ctx context.Context: 超时取消
//   - url string: 请求路径，示例：http://192.168.1.20:13000
//   - body sonRpcBody: any, 请求体
//   - headers map[string]string: 请求头
//   - tr http.Transport:
//
// Returns:
//   - []byte: 响应内容
//   - error: 错误
func Post[T any](ctx context.Context, url string, jsonRpcBody *JsonRpcBody, headers map[string]string, tr http.RoundTripper) (*JsonRpcResponse[*T], error) {
	response, err := rawPost(ctx, url, jsonRpcBody, headers, tr)
	if err != nil {
		return nil, err
	}

	var t JsonRpcResponse[*T]
	if err = json.Unmarshal(response, &t); err != nil {
		log.Error().Err(err).Msg("Failed to unmarshal response body")
		return nil, err
	}

	return &t, nil
}

func rawPost(ctx context.Context, baseUrl string, jsonRpcBody *JsonRpcBody, headers map[string]string, tr http.RoundTripper) ([]byte, error) {
	log.Debug().Msgf("开始发送JsonRpc请求，url: %s, body: %+v", baseUrl, jsonRpcBody)
	bodyBytes, err := json.Marshal(jsonRpcBody)
	if err != nil {
		return nil, err
	}
	body := strings.NewReader(string(bodyBytes))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, baseUrl, body)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create http request")
		return nil, err
	}

	if len(headers) != 0 {
		for key, value := range headers {
			request.Header.Set(key, value)
		}
	}
	client := &http.Client{Transport: tr}
	request.TransferEncoding = []string{}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send http request")
		return nil, fmt.Errorf("请求节点失败：%s", err)
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close response body")
		}
	}(response.Body)

	if res, err := io.ReadAll(response.Body); err != nil {
		log.Error().Err(err).Msg("Failed to read response body")
		return nil, err
	} else {
		return res, nil
	}
}
