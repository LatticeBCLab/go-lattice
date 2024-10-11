package client

import (
	"context"
	"strconv"

	"github.com/LatticeBCLab/go-lattice/common/types"
)

func (api *httpApi) GetSubchain(ctx context.Context, subchainId string) (*types.Subchain, error) {
	response, err := Post[types.Subchain](ctx, api.NodeUrl, NewJsonRpcBody("latc_latcInfo"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetCreatedSubchain(ctx context.Context) ([]uint64, error) {
	response, err := Post[[]uint64](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_getCreatedAllChains"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetJoinedSubchain(ctx context.Context) ([]uint64, error) {
	response, err := Post[[]uint64](ctx, api.NodeUrl, NewJsonRpcBody("node_getAllChainId"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetSubchainRunningStatus(ctx context.Context, subchainID string) (*types.SubchainRunningStatus, error) {
	response, err := Post[string](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_getChainStatus"), api.newHeaders(subchainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	status := *response.Result
	return &types.SubchainRunningStatus{
		Status:  status,
		Running: status == types.SubchainStatusRUNNING,
	}, nil
}

func (api *httpApi) JoinSubchain(ctx context.Context, subchainId, networkId uint64, inode string) error {
	response, err := Post[any](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_selfJoinChain", subchainId, networkId, inode), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) StartSubchain(ctx context.Context, subchainId string) error {
	response, err := Post[any](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_startSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) StopSubchain(ctx context.Context, subchainId string) error {
	response, err := Post[any](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_stopSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) DeleteSubchain(ctx context.Context, subchainId string) error {
	response, err := Post[any](ctx, api.NodeUrl, NewJsonRpcBody("cbyc_delSelfChain"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}

func (api *httpApi) GetSubchainPeers(ctx context.Context, subchainId string) (map[string]*types.SubchainPeer, error) {
	response, err := Post[map[string]*types.SubchainPeer](ctx, api.NodeUrl, NewJsonRpcBody("latc_peers"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetSubchainBriefInfo(ctx context.Context, subchainID string) ([]*types.SubchainBriefInfo, error) {
	subchainIdAsInt := int64(-1)
	if subchainID != "" {
		var err error
		subchainIdAsInt, err = strconv.ParseInt(subchainID, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	response, err := Post[[]*types.SubchainBriefInfo](ctx, api.NodeUrl, NewJsonRpcBody("latc_othersLatcInfo", subchainIdAsInt), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}
