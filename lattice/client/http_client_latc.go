package client

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strconv"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func (api *httpApi) GetLatcInfo(ctx context.Context, chainId string) (*types.NodeProtocolConfig, error) {
	response, err := Post[types.NodeProtocolConfig](ctx, api.NodeUrl, NewJsonRpcBody("latc_latcInfo"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeCertificate(ctx context.Context) (*types.NodeCertificate, error) {
	response, err := Post[x509.Certificate](ctx, api.NodeUrl, NewJsonRpcBody("latc_getOwnerCert"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}

	return x509CertificateToNodeCertificate(response.Result)
}

func x509CertificateToNodeCertificate(certificate *x509.Certificate) (*types.NodeCertificate, error) {
	subjectLocality := certificate.Subject.Locality
	blockHeight, err := strconv.ParseUint(subjectLocality[0], 10, 64)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	block := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificate.Raw,
	}

	return &types.NodeCertificate{
		Certificate:        certificate,
		BlockHeightAtIssue: blockHeight,
		Type:               types.NodeCertificateType(subjectLocality[1]),
		OwnerAddress:       subjectLocality[2],
		PEMCertificate:     string(pem.EncodeToMemory(block)),
	}, err
}

func (api *httpApi) GetPeerNodeCertificate(ctx context.Context, serialNumber string) (*types.NodeCertificate, error) {
	response, err := Post[x509.Certificate](ctx, api.NodeUrl, NewJsonRpcBody("latc_getCert", serialNumber), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}

	return x509CertificateToNodeCertificate(response.Result)
}

func (api *httpApi) GetPeerNodeCertificateByAddress(ctx context.Context, nodeAddress string) (*types.NodeCertificate, error) {
	peersMap, err := api.GetLatcPeers(ctx, emptyChainId)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	peers := lo.MapToSlice(peersMap, func(key string, value *types.SubchainPeer) *types.SubchainPeer { return value })
	filtered := lo.Filter(peers, func(peer *types.SubchainPeer, _ int) bool { return peer.Saint == nodeAddress })
	if len(filtered) == 0 {
		log.Error().Err(err)
		return nil, fmt.Errorf("no peer found with address %s", nodeAddress)
	}

	return api.GetPeerNodeCertificate(ctx, filtered[0].CertificateSN.String())
}

func (api *httpApi) GetNodeConfiguration(ctx context.Context) (*types.NodeConfiguration, error) {
	response, err := Post[types.NodeConfiguration](ctx, api.NodeUrl, NewJsonRpcBody("latc_getConfig"), api.newHeaders(emptyChainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) LoadNodeConfiguration(ctx context.Context, chainId string) (*types.NodeConfiguration, error) {
	response, err := Post[types.NodeConfiguration](ctx, api.NodeUrl, NewJsonRpcBody("latc_loadConfig"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeProtocol(ctx context.Context, chainId string) (*types.NodeProtocol, error) {
	response, err := Post[types.NodeProtocol](ctx, api.NodeUrl, NewJsonRpcBody("latc_getProtocols"), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetNodeConfig(ctx context.Context, chainID string) (*types.NodeConfiguration, error) {
	response, err := Post[types.NodeConfiguration](ctx, api.NodeUrl, NewJsonRpcBody("latc_getConfig"), api.newHeaders(chainID), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}

func (api *httpApi) GetLatcPeers(ctx context.Context, subchainId string) (map[string]*types.SubchainPeer, error) {
	response, err := Post[map[string]*types.SubchainPeer](ctx, api.NodeUrl, NewJsonRpcBody("latc_peers"), api.newHeaders(subchainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return *response.Result, nil
}

func (api *httpApi) GetBalanceWithPending(ctx context.Context, chainId, accountAddress string) (*types.Balance, error) {
	response, err := Post[types.Balance](ctx, api.NodeUrl, NewJsonRpcBody("latc_getBalanceWithPending", accountAddress), api.newHeaders(chainId), api.transport)
	if err != nil {
		return nil, err
	}
	if response.Error != nil {
		return nil, response.Error.Error()
	}
	return response.Result, nil
}
