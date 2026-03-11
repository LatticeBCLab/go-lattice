package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/samber/lo"
)

func NewIdentityContract() IdentityContract {
	return &identityContract{
		abi: abi.NewAbi(IdentityBuiltinContract.AbiString),
	}
}

type IdentityContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	// ChangeIdentity 更换身份
	ChangeIdentity(oldAddress, newAddress, data string) (string, error)
	// CreateDID 创建DID
	CreateDID(oldAddress, doc string) (string, error)
}

type identityContract struct {
	abi abi.LatticeAbi
}

func (c *identityContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *identityContract) ContractAddress() string {
	return IdentityBuiltinContract.Address
}

func (c *identityContract) ChangeIdentity(oldAddress, newAddress, data string) (string, error) {
	function, err := c.abi.GetLatticeFunction("changeIdentity", oldAddress, newAddress, lo.Ternary(data == "", "identity", data))
	if err != nil {
		return "", err
	}
	return function.Encode()
}

func (c *identityContract) CreateDID(oldAddress, doc string) (string, error) {
	function, err := c.abi.GetLatticeFunction("createDID", oldAddress, doc)
	if err != nil {
		return "", err
	}
	return function.Encode()
}

var IdentityBuiltinContract = Contract{
	Description: "身份合约",
	Address:     "zltc_aQdmesGLjoJ5FJ65t2F7Nf9tTAT2C3dxA",
	AbiString: `[
	{
		"inputs": [
			{
				"internalType": "address",
				"name": "_oldAddress",
				"type": "address"
			},
			{
				"internalType": "address",
				"name": "_newAddress",
				"type": "address"
			},
			{
				"internalType": "string",
				"name": "_data",
				"type": "string"
			}
		],
		"name": "changeIdentity",
		"outputs": [],
		"stateMutability": "pure",
		"type": "function"
	},{
		"inputs": [
			{
				"internalType": "address",
				"name": "_oldAddress",
				"type": "address"
			},
			{
				"internalType": "string",
				"name": "_doc",
				"type": "string"
			}
		],
		"name": "createDID",
		"outputs": [
			{
				"internalType": "string",
				"name": "_idbHash",
				"type": "string"
			},
			{
				"internalType": "string",
				"name": "_diddocHash",
				"type": "string"
			}
		],
		"stateMutability": "pure",
		"type": "function"
}
]`,
}
