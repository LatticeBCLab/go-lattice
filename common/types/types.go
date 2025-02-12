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
//   - Saint
//   - ChainId
//   - DaemonBlockHeight
//   - DaemonBlockHash
//   - GinHttpPort       文件传输服务的端口
//   - DFS 		   		 是否使用文件传输服务
//   - INode       		 节点的INode信息
type SubchainPeer struct {
	Id                string   `json:"Id"`
	Saint             string   `json:"Saint"`
	ChainId           uint32   `json:"ChainId"`
	DaemonBlockHeight *big.Int `json:"DBNum"`
	DaemonBlockHash   string   `json:"DBHash"`
	GinHttpPort       int      `json:"GinHttpPort"`
	DFS               bool     `json:"DFS"`
	INode             string   `json:"inode"`
}

// NodeConfig 节点配置信息
type NodeConfig struct {
	Lattice struct {
		NetworkIDGroup []int `json:"networkIDGroup"`
	} `json:"latc"`
	Node struct {
		Name    string `json:"name"`
		DataDir string `json:"dataDir"`
		Network struct {
			Host      string   `json:"host,omitempty"`
			HTTPPort  int      `json:"HTTPPort,omitempty"`
			WSPort    int      `json:"WSPort,omitempty"`
			P2PPort   int      `json:"P2PPort,omitempty"`
			Bootstrap []string `json:"bootstrap"`
			JWTEnable bool     `json:"JWTEnable"`
			JWTSecret string   `json:"JWTSecret"`
		}
		DistributedStorage struct {
			Enable      bool `json:"Enable,omitempty"`
			GinHTTPPort int  `json:"GinHTTPPort,omitempty"`
		}
		Freeze struct {
			Enable   bool `json:"Enable,omitempty"`
			Interval int  `json:"FreezeInterval,omitempty"`
		}
		NVM struct {
			MultilingualContractFileStoragePath     string `json:"ContractFileDir,omitempty"`
			MultilingualContractPattern             string `json:"Pattern,omitempty"`    // 生成合约执行路径pattern
			MultilingualContractRedundantDeployment bool   `json:"Redundancy,omitempty"` // 是否支持冗余部署，即对一个合约文件部署多次，生成多个合约地址
		} `json:"Nvm,omitempty"`
	}
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
// - LatcID                          链ID
// - Name                            名称
// - Desc                            描述
// - LatcGodAddr                     守护链的地址
// - LatcSaints                      共识节点列表
// - Consensus                       共识
// - Epoch                           重置投票和检查点的纪元长度
// - Tokenless                       false:有通证 true:无通证
// - Period                          出块间隔
// - EnableNoTxDelayedMining         是否不允许无交易时快速出空块，无交易时延迟出块
// - NoTxDelayedMiningPeriodMultiple 无交易时的延迟出块间隔倍数
// - IsGM                            是否使用了Sm2p256v1曲线
// - RootPublicKey                   中心化CA根证书公钥
// - EnableContractLifecycle         是否开启合约生命周期
// - EnableVotingDictatorship        是否开启投票(合约生命周期)时盟主一票制度
// - ContractDeploymentVotingRule    合约部署的投票规则
// - EnableContractManagement        是否开启合约管理
// - ChainByChainVotingRule          以链建链投票规则
// - ProposalExpirationDays          提案的过期天数，默认7天
// - ConfigurationModifyVotingRule   配置修改的投票规则
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
	Latc latcConfig `json:"Latc"` //Lattice config
	Node nodeConfig `json:"Node"` //Node config
}

type latcConfig struct {
	NetworkIDGroup []int  `json:"NetworkIDGroup"`
	Name           string `json:"Name"`
	Desc           string `json:"Desc"`
	LevelDB        struct {
		OpenFilesCacheCapacity int `json:"OpenFilesCacheCapacity"`
		MemoryCache            int `json:"MemoryCache"`
	} `json:"LevelDB"`
	AutoWitness        bool   `json:"AutoWitness"`
	ResendSigInterval  int    `json:"ResendSigInterval"`
	ViewChangeTimeout  int    `json:"ViewChangeTimeout"`
	NoRecursion        bool   `json:"NoRecursion"`
	TrieCleanLimit     int    `json:"TrieCleanLimit"`
	MinReadySize       int    `json:"MinReadySize"`
	StartInterval      int    `json:"StartInterval"`
	ReadyInterval      int    `json:"ReadyInterval"`
	StoreNeed          bool   `json:"StoreNeed"`
	EviEnable          bool   `json:"EviEnable"`
	EviLevel           int    `json:"EviLevel"`
	WorldStateCompress bool   `json:"WorldStateCompress"`
	TrieNumberLimit    int    `json:"TrieNumberLimit"`
	RaftDebug          bool   `json:"RaftDebug"`
	ConfigsDir         string `json:"ConfigsDir"`
	DebugPWD           string `json:"DebugPWD"`
}

