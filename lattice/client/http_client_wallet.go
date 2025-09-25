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

// ProxyReEncryption 代理重加密
// Args:
//   - ciphertext string: 密文，和重加密 wasm 保持一致，hex字符串(不带0x前缀)
//   - businessAddress string: 业务合约地址
//   - initiator string: 发起者
//   - whitelist string: 数据接收者
//
// Returns:
//   - string: 重加密后的密文，可能的error(16进制)
//   - error
func (api *httpApi) ProxyReEncryption(ctx context.Context, chainId string, ciphertext, businessAddress, initiator, whitelist string) (string, error) {
	response, err := Post[string](
		ctx,
		api.NodeUrl,
		NewJsonRpcBody("wallet_proxyRecrypt", map[string]string{"ciphertext": ciphertext, "businessId": businessAddress, "initiator": initiator, "whiteList": whitelist}),
		api.newHeaders(chainId),
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

// ImportCertificates 导入证书
func (api *httpApi) ImportCertificates(ctx context.Context, chainId string, pemCertificates []string) error {
	response, err := Post[any](
		ctx,
		api.NodeUrl,
		NewJsonRpcBody("wallet_importCert", pemCertificates),
		api.newHeaders(chainId),
		api.transport,
	)
	if err != nil {
		return err
	}
	if response.Error != nil {
		return response.Error.Error()
	}
	return nil
}
