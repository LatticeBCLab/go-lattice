package types

import (
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/LatticeBCLab/go-lattice/common/constant"
)

const (
	AddressVersion = 1
	AddressLength  = 20 // 20 byte
	AddressTitle   = "zltc"
	HashLength     = 32 // 32 byte
)

type Number string

func (n Number) MustToBigInt() *big.Int {
	num := new(big.Int)
	num.SetString(string(n), 10)
	return num
}

// UploadFileResponse 文件上传到链上的返回结果
//
//   - CID 文件唯一标识，示例：GPK3PveRaWoK6S2b53D3ZeJTm4nBvv2vSVjRStRLQcyX
//   - FilePath 文件存储地址，示例：JG-DFS/tempFileDir/20240816/1723793768943848748_avatar.svg"
//   - Message 返回信息，示例：success
//   - OccupiedStorageByte 文件占用的存储字节数，单位为byte，示例：255686
//   - StorageAddress 需要冗余存储文件的节点地址，示例：DFS_beforeSign||zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi;zltc_nBGgKoo1rzd4thjfauEN6ULj7jp1zXhxE;
type UploadFileResponse struct {
	CID                 string `json:"cid,omitempty"`
	FilePath            string `json:"filePath,omitempty"`
	Message             string `json:"message,omitempty"`
	OccupiedStorageByte int64  `json:"needStorageSize,omitempty"`
	StorageAddress      string `json:"storageAddress,omitempty"`
}

// NodeInfo 节点
//
//   - ID 示例：16Uiu2HAmQ7Da6iuScYSYs8XGJs95hiKdS6tgmbqUUuKC62Xh3s4V
//   - Name 示例：ZLTC2_1
//   - Version
//   - INode 节点连接信息，示例：/ip4/192.168.1.185/tcp/13801/p2p/16Uiu2HAmQ7Da6iuScYSYs8XGJs95hiKdS6tgmbqUUuKC62Xh3s4V
//   - Inr
//   - IP
//   - Ports
//   - ListenAddress
type NodeInfo struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
	INode   string `json:"inode,omitempty"`
	Inr     string `json:"inr,omitempty"`
	IP      string `json:"ip,omitempty"`
	Ports   struct {
		P2PPort       uint16 `json:"discovery,omitempty"`
		WebsocketPort uint16 `json:"listener,omitempty"`
		HTTPPort      uint16 `json:"httpPort,omitempty"`
	}
	ListenAddress string `json:"listenAddr,omitempty"`
}

// Subchain 子链信息
//   - ID                              链ID
//   - Name                            名称
//   - Desc                            描述
//   - LatcGodAddr                     守护链的地址
//   - LatcSaints                      共识节点列表
//   - Consensus                       共识
//   - Epoch                           重置投票和检查点的纪元长度
//   - Tokenless                       false:有通证 true:无通证
//   - Period                          出块间隔
//   - EnableNoTxDelayedMining         是否不允许无交易时快速出空块，无交易时延迟出块
//   - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
//   - IsGM                            是否使用了Sm2p256v1曲线
//   - RootPublicKey                   中心化CA根证书公钥
//   - EnableContractLifecycle         是否开启合约生命周期
//   - EnableVotingDictatorship        是否开启投票(合约生命周期)时盟主一票制度
//   - ContractDeploymentVotingRule    合约部署的投票规则
//   - EnableContractManagement        是否开启合约管理
//   - ChainByChainVotingRule          以链建链投票规则
//   - ProposalExpirationDays          提案的过期天数，默认7天
//   - ConfigurationModifyVotingRule   配置修改的投票规则
type Subchain struct {
	ID                              uint64     `json:"latcId,omitempty"`
	Name                            string     `json:"name,omitempty"`
	Desc                            string     `json:"desc,omitempty"`
	LatcGodAddr                     string     `json:"latcGod,omitempty"`
	LatcSaints                      []string   `json:"latcSaints,omitempty"`
	Consensus                       string     `json:"consensus,omitempty"`
	Epoch                           uint       `json:"epoch,omitempty"`
	Tokenless                       bool       `json:"tokenless,omitempty"`
	Period                          uint       `json:"period,omitempty"`
	EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor,omitempty"`
	NoTxDelayedMiningPeriodMultiple uint32     `json:"emptyAnchorPeriodMul,omitempty"`
	IsGM                            bool       `json:"GM,omitempty"`
	RootPublicKey                   string     `json:"rootPublicKey,omitempty"`
	EnableContractLifecycle         bool       `json:"isContractVote,omitempty"`
	EnableVotingDictatorship        bool       `json:"isDictatorship,omitempty"`
	ContractDeploymentVotingRule    VotingRule `json:"deployRule,omitempty"`
	EnableContractManagement        bool       `json:"contractPermission,omitempty"`
	ChainByChainVotingRule          VotingRule `json:"chainByChainVote,omitempty"`
	ProposalExpirationDays          uint       `json:"ProposalExpireTime,omitempty"`
	ConfigurationModifyVotingRule   VotingRule `json:"configModifyRule,omitempty"`
}

