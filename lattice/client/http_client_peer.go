package client

import (
	"context"
	"encoding/json"
)

func (api *httpApi) ConnectPeerAsync(ctx context.Context, id string) error {
	response, err := Post[json.RawMessage](ctx, api.NodeUrl, NewJsonRpcBody("latc_connectPeer", id), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) DisconnectPeerAsync(ctx context.Context, id string) error {
	response, err := Post[json.RawMessage](ctx, api.NodeUrl, NewJsonRpcBody("latc_disconnectPeer", id), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}
