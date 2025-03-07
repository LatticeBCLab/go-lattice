---
title: Getting Started
---

<h1 align="center" class="text-blue-700">快速开始</h1>

## 1.前置条件
你已经有一个运行中的 **节点** 。
+ 获取节点的 HTTP 连接信息：`http://127.0.0.1:8080`；
+ 获取一个账户的私钥，用于转账、合约部署和调用。

## 2.安装

### 2.1 安装 go-lattice

```bash
go install github.com/LatticeBCLab/go-lattice
```

### 2.2 初始化 ZLattice 客户端

现在你可以初始化一个新的 ZLattice 客户端，如下所示：

```go
package main
import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/LatticeBCLab/go-lattice/abi"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/crypto"
	zlattice "github.com/LatticeBCLab/go-lattice/lattice"
	"github.com/LatticeBCLab/go-lattice/lattice/client"
	"github.com/LatticeBCLab/go-lattice/wallet"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"time"
)
type Config struct {
	Node struct {
		IP                  string
		HttpPort            int
		WebsocketPort       int
		SecureConfiguration struct {
			JwtSecret             string
			JwtExpirationDuration time.Duration
		}
		Curve types.Curve
	}
}
func New(config *Config) *ZLatticeClient {
	return &ZLatticeClient{
		config: config,
		client: initZLatticeClient(config),
	}
}
type ZLatticeClient struct {
	config *Config
	client zlattice.Lattice
}
func initZLatticeClient(config *Config) zlattice.Lattice {
	return zlattice.NewLattice(
		&zlattice.ChainConfig{
			Curve: config.Node.Curve,
		},
		&zlattice.ConnectingNodeConfig{
			Ip:                         config.Node.IP,
			HttpPort:                   uint16(config.Node.HttpPort),
			JwtSecret:                  config.Node.SecureConfiguration.JwtSecret,
			JwtTokenExpirationDuration: config.Node.SecureConfiguration.JwtExpirationDuration,
		},
		zlattice.NewMemoryBlockCache(10*time.Second, time.Minute, time.Minute),
		zlattice.NewAccountLock(),
		&zlattice.Options{MaxIdleConnsPerHost: 200},
	)
}
// Client 晶格链节点的客户端
func (c *ZLatticeClient) Client() zlattice.Lattice {
	return c.client
}
// HttpApi 晶格链节点的 HTTP API 客户端
func (c *ZLatticeClient) HttpApi() client.HttpApi {
	return c.client.HttpApi()
}
```

### 2.3 创建Credentials
创建凭证
```go
import (
    zlattice "github.com/LatticeBCLab/go-lattice/lattice"
)

const (
    chainId    = "1"
    passphrase = "9snjka823njk"
    fileKey    = `{
        "uuid": "bb889ee6-5d1d-474e-9514-5bbf412a42ec",
        "address": "zltc_iEUCcfMhVYy3zcpp8zLjoaTAeN6PZfMBL",
        "cipher": {
            "aes": {
                "cipher": "aes-128-ctr",
                "iv": "23b4ddcd8cfea7e37b3c69bbb600934f"
            },
            "kdf": {
                "kdf": "scrypt",
                "kdfParams": {
                    "DKLen": 32,
                    "n": 262144,
                    "p": 1,
                    "r": 8,
                    "salt": "87cf307be225ce2eaf255d602233852200195d838b5d98c4078ceb6235ec46e4"
                }
            },
            "cipherText": "672b3de4784fc0d17941ae257908672dd4984a43c616147366a42bc2e9ef2d8a",
            "mac": "bd6ac051c41f4d0238464a66df004de357baf2f3f03ced8ccba0a497e14044bd"
        },
        "isGM": true
    }`
    accountAddress = "zltc_iEUCcfMhVYy3zcpp8zLjoaTAeN6PZfMBL"
    linkerAddress  = "zltc_jF4U7umzNpiE8uU35RCBp9f2qf53H5CZZ"
    privateKey     = "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"
)

// NewCredentials 推荐
func NewCredentials() *zlattice.Credentials {
    return &zlattice.Credentials{
        AccountAddress: accountAddress,
        PrivateKey:     privateKey,
    }
}

func NewCredentialsWithPassphrase() *zlattice.Credentials {
    return &zlattice.Credentials{
        AccountAddress: accountAddress,
        FileKey:        fileKey,
        Passphrase:     passphrase,
    }
}
```

## 3.转账
::: tip
**什么是转账？**

转账是区块链中最基础的操作，指的是从一个账户向另一个账户发送代币或原生货币的过程。在晶格链中，每笔转账都会被记录在区块中，具有不可篡改性和可追溯性。
:::


