package types

import (
	"fmt"
	"github.com/samber/lo"
	"math/big"
)

type Proposal[T ContractLifecycleProposal | ModifyChainConfigProposal | SubchainProposal | NodeCertificateProposal] struct {
	Type    ProposalType `json:"proposalType,omitempty"`
	Content *T           `json:"proposalContent,omitempty"`
}

// ContractLifecycleProposal 合约生命周期提案
// State      ProposalState
// IsRevoke   0-吊销、1-冻结、2-解冻. 【已废弃】
// Period     0-部署、1-升级、2-吊销、3-冻结、4-解冻. 【已废弃】
// CreatedAt  提案创建时间戳
// ModifiedAt 提案修改时间戳
// DBHeight   提案结束时的守护区块高度
// ContractManagerBits 合约状态state：合约信息中ContractManagerBits[0] 的值，freezeState: 合约信息中ContractManagerBits[1] 的值
type ContractLifecycleProposal struct {
	Id                  string        `json:"proposalId"`
	State               ProposalState `json:"proposalState"`
	Nonce               uint64        `json:"nonce"`
	Launcher            string        `json:"launcher"`
	TxHash              string        `json:"txHash"`
	ContractAddress     string        `json:"contractAddress"`
	CreatedAt           int64         `json:"createAt,omitempty"`
	ModifiedAt          int64         `json:"modifiedAt,omitempty"`
	DBHeight            uint64        `json:"dbNumber"`
	ContractManagerBits []byte        `json:"contractManagerBits"`
}

// GetContractStates 获取合约状态集合
func (c *ContractLifecycleProposal) GetContractStates() []ContractLifecycleState {
	var states []ContractLifecycleState
	part := fmt.Sprintf("%08b", c.ContractManagerBits[0])
	freezeBit := part[5]
	upgradeBit := part[6]
	deployBit := part[7]

	if freezeBit == '1' {
		states = append(states, ContractLifecycleStateFREEZE)
	}
	if upgradeBit == '1' {
		states = append(states, ContractLifecycleStateUPGRADE)
	}
	if deployBit == '1' {
		states = append(states, ContractLifecycleStateDEPLOY)
	}

	return lo.Ternary(len(states) == 0, []ContractLifecycleState{ContractLifecycleStateUNSPECIFIED}, states)
}

// ModifyChainConfigProposal 修改链配置提案
// State ProposalState
type ModifyChainConfigProposal struct {
	Id                string   `json:"proposalId"`
	State             uint8    `json:"proposalState"`
	Nonce             uint64   `json:"nonce"`
	Type              uint8    `json:"modifyType"`
	Period            uint32   `json:"period"`
	IsDictatorship    bool     `json:"isDictatorship"`
	NoEmptyAnchor     bool     `json:"noEmptyAnchor"`
	DeployRule        uint8    `json:"deployRule"`
	LatcSaint         []string `json:"latcSaint"`
	Consensus         string   `json:"consensus"`
	RelatedProposalId string   `json:"nodeCertProposal"`
}

// SubchainProposal 以链建链的提案内容
// State ProposalState
type SubchainProposal struct {
	Id          string `json:"proposalId,omitempty"`
	State       uint8  `json:"proposalState,omitempty"`
	ChainConfig struct {
		// create
		NewChain struct {
			Id                              uint32     `json:"chainId,omitempty"`
			Name                            string     `json:"name,omitempty"`
			Period                          uint32     `json:"period,omitempty"`
			Tokenless                       bool       `json:"tokenless,omitempty"`
			EnableNoTxDelayedMining         bool       `json:"noEmptyAnchor,omitempty"`
			NoTxDelayedMiningPeriodMultiple uint32     `json:"emptyAnchorPeriodMul,omitempty"`
			EnableContractLifecycle         bool       `json:"isContractVote,omitempty"`
			EnableVotingDictatorship        bool       `json:"isDictatorship,omitempty"`
			ContractDeploymentVotingRule    VotingRule `json:"deployRule,omitempty"`
			EnableContractManagement        bool       `json:"contractPermission,omitempty"`
		} `json:"newChain,omitempty"`
		Consensus              string `json:"consensus,omitempty"`
		Timestamp              int64  `json:"timestamp,omitempty"`
		ParentHash             string `json:"parentHash,omitempty"`
		JoinSubchainProposalId string `json:"joinProposalId,omitempty"` // the proposal id that apply for join subchain
	} `json:"ChainConfig,omitempty"`
}

type NodeCertificateProposal struct {
	Id                string     `json:"proposalId"`
	State             uint8      `json:"proposalState"`
	Nonce             uint64     `json:"nonce"`
	Receipt           string     `json:"receipt"`
	Launcher          string     `json:"launcher"`
	CreatedAt         uint64     `json:"createAt"`
	ModifiedAt        uint64     `json:"modifiedAt"`
	TxHash            string     `json:"txHash"`
	PreDBNumber       uint64     `json:"dbNumber"`
	ProposalAddress   string     `json:"proposalAddress"`
	OrganizationName  string     `json:"orgName"`
	NodeCerts         []NodeCert `json:"nodeCertParam"`
	RelatedProposalId string     `json:"configModifyProposal"`
}

