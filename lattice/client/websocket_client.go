package client

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"reflect"
	"time"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/buger/jsonparser"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// 订阅结果，支持读取订阅数据和主动结束订阅
type Subscribe[T any] interface {
	// Read: 读取一条订阅数据
	Read() (*T, error)
	// Close: 主动结束订阅
	Close() error
}

type WebSocketApi interface {
	// Subscribe 订阅任意数据，测试用
	Subscribe(ctx context.Context, method string, params ...any) (Subscribe[map[string]any], error)

	// Workflow 订阅工作流
	Workflow(ctx context.Context, cond *types.WorkflowSubscribeCondition) (Subscribe[any], error)
	// WorkflowApi 订阅工作流-接口调用
	WorkflowApi(ctx context.Context, level types.WorkflowLevel, chainId *big.Int, apiScope string) (Subscribe[types.WorkflowApi], error)
	// WorkflowNodeConnect 订阅工作流-节点连接
	WorkflowNodeConnect(ctx context.Context, level types.WorkflowLevel) (Subscribe[types.WorkflowNodeConnect], error)
	// WorkflowHandshake 订阅工作流-节点握手
	WorkflowHandshake(ctx context.Context, level types.WorkflowLevel) (Subscribe[types.WorkflowHandshake], error)
	// WorkflowTransaction 订阅工作流-交易执行
	WorkflowTransaction(ctx context.Context, level types.WorkflowLevel, chainId *big.Int, hash *common.Hash) (Subscribe[types.WorkflowTransaction], error)
	// WorkflowDaemonBlock 订阅工作流-守护区块
	WorkflowDaemonBlock(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowDaemonBlock], error)
	// WorkflowConsensus 订阅工作流-参与共识
	WorkflowConsensus(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowConsensus], error)
	// WorkflowChainByChain 订阅工作流-通道管理
	WorkflowChainByChain(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowChainByChain], error)
}

type WebSocketApiInitParam struct {
	WebSocketUrl string // 节点的URL

	HandshakeTimeout time.Duration
	ReadBufferSize   int
	WriteBufferSize  int
}

// NewWebSocketApi creates a new WebSocket API for the Lattice node.
func NewWebSocketApi(args *WebSocketApiInitParam) WebSocketApi {
	if args.HandshakeTimeout == 0 {
		args.HandshakeTimeout = 30 * time.Second
	}
	return &webSocketApi{
		webSocketUrl: args.WebSocketUrl,
		dialer: &websocket.Dialer{
			Proxy:            http.ProxyFromEnvironment,
			HandshakeTimeout: args.HandshakeTimeout,
			ReadBufferSize:   args.ReadBufferSize,
			WriteBufferSize:  args.WriteBufferSize,
		},
	}
}

// 订阅消息
type SubscribeResponse struct {
	JsonRpc string `json:"jsonrpc,omitempty"`
	Method  string `json:"method,omitempty"` // 方法名, latc_subscription, node_subscription
	Params  struct {
		SubApi string          `json:"subapi,omitempty"` // subapi
		Result json.RawMessage `json:"result,omitempty"` // 订阅数据
	} `json:"params,omitempty"` // 返回结果
}

type webSocketApi struct {
	webSocketUrl string
	dialer       *websocket.Dialer
}

// 订阅任意类型的工作流
type workflowResult struct {
	subscribeRawResult
}

func unmarshalJSON[T any](b []byte) (any, error) {
	var result T
	err := json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *workflowResult) Read() (*any, error) {
	b, err := r.subscribeRawResult.Read()
	if err != nil {
		return nil, err
	}
	t, err := jsonparser.GetInt(b, "flowType")
	if err != nil {
		return nil, err
	}
	var result any
	switch types.WorkflowType(t) {
	case types.WorkflowType_API:
		result, err = unmarshalJSON[types.WorkflowApi](b)
	case types.WorkflowType_CHAIN_BY_CHAIN:
		result, err = unmarshalJSON[types.WorkflowChainByChain](b)
	case types.WorkflowType_CONSENSUS:
		result, err = unmarshalJSON[types.WorkflowConsensus](b)
	case types.WorkflowType_DAEMON_BLOCK:
		result, err = unmarshalJSON[types.WorkflowDaemonBlock](b)
	case types.WorkflowType_HANDSHAKE:
		result, err = unmarshalJSON[types.WorkflowHandshake](b)
	case types.WorkflowType_NODE_CONNECT:
		result, err = unmarshalJSON[types.WorkflowNodeConnect](b)
	case types.WorkflowType_TRANSACTION:
		result, err = unmarshalJSON[types.WorkflowTransaction](b)
	default:
		return nil, fmt.Errorf("unknown workflow type %d", t)
	}
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// Workflow implements WebSocketApi.
func (w *webSocketApi) Workflow(ctx context.Context, cond *types.WorkflowSubscribeCondition) (Subscribe[any], error) {
	r, err := w.subscribe(ctx, "node_subscribe", "workflow", cond)
	if err != nil {
		return nil, err
	}
	return &workflowResult{*r}, nil
}

// 订阅具体类型的工作流时检查返回的工作流类型是否与订阅的工作流类型相同
type workflowCheckResult[T any] struct {
	workflowResult
}

func (r *workflowCheckResult[T]) Read() (*T, error) {
	result, err := r.workflowResult.Read()
	if err != nil {
		return nil, err
	}
	checkResult, ok := (*result).(*T)
	if !ok {
		return nil, fmt.Errorf("workflow type error, can not convert %s to %s", reflect.TypeOf(*result).Name(), reflect.TypeOf((*T)(nil)).Name())
	}
	return checkResult, nil
}

func workflow[T any](w *webSocketApi, ctx context.Context, cond *types.WorkflowSubscribeCondition) (Subscribe[T], error) {
	r, err := w.Workflow(ctx, cond)
	if err != nil {
		return nil, err
	}
	return &workflowCheckResult[T]{*(r.(*workflowResult))}, nil
}

// WorkflowApi implements WebSocketApi.
func (w *webSocketApi) WorkflowApi(ctx context.Context, level types.WorkflowLevel, chainId *big.Int, apiScope string) (Subscribe[types.WorkflowApi], error) {
	return workflow[types.WorkflowApi](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_API, Level: level, ChainId: chainId, ApiScope: apiScope})
}

