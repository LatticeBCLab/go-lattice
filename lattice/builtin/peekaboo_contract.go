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
	method := "delPayload"
	if visibility {
		method = "addPayload"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *peekabooContract) ToggleHash(hash string, visibility bool) (string, error) {
	method := "delHash"
	if visibility {
		method = "addHash"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *peekabooContract) ToggleCode(hash string, visibility bool) (string, error) {
	method := "delCode"
	if visibility {
		method = "addCode"
	}

	fn, err := c.abi.GetLatticeFunction(method, hash)
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
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addPayload",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delPayload",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addHash",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "addCode",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delHash",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"constant": false,
			"inputs": [
				{
					"internalType": "bytes32",
					"name": "_hash",
					"type": "bytes32"
				}
			],
			"name": "delCode",
			"outputs": [],
			"payable": false,
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
