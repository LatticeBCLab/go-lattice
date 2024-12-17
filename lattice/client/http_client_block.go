package client

import (
	"context"

	"github.com/LatticeBCLab/go-lattice/common/types"
)

func (api *httpApi) GetLatestBlock(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getCurrentTBDB", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetLatestBlockWithPending(ctx context.Context, chainId, accountAddress string) (*types.LatestBlock, error) {
	response, err := Post[types.LatestBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getPendingTBDB", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetLatestDaemonBlock(ctx context.Context, chainID string) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getCurrentDBlock"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetDaemonBlockByHash(ctx context.Context, chainId, hash string) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getDBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetDaemonBlockByHeight(ctx context.Context, chainId string, height uint64) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getDBlockByNumber", height), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTransactionBlockByHash(ctx context.Context, chainId, hash string) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTransactionsPagination(ctx context.Context, chainId string, startDaemonBlockHeight uint64, pageSize uint16) (*types.TransactionsPagination, error) {
	response, err := Post[types.TransactionsPagination](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockPagesByDNumber", startDaemonBlockHeight, pageSize), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
