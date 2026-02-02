package sm2p256v1

import (
	"fmt"
	"testing"

	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestSm2p256v1Api_GenerateKeyPair(t *testing.T) {
	crypto := New()
	sk, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	t.Log(sk)
	skHexString, err := crypto.SKToHexString(sk)
	if err != nil {
		t.Error(err)
	}
	pkHexString, err := crypto.PKToHexString(&sk.PublicKey)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(skHexString)
	fmt.Println(pkHexString)
}

func TestSm2p256v1Api_HexToSK(t *testing.T) {
	sk := "0xb797270aab585e237498f203b6b7c85f812db111303330b7b306cebcfc08f913"
	crypto := New()
	priv, err := crypto.HexToSK(sk)
	if err != nil {
		t.Error(err)
	}
	crypto.Verify(
		hexutil.MustDecode("0x0d93d4aaaea5dc74818b6f33f95672a362e553b7272e646fa5b4a44cc6e9b833"),
		hexutil.MustDecode("0x24e2393dc5aba1b35616401d022a753650331d91bb6c26384dee9cdb3f8d522fef18a7d5031da24d6c60c16ca632674158cbbe5bf65ea0613db4ce57654be7850168b5f899747b1dfcd2f76118d8437713df31f250516a8de1473042228c0d6529"),
		&priv.PublicKey)
	fmt.Println(crypto.PKToHexString(&priv.PublicKey))
}

func TestSm2p256v1Api_Sign(t *testing.T) {
	crypto := New()
	sk, err := crypto.GenerateKeyPair()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(crypto.SKToHexString(sk))
	fmt.Println(crypto.PKToHexString(&sk.PublicKey))

	hash := []byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
	signature, err := crypto.Sign(hash, sk)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(hexutil.Encode(signature))
	}
	pass := crypto.Verify(hash, signature, &sk.PublicKey)
	fmt.Println(pass)
}

func TestSm2p256v1Api_PKToAddress(t *testing.T) {
	crypto := New()
	sk, err := crypto.BytesToSK(hexutil.MustDecode("0xca51dfee6b7337bd26c716931fa4a5c31eb7d91fa44bd254bad453d2bd0b815a"))
	assert.Nil(t, err)
	fmt.Println(crypto.PKToHexString(&sk.PublicKey))
	addr, err := crypto.PKToAddress(&sk.PublicKey)
	assert.Nil(t, err)
	t.Log("ETH Address:", addr.Hex())
	t.Log("ZLTC Address:", convert.AddressToZltc(addr))
	assert.Len(t, addr, 20)
	// 0x04376db6870f8ca937c94a49761db33fc71c5de643f07cb1501504644ef86360f7fb7974b11058c76a56c03bee897c0b5f640613cb6a3ff41fb23426d2b5e17cbb
}
