package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeekabooContract_Encode(t *testing.T) {
	contract := NewPeekabooContract()

	t.Run("Hide code", func(t *testing.T) {
		hash := "0x85d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
		actual, err := contract.ToggleCode(hash, true)
		assert.NoError(t, err)
		expect := "0xfdb2787485d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
		assert.Equal(t, expect, actual)
	})

	t.Run("Show code", func(t *testing.T) {
		hash := "0x85d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
		actual, err := contract.ToggleCode(hash, false)
		assert.NoError(t, err)
		expect := "0xbd63758085d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
		assert.Equal(t, expect, actual)
	})

	t.Run("Encode batch dele payload", func(t *testing.T) {
		hashes := []string{"0x85d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"}
		code, err := contract.BatchTogglePayload(hashes, true)
		expect := "0x8679bb350000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000185d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
		assert.NoError(t, err)
		assert.Equal(t, expect, code)
	})
}
