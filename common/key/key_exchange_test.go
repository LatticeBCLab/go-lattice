package key

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestECDHEExchange_Exchange(t *testing.T) {
	// Client
	instance := NewECDHEExchange()
	clientSK, clientRandom, err := instance.GenerateSharedParams()
	assert.Nil(t, err)
	result, err := instance.Exchange(&ExchangeParams{
		PK:     clientSK.PublicKey(),
		Random: clientRandom,
	})
	assert.Nil(t, err)
	t.Logf("ID: %s", result.AccessKey.ID)
	t.Logf("Secret: %s", result.AccessKey.Secret)

	// Server
	server := &ExchangeParams{
		PK:     result.PK,
		Random: result.Random,
	}
	ak, err := instance.ClientExchange(clientSK, clientRandom, server)
	assert.Nil(t, err)

	// Assert
	assert.Equal(t, result.AccessKey.ID, ak.ID)
	assert.Equal(t, result.AccessKey.Secret, ak.Secret)
}
