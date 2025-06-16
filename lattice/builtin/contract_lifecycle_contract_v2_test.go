package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContractLifecycleContractV2(t *testing.T) {
	contract := NewContractLifecycleContractV2()

	t.Run("Freeze contract", func(t *testing.T) {
		code, err := contract.Freeze("zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6")
		assert.NoError(t, err)
		expectedCode := "0xb38c170b0000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b0000000000000000000000000000000000000000000000000000000000000001"
		assert.Equal(t, expectedCode, code)
	})

	t.Run("Unfreeze contract", func(t *testing.T) {
		code, err := contract.Unfreeze("0x9293c604c644BfAc34F498998cC3402F203d4D6B")
		assert.NoError(t, err)
		expectedCode := "0xb38c170b0000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b0000000000000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expectedCode, code)
	})
}
