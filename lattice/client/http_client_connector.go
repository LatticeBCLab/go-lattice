package client

import (
	"context"

	"github.com/LatticeBCLab/go-lattice/common/types"
)

func (api *httpApi) GetCreateContractSolidity(
	ctx context.Context,
	params *types.CreateDataContractParams,
) (string, error) {
	response, err := Post[string](
		ctx,
		api.NodeUrl,
		NewJsonRpcBody("wallet_getCreateContractSolidity", params),
		api.newHeaders(""),
		api.transport,
	)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}