### 3.1 异步转账
```go
func (c *ZLatticeClient) Transfer() {
    payload := "0x01"
    hash, err := c.Client().Transfer(context.Background(), NewCredentials(), chainId, linkerAddress, payload, 0, 0)
    if err != nil {
        log.Fatal().Err(err)
    }
    log.Info().Msgf("转账交易哈希: %s", hash.String())
}
```

### 3.2 同步转账

```go
func (c *ZLatticeClient) TransferWaitReceipt() {
    payload := "0x01"
    hash, receipt, err := c.Client().TransferWaitReceipt(context.Background(), NewCredentials(), chainId, linkerAddress, payload, 0, 0, zlattice.DefaultBackOffRetryStrategy())
    if err != nil {
        log.Fatal().Err(err)
    }
    log.Info().Msgf("转账交易哈希: %s", hash.String())
    log.Error().Msgf("转账交易回执：%v", receipt)
}
```

## 4.智能合约
::: tip
智能合约是一段存储在区块链上的代码，它会在满足触发条件时自动执行，就像一台“不可篡改的自动售货机”。是一个去中心化、不可篡改的程序，它存储在区块链中，并在交易或事件触发时执行特定的逻辑。
:::

### 4.1 部署合约
```go
func (c *ZLatticeClient) Deploy() {
    abiString := `[
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
            }
        ],
        "name": "changeAddress",
        "outputs": [],
        "stateMutability": "pure",
        "type": "function"
    }
]`
    bytecode := "0x01091821020128128108201281"
    latticeAbi := abi.NewAbi(abiString)
    method := latticeAbi.GetConstructor()
    if err != nil {
        log.Fatal().Err(err)
    }
    data, err := method.Encode()
    if err != nil {
        log.Fatal().Err(err).Msgf("编码合约参数错误")
    }
    hash, receipt, err := c.Client().DeployContractWaitReceipt(context.Background(), NewCredentials(), chainId, bytecode+data, "0x", 0, 0, zlattice.DefaultBackOffRetryStrategy())
    if err != nil {
        log.Fatal().Err(err).Msg("部署合约失败")
    }
    if !receipt.Success {
        log.Fatal().Str("错误", receipt.ContractRet).Msgf("部署合约失败")
    }
    log.Info().Msgf("转账交易哈希: %s", hash.String())
    log.Error().Msgf("转账交易回执：%v", receipt)
}
```

### 4.2 调用合约
#### 4.2.1 异步调用
```go
func (c *ZLatticeClient) CallContract() {
    contractAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    abiString := `[
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
            }
        ],
        "name": "changeAddress",
        "outputs": [],
        "stateMutability": "pure",
        "type": "function"
    }
]`
    oldAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    newAddress := "zltc_VmGey8EMn7MdyXMoPx9a3sJhhkaJykv3Z"
    latticeAbi := abi.NewAbi(abiString)
    method, err := latticeAbi.GetLatticeFunction("changeAddress", oldAddress, newAddress)
    if err != nil {
        log.Fatal().Err(err)
    }
    data, err := method.Encode()
    if err != nil {
        log.Fatal().Err(err).Msgf("编码合约参数错误")
    }
    hash, err := c.Client().CallContract(context.Background(), NewCredentials(), chainId, contractAddress, data, "0x", 0, 0)
    if err != nil {
        log.Fatal().Err(err)
    }
    log.Info().Msgf("调用合约交易哈希: %s", hash.String())
}
```
### 4.2.2 同步调用
```go
func (c *ZLatticeClient) CallContractWaitReceipt() {
    contractAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    abiString := `[
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
            }
        ],
        "name": "changeAddress",
        "outputs": [],
        "stateMutability": "pure",
        "type": "function"
    }
]`
    oldAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    newAddress := "zltc_VmGey8EMn7MdyXMoPx9a3sJhhkaJykv3Z"
    latticeAbi := abi.NewAbi(abiString)
    method, err := latticeAbi.GetLatticeFunction("changeAddress", oldAddress, newAddress)
    if err != nil {
        log.Fatal().Err(err)
    }
    data, err := method.Encode()
    if err != nil {
        log.Fatal().Err(err).Msgf("编码合约参数错误")
    }
    hash, receipt, err := c.Client().CallContractWaitReceipt(context.Background(), NewCredentials(), chainId, contractAddress, data, "0x", 0, 0, zlattice.DefaultBackOffRetryStrategy())
    if err != nil {
        log.Fatal().Err(err)
    }
    if !receipt.Success {
        log.Fatal().Str("错误", receipt.ContractRet).Msgf("执行合约失败")
    }
    log.Info().Msgf("调用合约交易哈希: %s", hash.String())
    log.Error().Msgf("调用交易回执：%v", receipt)
}
```

