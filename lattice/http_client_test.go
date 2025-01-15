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
	connectingNodeConfig := &ConnectingNodeConfig{Ip: "192.168.3.51", HttpPort: 13000}
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
