package client

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	defaultChannelSize = 1024
)

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

type WebSocketApi interface {
	Subscribe(ctx context.Context, name string) (chan any, error)
}

type webSocketApi struct {
	webSocketUrl string
	dialer       *websocket.Dialer
}

func (api *webSocketApi) Subscribe(ctx context.Context, name string) (chan any, error) {
	conn, _, err := api.dialer.DialContext(ctx, api.webSocketUrl, nil)
	if err != nil {
		return nil, err
	}
	err = conn.WriteJSON(NewJsonRpcBody("latc_subscribe", name))
	if err != nil {
		conn.Close()
		return nil, err
	}
	output := make(chan any, defaultChannelSize)
	go func() {
		defer close(output)
		defer conn.Close()
		for {
			var resp JsonRpcResponse[any]
			err := conn.ReadJSON(&resp)
			if err != nil {
				log.Error().Err(err).Msg("Subscribe:读取订阅信息失败")
				output <- err
				return
			} else if resp.Error != nil {
				log.Info().Err(resp.Error.Error()).Msg("Subscribe:订阅消息返回错误")
				output <- resp.Error.Error()
			} else {
				output <- resp.Result
			}
		}
	}()
	return output, nil
}
