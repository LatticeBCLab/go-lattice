package abi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatticeFunction_ConvertArguments(t *testing.T) {
	abi := NewAbi(complexABI)

	t.Run("construction", func(t *testing.T) {
		fn := abi.GetConstructor(`[1001,Jack,0x9293c604c644bfac34f498998cc3402f203d4d6b]`)
		actual, err := fn.Encode()
		assert.NoError(t, err)
		t.Log(actual)
		expected := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000003e900000000000000000000000000000000000000000000000000000000000000600000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b00000000000000000000000000000000000000000000000000000000000000044a61636b00000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})

	t.Run("add student", func(t *testing.T) {
		fn, _ := abi.GetLatticeFunction("addStudent", "1001", "Jack", "0x9293c604c644bfac34f498998cc3402f203d4d6b")
		actual, err := fn.Encode()
		assert.NoError(t, err)
		expected := "0xe7d0f6ee00000000000000000000000000000000000000000000000000000000000003e900000000000000000000000000000000000000000000000000000000000000600000000000000000000000009293c604c644bfac34f498998cc3402f203d4d6b00000000000000000000000000000000000000000000000000000000000000044a61636b00000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})

	t.Run("get student slice", func(t *testing.T) {
		fn, _ := abi.GetLatticeFunction("getStudentsSlice", "[1,2,3]")
		actual, err := fn.Encode()
		assert.NoError(t, err)
		expected := "0xc6fc9df000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000003"
		assert.Equal(t, expected, actual)
	})

	t.Run("update student name", func(t *testing.T) {
		fn, _ := abi.GetLatticeFunction("updateStudentName", "1002", "Tom")
		actual, err := fn.Encode()
		assert.NoError(t, err)
		expected := "0x440dba6c00000000000000000000000000000000000000000000000000000000000003ea00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000003546f6d0000000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})
}

func TestLatticeFunction_EncodeConstructor(t *testing.T) {
	abi := `[{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"get","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"index","type":"uint256"}],"name":"getPeople","outputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getPeoples","outputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person[]","name":"","type":"tuple[]"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"people","outputs":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_name","type":"string"}],"name":"set","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"string","name":"name","type":"string"},{"internalType":"uint256","name":"age","type":"uint256"}],"internalType":"struct HelloWorld.Person","name":"_person","type":"tuple"}],"name":"setPeople","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	data, _ := NewAbi(abi).GetConstructor("jack").Encode()
	expectData := "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000046a61636b00000000000000000000000000000000000000000000000000000000"
	assert.Equal(t, expectData, data)
}

var complexABI = `[
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "uint256",
                        "name": "id",
                        "type": "uint256"
                    },
                    {
                        "internalType": "string",
                        "name": "name",
                        "type": "string"
                    },
                    {
                        "internalType": "address",
                        "name": "wallet",
                        "type": "address"
                    }
                ],
                "internalType": "struct StudentManager.Student",
                "name": "_firstStudent",
                "type": "tuple"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "constructor"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "id",
                "type": "uint256"
            },
            {
                "indexed": false,
                "internalType": "string",
                "name": "name",
                "type": "string"
            },
            {
                "indexed": false,
                "internalType": "address",
                "name": "wallet",
                "type": "address"
            }
        ],
        "name": "StudentAdded",
        "type": "event"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_id",
                "type": "uint256"
            },
            {
                "internalType": "string",
                "name": "_name",
                "type": "string"
            },
            {
                "internalType": "address",
                "name": "_wallet",
                "type": "address"
            }
        ],
        "name": "addStudent",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256[]",
                "name": "_ids",
                "type": "uint256[]"
            },
            {
                "internalType": "string[]",
                "name": "_names",
                "type": "string[]"
            },
            {
                "internalType": "address[]",
                "name": "_wallets",
                "type": "address[]"
            }
        ],
        "name": "addStudentsBatch",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_id",
                "type": "uint256"
            }
        ],
        "name": "getStudentData",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "bytes32",
                "name": "_hashedId",
                "type": "bytes32"
            }
        ],
        "name": "getStudentWallet",
        "outputs": [
            {
                "internalType": "address",
                "name": "",
                "type": "address"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256[]",
                "name": "_ids",
                "type": "uint256[]"
            }
        ],
        "name": "getStudentsSlice",
        "outputs": [
            {
                "components": [
                    {
                        "internalType": "uint256",
                        "name": "id",
                        "type": "uint256"
                    },
                    {
                        "internalType": "string",
                        "name": "name",
                        "type": "string"
                    },
                    {
                        "internalType": "address",
                        "name": "wallet",
                        "type": "address"
                    }
                ],
                "internalType": "struct StudentManager.Student[]",
                "name": "",
                "type": "tuple[]"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_id",
                "type": "uint256"
            },
            {
                "internalType": "bytes",
                "name": "_data",
                "type": "bytes"
            }
        ],
        "name": "storeStudentData",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "name": "students",
        "outputs": [
            {
                "internalType": "uint256",
                "name": "id",
                "type": "uint256"
            },
            {
                "internalType": "string",
                "name": "name",
                "type": "string"
            },
            {
                "internalType": "address",
                "name": "wallet",
                "type": "address"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_id",
                "type": "uint256"
            },
            {
                "components": [
                    {
                        "internalType": "uint256",
                        "name": "id",
                        "type": "uint256"
                    },
                    {
                        "internalType": "string",
                        "name": "name",
                        "type": "string"
                    },
                    {
                        "internalType": "address",
                        "name": "wallet",
                        "type": "address"
                    }
                ],
                "internalType": "struct StudentManager.Student",
                "name": "_newInfo",
                "type": "tuple"
            }
        ],
        "name": "updateStudentInfo",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_id",
                "type": "uint256"
            },
            {
                "internalType": "string",
                "name": "_newName",
                "type": "string"
            }
        ],
        "name": "updateStudentName",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`