// SubchainBriefInfo 子链的简要信息
type SubchainBriefInfo struct {
	ID          uint64 `json:"LatcID,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Desc,omitempty"`
}

// ConsensusNodeStatus 共识节点的状态
//   - Address						   节点地址
//   - WitnessedBlockCount             见证区块数量
//   - FailureWitnessedBlockCount 	   见证失败的区块数量
//   - ShouldWitnessedBlockCount	   应当见证的区块数量
type ConsensusNodeStatus struct {
	Address                    string `json:"Addr,omitempty"`
	WitnessedBlockCount        uint64 `json:"SignatureCount,omitempty"`
	FailureWitnessedBlockCount uint64 `json:"SignatureFailCount,omitempty"`
	ShouldWitnessedBlockCount  uint64 `json:"ShouldSignatureCount,omitempty"`
	Online                     bool   `json:"online,omitempty"`
}

// NodePeer peer
type NodePeer struct {
	INode     string            `json:"inode"`
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Caps      []string          `json:"caps"`
	Network   NodePeerNetwork   `json:"network"`
	Protocols NodePeerProtocols `json:"protocols"`
}

type NodePeerNetwork struct {
	LocalAddress  string `json:"localAddress"`
	RemoteAddress string `json:"remoteAddress"`
	Inbound       bool   `json:"inbound"`
	Trusted       bool   `json:"trusted"`
	Static        bool   `json:"static"`
}

type NodePeerProtocols struct {
	Latc string `json:"latc"`
}

// SubchainPeer 子链（通道）的节点对等节点信息
//
//   - Id
//   - Saint			 节点地址
//   - ChainId			 链ID
//   - IsBootstrap       节点是否为种子节点
//   - CertificateSN     节点证书序列号
//   - DaemonBlockHeight 守护区块高度
//   - DaemonBlockHash	 守护区块哈希
//   - GinHttpPort       文件传输服务的端口
//   - DFS 		   		 是否使用文件传输服务
//   - INode       		 节点的INode信息
type SubchainPeer struct {
	Id                string   `json:"Id,omitempty"`
	Saint             string   `json:"Saint,omitempty"`
	ChainId           uint32   `json:"ChainId,omitempty"`
	IsBootstrap       bool     `json:"IsBoot,omitempty"`
	CertificateSN     *big.Int `json:"CertSerialNumber,omitempty"`
	DaemonBlockHeight *big.Int `json:"DBNum,omitempty"`
	DaemonBlockHash   string   `json:"DBHash,omitempty"`
	GinHttpPort       int      `json:"GinHttpPort,omitempty"`
	DFS               bool     `json:"DFS,omitempty"`
	INode             string   `json:"inode,omitempty"`
}