type NodeCert struct {
	Address      string   `json:"address,omitempty"`      // 申请的节点地址/撤销的节点地址
	SerialNumber *big.Int `json:"serialNumber,omitempty"` // 撤销证书的序列号
	CertDigest   []byte   `json:"certDigest,omitempty"`   // 证书摘要
	Signs        []byte   `json:"signs,omitempty"`        // 共识节点签名
	Revoked      bool     `json:"revoked,omitempty"`      // 是否已经撤销
}

// ProposalState 提案状态
//   - ProposalStateNONE 	 空值
//   - ProposalStateINITIAL  提案正在进行投票
//   - ProposalStateSUCCESS  提案投票通过
//   - ProposalStateFAILED 	 提案投票未通过
//   - ProposalStateEXPIRED  提案已过期
//   - ProposalStateERROR 	 提案执行错误
//   - ProposalStateCANCEL   提案已取消
//   - ProposalStateNOTSTART 提案未开始
type ProposalState uint8

const (
	ProposalStateNONE ProposalState = iota
	ProposalStateINITIAL
	ProposalStateSUCCESS
	ProposalStateFAILED
	ProposalStateEXPIRED
	ProposalStateERROR
	ProposalStateCANCEL
	ProposalStateNOTSTART
)

// ProposalType 提案类型
//   - ProposalTypeNone						None
//   - ProposalTypeContractManagement		合约内部管理
//   - ProposalTypeContractLifecycle		合约生命周期
//   - ProposalTypeModifyChainConfiguration 修改链配置
//   - ProposalTypeChainByChain				以链建链
//   - ProposalTypeNodeCertificate          节点证书
type ProposalType uint8

const (
	ProposalTypeNone ProposalType = iota
	ProposalTypeContractManagement
	ProposalTypeContractLifecycle
	ProposalTypeModifyChainConfiguration
	ProposalTypeChainByChain
	ProposalTypeNodeCertificate
)

// ModifyChainConfigurationType 修改链配置类型
//   - ModifyChainConfigurationTypeUpdatePeriod 								更新出块时间
//   - ModifyChainConfigurationTypeISEnableContractLifecycleVotingDictatorship  是否开启合约生命周期投票的盟主独裁机制，否则为共识投票
//   - ModifyChainConfigurationTypeAddConsensusNodes 							添加共识节点
//   - ModifyChainConfigurationTypeDeleteConsensusNodes 						删除共识节点
//   - ModifyChainConfigurationTypeUpdateConsensus   							更换共识
//   - ModifyChainConfigurationTypeUpdateContractDeploymentVotingRule 			更新合约部署的投票规则
//   - ModifyChainConfigurationTypeEnableNoTxDelayedMining						无交易时是否延迟出块
//   - ModifyChainConfigurationTypeEnableContractLifecycle 						是否开启合约生命周期
//   - ModifyChainConfigurationTypeEnableContractManagement 					是否启用合约内部权限管理
//   - ModifyChainConfigurationTypeReplaceConsensusNodes 						替换共识节点
//   - ModifyChainConfigurationTypeUpdateNoTxDelayedMiningPeriodMultiple		更新无交易时延迟出块的阶段倍数
//   - ModifyChainConfigurationTypeUpdateProposalExpirationDays					更新提案的过期天数
//   - ModifyChainConfigurationTypeUpdateChainByChainVotingRule 				更新以链建链的投票规则（修改通道管理规则）
type ModifyChainConfigurationType uint8

const (
	ModifyChainConfigurationTypeUpdatePeriod ModifyChainConfigurationType = iota
	ModifyChainConfigurationTypeISEnableContractLifecycleVotingDictatorship
	ModifyChainConfigurationTypeAddConsensusNodes
	ModifyChainConfigurationTypeDeleteConsensusNodes
	ModifyChainConfigurationTypeUpdateConsensus
	ModifyChainConfigurationTypeUpdateContractDeploymentVotingRule
	ModifyChainConfigurationTypeEnableNoTxDelayedMining
	ModifyChainConfigurationTypeEnableContractLifecycle
	ModifyChainConfigurationTypeEnableContractManagement
	ModifyChainConfigurationTypeReplaceConsensusNodes
	ModifyChainConfigurationTypeUpdateNoTxDelayedMiningPeriodMultiple
	ModifyChainConfigurationTypeUpdateProposalExpirationDays
	ModifyChainConfigurationTypeUpdateChainByChainVotingRule
)

// VoteSuggestion 投票建议
//   - VoteSuggestionDISAPPROVE 反对
//   - VoteSuggestionAPPROVE	同意
type VoteSuggestion uint8

const (
	VoteSuggestionDISAPPROVE VoteSuggestion = iota
	VoteSuggestionAPPROVE
)

// VoteDetails 投票详情
//   - VoteId
//   - ProposalId
//   - VoteSuggestion
//   - Address
//   - ProposalType
//   - Nonce
//   - CreatedAt
type VoteDetails struct {
	VoteId         string         `json:"voteId"`
	ProposalId     string         `json:"proposalId"`
	VoteSuggestion VoteSuggestion `json:"voteSuggestion"`
	Address        string         `json:"address"`
	ProposalType   ProposalType   `json:"proposalType"`
	Nonce          uint64         `json:"nonce"`
	CreatedAt      uint64         `json:"createAt"`
}
