package builtin

import (
	"testing"

	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/stretchr/testify/assert"
)

func TestCredibilityContractEncode(t *testing.T) {
	contract := NewCredibilityContract()
	t.Run("Toggle data visibility", func(t *testing.T) {
		dataId := "1"
		businessAddress := "0x9293c604c644bfac34f498998cc3402f203d4d6b"
		data, err := contract.ToggleVisibility(dataId, businessAddress)
		assert.NoError(t, err)
		expected := "0xa2ec965700000000000000000000000000000000000000000000000000000000000000400000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b00000000000000000000000000000000000000000000000000000000000000013100000000000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, data)
	})
	t.Run("Batch toggle data visibility", func(t *testing.T) {
		actual, err := contract.BatchToggleVisibility([]ToggleVisibilityParam{
			{"1", convert.ZltcMustToAddress("zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6")},
			{"2", convert.ZltcMustToAddress("zltc_dhdfbm9JEoyDvYoCDVsABiZj52TAo9Ei6")},
		})
		assert.NoError(t, err)
		expected := "0x6726b1f800000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000400000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b0000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b00000000000000000000000000000000000000000000000000000000000000013200000000000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})

	t.Run("Create protocol", func(t *testing.T) {
		actual, err := contract.CreateProtocol(1, []byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}"))
		assert.NoError(t, err)
		expected := "0xef7e985800000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000373796e746178203d202270726f746f33223b0a0a6d6573736167652053747564656e74207b0a09737472696e67206964203d20313b0a09737472696e67206e616d65203d20323b0a7d0000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})

	t.Run("Batch create protocol", func(t *testing.T) {
		request := make([]CreateProtocolRequest, 2)
		request[0] = CreateProtocolRequest{
			10,
			convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
		}
		request[1] = CreateProtocolRequest{
			10,
			convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
		}
		data, err := contract.BatchCreateProtocol(request)
		assert.NoError(t, err)
		expectData := "0x589960aa0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000373796e746178203d202270726f746f33223b0a0a6d6573736167652053747564656e74207b0a09737472696e67206964203d20313b0a09737472696e67206e616d65203d20323b0a7d0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000373796e746178203d202270726f746f33223b0a0a6d6573736167652053747564656e74207b0a09737472696e67206964203d20313b0a09737472696e67206e616d65203d20323b0a7d0000000000000000000000000000000000000000000000"
		assert.Equal(t, expectData, data)
	})
	t.Run("Batch update protocol", func(t *testing.T) {
		request := make([]UpdateProtocolRequest, 1)
		request[0] = UpdateProtocolRequest{
			ProtocolUri: 8589934595,
			Data:        convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
		}
		data, err := contract.BatchUpdateProtocol(request)
		assert.NoError(t, err)
		expectData := "0x344b37ce00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000002000000030000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000373796e746178203d202270726f746f33223b0a0a6d6573736167652053747564656e74207b0a09737472696e67206964203d20313b0a09737472696e67206e616d65203d20323b0a7d0000000000000000000000000000000000000000000000"
		assert.Equal(t, expectData, data)
	})

	t.Run("Write data", func(t *testing.T) {
		addr, _ := convert.ZltcToAddress("zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D")
		println("addr:", addr.Hex())
		data, err := contract.Write(&WriteLedgerRequest{
			ProtocolUri: 8589934595,
			Hash:        "1",
			Data:        convert.BytesToBytes32Arr([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
			Address:     addr,
		})
		println("Data:", convert.BytesToBytes32Arr([]byte{1, 2, 3, 4, 5, 6, 7, 8}))
		assert.NoError(t, err)
		expectData := "0x4131ff530000000000000000000000000000000000000000000000000000000200000003000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000561717f7922a233720ae38acaa4174cda0bf17660000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010102030405060708000000000000000000000000000000000000000000000000"
		assert.Equal(t, expectData, data)
	})
	t.Run("Batch write data", func(t *testing.T) {
		addr, _ := convert.ZltcToAddress("zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D")
		batchRequest := make([]WriteLedgerRequest, 1)
		batchRequest[0] = WriteLedgerRequest{
			ProtocolUri: 8589934595,
			Hash:        "1",
			Data:        convert.BytesToBytes32Arr([]byte{1, 2, 3, 4, 5, 6, 7, 8}),
			Address:     addr,
		}
		data, err := contract.BatchWrite(batchRequest)
		assert.NoError(t, err)
		expectData := "0x77b34b730000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000200000003000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000561717f7922a233720ae38acaa4174cda0bf17660000000000000000000000000000000000000000000000000000000000000001310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010102030405060708000000000000000000000000000000000000000000000000"
		assert.Equal(t, expectData, data)
	})
}
