package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeekabooContract_ToggleCode(t *testing.T) {
	contract := NewPeekabooContract()
	hash := "0x85d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
	code1, err := contract.ToggleCode(hash, true)
	assert.NoError(t, err)
	expect1 := "0xbd63758085d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
	assert.Equal(t, expect1, code1)

	code2, err := contract.ToggleCode(hash, false)
	assert.NoError(t, err)
	expect2 := "0xfdb2787485d6713b14264b030b7c4f71224691d866c0058e9cdba5e0cc417b47d6003ecd"
	assert.Equal(t, expect2, code2)
}
