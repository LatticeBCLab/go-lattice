package lattice

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/client"
	"github.com/stretchr/testify/assert"
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

	t.Run("connect node async", func(t *testing.T) {
		inode := "/ip4/172.22.0.23/tcp/39901/p2p/16Uiu2HAkyXyfmor1mdzWR8qpH36GRppMHPyDJdNrRo7XQLXGZqyU"
		err := httpApi.ConnectNodeAsync(ctx, inode)
		assert.NoError(t, err)
	})

	t.Run("connect peer async", func(t *testing.T) {
		id := "16Uiu2HAkyXyfmor1mdzWR8qpH36GRppMHPyDJdNrRo7XQLXGZqyU"
		err := httpApi.ConnectPeerAsync(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("disconnect node async", func(t *testing.T) {
		id := "16Uiu2HAkyXyfmor1mdzWR8qpH36GRppMHPyDJdNrRo7XQLXGZqyU"
		err := httpApi.DisconnectPeerAsync(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("get current transaction block", func(t *testing.T) {
		accountAddress := "zltc_j5yLhxm8fkwJkuhapqmqmJ1vYY2gLfPLy"
		block, err := httpApi.GetCurrentTBlock(ctx, chainId, accountAddress)
		assert.NoError(t, err)
		t.Log(block)
	})

	t.Run("get tblock by height", func(t *testing.T) {
		accountAddress := "zltc_j5yLhxm8fkwJkuhapqmqmJ1vYY2gLfPLy"
		block, err := httpApi.GetTBlockByHeight(ctx, chainId, accountAddress, 1)
		assert.NoError(t, err)
		t.Log(block)
	})

	t.Run("get balance with pending", func(t *testing.T) {
		accountAddress := "zltc_j5yLhxm8fkwJkuhapqmqmJ1vYY2gLfPLy"
		balance, err := httpApi.GetBalanceWithPending(ctx, chainId, accountAddress)
		assert.NoError(t, err)
		t.Log(balance)
	})

	t.Run("get genesis block", func(t *testing.T) {
		block, err := httpApi.GetGenesisBlock(ctx, chainId)
		assert.NoError(t, err)
		t.Log(block)
	})

	t.Run("get tblocks by heights", func(t *testing.T) {
		txs, err := httpApi.GetTBlocksByHeights(ctx, chainId, "zltc_oLW2u5m7zLVYRhipEtqix1aRkPn18ommF", []uint64{1})
		assert.NoError(t, err)
		t.Log(txs)
	})

	t.Run("get dblocks by heights", func(t *testing.T) {
		blocks, err := httpApi.GetDBlocksByHeights(ctx, chainId, []uint64{1})
		assert.NoError(t, err)
		t.Log(blocks)
	})

	t.Run("get receipts", func(t *testing.T) {
		receipts, err := httpApi.GetReceipts(ctx, chainId, []string{"0xc043e226b626580cf0fa9d2245999f6cd9ae669ab1443caf7bffe94e6019530e"})
		assert.NoError(t, err)
		t.Log(receipts)
	})

	t.Run("get genesis", func(t *testing.T) {
		block, err := httpApi.GetGenesisBlock(ctx, chainId)
		assert.NoError(t, err)
		t.Log(block)
	})

	t.Run("proxy re-encryption", func(t *testing.T) {
		result, err := httpApi.ProxyReEncryption(
			ctx,
			chainId,
			"0102",
			"zltc_Vsv3TDAxHpiKQxZAS2ctaJr6UnGkp7mkT",
			"zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi",
			"zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi",
		)
		assert.NoError(t, err)
		t.Log(result)
	})
}

func TestConfig(t *testing.T) {
	httpApi := setupHttpClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	t.Run("get node configuration", func(t *testing.T) {
		config, err := httpApi.GetNodeConfig(ctx, "")
		assert.Nil(t, err)
		t.Logf("%+v", config)
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
