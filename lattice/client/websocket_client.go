package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/buger/jsonparser"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type SubscriptionResult struct {
	ID     string          `json:"subapi,omitempty"` // id
	Result json.RawMessage `json:"result,omitempty"` // 数据
}

// 订阅消息
type SubscribeResponse struct {
	JsonRpc string             `json:"jsonrpc,omitempty"`
	Method  string             `json:"method,omitempty"` // 方法名, latc_subscription, node_subscription
	Params  SubscriptionResult `json:"params,omitempty"` // 返回结果
}

// 订阅结果，支持读取订阅数据和主动结束订阅
type Subscribe[T any] interface {
	// ID 返回订阅ID
	ID() string
	// Read 读取一条订阅数据
	Read() (T, error)
	// Close 主动结束订阅
	Close() error
}

type WebSocketApi interface {
	// Subscribe 订阅任意数据，测试用
	Subscribe(ctx context.Context, method string, params ...any) (Subscribe[map[string]any], error)

	// Workflow 订阅工作流
	Workflow(ctx context.Context, cond *types.WorkflowSubscribeCondition) (Subscribe[types.Workflow], error)
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

type webSocketApi struct {
	webSocketUrl string
	dialer       *websocket.Dialer
}

// 订阅任意类型的工作流
type workflowResult struct {
	subscribeRawResult
}

func unmarshalWorkflow[T any](b []byte) (*T, error) {
	var result T
	err := json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r workflowResult) Read() (types.Workflow, error) {
	raw, err := r.subscribeRawResult.Read()
	if err != nil {
		return nil, err
	}
	b := raw.Params.Result
	t, err := jsonparser.GetInt(b, "flowType")
	if err != nil {
		return nil, err
	}
	var result types.Workflow
	switch types.WorkflowType(t) {
	case types.WorkflowType_API:
		result, err = unmarshalWorkflow[types.WorkflowApi](b)
	case types.WorkflowType_CHAIN_BY_CHAIN:
		result, err = unmarshalWorkflow[types.WorkflowChainByChain](b)
	case types.WorkflowType_CONSENSUS:
		result, err = unmarshalWorkflow[types.WorkflowConsensus](b)
	case types.WorkflowType_DAEMON_BLOCK:
		result, err = unmarshalWorkflow[types.WorkflowDaemonBlock](b)
	case types.WorkflowType_HANDSHAKE:
		result, err = unmarshalWorkflow[types.WorkflowHandshake](b)
	case types.WorkflowType_NODE_CONNECT:
		result, err = unmarshalWorkflow[types.WorkflowNodeConnect](b)
	case types.WorkflowType_TRANSACTION:
		result, err = unmarshalWorkflow[types.WorkflowTransaction](b)
	default:
		return nil, fmt.Errorf("unknown workflow type %d", t)
	}
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Workflow implements WebSocketApi.
func (w *webSocketApi) Workflow(ctx context.Context, cond *types.WorkflowSubscribeCondition) (Subscribe[types.Workflow], error) {
	if cond == nil {
		cond = &types.WorkflowSubscribeCondition{}
	}
	r, err := w.subscribe(ctx, "node_subscribe", "workflow", cond)
	if err != nil {
		return nil, err
	}
	return &workflowResult{r}, nil
}

type subscribeResult[T any] struct {
	subscribeRawResult
}

func (r subscribeResult[T]) Read() (result T, err error) {
	b, err := r.subscribeRawResult.Read()
	if err != nil {
		return
	}
	err = json.Unmarshal(b.Params.Result, &result)
	if err != nil {
		return
	}
	return
}

func subscribe[T any](w *webSocketApi, ctx context.Context, method string, params ...any) (Subscribe[T], error) {
	raw, err := w.subscribe(ctx, method, params...)
	if err != nil {
		return nil, err
	}
	return &subscribeResult[T]{raw}, nil
}

// Subscribe implements WebSocketApi.
func (w *webSocketApi) Subscribe(ctx context.Context, method string, params ...any) (Subscribe[map[string]any], error) {
	return subscribe[map[string]any](w, ctx, method, params...)
}

type subscribeRawResult struct {
	conn *websocket.Conn
	id   string
}

func (r subscribeRawResult) ID() string {
	return r.id
}

func (r subscribeRawResult) Read() (*SubscribeResponse, error) {
	var resp SubscribeResponse
	err := r.conn.ReadJSON(&resp)
	if err != nil {
		return nil, err
	}
	// 检验订阅ID是否匹配
	if resp.Params.ID != r.id {
		return nil, fmt.Errorf("订阅ID不匹配，got %s, expect %s", resp.Params.ID, r.id)
	}
	return &resp, nil
}

func (r subscribeRawResult) Close() error {
	return r.conn.Close()
}

func (w *webSocketApi) subscribe(ctx context.Context, method string, params ...any) (result subscribeRawResult, err error) {
	// 建立连接
	conn, _, err := w.dialer.DialContext(ctx, w.webSocketUrl, nil)
	if err != nil {
		return
	}
	// 发送订阅请求
	req := NewJsonRpcBody(method, params...)
	debugReq, _ := json.Marshal(req)
	log.Debug().Msgf("Subscribe:开始订阅, req=%s", debugReq)
	err = conn.WriteJSON(req)
	if err != nil {
		conn.Close()
		return
	}
	// 读取订阅请求响应
	var resp JsonRpcResponse[string]
	err = conn.ReadJSON(&resp)
	if err != nil {
		log.Error().Err(err).Msgf("Subscribe:订阅失败, req=%v", req)
		conn.Close()
		return
	}
	// 开始读取订阅数据
	result = subscribeRawResult{
		conn: conn,
		id:   resp.Result,
	}
	return
}
