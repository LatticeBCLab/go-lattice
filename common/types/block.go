package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type LatestBlock struct {
	Height          uint64      `json:"currentTBlockNumber"` // 账户最新的高度
	Hash            common.Hash `json:"currentTBlockHash"`   // 账户最新的一笔交易哈希
	DaemonBlockHash common.Hash `json:"currentDBlockHash"`   // 守护区块的哈希
}

// IncrHeight 增长高度
func (b *LatestBlock) IncrHeight() {
	b.Height++
}

// DaemonBlock 守护区块
//   - Hash			守护区块哈希
//   - ParentHash
//   - LedgerHash
//   - ReceiptsHash 回执哈希
//   - Coinbase
//   - Signer
//   - Contracts
//   - Difficulty
//   - Height	   守护区块高度
//   - Extra
//   - Reward
//   - Pow
//   - Timestamp
//   - Size
//   - TD
//   - TTD
//   - Version
//   - TxHashes
//   - Receipts
//   - Anchors
type DaemonBlock struct {
	Hash         common.Hash   `json:"hash"`
	ParentHash   common.Hash   `json:"parentHash"`
	LedgerHash   common.Hash   `json:"ledgerHash"`
	ReceiptsHash common.Hash   `json:"receiptsHash"`
	Coinbase     string        `json:"coinbase"`
	Signer       string        `json:"signer"`
	Contracts    []string      `json:"contracts"`
	Difficulty   *big.Int      `json:"difficulty"`
	Height       *big.Int      `json:"number"`
	LatestHeight *big.Int      `json:"lastedDBNumber"`
	Extra        string        `json:"extra"`
	Reward       *big.Int      `json:"reward"`
	Pow          *big.Int      `json:"pow"`
	Timestamp    uint64        `json:"timestamp"`
	Size         uint32        `json:"size"`
	TD           uint64        `json:"td"`
	TTD          uint64        `json:"ttd"`
	Version      uint32        `json:"version"`
	TxHashes     []common.Hash `json:"txHashList"`
	Receipts     []*Receipt    `json:"receipts"`
	Anchors      []*Anchor     `json:"anchors"`
}

type Anchor struct {
	Height *big.Int `json:"number"`
	Hash   string   `json:"hash"`
	Owner  string   `json:"owner"`
}

// TransactionBlock 交易区块
type TransactionBlock struct {
	Height     *big.Int    `json:"number"`
	Hash       common.Hash `json:"hash,omitempty"`
	ParentHash common.Hash `json:"parentHash"`
	DaemonHash common.Hash `json:"daemonHash"`
	Payload    string      `json:"payload"`
	Hub        []string    `json:"hub"`
	Timestamp  uint64      `json:"timestamp"`
	Type       string      `json:"type"`
	Owner      string      `json:"owner"`
	Linker     string      `json:"linker"`
	Code       string      `json:"code"`
	CodeHash   string      `json:"codeHash"`
	Amount     string      `json:"amount"`
	Joule      uint64      `json:"joule"`
	Sign       string      `json:"sign"`
	Pow        string      `json:"proofOfWork"`
	Size       uint64      `json:"size"`
}

// TransactionsPagination 查询交易区块的分页结果
//   - NextStartingHeight 下一次查询的起始高度
//   - Count			  查询到的交易数量
//   - Transactions		  查询到的交易集合
type TransactionsPagination struct {
	NextStartingHeight *big.Int            `json:"height"`
	Count              uint32              `json:"count"`
	Transactions       []*TransactionBlock `json:"tBlocks"`
}

type Signature struct {
	From   string      `json:"from"`
	Owner  string      `json:"owner"`
	Hash   common.Hash `json:"hash"`
	Number *big.Int    `json:"number"`
	Status uint        `json:"status"`
	Sign   []byte      `json:"sign"`
}

type WitnessProof struct {
	Hash         common.Hash  `json:"Hash"`
	Owner        string       `json:"Owner"`
	DaemonHash   common.Hash  `json:"DaemonHash"`
	Number       *big.Int     `json:"Number"`
	EndNumber    *big.Int     `json:"EndNumber"`
	ExpectNumber *big.Int     `json:"Expect"` //预料中有其证明的DBlock高度
	Signers      []*Signature `json:"Signers"`
}
