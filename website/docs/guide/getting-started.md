---
title: Getting Started
---

<h1 align="center" class="text-blue-700">快速开始</h1>

## 1.前置条件
你已经有一个运行中的 **节点** 。
+ 获取节点的 HTTP 连接信息：`http://127.0.0.1:8080`；
+ 获取一个账户的私钥，用于转账、合约部署和调用。

## 2.安装

```bash
go install github.com/LatticeBCLab/go-lattice
```
现在你可以初始化一个新的 ZLattice 客户端，如下所示：

```go
import (
    zlattice "github.com/LatticeBCLab/go-lattice/lattice"
    "github.com/LatticeBCLab/go-lattice/lattice/client"
)

type Curve string

const (
    Secp256k1 Curve = "secp256k1"
	Sm2p256v1 Curve = "sm2p256v1"
)

type Config struct {
    Node struct {
        IP string
        HttpPort int
        WebsocketPort int
        SecureConfiguration struct {
            JwtSecret string
            JwtExpirationDuration time.Duration
        }
        Curve Curve
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

func initZLatticeClient(config *Config) *ZLatticeClient {
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

// 晶格链节点的客户端
func (c *ZLatticeClient) Client() zlattice.Lattice {
    return c.client
}

// 晶格链节点的 HTTP API 客户端
func (c *ZLatticeClient) HttpApi() client.HttpApi {
    return c.client.HttpApi()
}
```

创建Credentials
```go
import (
    zlattice "github.com/LatticeBCLab/go-lattice/lattice"
)

const (
    passphrase = "9snjka823njk"
    fileKey = `{
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
    accountAddress = "0x0000000000000000000000000000000000000000"
    privateKey = "0x23d5b2a2eb0a9c8b86d62cbc3955cfd1fb26ec576ecc379f402d0f5d2b27a7bb"
)

// 推荐
func NewCredentials() *zlattice.Credentials {
    return &zlattice.Credentials{
        AccountAddress: accountAddress,
        PrivateKey: privateKey,
    }
}

func NewCredentialsWithPassphrase() *zlattice.Credentials {
    return &zlattice.Credentials{
        AccountAddress: accountAddress,
        FileKey: fileKey,
        Passphrase: passphrase,
    }
}

```

## 3.转账
::: tip
**什么是转账？**

转账是区块链中最基础的操作，指的是从一个账户向另一个账户发送代币或原生货币的过程。在晶格链中，每笔转账都会被记录在区块中，具有不可篡改性和可追溯性。
:::

- From: 转账发起方
- To: 转账接收方
- Amount: 转账金额
- GasLimit: 转账消耗的Gas
- GasPrice: 转账消耗的Gas价格

### 3.1 异步转账

```go
latticeClient := initZLatticeClient(cfg).Client()
credentials := NewCredentials()
chainId := "1"
to := "zltc_Z1pnS94bP4hQSYLs4aP4UwBP9pH8bEvhi"
payload := "0x" // hex string
amount := 100
joule := 1

cancelCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
hash, err := latticeClient.Transfer(cancelCtx, credentials, chainId, to, payload, amount, joule)
if err != nil {
    log.Fatal(err)
}
fmt.Println(hash)
```

### 3.2 同步转账
```go
cancelCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
hash, receipt, err := latticeClient.TransferWaitReceipt(cancelCtx, credentials, chainId, to, payload, amount, joule, zlattice.DefaultFixedRetryStrategy())
if err != nil {
    log.Fatal(err)
}
fmt.Println(hash)
```

## 4.合约
### 部署合约

### 调用合约

### 预调用合约

### 内置合约

## Http

### 获取收据

### 获取交易块

### 获取守护块

