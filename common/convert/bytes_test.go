package convert

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestZltcToAddress(t *testing.T) {
	zltcAddr := "zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6"
	ethAddr := "0x9293c604c644BfAc34F498998cC3402F203d4D6B"
	addr, err := ZltcToAddress(zltcAddr)
	assert.Nil(t, err)
	assert.Equal(t, ethAddr, addr.Hex())
}

func TestAddressToZltc(t *testing.T) {
	zltcAddr := "zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6"
	ethAddr := "0x9293c604c644BfAc34F498998cC3402F203d4D6B"
	addr := AddressToZltc(common.HexToAddress(ethAddr))
	assert.Equal(t, zltcAddr, addr)
}
