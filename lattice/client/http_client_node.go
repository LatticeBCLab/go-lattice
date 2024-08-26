package client

import (
	"context"
	"github.com/wylu1037/lattice-go/common/types"
)

func (api *httpApi) GetNodeInfo(_ context.Context) (*types.NodeInfo, error) {
	response, err := Post[types.NodeInfo](api.Url, NewJsonRpcBody("node_nodeInfo"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetConsensusNodesStatus(_ context.Context, chainID string) ([]*types.ConsensusNodeStatus, error) {
	response, err := Post[[]*types.ConsensusNodeStatus](api.Url, NewJsonRpcBody("witness_nodeList"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetGenesisNodeAddress(_ context.Context, chainID string) (string, error) {
	response, err := Post[string](api.Url, NewJsonRpcBody("wallet_getGenesisNode"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodePeers(_ context.Context, chainID string) ([]*types.NodePeer, error) {
	response, err := Post[[]*types.NodePeer](api.Url, NewJsonRpcBody("node_peers"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodeConfig(_ context.Context, chainID string) (*types.NodeConfig, error) {
	response, err := Post[*types.NodeConfig](api.Url, NewJsonRpcBody("latc_getConfig"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}