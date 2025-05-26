package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/samber/lo"
)

func NewProxyReEncryptionContract() ProxyReEncryptionContract {
	return &proxyReEncryptionContract{
		abi: abi.NewAbi(ProxyReEncryptionBuiltinContract.AbiString),
	}
}

type ProxySecret struct {
	Whitelist string `json:"WhiteList"` // 白名单用户
	Cipher    string `json:"Sk"`        // 代理重加密密钥片段的密文
	Proxy     string `json:"Proxy"`     // 代理节点地址
}

type ProxyReEncryptionContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	// Sharding 返回某个白名单的所有片段
	Sharding(businessId, initiator, whitelist string) (string, error)
	// RemoveProxySecret 删除代理重加密密钥，会删除所有的代理片段，支持批量删除
	RemoveProxySecret(businessId, initiator string, whitelists []string) (string, error)
	// SelectProxySecret 查找某个代理节点的某个白名单的片段
	SelectProxySecret(proxy, businessId, initiator, whitelist string) (string, error)
	// StoreProxySecret 存储代理重加密密钥片段，支持批量新增
	StoreProxySecret(secrets []ProxySecret, businessId, initiator string) (string, error)
	// UpdateProxySecret 更新所有代理的片段，本质是先调用删除 RemoveProxySecret，然后新增 1 StoreProxySecret
	UpdateProxySecret(secrets []ProxySecret, businessId, initiator string, whitelists []string) (string, error)
}

type proxyReEncryptionContract struct {
	abi abi.LatticeAbi
}

func (c *proxyReEncryptionContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *proxyReEncryptionContract) ContractAddress() string {
	return ProxyReEncryptionBuiltinContract.Address
}

func (c *proxyReEncryptionContract) Sharding(businessId, initiator, whitelist string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("sharding", businessId, initiator, whitelist)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *proxyReEncryptionContract) RemoveProxySecret(businessId, initiator string, whitelists []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("removeProxySecret", businessId, initiator, whitelists)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *proxyReEncryptionContract) SelectProxySecret(proxy, businessId, initiator, whitelist string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("selectProxySecret", proxy, businessId, initiator, whitelist)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *proxyReEncryptionContract) StoreProxySecret(secrets []ProxySecret, businessId, initiator string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("storeProxySecret",
		lo.Map(secrets, func(item ProxySecret, index int) []interface{} {
			return []interface{}{item.Whitelist, item.Cipher, item.Proxy}
		}),
		businessId,
		initiator,
	)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

func (c *proxyReEncryptionContract) UpdateProxySecret(secrets []ProxySecret, businessId, initiator string, whitelists []string) (string, error) {
	fn, err := c.abi.GetLatticeFunction("updateProxySecret",
		lo.Map(secrets, func(item ProxySecret, index int) []interface{} {
			return []interface{}{item.Whitelist, item.Cipher, item.Proxy}
		}),
		businessId,
		initiator,
		whitelists,
	)
	if err != nil {
		return "", err
	}
	return fn.Encode()
}

var ProxyReEncryptionBuiltinContract = &Contract{
	Description: "代理重加密预置合约",
	Address:     "zltc_QLbz7JHiBTspVQ2d8UCmyL1PGWzZnrVHJ",
	AbiString: `[
    {
        "inputs": [
            {
                "internalType": "string",
                "name": "businessId",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "initiator",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "whiteList",
                "type": "string"
            }
        ],
        "name": "sharding",
        "outputs": [
            {
                "internalType": "string[]",
                "name": "",
                "type": "string[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "string",
                "name": "businessId",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "initiator",
                "type": "string"
            },
            {
                "internalType": "string[]",
                "name": "whiteList",
                "type": "string[]"
            }
        ],
        "name": "removeProxySecret",
        "outputs": [
            {
                "internalType": "string[]",
                "name": "",
                "type": "string[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "string",
                "name": "proxy",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "businessId",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "initiator",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "whiteList",
                "type": "string"
            }
        ],
        "name": "selectProxySecret",
        "outputs": [
            {
                "internalType": "string",
                "name": "",
                "type": "string"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "string",
                        "name": "WhiteList",
                        "type": "string"
                    },
                    {
                        "internalType": "string",
                        "name": "Sk",
                        "type": "string"
                    },
                    {
                        "internalType": "string",
                        "name": "Proxy",
                        "type": "string"
                    }
                ],
                "internalType": "struct TpreRelation[]",
                "name": "sec",
                "type": "tuple[]"
            },
            {
                "internalType": "string",
                "name": "businessId",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "initiator",
                "type": "string"
            }
        ],
        "name": "storeProxySecret",
        "outputs": [
            {
                "internalType": "string[]",
                "name": "",
                "type": "string[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "components": [
                    {
                        "internalType": "string",
                        "name": "WhiteList",
                        "type": "string"
                    },
                    {
                        "internalType": "string",
                        "name": "Sk",
                        "type": "string"
                    },
                    {
                        "internalType": "string",
                        "name": "Proxy",
                        "type": "string"
                    }
                ],
                "internalType": "struct TpreRelation[]",
                "name": "sec",
                "type": "tuple[]"
            },
            {
                "internalType": "string",
                "name": "businessId",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "initiator",
                "type": "string"
            },
            {
                "internalType": "string[]",
                "name": "whiteList",
                "type": "string[]"
            }
        ],
        "name": "updateProxySecret",
        "outputs": [
            {
                "internalType": "string[]",
                "name": "",
                "type": "string[]"
            }
        ],
        "stateMutability": "nonpayable",
        "type": "function"
    }
]`,
}
