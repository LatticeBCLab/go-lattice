package types

type Proposal[T ContractLifecycleProposal | ModifyChainConfigProposal | SubchainProposal] struct {
	Type    uint8 `json:"proposalType,omitempty"`
	Content *T    `json:"proposalContent,omitempty"`
}

// ContractLifecycleProposal 合约生命周期提案
// State      ProposalState
// IsRevoke   0-吊销、1-冻结、2-解冻
// Period     0-部署、1-升级、2-吊销、3-冻结、4-解冻
// CreatedAt  提案创建时间戳
// ModifiedAt 提案修改时间戳
// DBHeight   提案结束时的守护区块高度,
type ContractLifecycleProposal struct {
	Id              string        `json:"proposalId"`
	State           ProposalState `json:"proposalState"`
	Nonce           uint64        `json:"nonce"`
	ContractAddress string        `json:"contractAddress"`
	IsRevoke        uint32        `json:"isRevoke"`
	Period          uint8         `json:"period"`
	CreatedAt       int64         `json:"createAt,omitempty"`
	ModifiedAt      int64         `json:"modifiedAt,omitempty"`
	DBHeight        uint64        `json:"dbNumber"`
}

// ModifyChainConfigProposal 修改链配置提案
// State ProposalState
type ModifyChainConfigProposal struct {
	Id             string   `json:"proposalId"`
	State          uint8    `json:"proposalState"`
	Nonce          uint64   `json:"nonce"`
	Type           uint8    `json:"modifyType"`
	Period         uint32   `json:"period"`
	IsDictatorship bool     `json:"isDictatorship"`
	NoEmptyAnchor  bool     `json:"noEmptyAnchor"`
	DeployRule     uint8    `json:"deployRule"`
	LatcSaint      []string `json:"latcSaint"`
	Consensus      string   `json:"consensus"`
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
type ProposalType uint8

const (
	ProposalTypeNone ProposalType = iota
	ProposalTypeContractManagement
	ProposalTypeContractLifecycle
	ProposalTypeModifyChainConfiguration
	ProposalTypeChainByChain
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
