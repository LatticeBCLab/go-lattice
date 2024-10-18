package client

import (
	"context"
	"math/big"

	"github.com/LatticeBCLab/go-lattice/common/types"
)

func (api *httpApi) Freeze(ctx context.Context, chainId string, dblockNumber *big.Int) (uint64, error) {
	response, err := Post[uint64](ctx, api.NodeUrl, NewJsonRpcBody("latc_freeze", dblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return 0, err
	}
	if response.Error != nil {
		return 0, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetFreezeDBlockByHash(ctx context.Context, chainId string, hash string) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeDBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetFreezeDBlockByNumber(ctx context.Context, chainId string, dblockNumber *big.Int) (*types.DaemonBlock, error) {
	response, err := Post[types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeDBlockByNumber", dblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetFreezeTBlockByHash(ctx context.Context, chainId string, hash string) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeTBlockByHash", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetFreezeTBlockByNumber(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeTBlockByNumber", address, tblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetFreezeReceipt(ctx context.Context, chainId string, address string, tblockNumber *big.Int) (*types.Receipt, error) {
	response, err := Post[types.Receipt](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeReceipt", address, tblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetFreezeSaveSpace(ctx context.Context, chainId string) (*types.FreezeSaveSpace, error) {
	response, err := Post[types.FreezeSaveSpace](ctx, api.NodeUrl, NewJsonRpcBody("latc_getFreezeSaveSpace"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
