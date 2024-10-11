package crypto

import (
	"fmt"
	"testing"

	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/stretchr/testify/assert"
)

func TestGenerateKeyPairForSm2p256v1(t *testing.T) {
	c := NewCrypto(types.Sm2p256v1)
	privateKey, err := c.GenerateKeyPair()
	assert.Nil(t, err)
	privateKeyHexString, _ := c.SKToHexString(privateKey)
	fmt.Println(privateKeyHexString)
	address, _ := c.PKToAddress(&privateKey.PublicKey)
	fmt.Println(convert.AddressToZltc(address))
}

func TestGenerateAddressFromPrivateKey(t *testing.T) {
	api := NewCrypto(types.Sm2p256v1)
	sk := "0x72ffdd7245e0ad7cffd533ad99f54048bf3fa6358e071fba8c2d7783d992d997"
	privateKey, _ := api.HexToSK(sk)
	address, _ := api.PKToAddress(&privateKey.PublicKey)
	zltc := convert.AddressToZltc(address)
	expect := "zltc_jF4U7umzNpiE8uU35RCBp9f2qf53H5CZZ"
	assert.Equal(t, zltc, expect)
}

func TestConvertPublicKeyToAddressForSm2p256v1(t *testing.T) {
	api := NewCrypto(types.Sm2p256v1)
	sk := "0x9860956de90cc61a05447ea067197be1fa08d712c4a5088c9cb62182bdca0f92"

	privateKey, _ := api.HexToSK(sk)

	address, _ := api.PKToAddress(&privateKey.PublicKey)
	actual := convert.AddressToZltc(address)
	expected := "zltc_oJCrxCx6X23m5xVZFLjexi8GGaib6Zwff"
	assert.Equal(t, expected, actual)
}

func TestConvertPublicKeyToAddressForSecp256k1(t *testing.T) {
	api := NewCrypto(types.Secp256k1)
	sk := "0xd2c784688ab85d689e358a7b030c9f26b8ee45e66e89d8842fa88da3b9637955"

	privateKey, _ := api.HexToSK(sk)
	fmt.Println(api.PKToHexString(&privateKey.PublicKey))
	address, _ := api.PKToAddress(&privateKey.PublicKey)
	actual := convert.AddressToZltc(address)
	expected := "zltc_cWAvRSgCKgfyp5Rz5TH8srmrZsH5fVYpg"
	assert.Equal(t, expected, actual)
}
