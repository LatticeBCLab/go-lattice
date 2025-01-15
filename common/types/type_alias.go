package types

// Curve Elliptic curve
//   - Secp256k1 NIST
//   - Sm2p256v1 GM
type Curve string

const (
	Secp256k1 Curve = "secp256k1"
	Sm2p256v1 Curve = "sm2p256v1"
)

// Fee 费用
//
// 最小单位为 wei
// 1eth = 1e18wei
// 1eth = 1e9gwei
type Fee string

// VotingRule 投票规则类型
//   - VotingRuleNO        不需要投票
//   - VotingRuleLEADER    盟主一票制
//   - VotingRuleCONSENSUS 共识投票
type VotingRule uint8

const (
	VotingRuleNO VotingRule = iota
	VotingRuleLEADER
	VotingRuleCONSENSUS
)

// Consensus 共识类型
//   - ConsensusPOA  poa共识
//   - ConsensusPBFT pbft共识
//   - ConsensusRAFT raft共识
type Consensus string

const (
	ConsensusPOA  Consensus = "PoA"
	ConsensusPBFT Consensus = "PBFT"
	ConsensusRAFT Consensus = "RAFT"
)

// NodeType 节点类型
//   - NodeTypeGENESIS   创世节点
//   - NodeTypeCONSENSUS 共识节点
//   - NodeTypeWITNESS   见证节点
//   - NodeTypeUNKNOWN   未知节点
type NodeType uint8

const (
	NodeTypeGENESIS NodeType = iota
	NodeTypeCONSENSUS
	NodeTypeWITNESS
	NodeTypeUNKNOWN
)

// ContractState 合约状态
//   - ContractStatePROHIBITED  禁止执行合约
//   - ContractStateALLOWABLE   允许执行合约
//   - ContractStateUNAVAILABLE 合约不可调用，处于冻结状态
type ContractState uint8

const (
	ContractStatePROHIBITED ContractState = iota
	ContractStateALLOWABLE
	ContractStateUNAVAILABLE
)

// ContractManagementMode 合约管理模式
//   - ContractManagementModeWHITELIST 白名单模式
//   - ContractManagementModeBLACKLIST 黑名单模式
type ContractManagementMode uint8

const (
	ContractManagementModeWHITELIST ContractManagementMode = iota
	ContractManagementModeBLACKLIST
)

// ContractLang 合约语言
//
//   - ContractLangGo	Go
//   - ContractLangJava Java.
type ContractLang string

const (
	ContractLangGo   ContractLang = "Go"
	ContractLangJava ContractLang = "Java"
)

// EvidenceType 留痕类型
//   - EvidenceTypeVOTING   投票
//   - EvidenceTypeTBLOCK   账户交易
//   - EvidenceTypeDBLOCK   守护区块
//   - EvidenceTypeSIGN     签名
//   - EvidenceTypePRECALL  预执行合约
//   - EvidenceTypeONCHAIN  发布合约交易
//   - EvidenceTypeEXECUTE  执行合约交易
//   - EvidenceTypeUPDATE   合约升级
//   - EvidenceTypeUPGRADE  升级合约的账户交易
//   - EvidenceTypeDEPLOY   合约部署
//   - EvidenceTypeCALL     合约调用
//   - EvidenceTypeREVOKE   合约吊销
//   - EvidenceTypeFREEZE   合约冻结
//   - EvidenceTypeUNFREEZE 合约解冻
//   - EvidenceTypeERROR    error错误
//   - EvidenceTypeCRITICAL crit错误
//   - EvidenceTypeADDED    增加账户
//   - EvidenceTypeDELETED  删除账户
//   - EvidenceTypeLOCKED   锁定账户
//   - EvidenceTypeUNLOCKED 解锁账户
//   - EvidenceTypeORACLE   预言机
//   - EvidenceTypeMODIFY_CONFIG  链配置修改
//   - EvidenceTypeCHAIN_BY_CHAIN 子链
//   - EvidenceTypeHANDSHAKE      握手
type EvidenceType string