// WorkflowChainByChain implements WebSocketApi.
func (w *webSocketApi) WorkflowChainByChain(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowChainByChain], error) {
	return workflow[types.WorkflowChainByChain](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_CHAIN_BY_CHAIN, Level: level, ChainId: chainId})
}

// WorkflowConsensus implements WebSocketApi.
func (w *webSocketApi) WorkflowConsensus(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowConsensus], error) {
	return workflow[types.WorkflowConsensus](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_CONSENSUS, Level: level, ChainId: chainId})
}

// WorkflowDaemonBlock implements WebSocketApi.
func (w *webSocketApi) WorkflowDaemonBlock(ctx context.Context, level types.WorkflowLevel, chainId *big.Int) (Subscribe[types.WorkflowDaemonBlock], error) {
	return workflow[types.WorkflowDaemonBlock](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_DAEMON_BLOCK, Level: level, ChainId: chainId})
}

// WorkflowHandshake implements WebSocketApi.
func (w *webSocketApi) WorkflowHandshake(ctx context.Context, level types.WorkflowLevel) (Subscribe[types.WorkflowHandshake], error) {
	return workflow[types.WorkflowHandshake](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_HANDSHAKE, Level: level})
}

// WorkflowNodeConnect implements WebSocketApi.
func (w *webSocketApi) WorkflowNodeConnect(ctx context.Context, level types.WorkflowLevel) (Subscribe[types.WorkflowNodeConnect], error) {
	return workflow[types.WorkflowNodeConnect](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_NODE_CONNECT, Level: level})
}

// WorkflowTransaction implements WebSocketApi.
func (w *webSocketApi) WorkflowTransaction(ctx context.Context, level types.WorkflowLevel, chainId *big.Int, hash *common.Hash) (Subscribe[types.WorkflowTransaction], error) {
	return workflow[types.WorkflowTransaction](w, ctx, &types.WorkflowSubscribeCondition{Type: types.WorkflowType_TRANSACTION, Level: level, ChainId: chainId, Hash: hash})
}

type subscribeResult[T any] struct {
	subscribeRawResult
}

func (r *subscribeResult[T]) Read() (*T, error) {
	b, err := r.subscribeRawResult.Read()
	if err != nil {
		return nil, err
	}
	var result T
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func subscribe[T any](w *webSocketApi, ctx context.Context, method string, params ...any) (Subscribe[T], error) {
	raw, err := w.subscribe(ctx, method, params...)
	if err != nil {
		return nil, err
	}
	return &subscribeResult[T]{*raw}, nil
}

// Subscribe implements WebSocketApi.
func (w *webSocketApi) Subscribe(ctx context.Context, method string, params ...any) (Subscribe[map[string]any], error) {
	return subscribe[map[string]any](w, ctx, method, params...)
}

type subscribeRawResult struct {
	conn *websocket.Conn
}

func (r *subscribeRawResult) Read() ([]byte, error) {
	var resp SubscribeResponse
	err := r.conn.ReadJSON(&resp)
	if err != nil {
		return nil, err
	}
	return resp.Params.Result, nil
}

func (r *subscribeRawResult) Close() error {
	return r.conn.Close()
}

func (w *webSocketApi) subscribe(ctx context.Context, method string, params ...any) (*subscribeRawResult, error) {
	// 建立连接
	conn, _, err := w.dialer.DialContext(ctx, w.webSocketUrl, nil)
	if err != nil {
		return nil, err
	}
	// 发送订阅请求
	req := NewJsonRpcBody(method, params...)
	log.Debug().Msgf("Subscribe:开始订阅, req=%v", req)
	err = conn.WriteJSON(req)
	if err != nil {
		conn.Close()
		return nil, err
	}
	// 读取订阅请求响应
	var resp JsonRpcResponse[string]
	err = conn.ReadJSON(&resp)
	if err != nil {
		log.Error().Err(err).Msgf("Subscribe:订阅失败, req=%v", req)
		conn.Close()
		return nil, err
	}
	// 开始读取订阅数据
	output := subscribeRawResult{
		conn: conn,
	}
	return &output, nil
}