// ContractInformation 合约信息
//
//   - ContractAddress 合约地址
//   - Owner           合约的部署者地址
//   - State           合约的状态
//   - Version         合约的版本
//   - ProposalId      合约的提案ID，包括 部署、升级、吊销
//   - CreatedAt       合约的部署时间戳(s)
//   - UpdatedAt       合约的修改时间戳(s)
type ContractInformation struct {
	ContractAddress string        `json:"address"`
	Owner           string        `json:"deploymentAddress"`
	State           ContractState `json:"state"`
	Version         uint8         `json:"version"`
	ProposalId      string        `json:"votingProposalId,omitempty"`
	CreatedAt       int64         `json:"createAt"`
	UpdatedAt       int64         `json:"modifiedAt"`
}

// ContractManagement 合约管理信息
//
//   - Mode           合约管理模式，白名单 or 黑名单
//   - Threshold      投票通过的阈值，大于10则按照权重加和，小于等于10则按照百分比
//   - Whitelist	  合约白名单
//   - Blacklist	  合约黑名单
//   - Administrators 合约管理员：`{"zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi": 10}`
type ContractManagement struct {
	Mode           ContractManagementMode `json:"permissionMode"`
	Threshold      uint64                 `json:"threshold"`
	Whitelist      []string               `json:"whiteList"`
	Blacklist      []string               `json:"blackList"`
	Administrators map[string]uint8       `json:"managerList"`
}

// DeployMultilingualContractCode 部署多语言智能合约的代码
//   - FileName 上传到链上的合约文件名
type DeployMultilingualContractCode struct {
	FileName string `json:"contractName,omitempty"`
}

// UpgradeMultilingualContractCode 升级多语言智能合约的代码
//   - FileName 上传到链上的合约文件名
type UpgradeMultilingualContractCode struct {
	FileName string `json:"contractName,omitempty"`
}

// CallMultilingualContractCode 调用多语言智能合约的代码
//   - Method	 调用的合约方法名，示例：`double`
//   - Arguments 调用的合约方法参数，示例：`{"number":[56,50,55]}`
type CallMultilingualContractCode struct {
	Method    string            `json:"methodName,omitempty"`
	Arguments map[string][]byte `json:"methodArgs,omitempty"`
}

