package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNodeCertManagerContract_Encode(t *testing.T) {
	contract := NewNodeCertManagerContract()

	t.Run("Revoke cert", func(t *testing.T) {
		sn := "109313213476126583182290498689290247244"
		actual, err := contract.Revoke([]string{sn})
		assert.NoError(t, err)
		expect := "0x420b3bd10000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000523cf58890fa44499270fb81ac549c4c"
		assert.Equal(t, expect, actual)
	})
}
