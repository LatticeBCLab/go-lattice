---
title: Getting Started
---

# 快速开始

## 前置条件
你已经有一个运行中的 **节点** 。并且获取节点的 HTTP 连接信息：`http://127.0.0.1:8080`。

## 安装

```bash
go install github.com/LatticeBCLab/go-lattice
```
现在你可以初始化一个新的 ZLattice 客户端，如下所示：

```go
import (
    zlattice "github.com/LatticeBCLab/go-lattice/lattice"
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
```


## 转账

## 合约
### 部署合约

### 调用合约

### 预调用合约

### 内置合约

## Http

### 获取收据

### 获取交易块

### 获取守护块