func (c *DeployMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

func (c *UpgradeMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

func (c *CallMultilingualContractCode) Encode() string {
	if bytes, err := json.Marshal(c); err != nil {
		return constant.HexPrefix
	} else {
		return constant.HexPrefix + hex.EncodeToString(bytes)
	}
}

// NodeProtocol 节点的网络协议信息
type NodeProtocol struct {
	Genesis        string              `json:"genesis"`
	NetWorkIdGroup []int               `json:"netWorkIdGroup"`
	Config         *NodeProtocolConfig `json:"config"`
}

// NodeProtocolConfig
//   - LatcID                          链ID
//   - Name                            名称
//   - Desc                            描述
//   - LatcGodAddr                     守护链的地址
//   - LatcSaints                      共识节点列表
//   - Consensus                       共识
//   - Epoch                           重置投票和检查点的纪元长度
//   - Tokenless                       false:有通证 true:无通证
//   - Period                          出块间隔
//   - EnableNoTxDelayedMining         是否不允许无交易时快速出空块，无交易时延迟出块
//   - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
//   - IsGM                            是否使用了Sm2p256v1曲线
//   - RootPublicKey                   中心化CA根证书公钥
//   - EnableContractLifecycle         是否开启合约生命周期
//   - EnableVotingDictatorship        是否开启投票(合约生命周期)时盟主一票制度
//   - ContractDeploymentVotingRule    合约部署的投票规则
//   - EnableContractManagement        是否开启合约管理
//   - ChainByChainVotingRule          以链建链投票规则
//   - ProposalExpirationDays          提案的过期天数，默认7天
//   - ConfigurationModifyVotingRule   对于配置修改的投票规则来说，0、2都是共识投票，1是盟主一票制
type NodeProtocolConfig struct {
	LatcID                          *big.Int   `json:"latcId,omitempty"`
	Name                            string     `json:"Name,omitempty"`
	Desc                            string     `json:"Desc,omitempty"`
	LatcGodAddr                     string     `json:"latcGod,omitempty"`
	LatcSaints                      []string   `json:"latcSaints,omitempty"`
	Consensus                       string     `json:"consensus,omitempty"`
	Epoch                           uint       `json:"epoch,omitempty"`
	Tokenless                       bool       `json:"tokenless,omitempty"`
	Period                          uint       `json:"period,omitempty"`
	EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor,omitempty"`
	NoTxDelayedMiningPeriodMultiple uint64     `json:"emptyAnchorPeriodMul,omitempty"`
	IsGM                            bool       `json:"GM,omitempty"`
	RootPublicKey                   string     `json:"rootPublicKey,omitempty"`
	EnableContractLifecycle         bool       `json:"isContractVote,omitempty"`
	EnableVotingDictatorship        bool       `json:"isDictatorship,omitempty"`
	ContractDeploymentVotingRule    VotingRule `json:"deployRule,omitempty"`
	EnableContractManagement        bool       `json:"contractPermission,omitempty"`
	ChainByChainVotingRule          VotingRule `json:"chainByChainVote,omitempty"`
	ProposalExpirationDays          uint       `json:"ProposalExpireTime,omitempty"`
	ConfigurationModifyVotingRule   VotingRule `json:"configModifyRule,omitempty"`
}

// NodeConfirmedConfiguration 节点确认的配置信息
//   - LatcId
//   - LatcGod
//   - LatcSaints
//   - Consensus
//   - Epoch
//   - Tokenless
//   - Period
//   - EnableNoTxDelayedMining		   是否不允许无交易时快速出空块，无交易时延迟出块
//   - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
//   - IsGM							   是否使用了Sm2p256v1曲线
//   - RootPublicKey
//   - EnableContractLifecycle		   是否开启合约生命周期
//   - EnableVotingDictatorship		   是否开启投票(合约生命周期)时盟主一票制度
//   - ContractDeploymentVotingRule    合约部署的投票规则
//   - EnableContractManagement		   是否开启合约管理
//   - ChainVote
type NodeConfirmedConfiguration struct {
	LatcId                          uint32     `json:"latcId"`
	LatcGod                         string     `json:"latcGod"`
	LatcSaints                      []string   `json:"latcSaints"`
	Consensus                       string     `json:"consensus"`
	Epoch                           uint32     `json:"epoch"`
	Tokenless                       bool       `json:"tokenless"`
	Period                          uint32     `json:"period"`
	EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor"`
	NoTxDelayedMiningPeriodMultiple uint32     `json:"emptyAnchorPeriodMul"`
	IsGM                            bool       `json:"gm"`
	RootPublicKey                   string     `json:"rootPublicKey"`
	EnableContractLifecycle         bool       `json:"isContractVote"`
	EnableVotingDictatorship        bool       `json:"isDictatorship"`
	ContractDeploymentVotingRule    VotingRule `json:"deployRule"`
	EnableContractManagement        bool       `json:"contractPermission"`
	ChainVote                       bool       `json:"chainVote"`
}

// NodeVersion 节点的版本信息
//   - IncentiveInfo  example: 无激励
//   - ConsortiumInfo example: 联盟链
//   - BuildDateTime  2024-05-28 07:47:18
//   - Version		  v2.0.0
type NodeVersion struct {
	IncentiveInfo  string `json:"IsIncentive"`
	ConsortiumInfo string `json:"IsConsortium"`
	BuildDateTime  string `json:"BuildDate"`
	Version        string `json:"Version"`
}

// NodeConfiguration 节点配置信息
type NodeConfiguration struct {
	Latc LatcConfig `json:"Latc" yaml:"latc"` //Lattice config
	Node NodeConfig `json:"Node" yaml:"node"` //Node config
}
type LatcConfig struct {
	NetworkIDGroup  []int           `json:"NetworkIDGroup" yaml:"networkIDGroup"`
	Name            string          `json:"Name" yaml:"name"`
	Desc            string          `json:"Desc" yaml:"desc"`
	ConfigsDir      string          `json:"ConfigsDir" yaml:"configsDir"`
	RaftDebug       bool            `json:"RaftDebug" yaml:"raftDebug"`
	ConsensusConfig ConsensusConfig `json:"ConsensusConfig" yaml:"consensus"`
	TxPoolConfig    TxPoolConfig    `json:"TxPoolConfig" yaml:"txPool"`
	StateConfig     StateConfig     `json:"StateConfig" yaml:"state"`
	EvidenceConfig  EvidenceConfig  `json:"EvidenceConfig" yaml:"evidence"`
}
type ConsensusConfig struct {
	ResendSigInterval int `json:"ResendSigInterval" yaml:"resendsiginterval"`
	ViewChangeTimeout int `json:"ViewChangeTimeout" yaml:"viewchangetimeout"`
}
type TxPoolConfig struct {
	MinReadySize  int  `json:"MinReadySize" yaml:"minReadySize"`
	StartInterval int  `json:"StartInterval" yaml:"startInterval"`
	ReadyInterval int  `json:"ReadyInterval" yaml:"readyInterval"`
	StoreNeed     bool `json:"StoreNeed" yaml:"storeNeed"`
}
type StateConfig struct {
	WorldStateCompress bool `json:"WorldStateCompress" yaml:"worldStateCompress"`
	TrieNumberLimit    int  `json:"TrieNumberLimit" yaml:"trieNumberLimit"`
	TrieCleanLimit     int  `json:"TrieCleanLimit" yaml:"trieCleanLimit"`
}
type EvidenceConfig struct {
	EviEnable bool `json:"EviEnable" yaml:"evienable"`
	EviLevel  int  `json:"EviLevel" yaml:"evilevel"`
}
type NodeConfig struct {
	DebugPWD        string `json:"DebugPWD" yaml:"debugPWD"`
	DebugDB         bool   `json:"DebugDB" yaml:"debugDB"`
	Name            string `json:"Name" yaml:"name"`
	Version         string `json:"Version" yaml:"version"`
	IteratorVersion int    `json:"IteratorVersion" yaml:"iteratorVersion"`
	DataDir         string `json:"DataDir" yaml:"dataDir"`
	OffChainDir     string `json:"OffChainDir" yaml:"offChainDir"`   // 业务目录，目前是 存证溯源业务 的数据库
	TransientDir    string `json:"TransientDir" yaml:"transientDir"` // 非必要数据目录，可删除的数据目录
	FileKeyDir      string `json:"FileKeyDir" yaml:"fileKeyDir"`
	ConfigsDir      string `json:"ConfigsDir" yaml:"configsDir"`
	Network         struct {
		Host                 string   `json:"Host" yaml:"host"`
		AuthVirtualHosts     []string `json:"AuthVirtualHosts" yaml:"authVirtualHosts"`
		AuthPort             int      `json:"AuthPort" yaml:"authPort"`
		HTTPPort             int      `json:"HTTPPort" yaml:"HTTPPort"`
		HTTPCors             []string `json:"HTTPCors" yaml:"HTTPCors"`
		HTTPVirtualHosts     []string `json:"HTTPVirtualHosts" yaml:"HTTPVirtualHosts"`
		HTTPModules          []string `json:"HTTPModules" yaml:"HTTPModules"`
		HTTPPathPrefix       string   `json:"HTTPPathPrefix" yaml:"HTTPPathPrefix"`
		BatchRequestLimit    int      `json:"BatchRequestLimit" yaml:"batchRequestLimit"`
		BatchResponseMaxSize int      `json:"BatchResponseMaxSize" yaml:"batchResponseMaxSize"`
		WSPort               int      `json:"WSPort" yaml:"WSPort"`
		WSPathPrefix         string   `json:"WSPathPrefix" yaml:"WSPathPrefix"`
		WSOrigins            []string `json:"WSOrigins" yaml:"WSOrigins"`
		WSModules            []string `json:"WSModules" yaml:"WSModules"`
		P2PPort              int      `json:"P2PPort" yaml:"P2PPort"`
		Bootstrap            []string `json:"Bootstrap" yaml:"bootstrap"`
		NoDiscovery          bool     `json:"NoDiscovery" yaml:"noDiscovery"`
		MaxPeers             int      `json:"MaxPeers" yaml:"maxPeers"`
		StartNTP             bool     `json:"StartNTP" yaml:"startNTP"`
		NTPHost              string   `json:"NTPHost" yaml:"NTPHost"`
		JWTEnable            bool     `json:"JWTEnable" yaml:"JWTEnable"`
		JWTExpiryTimeout     int      `json:"JWTExpiryTimeout" yaml:"JWTExpiryTimeout"`
		JWTSecret            string   `json:"JWTSecret" yaml:"JWTSecret"`
	} `json:"Network" yaml:"network"`
	DatabaseConfig struct {
		DBType                 string `json:"DBType" yaml:"dbType"`
		OpenFilesCacheCapacity int    `json:"OpenFilesCacheCapacity" yaml:"openFilesCacheCapacity"`
		MemoryCache            int    `json:"MemoryCache" yaml:"memoryCache"`
		CouchDB                struct {
			Address  string `json:"Address" yaml:"address"`
			Username string `json:"Username" yaml:"username"`
			Password string `json:"Password" yaml:"password"`
		} `json:"CouchDB" yaml:"couchDB"`
	} `json:"DatabaseConfig" yaml:"database"`
	Log struct {
		LogLevel        int  `json:"LogLevel" yaml:"logLevel"`
		Console         bool `json:"Console" yaml:"console"`
		LogKeep         int  `json:"LogKeep" yaml:"logKeep"`
		LogMaxSize      int  `json:"LogMaxSize" yaml:"logMaxSize"`
		UseZhLog        bool `json:"UseZhLog" yaml:"useZLLog"`
		FileSliceSize   int  `json:"FileSliceSize" yaml:"fileSliceSize"`
		RetryChannelNum int  `json:"RetryChannelNum" yaml:"retryChannelNum"`
		UploadTimeout   int  `json:"UploadTimeout" yaml:"uploadTimeout"`
		GoroutineMax    int  `json:"GoroutineMax" yaml:"goroutineMax"`
	} `json:"Log" yaml:"log"`
	Nvm struct {
		BaseDir         string `json:"BaseDir" yaml:"baseDir"`
		Pattern         string `json:"Pattern" yaml:"pattern"`
		ContractFileDir string `json:"ContractFileDir" yaml:"contractFileDir"`
		Redundancy      bool   `json:"Redundancy" yaml:"redundancy"`
		VM              struct {
			Native struct {
				StopTimeout int `json:"StopTimeout" yaml:"stopTimeout"`
				Docker      struct {
					Enable    bool    `json:"Enable" yaml:"enable"`
					ImageName string  `json:"ImageName" yaml:"imageName"`
					Cpus      float64 `json:"Cpus" yaml:"cpus"`
					Memory    string  `json:"Memory" yaml:"memory"`
				} `json:"Docker" yaml:"docker"`
			} `json:"Native" yaml:"native"`
		} `json:"VM" yaml:"vm"`
	} `json:"Nvm" yaml:"nvm"`
	DistributedStorage struct {
		Enable      bool `json:"Enable" yaml:"enable"`
		GinHTTPPort int  `json:"GinHTTPPort" yaml:"ginHTTPPort"`
		ChunkSize   int  `json:"ChunkSize" yaml:"chunkSize"`
		NodeAmount  int  `json:"NodeAmount" yaml:"nodeAmount"`
	} `json:"DistributedStorage" yaml:"distributedStorage"`
	Freeze struct {
		Enable         bool `json:"Enable" yaml:"enable"`
		FreezeInterval int  `json:"FreezeInterval" yaml:"freezeInterval"`
		Threshold      int  `json:"Threshold" yaml:"threshold"`
		OneHand        int  `json:"OneHand" yaml:"oneHand"`
	} `json:"Freeze" yaml:"freeze"`
	MessageQueue struct {
		Enable               bool     `json:"Enable" yaml:"enable"`
		KafkaIP              string   `json:"KafkaIP" yaml:"kafkaIP"`
		KafkaPort            int      `json:"KafkaPort" yaml:"kafkaPort"`
		KafkaBrokers         []string `json:"KafkaBrokers" yaml:"KafkaBrokers"`
		DBlockWithDeposit    bool     `json:"DBlockWithDeposit" yaml:"DBlockWithDeposit"`
		DBlockWithReceipts   bool     `json:"DBlockWithReceipts" yaml:"DBlockWithReceipts"`
		SASL                 bool     `json:"SASL" yaml:"SASL"`
		SASLUser             string   `json:"SASLUser" yaml:"SASLUser"`
		SASLPassword         string   `json:"SASLPassword" yaml:"SASLPassword"`
		MQTopicNumPartitions int      `json:"MQTopicNumPartitions" yaml:"MQTopicNumPartitions"`
		MQNeedRelay          bool     `json:"MQNeedRelay" yaml:"MQNeedRelay"`
		CanRelayToMQ         bool     `json:"CanRelayToMQ" yaml:"canRelayToMQ"`
	} `json:"MessageQueue" yaml:"messageQueue"`
	PProfOn        bool   `json:"PProfOn" yaml:"pProfOn"`
	PProfPort      int    `json:"PProfPort" yaml:"pProfPort"`
	IPCPath        string `json:"IPCPath" yaml:"iPCPath"`
	NetworkIDGroup []int  `json:"NetworkIDGroup" yaml:"networkIDGroup"`
}

// FreezeSaveSpace 区块压缩节省的空间
type FreezeSaveSpace struct {
	SaveSpace uint64 `json:"saveSpace"`
	Unit      string `json:"unit"` // byte
}

// FreezeInterval 区块压缩间隔
type FreezeInterval struct {
	FreezeInterval uint64 `json:"freezeInterval"`
	Unit           string `json:"unit"` // minute
}

func (fi *FreezeInterval) String() string {
	return fmt.Sprintf("%d分钟/次", fi.FreezeInterval)
}

// NodeCertificate node digital certificate
type NodeCertificate struct {
	Certificate        *x509.Certificate   `json:"certificate,omitempty"`
	Type               NodeCertificateType `json:"certificateType,omitempty"`    // Certificate type
	OwnerAddress       string              `json:"ownerAddress,omitempty"`       // Address of the certificate owner
	BlockHeightAtIssue uint64              `json:"blockHeightAtIssue,omitempty"` // Block height at the issue certificate
	PEMCertificate     string              `json:"pemCertificate,omitempty"`
}

// LatcConsensus 前置共识和后置共识
type LatcConsensus struct {
	DaemonConsensus string `json:"daemonConsensus"` // 后置共识: Witness, pbft, raft
	FrontConsensus  string `json:"frontConsensus"`  // 前置共识: PoA, PoAP, View
}

// SyncStatus 当前区块同步状态
type SyncStatus struct {
	PackSize  int  `json:"packSize"`
	Syncing   int  `json:"syncing"`
	Reconnect bool `json:"reconnect"`
}
