package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CreateBusinessContractAddress 创建存证业务的业务合约地址
const CreateBusinessContractAddress = "zltc_QLbz7JHiBTspS9WTWJUrbNsB5wbENMweQ"

func NewCredibilityContract() CredibilityContract {
	return &credibilityContract{
		abi: abi.NewAbi(CredibilityBuiltinContract.AbiString),
	}
}

// WriteLedgerRequest 存证数据的请求结构体
type WriteLedgerRequest struct {
	ProtocolUri uint64         `json:"protocolUri"` // uri:协议号
	Hash        string         `json:"hash"`        // dataId:数据ID
	Data        [][32]byte     `json:"data"`        // data:存证的数据
	Address     common.Address `json:"address"`     // businessContractAddress:业务合约地址
}

// CreateProtocolRequest 创建协议的结构体
//
//   - ProtocolSuite 协议簇（行业号）
//   - Data 协议内容
type CreateProtocolRequest struct {
	ProtocolSuite uint64     `json:"protocolSuite"`
	Data          [][32]byte `json:"data"`
}

// UpdateProtocolRequest 更新协议的结构体
//
//   - ProtocolUri 协议号
//   - Data 更新的内容
type UpdateProtocolRequest struct {
	ProtocolUri uint64     `json:"protocolUri"`
	Data        [][32]byte `json:"data"`
}

type CredibilityContract interface {

	// MyAbi 返回存证合约的ABI对象
	//
	// Returns:
	//   - *myabi.ABI
	MyAbi() *myabi.ABI

	// ContractAddress 获取以链建链的合约地址
	//
	// Returns:
	//   - string: 合约地址，zltc_QLbz7JHiBTspUvTPzLHy5biDS9mu53mmv
	ContractAddress() string

	// GetCreateBusinessContractAddress 获取创建业务合约的合约地址
	//
	// Returns:
	//   - string: 创建业务合约的合约地址，zltc_QLbz7JHiBTspS9WTWJUrbNsB5wbENMweQ
	GetCreateBusinessContractAddress() string

	// CreateBusiness 创建业务合约地址
	//
	// Returns:
	//   - data string
	//   - err error
	CreateBusiness() (data string, err error)

	// CreateProtocol 创建协议
	//
	// Parameters:
	//   - tradeNumber uint64
	//   - message []byte
	//
	// Returns:
	//   - data string
	//   - err error
	CreateProtocol(tradeNumber uint64, message []byte) (data string, err error)

	// BatchCreateProtocol 批量创建协议
	//
	// Parameters:
	//   - request []CreateProtocolRequest
	//
	// Returns:
	//   - data string
	//   - err error
	BatchCreateProtocol(request []CreateProtocolRequest) (data string, err error)

	// ReadProtocol 读取协议
	//
	// Parameters:
	//   - uri uint64
	//
	// Returns:
	//   - data string
	//   - err error
	ReadProtocol(uri uint64) (data string, err error)

	// UpdateProtocol 更新协议
	//
	// Parameters:
	//   - uri uint64
	//   - message []byte
	//
	// Returns:
	//   - data string
	//   - err error
	UpdateProtocol(uri uint64, message []byte) (data string, err error)

	// BatchUpdateProtocol 批量更新协议
	//
	// Parameters:
	//   - request []UpdateProtocolRequest
	//
	// Returns:
	//   - data string
	//   - err error
	BatchUpdateProtocol(request []UpdateProtocolRequest) (data string, err error)

	// Write 写入存证数据
	//
	// Parameters:
	//   - request *WriteLedgerRequest
	//
	// Returns:
	//   - data string
	//   - err error
	Write(request *WriteLedgerRequest) (data string, err error)

	// BatchWrite 批量写入存证数据
	//
	// Parameters:
	//   - request []*WriteLedgerRequest
	//
	// Returns:
	//   - data string
	//   - err error
	BatchWrite(request []WriteLedgerRequest) (data string, err error)

	// Read 读取存证数据
	//
	// Parameters:
	//   - dataId string
	//   - businessContractAddress string
	//
	// Returns:
	//   - data string
	//   - err error
	Read(dataId, businessContractAddress string) (data string, err error)
}

type credibilityContract struct {
	abi abi.LatticeAbi
}

func (c *credibilityContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *credibilityContract) ContractAddress() string {
	return CredibilityBuiltinContract.Address
}

func (c *credibilityContract) GetCreateBusinessContractAddress() string {
	return CreateBusinessContractAddress
}

func (c *credibilityContract) CreateBusiness() (data string, err error) {
	return hexutil.Encode([]byte{49}), nil
}

