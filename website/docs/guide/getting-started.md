---
title: Getting Started
---

# Getting Started

## Prerequisites
You already have a running **node** of `ZLattice`.

## Install

```bash
go install github.com/LatticeBCLab/go-lattice
```
Now you can initialize a new ZLattice client as follows:

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


## Transfer

## Contract
### Deploy Contract

### Call Contract

### Precall Contract

### Builtin Contract

## Http

### Get Receipt

### Get TBlock

### Get DBlokc

