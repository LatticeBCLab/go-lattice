package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func NewTraceabilityContract() TraceabilityContract {
	return &traceabilityContract{
		abi: abi.NewAbi(TraceabilityBuiltinContract.AbiString),
	}
}

type TraceabilityContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	// Write traceability data
	Write(request *WriteTraceabilityRequest) (data string, err error)
	// Read traceability data
	Read(id string) (data string, err error)
}

type WriteTraceabilityRequest struct {
	TraceabilityId string         `json:"traceabilityId"`
	ProtocolUri    uint64         `json:"protocolUri"`
	Hash           string         `json:"hash"`
	Data           [][32]byte     `json:"data"`
	Address        common.Address `json:"address"`
}

type traceabilityContract struct {
	abi abi.LatticeAbi
}

func (c *traceabilityContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *traceabilityContract) ContractAddress() string {
	return TraceabilityBuiltinContract.Address
}

func (c *traceabilityContract) Write(request *WriteTraceabilityRequest) (data string, err error) {
	code, err := c.abi.RawAbi().Pack("writeTraceability", request.TraceabilityId, request.ProtocolUri, request.Hash, request.Data, request.Address)
	if err != nil {
		return "", err
	}

	return hexutil.Encode(code), nil
}

func (c *traceabilityContract) Read(id string) (data string, err error) {
	fn, err := c.abi.GetLatticeFunction("getTraceability", id)
	if err != nil {
		return "", err
	}

	return fn.Encode()
}

var TraceabilityBuiltinContract = Contract{
	Description: "溯源合约",
	Address:     "zltc_QLbz7JHiBTszrxNq8kcAzz8TYKsgttvgf",
	AbiString: `[
	{
			"inputs": [
				{
					"internalType": "string",
					"name": "traceabilityId",
					"type": "string"
				}
			],
			"name": "getTraceability",
			"outputs": [
				{
					"internalType": "string[]",
					"name": "",
					"type": "string[]"
				}
			],
			"stateMutability": "view",
			"type": "function"
		},
		{
			"inputs": [
				{
					"internalType": "string",
					"name": "traceabilityId",
					"type": "string"
				},
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
		}
	]`,
}