func (c *credibilityContract) CreateProtocol(tradeNumber uint64, message []byte) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("addProtocol", tradeNumber, convert.BytesToBytes32Arr(message))
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) BatchCreateProtocol(request []CreateProtocolRequest) (data string, err error) {
	code, err := c.abi.RawAbi().Pack("addProtocolBatch", request)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *credibilityContract) ReadProtocol(uri uint64) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("getAddress", uri)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) UpdateProtocol(uri uint64, message []byte) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("updateProtocol", uri, convert.BytesToBytes32Arr(message))
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

func (c *credibilityContract) BatchUpdateProtocol(request []UpdateProtocolRequest) (data string, err error) {
	code, err := c.abi.RawAbi().Pack("updateProtocolBatch", request)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *credibilityContract) Write(request *WriteLedgerRequest) (data string, err error) {
	code, err := c.abi.RawAbi().Pack("writeTraceability", request.ProtocolUri, request.Hash, request.Data, request.Address)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *credibilityContract) BatchWrite(request []WriteLedgerRequest) (data string, err error) {
	code, err := c.abi.RawAbi().Pack("writeTraceabilityBatch", request)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *credibilityContract) Read(dataId, businessContractAddress string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("getTraceability", dataId, businessContractAddress)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var CredibilityBuiltinContract = Contract{
	Description: "存证溯源合约",
	Address:     "zltc_QLbz7JHiBTspUvTPzLHy5biDS9mu53mmv",
	AbiString: `[
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolSuite",
					"type": "uint64"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				}
			],
			"name": "addProtocol",
			"outputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
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
							"internalType": "uint64",
							"name": "ProtocolSuite",
							"type": "uint64"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct ProtocolSuiteParam[]",
					"name": "protocols",
					"type": "tuple[]"
				}
			],
			"name": "addProtocolBatch",
			"outputs": [
				{
					"internalType": "uint64[]",
					"name": "",
					"type": "uint64[]"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				}
			],
			"name": "getAddress",
			"outputs": [
				{
					"components": [
						{
							"internalType": "address",
							"name": "updater",
							"type": "address"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct credibilidity.Protocol[]",
					"name": "protocol",
					"type": "tuple[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				}
			],
			"name": "updateProtocol",
			"outputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
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
							"internalType": "uint64",
							"name": "ProtocolUri",
							"type": "uint64"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct ProtocolParam[]",
					"name": "protocols",
					"type": "tuple[]"
				}
			],
			"name": "updateProtocolBatch",
			"outputs": [
				{
					"internalType": "uint64[]",
					"name": "",
					"type": "uint64[]"
				}
			],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "getTraceability",
			"outputs": [
				{
					"components": [
						{
							"internalType": "uint64",
							"name": "number",
							"type": "uint64"
						},
						{
							"internalType": "uint64",
							"name": "protocol",
							"type": "uint64"
						},
						{
							"internalType": "address",
							"name": "updater",
							"type": "address"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						}
					],
					"internalType": "struct credibilidity.Evidence[]",
					"name": "evi",
					"type": "tuple[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "setDataSecret",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				},
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "writeTraceability",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "uint64",
					"name": "protocolUri",
					"type": "uint64"
				},
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				},
				{
					"internalType": "bytes32[]",
					"name": "data",
					"type": "bytes32[]"
				},
				{
					"internalType": "address",
					"name": "address",
					"type": "address"
				}
			],
			"name": "writeTraceabilityWithStatus",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"components": [
						{
							"internalType": "uint64",
							"name": "protocolUri",
							"type": "uint64"
						},
						{
							"internalType": "string",
							"name": "hash",
							"type": "string"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						},
						{
							"internalType": "address",
							"name": "address",
							"type": "address"
						}
					],
					"internalType": "struct Business.batch[]",
					"name": "bt",
					"type": "tuple[]"
				}
			],
			"name": "writeTraceabilityBatch",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"components": [
						{
							"internalType": "uint64",
							"name": "protocolUri",
							"type": "uint64"
						},
						{
							"internalType": "string",
							"name": "hash",
							"type": "string"
						},
						{
							"internalType": "bytes32[]",
							"name": "data",
							"type": "bytes32[]"
						},
						{
							"internalType": "address",
							"name": "address",
							"type": "address"
						}
					],
					"internalType": "struct Business.batch[]",
					"name": "bt",
					"type": "tuple[]"
				}
			],
			"name": "writeTraceabilityBatchWithStatus",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				}
			],
			"name": "quickWriteTraceability",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "hash",
					"type": "string"
				}
			],
			"name": "getQuickTraceability",
			"outputs": [],
			"stateMutability": "nonpayable",
			"type": "function"
		}
	]`,
}