const (
	EvidenceTypeVOTING         EvidenceType = "vote"
	EvidenceTypeTBLOCK         EvidenceType = "tblock"
	EvidenceTypeDBLOCK         EvidenceType = "dblock"
	EvidenceTypeSIGN           EvidenceType = "sign"
	EvidenceTypePRECALL        EvidenceType = "pre"
	EvidenceTypeONCHAIN        EvidenceType = "onChain"
	EvidenceTypeEXECUTE        EvidenceType = "execute"
	EvidenceTypeUPDATE         EvidenceType = "update"
	EvidenceTypeUPGRADE        EvidenceType = "upgrade"
	EvidenceTypeDEPLOY         EvidenceType = "deploy"
	EvidenceTypeCALL           EvidenceType = "call"
	EvidenceTypeREVOKE         EvidenceType = "revoke"
	EvidenceTypeFREEZE         EvidenceType = "freeze"
	EvidenceTypeUNFREEZE       EvidenceType = "release"
	EvidenceTypeERROR          EvidenceType = "error"
	EvidenceTypeCRITICAL       EvidenceType = "crit"
	EvidenceTypeADDED          EvidenceType = "add"
	EvidenceTypeDELETED        EvidenceType = "del"
	EvidenceTypeLOCKED         EvidenceType = "lock"
	EvidenceTypeUNLOCKED       EvidenceType = "unlock"
	EvidenceTypeORACLE         EvidenceType = "oracle"
	EvidenceTypeMODIFY_CONFIG  EvidenceType = "modifyConfig"
	EvidenceTypeCHAIN_BY_CHAIN EvidenceType = "chainByChain"
	EvidenceTypeHANDSHAKE      EvidenceType = "handshake"
)

// EvidenceLevel 留痕级别
//   - EvidenceLevelEMPTY    不填则默认为执行日志
//   - EvidenceLevelNONE	 执行日志
//   - EvidenceLevelERROR    error级别的错误日志
//   - EvidenceLevelCRITICAL crit级别的错误日志
type EvidenceLevel string

const (
	EvidenceLevelEMPTY    EvidenceLevel = ""
	EvidenceLevelNONE     EvidenceLevel = "none"
	EvidenceLevelERROR    EvidenceLevel = "error"
	EvidenceLevelCRITICAL EvidenceLevel = "crit"
)

// ContractLifecyclePeriod 合约生命周期阶段
//   - ContractLifecyclePeriodDEPLOYMENT 部署
//   - ContractLifecyclePeriodUPGRADE    升级
//   - ContractLifecyclePeriodREVOKE     吊销
//   - ContractLifecyclePeriodFREEZE     冻结
//   - ContractLifecyclePeriodUNFREEZE   解冻
type ContractLifecyclePeriod uint8

const (
	ContractLifecyclePeriodDEPLOYMENT ContractLifecyclePeriod = iota
	ContractLifecyclePeriodUPGRADE
	ContractLifecyclePeriodREVOKE
	ContractLifecyclePeriodFREEZE
	ContractLifecyclePeriodUNFREEZE
)

type TBlockState string

// TBlockState 账户区块状态
//   - TBlockStateEMPTY       空
//   - TBlockStateNONEXISTENT 不存在
//   - TBlockStatePENDING     在交易池中
//   - TBlockStateWITNESSING  见证中
//   - TBlockStateONCHAIN     已上链
//   - TBlockStateDAEMONIZED  已被守护（执行、部署）
const (
	TBlockStateEMPTY       TBlockState = ""
	TBlockStateNONEXISTENT TBlockState = "notExist"
	TBlockStatePENDING     TBlockState = "inPool"
	TBlockStateWITNESSING  TBlockState = "witnessing"
	TBlockStateONCHAIN     TBlockState = "onChain"
	TBlockStateDAEMONIZED  TBlockState = "daemon"
)
