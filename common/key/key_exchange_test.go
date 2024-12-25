package key

import (
	"github.com/defiweb/go-eth/hexutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECDHEExchange_Exchange(t *testing.T) {
	akType := AKTypeBaaS
	// Client
	instance := NewECDHEExchange()
	clientSK, clientRandom, err := instance.GenerateSharedParams()
	assert.Nil(t, err)
	result, err := instance.Exchange(akType, &ExchangeParams{
		PK:     clientSK.PublicKey(),
		Random: clientRandom,
	})
	assert.Nil(t, err)
	t.Logf("ID: %s", result.AccessKey.ID)
	t.Logf("Secret: %s", result.AccessKey.Secret)

	// Client
	client := &ExchangeParams{
		SK:     clientSK,
		Random: clientRandom,
	}

	// Server
	server := &ExchangeParams{
		PK:     result.PK,
		Random: result.Random,
	}
	ak, err := instance.SelfExchange(akType, client, server)
	assert.Nil(t, err)

	err = instance.ConfirmAccessKeyIdOrigin(akType, ak.ID, ak.Secret)
	assert.Nil(t, err)

	// Assert
	assert.Equal(t, result.AccessKey.ID, ak.ID)
	assert.Equal(t, result.AccessKey.Secret, ak.Secret)

	assert.Equal(t, akType, instance.GetAKType(ak.ID))
}

func TestECDHEExchange_SelfExchange(t *testing.T) {
	localPrivateKeyHexString := "3f0283f538e7ef552b08321c1babadd37ddd98e5804ad473d738751b7bbb0005"
	localPublicKeyHexString := "047e670d3dd59beef947d1f18f1c419380d6fa9bb9ec861df705cc9313e0b9030071ca0a216ca60fed9ba7788744244987f3556233168a80a6c254f5924765f818"
	localRandom := "6561a1a82a0e7058a76ed9b340044e5700bb77173727dc75e9afcf9b6eb779fd"

	remotePrivateKeyHexString := "7f075479ac03e077804303f408d937153a0a45f2a430641326a7ebcc5f16daad"
	remotePublicKeyHexString := "04cf04189227c01b8c3a55a36dec979e927a42174e48d2b8274bc78a55b5b1c5f3f9c29b7e2b6adfc410e1283198e43fbf51e987958761682a7072473836d05639"
	remoteRandom := "37e086bfef8fdb5e89cc28318024cd6f982a8a8e50a5b55448eba39289fe29df"

	localSK, err := ECDHSKFromHex(localPrivateKeyHexString)
	assert.Nil(t, err)
	localPK, err := ECDHPKFromHex(localPublicKeyHexString)
	assert.Nil(t, err)
	client := &ExchangeParams{
		SK:     localSK,
		PK:     localPK,
		Random: hexutil.MustHexToBytes(localRandom),
	}

	// Server
	remoteSK, err := ECDHSKFromHex(remotePrivateKeyHexString)
	assert.Nil(t, err)
	remotePK, err := ECDHPKFromHex(remotePublicKeyHexString)
	assert.Nil(t, err)
	server := &ExchangeParams{
		SK:     remoteSK,
		PK:     remotePK,
		Random: hexutil.MustHexToBytes(remoteRandom),
	}

	instance := NewECDHEExchange()
	ak, err := instance.SelfExchange(AKTypeBaaS, client, server)
	assert.Nil(t, err)

	expectedAccessKeyId := "2Xi9K4G7u9E3FtyM4iGqEjH"
	expectedAccessKeySecret := "6EdHs9VtuG5YgPYs99tcKm"
	assert.Equal(t, expectedAccessKeyId, ak.ID)
	assert.Equal(t, expectedAccessKeySecret, ak.Secret)
}
