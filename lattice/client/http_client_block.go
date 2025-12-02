package client

import (
	"context"
	"math/big"

	"github.com/rs/zerolog/log"

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

func (api *httpApi) GetTBlockState(ctx context.Context, chainId, hash string) (types.TBlockState, error) {
	response, err := Post[types.TBlockState](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockState", hash), api.newHeaders(chainId), api.transport)
	if err != nil {
		return types.TBlockStateEMPTY, err
	}
	if response.Error != nil {
		return types.TBlockStateEMPTY, response.Error.Error()
	}
	return *response.Result, nil
}

// GetLastBatchDBlockNumber implements HttpApi.
func (api *httpApi) GetLastBatchDBlockNumber(ctx context.Context, chainId string) (*big.Int, error) {
	response, err := Post[big.Int](ctx, api.NodeUrl, NewJsonRpcBody("latc_getLastedBatchDBNumber"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// GetDBlockProof implements HttpApi.
func (api *httpApi) GetDBlockProof(ctx context.Context, chainId string, dblockNumber *big.Int) (*types.WitnessProof, error) {
	response, err := Post[types.WitnessProof](ctx, api.NodeUrl, NewJsonRpcBody("latc_getDBlockProof", dblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

// GetTBlockProof implements HttpApi.
func (api *httpApi) GetTBlockProof(ctx context.Context, chainId string, accountAddress string, tblockNumber *big.Int) (*types.WitnessProof, error) {
	response, err := Post[types.WitnessProof](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockProof", accountAddress, tblockNumber), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetCurrentTBlock(ctx context.Context, chainId string, accountAddress string) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getCurrentTBlock", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTBlockByHeight(ctx context.Context, chainId, accountAddress string, height uint64) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockByNumber", accountAddress, height), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetGenesisBlock(ctx context.Context, chainId string) (*types.TransactionBlock, error) {
	response, err := Post[types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getGenesis"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetTBlocksByHeights(ctx context.Context, chainId string, accountAddress string, heights []uint64) ([]*types.TransactionBlock, error) {
	response, err := Post[[]*types.TransactionBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockByNumberRange", accountAddress, heights), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetDBlocksByHeights(ctx context.Context, chainId string, heights []uint64) ([]*types.DaemonBlock, error) {
	response, err := Post[[]*types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getDBlockByNumberRange", heights), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetRecentDBlocks(ctx context.Context, chainId string, limit uint32) ([]*types.DaemonBlock, error) {
	response, err := Post[[]*types.DaemonBlock](ctx, api.NodeUrl, NewJsonRpcBody("latc_getRecentDBlocks", limit), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetTBlockCount(ctx context.Context, chainId string) (*types.TBlockCount, error) {
	response, err := Post[types.TBlockCount](ctx, api.NodeUrl, NewJsonRpcBody("latc_getTBlockCount"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		log.Error().Err(err).Msgf("Failed to get tblock count")
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
