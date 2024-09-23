package crypto

import (
	"fmt"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCrypto(t *testing.T) {
	c := NewCrypto(types.Sm2p256v1)
	sk := "0xb58ee7d18f8ea223e8f4ca11cd813d3122990a354355f7b25f4891aa1be0ff2b"
	pk := "0x043bfd529f0827940b4130fc700e17d17e4f40ba38fd0006cc6a6f923da8139e05393ab1699638f80a84d4b3478205c7d99d84c58d5e8ac71a9fa69b2d2736fcbb"

	data := []byte("Hello World")
	cipher, err := c.Encrypt(data, pk)
	assert.Nil(t, err)

	source, err := c.Decrypt(cipher, sk)
	assert.Nil(t, err)
	assert.Equal(t, data, source)
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
