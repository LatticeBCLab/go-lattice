package wallet

import (
	"encoding/json"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileKey(t *testing.T) {
	curve := types.Sm2p256v1

	t.Run("Generate file key", func(t *testing.T) {
		passphrase := "Aa123456"
		privateKey := "0xbd7ea728f7e6240507b321cb4a937a8d34ecfd39c275dbacf31ddb4793691dcc"
		fileKey, err := GenerateFileKey(privateKey, passphrase, curve)
		assert.Nil(t, err)
		bs, err := json.Marshal(fileKey)
		assert.Nil(t, err)
		t.Log(string(bs))
	})

	t.Run("Decrypt file key", func(t *testing.T) {
		fileKeyString := `{"uuid":"bd50b183-cc93-4d12-b820-0980d26f5de5","address":"zltc_j5yLhxm8fkwJkuhapqmqmJ1vYY2gLfPLy","cipher":{"aes":{"cipher":"aes-128-ctr","iv":"d5dfd36bd6447d5b9ce4875e5f13abd8"},"kdf":{"kdf":"scrypt","kdfParams":{"DKLen":32,"n":262144,"p":1,"r":8,"salt":"838286b54f1c5a3be40e94bcae5211f75d8128d72aec7b0953c8127b3bfd0e08"}},"cipherText":"86153ac036fd29858fffdf2d3fc2b2b05d35e2d6c38f15c4ccf6efd3122f8a6a","mac":"3377046f1a81aad71d9944b61430d253e22d2ba6435296ee3cfa0b9a1d1d6574"},"isGM":true}`
		password := "Root1234"
		sk, err := NewFileKey(fileKeyString).Decrypt(password)
		assert.Nil(t, err)
		skString, err := crypto.NewCrypto(curve).SKToHexString(sk)
		assert.Nil(t, err)
		assert.Equal(t, skString, "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb")
	})
}
