package client

import (
	"context"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/wallet"
	"math/big"
)

func (api *httpApi) GetNodeInfo(ctx context.Context) (*types.NodeInfo, error) {
	response, err := Post[types.NodeInfo](ctx, api.NodeUrl, NewJsonRpcBody("node_nodeInfo"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetConsensusNodesStatus(ctx context.Context, chainID string) ([]*types.ConsensusNodeStatus, error) {
	response, err := Post[[]*types.ConsensusNodeStatus](ctx, api.NodeUrl, NewJsonRpcBody("witness_nodeList"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetGenesisNodeAddress(ctx context.Context, chainID string) (string, error) {
	response, err := Post[string](ctx, api.NodeUrl, NewJsonRpcBody("wallet_getGenesisNode"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodePeers(ctx context.Context) ([]*types.NodePeer, error) {
	response, err := Post[[]*types.NodePeer](ctx, api.NodeUrl, NewJsonRpcBody("node_peers"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetNodeConfirmedConfiguration(ctx context.Context, chainId string) (*types.NodeConfirmedConfiguration, error) {
	response, err := Post[types.NodeConfirmedConfiguration](ctx, api.NodeUrl, NewJsonRpcBody("wallet_getConfirmConfig"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeVersion(ctx context.Context) (*types.NodeVersion, error) {
	response, err := Post[types.NodeVersion](ctx, api.NodeUrl, NewJsonRpcBody("node_nodeVersion"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeSaintKey(ctx context.Context) (*wallet.FileKey, error) {
	response, err := Post[wallet.FileKey](ctx, api.NodeUrl, NewJsonRpcBody("node_getSaintKey"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeWorkingDirectory(ctx context.Context) (string, error) {
	response, err := Post[string](ctx, api.NodeUrl, NewJsonRpcBody("node_getLocationPath"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetSnapshot(ctx context.Context, chainId string, daemonBlockHeight *big.Int) (*types.NodeProtocolConfig, error) {
	response, err := Post[types.NodeProtocolConfig](ctx, api.NodeUrl, NewJsonRpcBody("clique_getSnapshot", daemonBlockHeight), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
