package block

import (
	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/crypto"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

// TransactionType 交易类型别名
//
//   - TransactionTypeGenesis			  创世交易
//   - TransactionTypeCreate			  构造交易
//   - TransactionTypeSend				  发送交易
//   - TransactionTypeReceive			  接收交易
//   - TransactionTypeDeployContract	  部署Solidity合约
//   - TransactionTypeCallContract	      调用Solidity合约
//   - TransactionTypeUpgradeContract	  升级Solidity合约
//   - TransactionTypeDeployJavaContract  部署Java和语
//   - TransactionTypeUpgradeJavaContract 升级Java合约
//   - TransactionTypeCallJavaContract	  调用Java合约
//   - TransactionTypeDeployGoContract	  部署Go合约
//   - TransactionTypeUpgradeGoContract	  升级Go合约
//   - TransactionTypeCallGoContract	  调用Go合约
type TransactionType string

const (
	TransactionTypeGenesis             TransactionType = "genesis"
	TransactionTypeCreate              TransactionType = "create"
	TransactionTypeSend                TransactionType = "send"
	TransactionTypeReceive             TransactionType = "receive"
	TransactionTypeDeployContract      TransactionType = "contract"
	TransactionTypeCallContract        TransactionType = "execute"
	TransactionTypeUpgradeContract     TransactionType = "update"
	TransactionTypeDeployJavaContract  TransactionType = "createJava"
	TransactionTypeUpgradeJavaContract TransactionType = "updateJava"
	TransactionTypeCallJavaContract    TransactionType = "executeJava"
	TransactionTypeDeployGoContract    TransactionType = "createGo"
	TransactionTypeUpgradeGoContract   TransactionType = "updateGo"
	TransactionTypeCallGoContract      TransactionType = "executeGo"
)

const (
	Genesis = iota
	Create
	Send
	Receive
	Contract
	Execute
	Upgrade
	RevokeContract
	FreezeContract
	UnfreezeContract
	DeployGoContract
	DeployJavaContract
	ExecuteGoContract
	ExecuteJavaContract
	UpgradeGoContract
	UpgradeJavaContract
)

var TransactionTypeCode = map[TransactionType]uint8{
	TransactionTypeGenesis:             Genesis,
	TransactionTypeCreate:              Create,
	TransactionTypeSend:                Send,
	TransactionTypeReceive:             Receive,
	TransactionTypeDeployContract:      Contract,
	TransactionTypeCallContract:        Execute,
	TransactionTypeUpgradeContract:     Upgrade,
	TransactionTypeDeployJavaContract:  DeployJavaContract,
	TransactionTypeUpgradeJavaContract: UpgradeJavaContract,
	TransactionTypeCallJavaContract:    ExecuteJavaContract,
	TransactionTypeDeployGoContract:    DeployGoContract,
	TransactionTypeUpgradeGoContract:   UpgradeGoContract,
	TransactionTypeCallGoContract:      ExecuteGoContract,
}

// Transaction 构造交易的结构体
type Transaction struct {
	Height      uint64          `json:"number"`
	Type        TransactionType `json:"type"`
	ParentHash  common.Hash     `json:"parentHash"`
	Hub         []common.Hash   `json:"hub"`
	DaemonHash  common.Hash     `json:"daemonHash"`
	CodeHash    common.Hash     `json:"codeHash"`
	Owner       string          `json:"owner"`
	Linker      string          `json:"linker"`
	Amount      *big.Int        `json:"amount"`
	Joule       *big.Int        `json:"joule"`
	Difficulty  uint64          `json:"difficulty"`
	Pow         *big.Int        `json:"pow"`
	ProofOfWork *big.Int        `json:"proofOfWork"`
	Payload     string          `json:"payload"`
	Timestamp   uint64          `json:"timestamp"`
	Code        string          `json:"code"`
	Sign        string          `json:"sign"`
	Hash        string          `json:"hash"`
	Hash2       common.Hash     `json:"hash2"`
	Key         string          `json:"key"`
	DataHash    string          `json:"dataHash"`
	ApplyHash   string          `json:"applyHash"`
}

func (tx *Transaction) GetTypeCode() uint8 {
	return TransactionTypeCode[tx.Type]
}

// GetOwnerAddress 获取Owner地址
//
// Parameters:
//
// Returns:
//   - common.Address
func (tx *Transaction) GetOwnerAddress() common.Address {
	addr, err := convert.ZltcToAddress(tx.Owner)
	if err != nil {
		return common.Address{}
	}
	return addr
}

// GetLinkerAddress 获取Linker地址
//
// Parameters:
//
// Returns:
//   - common.Address
func (tx *Transaction) GetLinkerAddress() common.Address {
	addr, err := convert.ZltcToAddress(tx.Linker)
	if err != nil {
		return common.Address{}
	}
	return addr
}

// DecodePayload decode 16进制的payload
//
// Parameters:
//
// Returns:
//   - []byte
func (tx *Transaction) DecodePayload() []byte {
	return hexutil.MustDecode(tx.Payload)
}

func (tx *Transaction) DecodeSign() []byte {
	return hexutil.MustDecode(tx.Sign)
}

// RlpEncodeHash 对交易进行rlp编码并计算哈希
//
// Parameters:
//   - chainId *big.Int: 区块链ID
//   - curve types.Curve: 椭圆曲线
//
// Returns:
//   - common.Hash: 哈希
func (tx *Transaction) RlpEncodeHash(chainId uint64, curve types.Curve) (common.Hash, error) {
	var err error
	hash := crypto.NewCrypto(curve).EncodeHash(func(writer io.Writer) {
		err = rlp.Encode(writer, []interface{}{
			tx.Height,
			tx.GetTypeCode(),
			tx.ParentHash,
			tx.Hub,
			tx.DaemonHash,
			tx.CodeHash,
			tx.GetOwnerAddress(),
			tx.GetLinkerAddress(),
			tx.Amount,
			tx.Joule,
			tx.Difficulty,
			tx.ProofOfWork,
			tx.DecodePayload(),
			tx.Timestamp,
			chainId,
			uint(0),
			uint(0),
		})
	})
	if err != nil {
		return common.Hash{}, err
	}
	return hash, nil
}

// SignTx 签名交易
//
// Parameters:
//   - curve types.Curve: 椭圆曲线类型
//   - hash []byte: 哈希
//   - skHex string: 私钥
//
// Returns:
//   - []byte: 签名
//   - error
func (tx *Transaction) sign(curve types.Curve, hash []byte, skHex string) ([]byte, error) {
	cryptoInstance := crypto.NewCrypto(curve)

	sk, err := cryptoInstance.HexToSK(skHex)
	if err != nil {
		return nil, err
	}

	sign, err := cryptoInstance.Sign(hash, sk)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func (tx *Transaction) CalculateTransactionHash(curve types.Curve) (common.Hash, error) {
	var err error
	hash := crypto.NewCrypto(curve).EncodeHash(func(writer io.Writer) {
		err = rlp.Encode(writer, []interface{}{
			tx.Height,
			tx.GetTypeCode(),
			tx.ParentHash,
			tx.DaemonHash,
			tx.CodeHash,
			tx.GetOwnerAddress(),
			tx.GetLinkerAddress(),
			tx.Hub,
			tx.Amount,
			uint(0), // income
			tx.Joule,
			tx.Difficulty,
			tx.ProofOfWork,
			tx.DecodePayload(),
			tx.Timestamp,
			tx.DecodeSign(),
			types.TXVersionLATEST,
		})
	})
	if err != nil {
		return common.Hash{}, err
	}
	return hash, nil
}

// SignTX 签名交易
//
// Parameters:
//   - chainId *big.Int
//   - curve types.Curve
//   - skHex string
//
// Returns:
//   - error
func (tx *Transaction) SignTX(chainId uint64, curve types.Curve, skHex string) error {
	hash, err := tx.RlpEncodeHash(chainId, curve)
	if err != nil {
		return err
	}

	return tx.SignHash(hash, curve, skHex)
}

// SignHash 签名哈希
func (tx *Transaction) SignHash(hash common.Hash, curve types.Curve, skHex string) error {
	signature, err := tx.sign(curve, hash[:], skHex)
	if err != nil {
		return err
	}
	tx.Sign = hexutil.Encode(signature)

	return nil
}
