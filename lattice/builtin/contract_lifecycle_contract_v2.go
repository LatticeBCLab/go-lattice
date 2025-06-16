package builtin

import "github.com/LatticeBCLab/go-lattice/abi"

func NewContractLifecycleContractV2() ContractLifecycleContractV2 {
	return &contractLifecycleContractV2{
		abi: abi.NewAbi(ContractLifecycleBuiltinContractV2.AbiString),
	}
}

type ContractLifecycleContractV2 interface {
	ContractAddress() string
	Freeze(contractAddress string) (string, error)
	Unfreeze(contractAddress string) (string, error)
}

type contractLifecycleContractV2 struct {
	abi abi.LatticeAbi
}

func (c *contractLifecycleContractV2) ContractAddress() string {
	return ContractLifecycleBuiltinContractV2.Address
}

func (c *contractLifecycleContractV2) Freeze(contractAddress string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, true)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *contractLifecycleContractV2) Unfreeze(contractAddress string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("launch", contractAddress, false)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var ContractLifecycleBuiltinContractV2 = Contract{
	Description: "合约生命周期合约版本2",
	Address:     "zltc_ZQJjaw74CKMjqYJFMKdEDaNTDMq5QKi3T",
	AbiString: `[
    {
        "inputs": [
            {
                "internalType": "address",
                "name": "Address",
                "type": "address"
            },
            {
                "internalType": "bool",
                "name": "IsFreeze",
                "type": "bool"
            }
        ],
        "name": "launch",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    }
]`,
}
