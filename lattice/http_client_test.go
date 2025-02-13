package lattice

import (
	"context"
	"fmt"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func setupHttpClient() client.HttpApi {
	connectingNodeConfig := &ConnectingNodeConfig{Ip: "172.22.0.23", HttpPort: 10332}
	initHttpClientArgs := &client.HttpApiInitParam{
		NodeAddress:                fmt.Sprintf("%s:%d", connectingNodeConfig.Ip, connectingNodeConfig.HttpPort),
		HttpUrl:                    connectingNodeConfig.GetHttpUrl(),
		GinServerUrl:               connectingNodeConfig.GetGinServerUrl(),
		Transport:                  (&Options{MaxIdleConnsPerHost: 200}).GetTransport(),
		JwtSecret:                  connectingNodeConfig.JwtSecret,
		JwtTokenExpirationDuration: connectingNodeConfig.JwtTokenExpirationDuration,
	}
	httpApi := client.NewHttpApi(initHttpClientArgs)

	return httpApi
}

func TestHttpClientRequest(t *testing.T) {
	httpApi := setupHttpClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	t.Run("get node certificate", func(t *testing.T) {
		wrappedCertificate, err := httpApi.GetNodeCertificate(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, wrappedCertificate)
		t.Log(wrappedCertificate.Type.ToChinese())
		t.Log(wrappedCertificate.PEMCertificate)
	})

	t.Run("get peer node certificate by serial number", func(t *testing.T) {
		serialNumber := "55415987022000890681526953338523518501"
		wrappedCertificate, err := httpApi.GetPeerNodeCertificate(ctx, serialNumber)
		assert.NoError(t, err)
		assert.NotNil(t, wrappedCertificate)
		t.Log(wrappedCertificate.Type.ToChinese())
		t.Log(wrappedCertificate.PEMCertificate)
	})

	t.Run("get peer node certificate by address", func(t *testing.T) {
		nodeAddress := "zltc_mXhpBG3X1dezEPumnfSQya38Awm654dXT"
		wrappedCertificate, err := httpApi.GetPeerNodeCertificateByAddress(ctx, nodeAddress)
		assert.NoError(t, err)
		assert.NotNil(t, wrappedCertificate)
		t.Log(wrappedCertificate.Type.ToChinese())
		t.Log(wrappedCertificate.PEMCertificate)
	})
}

func TestBlock(t *testing.T) {
	blockchainId := "1"
	httpApi := setupHttpClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	t.Run("get tblock state", func(t *testing.T) {
		hash := "0x9c3699c06cad8e50c7d2b59b1e1eefde3ad7a650870fb14efdfc5e9f26645510"
		state, err := httpApi.GetTBlockState(ctx, blockchainId, hash)
		assert.Nil(t, err)
		assert.Equal(t, state, types.TBlockStateDAEMONIZED)
	})
}

func TestSubchain(t *testing.T) {
	httpApi := setupHttpClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	t.Run("get created subchain", func(t *testing.T) {
		subchainIds, err := httpApi.GetCreatedSubchain(ctx)
		assert.Nil(t, err)
		t.Log(subchainIds)
	})
}
