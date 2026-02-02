package client

import (
	"context"

	"github.com/LatticeBCLab/go-lattice/common/types"
)

func (api *httpApi) GetCurrentIDB(ctx context.Context, chainId, owner string) (*types.CurrentIDB, error) {
	response, err := Post[types.CurrentIDB](ctx, api.NodeUrl, NewJsonRpcBody("latc_getCurrentIDB", owner), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetDIDBByHash(ctx context.Context, chainId, hash, docHash string) (*types.DIDB, error) {
	response, err := Post[types.DIDB](ctx, api.NodeUrl, NewJsonRpcBody("latc_getDIDBByHash", hash, docHash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