### 预调用合约
```go
func (c *ZLatticeClient) PreCallContract() {
    contractAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    abiString := `[
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
            }
        ],
        "name": "changeAddress",
        "outputs": [],
        "stateMutability": "pure",
        "type": "function"
    }
]`
    oldAddress := "zltc_YBomBNykwMqxm719giBL3VtYV4ABT9a8D"
    newAddress := "zltc_VmGey8EMn7MdyXMoPx9a3sJhhkaJykv3Z"
    latticeAbi := abi.NewAbi(abiString)
    method, err := latticeAbi.GetLatticeFunction("changeAddress", oldAddress, newAddress)
    if err != nil {
        log.Fatal().Err(err)
    }
    data, err := method.Encode()
    if err != nil {
        log.Fatal().Err(err).Msgf("编码合约参数错误")
    }
    receipt, err := c.Client().PreCallContract(context.Background(), chainId, accountAddress, contractAddress, data, "0x")
    if err != nil {
        log.Fatal().Err(err)
    }
    if !receipt.Success {
        log.Fatal().Str("错误", receipt.ContractRet).Msgf("预执行合约失败")
    }
    log.Error().Msgf("预调用交易回执：%v", receipt)
}
```

## 5.账户
### 5.1 生成私钥、账户
```go
func GenerateAccount() {
    c := crypto.NewCrypto(types.Sm2p256v1)
    priKey, err := c.GenerateKeyPair()
    if err != nil {
        log.Fatal().Msgf("生成密钥对错误")
    }
    privateKeyHexString, _ := c.SKToHexString(priKey)
    publicKeyHexString, _ := c.PKToHexString(&priKey.PublicKey)
    account, _ := c.PKToAddress(&priKey.PublicKey)
    log.Info().Str("私钥", privateKeyHexString)
    log.Info().Str("公钥", publicKeyHexString)
    log.Info().Str("账户地址", account.String())
}
```

### 5.2 生成FileKey
```go
func GenerateFileKey() {
    fk, err := wallet.GenerateFileKey(privateKey, passphrase, types.Sm2p256v1)
    log.Fatal().Err(err)
    bytes, err := json.Marshal(fk)
    log.Fatal().Err(err)
    fileKeyString := string(bytes)
    log.Info().Str("生成FileKey", fileKeyString)
}
```

### 5.3 解密FileKey
```go
func DecryptFileKey() {
    fileKeyString := `{"uuid":"bb889ee6-5d1d-474e-9514-5bbf412a42ec","address":"zltc_iEUCcfMhVYy3zcpp8zLjoaTAeN6PZfMBL","cipher":{"aes":{"cipher":"aes-128-ctr","iv":"23b4ddcd8cfea7e37b3c69bbb600934f"},"kdf":{"kdf":"scrypt","kdfParams":{"DKLen":32,"n":262144,"p":1,"r":8,"salt":"87cf307be225ce2eaf255d602233852200195d838b5d98c4078ceb6235ec46e4"}},"cipherText":"672b3de4784fc0d17941ae257908672dd4984a43c616147366a42bc2e9ef2d8a","mac":"bd6ac051c41f4d0238464a66df004de357baf2f3f03ced8ccba0a497e14044bd"},"isGM":true}`
    password := "Aa123456"
    sk, err := wallet.NewFileKey(fileKeyString).Decrypt(password)
    log.Fatal().Err(err)
    skString, err := crypto.NewCrypto(types.Sm2p256v1).SKToHexString(sk)
    log.Fatal().Err(err)
    log.Info().Str("解密获取私钥", skString)
}
```

## 6.HTTP Client

### 6.1 获取收据
```go
chainId := "1"
hash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
zlatticeClient := initZLatticeClient(cfg).Client()
receipt, err := zlatticeClient.HttpApi().GetReceipt(context.Background(), chainId, hash)
if err != nil {
    log.Fatal(err)
}
fmt.Println(receipt)
```

### 6.2 获取交易块
```go
chainId := "1"
hash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
zlatticeClient := initZLatticeClient(cfg).Client()
block, err := zlatticeClient.HttpApi().GetTransactionBlockByHash(context.Background(), chainId, hash)
if err != nil {
    log.Fatal(err)
}
fmt.Println(block)
```

### 6.3 获取守护块
```go
chainId := "1"
daemonBlockHash := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
zlatticeClient := initZLatticeClient(cfg).Client()
block, err := zlatticeClient.HttpApi().GetDaemonBlockByHash(context.Background(), chainId, daemonBlockHash)
if err != nil {
    log.Fatal(err)
}
fmt.Println(block)
```