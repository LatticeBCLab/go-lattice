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

// This is the interface for the HTTP API of the Lattice node.
type HttpApi interface {
	// CanDial 测试是否可以连接到节点
	//
	// Parameters:
	//   - timeout time.Duration
	//
	// Returns:
	//   - error
	CanDial(timeout time.Duration) error

	// GetLatestBlock 获取当前账户的最新的区块信息，不包括pending中的交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainId string
	//   - accountAddress string: 账户地址，zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi
	//
	// Returns:
	//   - types.LatestBlock
	//   - error
	GetLatestBlock(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error)

	// GetLatestBlockWithPending 获取当前账户的最新的区块信息，包括pending中的交易
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainId string
	//   - accountAddress string: 账户地址，zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi
	//
	// Returns:
	//   - types.LatestBlock
	//   - error
	GetLatestBlockWithPending(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error)

	// SendSignedTransaction 发送已签名的交易
	//
	// Parameters:
	//    - ctx context.Context
	//    - chainId string
	//    - signedTX *block.Transaction
	//
	// Returns:
	//    - o
	SendSignedTransaction(ctx context.Context, chainId string, signedTX *block.Transaction) (*common.Hash, error)

	// PreCallContract 预执行合约
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainId string
	//   - unsignedTX *block.Transaction: 未签名的交易
	//
	// Returns:
	//   - *types.Receipt
	//   - error
	PreCallContract(ctx context.Context, chainId string, unsignedTX *block.Transaction) (*types.Receipt, error)

	// GetReceipt 获取交易回执
	//
	// Parameters:
	//    - ctx context.Context
	//    - chainId string
	//    - hash string
	//
	// Returns:
	//    - types.Receipt
	//    - error
	GetReceipt(ctx context.Context, chainId, hash string) (*types.Receipt, error)

	// GetContractLifecycleProposal 获取合约生命周期提案
	//
	// Parameters:
	//    - ctx context.Context
	//    - chainId string: 链ID
	//    - contractAddress string: 合约地址
	//    - state types.ProposalState: 提案状态
	//
	// Returns:
	//    - types.Proposal[types.ContractLifecycleProposal]
	//    - error
	GetContractLifecycleProposal(ctx context.Context, chainId, contractAddress string, state types.ProposalState) ([]types.Proposal[types.ContractLifecycleProposal], error)

	// UploadFile 上传文件到链上
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainId string: 链ID
	//   - filePath string: 文件路径
	//
	// Returns:
	//   - *types.UploadFileResponse
	//   - error
	UploadFile(ctx context.Context, chainId, filePath string) (*types.UploadFileResponse, error)

	// DownloadFile 从链上下载文件
	//
	// Parameters:
	//   - ctx context.Context
	//   - cid string: 要下载文件的cid
	//   - filePath string: 指定的临时存储路径
	//
	// Returns:
	//   - error
	DownloadFile(ctx context.Context, cid, filePath string) error

	// GetNodeInfo 获取节点信息
	//
	// Parameters:
	//   - ctx context.Context
	//
	// Returns:
	//   - *types.NodeInfo,
	//   - error
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
	//
	// Returns:
	//   - []*types.ConsensusNodeStatus
	//   - error
	GetConsensusNodesStatus(ctx context.Context, chainID string) ([]*types.ConsensusNodeStatus, error)

	// GetGenesisNodeAddress 获取共识节点地址
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainID string
	//
	// Returns:
	//   - string
	//   - error
	GetGenesisNodeAddress(ctx context.Context, chainID string) (string, error)

	// GetLatestDaemonBlock 获取最新的守护区块信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainID string
	//
	// Returns:
	//   - *types.DaemonBlock
	//   - error
	GetLatestDaemonBlock(ctx context.Context, chainID string) (*types.DaemonBlock, error)

	// GetNodePeers 获取在主链上的节点的对等节点信息
	//
	// Parameters:
	//   - ctx context.Context
	//
	// Returns:
	//   - []*types.NodePeer
	//   - error
	GetNodePeers(ctx context.Context) ([]*types.NodePeer, error)

	// GetSubchainPeers 获取子链（通道）的节点对等节点信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - subchainId string
	//
	// Returns:
	//   - map[string]*types.SubchainPeer
	//   - error
	GetSubchainPeers(ctx context.Context, subchainId string) (map[string]*types.SubchainPeer, error)

	// GetNodeConfig 查询节点的配置信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainID string
	//
	// Returns:
	//   - *types.NodeConfig
	//   - error
	GetNodeConfig(ctx context.Context, chainID string) (*types.NodeConfig, error)

	// GetContractInformation 获取合约的信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainID string
	//   - contractAddress string
	//
	// Returns:
	//   - *types.ContractInformation
	//   - error
	GetContractInformation(ctx context.Context, chainID, contractAddress string) (*types.ContractInformation, error)

	// GetContractManagement 获取合约管理信息
	//
	// Parameters:
	//   - ctx context.Context
	//   - chainID string: 合约所在链的ID
	//   - contractAddress string: 合约地址
	//
	// Returns:
	//   - *types.ContractManagement
	//   - error
	GetContractManagement(ctx context.Context, chainID, contractAddress string, daemonBlockHeight *big.Int) (*types.ContractManagement, error)

	// GetDaemonBlockByHash 根据守护区块哈希查询守护区块信息
	GetDaemonBlockByHash(ctx context.Context, chainId, hash string) (*types.DaemonBlock, error)

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

	// GetEvidences 获取留痕信息
	//
	// Parameters:
	//   - chainId string
	//   - date string: 格式20240511
	//   - evidenceType string
	//   - page int
	//   - pageSize int
	//
	// Returns:
	//   - o
	GetEvidences(ctx context.Context, chainId, date string, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error)

	// GetErrorEvidences 获取错误留痕
	//
	// Parameters:
	//   - chainId string
	//   - date string: 格式20240511
	//   - level string: 日志级别, "none":执行日志, "error":error级别的错误日志, "crit":crit级别的错误日志, 不填则默认为执行日志
	//   - eviType string: 日志类型, "vote":投票, "tblock":账户交易, "dblock":守护区块, "sign":签名, "pre":预执行合约, "onChain":发布合约交易, "execute":执行合约交易, "update":合约升级, "upgrade":升级合约的账户交易, "deploy":合约部署, "call":合约调用, "revoke":合约吊销, "freeze":合约冻结, "release":合约解冻, "error":error错误, "crit":crit错误, "add":增加账户, "del":删除账户, "lock":锁定账户, "unlock":解锁账户, "oracle":预言机
	//   - page int
	//   - pageSize int
	//
	// Returns:
	//   - o
	GetErrorEvidences(ctx context.Context, chainId, date string, evidenceLevel types.EvidenceLevel, evidenceType types.EvidenceType, page, pageSize int) (*types.Evidences, error)
	// GetNodeConfirmedConfiguration 获取节点确认的配置信息
	GetNodeConfirmedConfiguration(ctx context.Context, chainId string) (*types.NodeConfirmedConfiguration, error)
	// GetNodeVersion 获取节点程序的版本信息
	GetNodeVersion(ctx context.Context) (*types.NodeVersion, error)
	// GetNodeSaintKey 获取节点的SaintKey信息
	GetNodeSaintKey(ctx context.Context) (*wallet.FileKey, error)
	// GetNodeConfiguration 获取节点的配置
	GetNodeConfiguration(ctx context.Context) (*types.NodeConfiguration, error)
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
	Freeze(ctx context.Context, dblockNumber *big.Int) (uint64, error)
	// GetFreezeDBlockByHash 通过区块哈希查询压缩的区块
	GetFreezeDBlockByHash(ctx context.Context, chainId string, hash string) (*types.DaemonBlock, error)
	// GetFreezeDBlockByHash 通过区块高度查询压缩的区块
	GetFreezeDBlockByNumber(ctx context.Context, chainId string, dblockNumber *big.Int) (*types.DaemonBlock, error)
	// GetFreezeTBlockByHash 通过区块哈希查询压缩的交易
	GetFreezeTBlockByHash(ctx context.Context, chainId string, hash string) (*types.TransactionBlock, error)
	// GetFreezeTBlockByNumber 通过账户地址和区块高度查询压缩的交易
	GetFreezeTBlockByNumber(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.TransactionBlock, error)
	// GetFreezeReceipt 通过账户地址和区块高度查询压缩的收据
	GetFreezeReceipt(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.Receipt, error)
	// GetFreezeSaveSpace 获取当前区块压缩已节省的空间
	GetFreezeSaveSpace(ctx context.Context) (*types.FreezeSaveSpace, error)
	// GetFreezeInterval 获取区块压缩的时间间隔
	GetFreezeInterval(ctx context.Context) (*types.FreezeInterval, error)
}

type httpApi struct {
	HttpAddress  string            // 节点的IP:Port
	NodeUrl      string            // 节点的Http请求路径
	GinServerUrl string            // 节点的Gin服务请求路径
	transport    http.RoundTripper // http transport
	jwtApi       Jwt               // jwt api
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
	if err := json.Unmarshal(response, &t); err != nil {
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
