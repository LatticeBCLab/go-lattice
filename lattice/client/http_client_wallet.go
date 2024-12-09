package client

import "context"

func (api *httpApi) ImportFileKey(ctx context.Context, fileKey string) (string, error) {
	response, err := Post[string](ctx, api.NodeUrl, NewJsonRpcBody("wallet_importFileKey", fileKey), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return "", err
	}
	if response.Error != nil {
		return "", response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) ImportRawKey(ctx context.Context, privateKey, password string) (bool, error) {
	response, err := Post[bool](ctx, api.NodeUrl, NewJsonRpcBody("wallet_importRawKey", privateKey, password), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return false, err
	}
	if response.Error != nil {
		return false, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetAccounts(ctx context.Context, chainId string) ([]string, error) {
	response, err := Post[[]string](ctx, api.NodeUrl, NewJsonRpcBody("wallet_accountList"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}
