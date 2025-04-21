package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
)

func NewNodeCertManagerContract() NodeCertManagerContract {
	return &nodeCertManagerContract{
		abi: abi.NewAbi(NodeCertManagerBuiltinContract.AbiString),
	}
}

type NodeCertManagerContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	// Revoke witness node cert
	Revoke(serialNumbers []string) (string, error)
}

type nodeCertManagerContract struct {
	abi abi.LatticeAbi
}

func (c *nodeCertManagerContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *nodeCertManagerContract) ContractAddress() string {
	return NodeCertManagerBuiltinContract.Address
}

func (c *nodeCertManagerContract) Revoke(serialNumbers []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("revoke", serialNumbers)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *nodeCertManagerContract) UploadKey(pubKeys []string) (string, error) {
	params := make([]interface{}, len(pubKeys))
	for i, nodePubKey := range pubKeys {
		param := make([]interface{}, 1)
		param[0] = nodePubKey
		params[i] = param
	}
	fn, err := c.abi.GetLatticeFunction("uploadKey", params)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *nodeCertManagerContract) Apply(certType uint8, orgName string, addresses []string) (string, error) {
	params := make([]interface{}, len(addresses))
	for i, node := range addresses {
		param := make([]interface{}, 1)
		param[0] = node
		params[i] = param
	}
	fn, err := c.abi.GetLatticeFunction("apply", certType, orgName, params)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var NodeCertManagerBuiltinContract = Contract{
	Description: "节点证书管理合约",
	Address:     "zltc_QLbz7JHxYJDL9LAguz9rKrwNtmfY2UoAZ",
	AbiString: `[
    {
        "inputs": [
            {
                "internalType": "uint8",
                "name": "certType",
                "type": "uint8"
            },
            {
                "internalType": "string",
                "name": "orgName",
                "type": "string"
            },
            {
                "components": [
                    {
                        "internalType": "address",
                        "name": "address",
                        "type": "address"
                    }
                ],
                "internalType": "struct ApplyCert[]",
                "name": "nodes",
                "type": "tuple[]"
            }
        ],
        "name": "apply",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint8",
                "name": "certType",
                "type": "uint8"
            },
            {
                "internalType": "string",
                "name": "orgName",
                "type": "string"
            },
            {
                "components": [
                    {
                        "internalType": "address",
                        "name": "address",
                        "type": "address"
                    }
                ],
                "internalType": "struct ApplyCert[]",
                "name": "nodes",
                "type": "tuple[]"
            }
        ],
        "name": "apply",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256[]",
                "name": "serialNumber",
                "type": "uint256[]"
            }
        ],
        "name": "revoke",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    },
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "bytes",
                        "name": "publicKey",
                        "type": "bytes"
                    }
                ],
                "internalType": "struct UploadKeyParam[]",
                "name": "nodes",
                "type": "tuple[]"
            }
        ],
        "name": "uploadKey",
        "outputs": [
            {
                "internalType": "bytes",
                "name": "",
                "type": "bytes"
            }
        ],
        "stateMutability": "pure",
        "type": "function"
    },
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "string",
                        "name": "OrgName",
                        "type": "string"
                    },
                    {
                        "internalType": "uint8",
                        "name": "CertType",
                        "type": "uint8"
                    },
                    {
                        "components": [
                            {
                                "internalType": "bytes",
                                "name": "PublicKey",
                                "type": "bytes"
                            }
                        ],
                        "internalType": "struct nodeCert.UploadKeyParam[]",
                        "name": "Applies",
                        "type": "tuple[]"
                    }
                ],
                "internalType": "struct nodeCert.UpAndApplyParam",
                "name": "param",
                "type": "tuple"
            }
        ],
        "name": "upAndApply",
        "outputs": [],
        "stateMutability": "pure",
        "type": "function"
    }
]`,
}