type nodeConfig struct {
	DebugPWD        string `json:"DebugPWD"`
	DebugDB         bool   `json:"DebugDB"`
	Name            string `json:"Name"`
	Version         string `json:"Version"`
	IteratorVersion int    `json:"IteratorVersion"`
	DataDir         string `json:"DataDir"`
	OffChainDir     string `json:"OffChainDir"`  // 业务目录，目前是 存证溯源业务 的数据库
	TransientDir    string `json:"TransientDir"` // 非必要数据目录，可删除的数据目录
	FileKeyDir      string `json:"FileKeyDir"`
	ConfigsDir      string `json:"ConfigsDir"`
	Network         struct {
		Host                 string   `json:"Host"`
		AuthVirtualHosts     []string `json:"AuthVirtualHosts"`
		AuthPort             int      `json:"AuthPort"`
		HTTPPort             int      `json:"HTTPPort"`
		HTTPCors             []string `json:"HTTPCors"`
		HTTPVirtualHosts     []string `json:"HTTPVirtualHosts"`
		HTTPModules          []string `json:"HTTPModules"`
		HTTPPathPrefix       string   `json:"HTTPPathPrefix"`
		BatchRequestLimit    int      `json:"BatchRequestLimit"`
		BatchResponseMaxSize int      `json:"BatchResponseMaxSize"`
		WSPort               int      `json:"WSPort"`
		WSPathPrefix         string   `json:"WSPathPrefix"`
		WSOrigins            []string `json:"WSOrigins"`
		WSModules            []string `json:"WSModules"`
		P2PPort              int      `json:"P2PPort"`
		Bootstrap            []string `json:"Bootstrap"`
		NoDiscovery          bool     `json:"NoDiscovery"`
		MaxPeers             int      `json:"MaxPeers"`
		StartNTP             bool     `json:"StartNTP"`
		NTPHost              string   `json:"NTPHost"`
		JWTEnable            bool     `json:"JWTEnable"`
		JWTExpiryTimeout     int      `json:"JWTExpiryTimeout"`
		JWTSecret            string   `json:"JWTSecret"`
	} `json:"Network"`
	Log struct {
		LogLevel        int  `json:"LogLevel"`
		Console         bool `json:"Console"`
		LogKeep         int  `json:"LogKeep"`
		LogMaxSize      int  `json:"LogMaxSize"`
		UseZhLog        bool `json:"UseZhLog"`
		FileSliceSize   int  `json:"FileSliceSize"`
		RetryChannelNum int  `json:"RetryChannelNum"`
		UploadTimeout   int  `json:"UploadTimeout"`
		GoroutineMax    int  `json:"GoroutineMax"`
	} `json:"Log"`
	Nvm struct {
		BaseDir         string `json:"BaseDir"`
		Pattern         string `json:"Pattern"`
		ContractFileDir string `json:"ContractFileDir"`
		ContractConfig  string `json:"ContractConfig"`
		Redundancy      bool   `json:"Redundancy"`
		VM              struct {
			EnableUpgrade bool `json:"EnableUpgrade"`
			Native        struct {
				Enable      bool `json:"Enable"`
				StopTimeout int  `json:"StopTimeout"`
				Docker      struct {
					Enable    bool    `json:"Enable"`
					ImageName string  `json:"ImageName"`
					Cpus      float64 `json:"Cpus"`
					Memory    string  `json:"Memory"`
				} `json:"Docker"`
			} `json:"Native"`
		} `json:"VM"`
	} `json:"Nvm"`
	DistributedStorage struct {
		Enable      bool `json:"Enable"`
		GinHTTPPort int  `json:"GinHTTPPort"`
		ChunkSize   int  `json:"ChunkSize"`
		NodeAmount  int  `json:"NodeAmount"`
	} `json:"DistributedStorage"`
	Freeze struct {
		Enable         bool `json:"Enable"`
		FreezeInterval int  `json:"FreezeInterval"`
		Threshold      int  `json:"Threshold"`
		OneHand        int  `json:"OneHand"`
	} `json:"Freeze"`
	MessageQueue struct {
		Enable               bool   `json:"Enable"`
		KafkaIP              string `json:"KafkaIP"`
		KafkaPort            int    `json:"KafkaPort"`
		DBlockWithDeposit    bool   `json:"DBlockWithDeposit"`
		DBlockWithReceipts   bool   `json:"DBlockWithReceipts"`
		SASL                 bool   `json:"SASL"`
		SASLUser             string `json:"SASLUser"`
		SASLPassword         string `json:"SASLPassword"`
		MQTopicNumPartitions int    `json:"MQTopicNumPartitions"`
		MQNeedRelay          bool   `json:"MQNeedRelay"`
		CanRelayToMQ         bool   `json:"CanRelayToMQ"`
	} `json:"MessageQueue"`
	CouchDB struct {
		Enable   bool   `json:"Enable"`
		Address  string `json:"Address"`
		Username string `json:"Username"`
		Password string `json:"Password"`
	} `json:"CouchDB"`
	PProfOn        bool   `json:"PProfOn"`
	PProfPort      int    `json:"PProfPort"`
	IPCPath        string `json:"IPCPath"`
	NetworkIDGroup []int  `json:"NetworkIDGroup"`
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
