package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentityContract_Encode(t *testing.T) {
	contract := NewIdentityContract()
	actual, err := contract.ChangeIdentity("", "", "")
	assert.NoError(t, err)
	expected := ""
	assert.Equal(t, expected, actual)
}
