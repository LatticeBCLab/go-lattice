package abi

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/LatticeBCLab/go-lattice/common/constant"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestDecodeReturn(t *testing.T) {
	myabi := FromJson("[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllAuthorizationIds\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllIdsAndAddresses\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"authorizationId\",\"type\":\"uint256\"}],\"name\":\"getAuthorizationDetails\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getAuthorizationsByUser\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"authorizedUser\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"expirationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"uses\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"ownerName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"authorizedName\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"assetName\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"isUnlimitedUses\",\"type\":\"bool\"}],\"internalType\":\"struct AssetAuthorization.GrantAuthorizationParams\",\"name\":\"params\",\"type\":\"tuple\"}],\"name\":\"grantAuthorizationAndMintNFT\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"authorizationId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"newExpirationTime\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"additionalUses\",\"type\":\"uint256\"}],\"name\":\"renewAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"authorizationId\",\"type\":\"uint256\"}],\"name\":\"revokeAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"authorizationId\",\"type\":\"uint256\"}],\"name\":\"useAuthorization\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
	result, err := DecodeReturn(myabi, "getAllIdsAndAddresses", "0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000004000000000000000000000000ae68b248eff0244e710191e92a2bf9cd74815974000000000000000000000000ae68b248eff0244e710191e92a2bf9cd74815974000000000000000000000000ae68b248eff0244e710191e92a2bf9cd74815974000000000000000000000000ae68b248eff0244e710191e92a2bf9cd74815974")
	assert.NoError(t, err)
	t.Log(result)
}

func TestDecodeGetAddress(t *testing.T) {
	myabi := FromJson("[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"protocolUri\",\"type\":\"uint64\"}],\"name\":\"getAddress\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"updater\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"data\",\"type\":\"bytes32[]\"}],\"internalType\":\"struct credibilidity.Protocol[]\",\"name\":\"protocol\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]")
	result, err := DecodeReturn(myabi, "getAddress", "0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000001e0000000000000000000000000000000000000000000000000000000000000034000000000000000000000000000000000000000000000000000000000000004a000000000000000000000000038d7001c4c3bcaa81f9073173531e12882bb0a8d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000873796e746178203d202270726f746f33223b0d0a0d0a6f7074696f6e206a6176615f7061636b616765203d2022636f6d2e7a6b6a672e7574696c6974792e70726f746f6275662e67656e657261746564223b0d0a6f7074696f6e206a6176615f6d756c7469706c655f66696c6573203d2066616c73653b0d0a6f7074696f6e206a6176615f6f757465725f636c6173736e616d65203d202254436572744948617368446562756750726f746f636f6c223b0d0a0d0a6d65737361676520544365727449486173684465627567207b0d0a09737472696e672048617368203d20313b0d0a7d0d0a0d0a00000000000000000000000000000000000000000000000000000000000000000000000038d7001c4c3bcaa81f9073173531e12882bb0a8d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000873796e746178203d202270726f746f33223b0d0a0d0a6f7074696f6e206a6176615f7061636b616765203d2022636f6d2e7a6b6a672e7574696c6974792e70726f746f6275662e67656e657261746564223b0d0a6f7074696f6e206a6176615f6d756c7469706c655f66696c6573203d2066616c73653b0d0a6f7074696f6e206a6176615f6f757465725f636c6173736e616d65203d202254436572744948617368446562756750726f746f636f6c223b0d0a0d0a6d65737361676520544365727449486173684465627567207b0d0a09696e7433322048617368203d20313b0d0a7d0d0a0d0a0000000000000000000000000000000000000000000000000000000000000000000000000038d7001c4c3bcaa81f9073173531e12882bb0a8d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000873796e746178203d202270726f746f33223b0d0a0d0a6f7074696f6e206a6176615f7061636b616765203d2022636f6d2e7a6b6a672e7574696c6974792e70726f746f6275662e67656e657261746564223b0d0a6f7074696f6e206a6176615f6d756c7469706c655f66696c6573203d2066616c73653b0d0a6f7074696f6e206a6176615f6f757465725f636c6173736e616d65203d202254436572744948617368446562756750726f746f636f6c223b0d0a0d0a6d65737361676520544365727449486173684465627567207b0d0a09737472696e672048617368203d20313b0d0a7d0d0a0d0a00000000000000000000000000000000000000000000000000000000000000000000000038d7001c4c3bcaa81f9073173531e12882bb0a8d0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000873796e746178203d202270726f746f33223b0d0a0d0a6f7074696f6e206a6176615f7061636b616765203d2022636f6d2e7a6b6a672e7574696c6974792e70726f746f6275662e67656e657261746564223b0d0a6f7074696f6e206a6176615f6d756c7469706c655f66696c6573203d2066616c73653b0d0a6f7074696f6e206a6176615f6f757465725f636c6173736e616d65203d202254436572744948617368446562756750726f746f636f6c223b0d0a0d0a6d65737361676520544365727449486173684465627567207b0d0a09737472696e672048617368203d20313b0d0a7d0d0a0d0a000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	var arr []*ProtocolContent
	err = json.Unmarshal([]byte(result[0]), &arr)
	if err != nil {
		assert.NoError(t, err)
	}
	for i, elem := range arr {
		arr[i].Version = uint32(i)
		arr[i].Message = string(Bytes32sToBytes(BytesArrToBytes32s(elem.Data)))
		arr[i].Updater = ToZLTCAddress(arr[i].Updater)
	}

	for _, elem := range arr {
		t.Log(elem.Message)
	}
}

type ProtocolContent struct {
	Updater string   `json:"updater,omitempty"` // 更新者账户地址
	Data    [][]byte `json:"data,omitempty"`    // 链上的原始数据 bytes32[] 类型
	Message string   `json:"message,omitempty"` // 解析后的协议内容
	Version uint32   `json:"version,omitempty"` // 协议的版本号
}

func Bytes32sToBytes(bytes32Arr [][32]byte) []byte {
	rows := len(bytes32Arr)
	bytes := make([]byte, (rows-1)*32)

	for i := 0; i < rows-1; i++ {
		copy(bytes[i*32:(i+1)*32], bytes32Arr[i][:])
	}

	arr := bytes32Arr[rows-1]
	for i := len(arr) - 1; i >= 0; i-- {
		if arr[i] != 0 {
			bytes = append(bytes, arr[:i+1]...)
			break
		}
	}

	return bytes
}

func ToZLTCAddress(address string) string {
	if strings.HasPrefix(address, constant.HexPrefix) {
		return convert.AddressToZltc(common.HexToAddress(address))
	}

	return address
}

func BytesArrToBytes32s(bytesArr [][]byte) [][32]byte {
	result := make([][32]byte, len(bytesArr))
	for i := 0; i < len(bytesArr); i++ {
		if len(bytesArr[i]) > 32 {
			copy(result[i][:], bytesArr[i][:32])
		} else {
			copy(result[i][:], bytesArr[i])
		}
	}
	return result
}
