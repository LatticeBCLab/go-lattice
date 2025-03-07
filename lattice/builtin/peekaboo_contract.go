package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
)

func NewPeekabooContract() PeekabooContract {
	return &peekabooContract{
		abi: abi.NewAbi(PeekabooBuiltinContract.AbiString),
	}
}

type PeekabooContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	TogglePayload(hash string, visibility bool) (string, error)
	ToggleHash(hash string, visibility bool) (string, error)
	ToggleCode(hash string, visibility bool) (string, error)
	BatchTogglePayload(hashes []string, visibility bool) (string, error)
}

type peekabooContract struct {
	abi abi.LatticeAbi
}

func (c *peekabooContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *peekabooContract) ContractAddress() string {
	return PeekabooBuiltinContract.Address
}

func (c *peekabooContract) TogglePayload(hash string, visibility bool) (string, error) {
	// 进行隐藏
	method := "addPayload"
	if visibility {
		// 删除隐藏记录
		method = "delPayload"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *peekabooContract) ToggleHash(hash string, visibility bool) (string, error) {
	method := "addHash"
	if visibility {
		method = "delHash"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *peekabooContract) ToggleCode(hash string, visibility bool) (string, error) {
	method := "addCode"
	if visibility {
		method = "delCode"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *peekabooContract) BatchTogglePayload(hashes []string, visibility bool) (string, error) {
	method := "batchAddPayload"
	if visibility {
		method = "batchDelPayload"
	}

	fn, err := c.abi.GetLatticeFunction(method, hashes)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var PeekabooBuiltinContract = Contract{
	Description: "区块隐藏合约",
	Address:     "zltc_a8Nx2gcs2XHye7MKVWykdanumqDkWXqRH",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addCode",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addHash",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addPayload",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32[]",
					"name": "_hash",
					"type": "bytes32[]"
				}
			],
			"name": "batchAddPayload",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32[]",
					"name": "_hash",
					"type": "bytes32[]"
				}
			],
			"name": "batchDelPayload",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delCode",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delHash",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delPayload",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
