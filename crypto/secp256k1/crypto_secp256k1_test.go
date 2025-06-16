package secp256k1

import (
	"fmt"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSecp256k1(t *testing.T) {
	ecc := New()
	t.Run("Generate key", func(t *testing.T) {
		sk, err := ecc.GenerateKeyPair()
		if err != nil {
			t.Error(err)
		}
		t.Log(sk)
		privateKey, _ := ecc.SKToHexString(sk) // 0x2235f144e584bd44eaef4b24297e3c2cb8a000a02f03314400fd7777c9824864
		fmt.Println("Private key", privateKey)
		publicKey, _ := ecc.PKToHexString(&sk.PublicKey)
		fmt.Println("Public key", publicKey) // 0x041362bd30468a548040df7cec7c5523cbb904c2b958cb4a21b82f35fe057a271f6d4fc9f042db35dc7aa74d861cfa57cf0948ff7a9d445dfa24467555f367673b
		address, _ := ecc.PKToAddress(&sk.PublicKey)
		fmt.Println(convert.AddressToZltc(address)) // zltc_WRR5eQmxPy5F9vzzr9WFi57xHiTvww5mV
	})

	t.Run("Get address from public key", func(t *testing.T) {
		publicKey := "0x04aee72e169c300700f4f7121b2885cc79ff6cf5abe862f2c6a72891fb62365ee2cef1f91e9da30060e37c22f8048095cb934439bbe6473239d3070608a4ade7dc"
		pk, _ := ecc.BytesToPK(hexutil.MustHexToBytes(publicKey))
		address, _ := ecc.PKToAddress(pk)
		fmt.Println("address", convert.AddressToZltc(address))
	})

	t.Run("Verify signature use given public key", func(t *testing.T) {
		publicKey := "0x04aee72e169c300700f4f7121b2885cc79ff6cf5abe862f2c6a72891fb62365ee2cef1f91e9da30060e37c22f8048095cb934439bbe6473239d3070608a4ade7dc"
		signature := "0xa3a060c88c45f8b024a9dd0719bdce15828eadf9d77afb38f39e87867823b01b75785047bdf74e523440e29fc29f13cbe4b8830b9680bf16d5d496a48984f2f0"
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
		pk, _ := ecc.BytesToPK(hexutil.MustHexToBytes(publicKey))
		pass := ecc.Verify(data, hexutil.MustHexToBytes(signature), pk)
		assert.Equal(t, true, pass)
	})

	t.Run("Sign and verify", func(t *testing.T) {
		sk, err := ecc.GenerateKeyPair()
		assert.Nil(t, err)
		hash := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
		signature, err := ecc.Sign(hash, sk)
		fmt.Println(hexutil.BytesToHex(signature))
		assert.Nil(t, err)
		passed := ecc.Verify(hash, signature, &sk.PublicKey)
		assert.True(t, passed)
	})
}
