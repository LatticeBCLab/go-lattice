package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/samber/lo"
)

type nodeCertificateContract struct {
	abi abi.LatticeAbi
}

// NodeCertificateType 节点证书类型
type NodeCertificateType uint8

const (
	Unspecified NodeCertificateType = iota
	InitConsensus
	InitClient
	Consensus
	Client
)

type RevokeNodeCertificateParam struct {
	SerialNumber uint32 `json:"serialNumber"`
	Client       string `json:"client"`
}

func NewNodeCertificateContract() NodeCertificateContract {
	return &nodeCertificateContract{
		abi: abi.NewAbi(NodeCertificateBuiltinContract.AbiString),
	}
}

type NodeCertificateContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	// Apply 申请证书
	Apply(certificateType NodeCertificateType, orgName string, nodes []string) (string, error)
	// Revoke 吊销证书
	Revoke(params []*RevokeNodeCertificateParam) (string, error)
	// UploadPublicKey 上传公钥
	UploadPublicKey(publicKeys []string) (string, error)
	// UploadPublicKeyAndApplyCertificate 上传公钥并申请证书
	UploadPublicKeyAndApplyCertificate(certificateType NodeCertificateType, orgName string, publicKeys []string) (string, error)
}

func (c *nodeCertificateContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *nodeCertificateContract) ContractAddress() string {
	return NodeCertificateBuiltinContract.Address
}

func (c *nodeCertificateContract) Apply(certificateType NodeCertificateType, orgName string, nodes []string) (string, error) {
	code, err := c.abi.RawAbi().Pack("apply", certificateType, orgName, lo.Map(nodes, func(item string, index int) struct {
		Address common.Address
	} {
		return struct{ Address common.Address }{Address: convert.ZltcMustToAddress(item)}
	}))
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *nodeCertificateContract) Revoke(params []*RevokeNodeCertificateParam) (string, error) {
	fn, err := c.abi.GetLatticeFunction("revoke",
		lo.Map(params, func(param *RevokeNodeCertificateParam, index int) uint32 {
			return param.SerialNumber
		}),
		lo.Map(params, func(param *RevokeNodeCertificateParam, index int) string {
			return param.Client
		}),
	)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *nodeCertificateContract) UploadPublicKey(publicKeys []string) (string, error) {
	code, err := c.abi.RawAbi().Pack("uploadKey", lo.Map(publicKeys, func(item string, index int) struct {
		PublicKey []byte
	} {
		return struct{ PublicKey []byte }{PublicKey: hexutil.MustDecode(item)}
	}))
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *nodeCertificateContract) UploadPublicKeyAndApplyCertificate(certificateType NodeCertificateType, orgName string, publicKeys []string) (string, error) {
	code, err := c.abi.RawAbi().Pack("upAndApply", struct {
		CertType NodeCertificateType
		OrgName  string
		Applies  []struct {
			PublicKey []byte
		}
	}{
		CertType: certificateType,
		OrgName:  orgName,
		Applies: lo.Map(publicKeys, func(item string, index int) struct {
			PublicKey []byte
		} {
			return struct{ PublicKey []byte }{PublicKey: hexutil.MustDecode(item)}
		}),
	})
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

var NodeCertificateBuiltinContract = Contract{
	Description: "节点证书合约",
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
                "internalType": "uint256[]",
                "name": "serialNumber",
                "type": "uint256[]"
            },
            {
                "internalType": "address[]",
                "name": "clients",
                "type": "address[]"
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
                "name": "applies",
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
                                "internalType": "bytes",
                                "name": "publicKey",
                                "type": "bytes"
                            }
                        ],
                        "internalType": "struct nodeCert.UploadKeyParam[]",
                        "name": "applies",
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
