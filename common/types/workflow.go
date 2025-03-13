package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// 工作流订阅条件
type WorkflowSubscribeCondition struct {
	// 工作流类型
	Type WorkflowType `json:"logType,omitempty"`
	// 工作流级别
	Level WorkflowLevel `json:"logLevel,omitempty"`
	// 链ID
	ChainId *big.Int `json:"chainId,omitempty"`
	// 订阅的API范围，类型不是接口调用时无效
	ApiScope string `json:"apiScope,omitempty"`
	// 订阅的交易hash，类型不是交易执行时无效
	Hash *common.Hash `json:"hash,omitempty"`
}

// 定义工作流接口
type Workflow interface {
	// Type 工作流类型
	GetType() WorkflowType
	// Level 工作流级别
	GetLevel() WorkflowLevel
	// 链ID，当工作流跟链无关时可能为空
	GetChainId() *big.Int
	// Info 信息，一般级别时的信息
	GetInfo() string
	// Error 错误信息，错误级别时的信息
	GetError() string
	// Timestamp 时间戳
	GetTimestamp() int64
	// Phase 短暂的订阅的结束标志
	GetPhase() WorkflowPhase
}

// 工作流订阅推送的信息（公有）
type WorkflowCommon struct {
	// 工作流类型
	Type WorkflowType `json:"flowType,omitempty"`
	// 工作流级别
	Level WorkflowLevel `json:"flowLevel"`
	// 链ID，当工作流跟链无关时可能为空
	ChainId *big.Int `json:"chainId,omitempty"`
	// 信息，一般级别时的信息
	Info string `json:"info,omitempty"`
	// 错误信息，错误级别时的信息
	Error string `json:"error,omitempty"`
	// 时间戳
	Timestamp int64
	// 短暂的订阅的结束标志
	Phase WorkflowPhase `json:"phase,omitempty"`
}

// 实现工作流接口
func (w *WorkflowCommon) GetType() WorkflowType {
	return w.Type
}
func (w *WorkflowCommon) GetLevel() WorkflowLevel {
	return w.Level
}
func (w *WorkflowCommon) GetChainId() *big.Int {
	return w.ChainId
}
func (w *WorkflowCommon) GetInfo() string {
	return w.Info
}
func (w *WorkflowCommon) GetError() string {
	return w.Error
}
func (w *WorkflowCommon) GetTimestamp() int64 {
	return w.Timestamp
}
func (w *WorkflowCommon) GetPhase() WorkflowPhase {
	return w.Phase
}

// 工作流级别
type WorkflowLevel int

const (
	WorkflowLevel_NONE WorkflowLevel = iota
	// 一般级别
	WorkflowLevel_INFO
	// 错误级别
	WorkflowLevel_ERROR
)

// 工作流类型
type WorkflowType int

const (
	WorkflowType_NONE WorkflowType = iota
	// 接口调用，必须传递接口前缀（rpc接口method方法名，如订阅wallet_相关接口，前缀为wallet_），不能接受所有节点调用的工作流（没必要）
	WorkflowType_API
	// 节点连接，在订阅之前的连接工作流不会推送
	WorkflowType_NODE_CONNECT
	// 节点握手，在订阅之前的握手工作流不会推送
	WorkflowType_HANDSHAKE
	// 交易执行，可选参数：hash，订阅指定交易的执行工作流，短暂的订阅，交易落块后订阅结束
	WorkflowType_TRANSACTION
	// 守护区块
	WorkflowType_DAEMON_BLOCK
	// 共识流程
	WorkflowType_CONSENSUS
	// 通道管理
	WorkflowType_CHAIN_BY_CHAIN
)

// 短暂的订阅的结束标志
// 订阅交易执行的工作流，目的是监听某笔交易的生命周期，如交易落块，即完成了监听的过程。
// 此类工作流定义为短暂的工作流，有结束标志，在订阅的工作周期结束后，不再推送消息。
// 结束标志的字段为 Phase 它目前只有一个值 END ，客户端可以据此判断订阅是否已经结束。
type WorkflowPhase string

const (
	WorkflowPhase_END WorkflowPhase = "END"
)

// 接口调用
type WorkflowApi struct {
	WorkflowCommon
	// 接口方法
	Method string `json:"method,omitempty"`
	// 接口参数
	Params string `json:"params,omitempty"`
	// 接口返回值，与公有的info不同是return在整个流程成功后才会返回
	Return string `json:"return,omitempty"`
}

// 节点连接
type WorkflowNodeConnect struct {
	WorkflowCommon
	// 连接节点的信息/节点hash
	INodeInfo string `json:"iNodeInfo,omitempty"`
}

// 节点握手
type WorkflowHandshake struct {
	WorkflowCommon
	// 节点信息
	PeerId string `json:"peerId,omitempty"`
}

// 交易执行
type WorkflowTransaction struct {
	WorkflowCommon
	// 交易区块hash
	Hash string `json:"hash,omitempty"`
	// 交易区块高度
	Height *big.Int `json:"height,omitempty"`
}

// 守护区块
type WorkflowDaemonBlock struct {
	WorkflowCommon
	// 守护区块hash
	Hash string `json:"hash,omitempty"`
	// 守护区块高度
	Height *big.Int `json:"height,omitempty"`
}

// 共识流程
type WorkflowConsensus struct {
	WorkflowCommon
}

// 通道管理
type WorkflowChainByChain struct {
	WorkflowCommon
	// 正在创建或连接的链ID
	ChildChainId *big.Int `json:"childChainId,omitempty"`
	// 以链建链的操作类型：新建链、停止链...
	ChainByChainType string `json:"chainByChainType,omitempty"`
}
