package lattice

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/LatticeBCLab/go-lattice/abi"
	"github.com/LatticeBCLab/go-lattice/common/constant"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	"github.com/LatticeBCLab/go-lattice/common/types"
	"github.com/LatticeBCLab/go-lattice/lattice/builtin"
	"github.com/LatticeBCLab/go-lattice/lattice/protobuf"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	proto      = "syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}"
	counterAbi = `[{"inputs":[],"name":"decrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getCount","outputs":[{"internalType":"int256","name":"","type":"int256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"incrementCounter","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	chainId    = "1"
)

var latticeApi = NewLattice(
	&ChainConfig{Curve: types.Sm2p256v1},
	&ConnectingNodeConfig{Ip: "172.22.0.25", HttpPort: 31015, JwtSecret: "4EWRP6LZ", JwtTokenExpirationDuration: 10 * time.Hour},
	NewMemoryBlockCache(10*time.Second, time.Minute, time.Minute),
	NewAccountLock(),
	&Options{MaxIdleConnsPerHost: 200},
)

var credentials = &Credentials{
	AccountAddress: "zltc_j5yLhxm8fkwJkuhapqmqmJ1vYY2gLfPLy",
	PrivateKey:     "0x88d80c38a8a10e03b54c2c2234e90d9809030a78e4fd2f99a6a189629b530f90",
}

func TestAttestation(t *testing.T) {
	contract := builtin.NewCredibilityContract()
	t.Run("Create business", func(t *testing.T) {
		data, err := contract.CreateBusiness()
		assert.NoError(t, err)
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.GetCreateBusinessContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		t.Logf("业务合约地址：%s", receipt.ContractRet)
		t.Log(string(r))
	})

	t.Run("Create protocol", func(t *testing.T) {
		data, err := contract.CreateProtocol(2, []byte(proto))
		assert.NoError(t, err)
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		result, err := abi.DecodeReturn(contract.MyAbi(), "addProtocol", receipt.ContractRet)
		assert.NoError(t, err)
		t.Logf("创建协议返回：%s", result[0])
		t.Log(string(r))
	})

	t.Run("Batch create protocol", func(t *testing.T) {
		request := make([]builtin.CreateProtocolRequest, 2)
		request[0] = builtin.CreateProtocolRequest{
			ProtocolSuite: 10,
			Data:          convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
		}
		request[1] = builtin.CreateProtocolRequest{
			ProtocolSuite: 10,
			Data:          convert.BytesToBytes32Arr([]byte("syntax = \"proto3\";\n\nmessage Student {\n\tstring id = 1;\n\tstring name = 2;\n}")),
		}
		data, err := contract.BatchCreateProtocol(request)
		assert.NoError(t, err)
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		result, err := abi.DecodeReturn(contract.MyAbi(), "addProtocolBatch", receipt.ContractRet)
		assert.NoError(t, err)
		t.Logf("批量创建协议返回：%s", result)
		t.Log(string(r))
	})

	t.Run("Write", func(t *testing.T) {
		jsonData := `{"id":"1","name":"jack"}`
		fd := protobuf.MakeFileDescriptor(strings.NewReader(proto))
		dataBytes, err := protobuf.MarshallMessage(fd, jsonData)
		assert.NoError(t, err)
		data, err := contract.Write(&builtin.WriteLedgerRequest{
			ProtocolUri: 8589934593,
			Hash:        "2",
			Data:        convert.BytesToBytes32Arr(dataBytes),
			Address:     common.HexToAddress("0x561717f7922a233720ae38acaa4174cda0bf1766"),
		})
		assert.NoError(t, err)
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)
		t.Logf("结束数据存证，交易哈希为：%s", hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		t.Log(string(r))
	})

	t.Run("Read", func(t *testing.T) {
		data, _ := contract.Read("2", "0x561717f7922a233720ae38acaa4174cda0bf1766")
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, contract.ContractAddress(), data, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)
		t.Logf("结束数据存证，交易哈希为：%s", hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		t.Log(string(r))
	})
}

func TestLattice_Peekaboo(t *testing.T) {
	peekabooContract := builtin.NewPeekabooContract()
	t.Run("Hidden payload", func(t *testing.T) {
		code, err := peekabooContract.TogglePayload("0x6d5191f752ccc991ba033007e7afaeb624e4b726a7632ba26966e267ebac68e5", false)
		assert.NoError(t, err)

		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, peekabooContract.ContractAddress(), code, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)

		t.Log(hash)
		t.Log(receipt)
	})

	t.Run("Hidden code", func(t *testing.T) {
		code, err := peekabooContract.ToggleCode("0x6d5191f752ccc991ba033007e7afaeb624e4b726a7632ba26966e267ebac68e5", false)
		assert.NoError(t, err)

		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, peekabooContract.ContractAddress(), code, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)

		t.Log(hash)
		t.Log(receipt)
	})

	t.Run("Hidden hash", func(t *testing.T) {
		code, err := peekabooContract.ToggleHash("0x6d5191f752ccc991ba033007e7afaeb624e4b726a7632ba26966e267ebac68e5", true)
		assert.NoError(t, err)

		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, peekabooContract.ContractAddress(), code, constant.ZeroPayload, 0, 0, DefaultBackOffRetryStrategy())
		assert.NoError(t, err)

		t.Log(hash)
		t.Log(receipt)
	})
}

func TestTransfer(t *testing.T) {
	payload := "0x10"
	linker := "zltc_S5KXbs6gFkEpSnfNpBg3DvZHnB9aasa6Q"
	t.Run("transfer", func(t *testing.T) {
		hash, err := latticeApi.Transfer(context.Background(), credentials, chainId, linker, payload, 0, 0)
		assert.NoError(t, err)
		t.Log(hash.String())
	})

	t.Run("transfer wait for receipt", func(t *testing.T) {
		hash, receipt, err := latticeApi.TransferWaitReceipt(context.Background(), credentials, chainId, linker, payload, 0, 0, NewBackOffRetryStrategy(10, time.Second))
		assert.NoError(t, err)
		t.Log(hash.String())
		t.Log(receipt)
	})
}

func TestContract(t *testing.T) {
	t.Run("Deploy contract", func(t *testing.T) {
		data := "0x60806040526000805534801561001457600080fd5b50610278806100246000396000f3fe608060405234801561001057600080fd5b50600436106100415760003560e01c80635b34b96614610046578063a87d942c14610050578063f5c5ad831461006e575b600080fd5b61004e610078565b005b610058610093565b60405161006591906100d0565b60405180910390f35b61007661009c565b005b600160008082825461008a919061011a565b92505081905550565b60008054905090565b60016000808282546100ae91906101ae565b92505081905550565b6000819050919050565b6100ca816100b7565b82525050565b60006020820190506100e560008301846100c1565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610125826100b7565b9150610130836100b7565b9250817f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0383136000831215161561016b5761016a6100eb565b5b817f80000000000000000000000000000000000000000000000000000000000000000383126000831216156101a3576101a26100eb565b5b828201905092915050565b60006101b9826100b7565b91506101c4836100b7565b9250827f8000000000000000000000000000000000000000000000000000000000000000018212600084121516156101ff576101fe6100eb565b5b827f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff018213600084121615610237576102366100eb565b5b82820390509291505056fea2646970667358221220d841351625356129f6266ada896818d690dbc4b0d176774a97d745dfbe2fe50164736f6c634300080b0033"
		hash, receipt, err := latticeApi.DeployContractWaitReceipt(context.Background(), credentials, chainId, data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		t.Log(string(r))

		// 获取正在进行的提案
		proposal, err := latticeApi.HttpApi().GetContractLifecycleProposal(context.Background(), chainId, receipt.ContractAddress, types.ProposalStateINITIAL, "", "")
		assert.NoError(t, err)

		voteContract := builtin.NewProposalContract()
		approveData, err := voteContract.Approve(proposal[0].Content.Id)
		assert.Nil(t, err)

		hash, receipt, err = latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, builtin.ProposalBuiltinContract.Address, approveData, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		re, _ := json.Marshal(receipt)
		t.Log(string(re))
	})

	t.Run("Call contract", func(t *testing.T) {
		function, err := abi.NewAbi(counterAbi).GetLatticeFunction("incrementCounter")
		assert.NoError(t, err)
		data, err := function.Encode()
		assert.NoError(t, err)
		hash, receipt, err := latticeApi.CallContractWaitReceipt(context.Background(), credentials, chainId, "zltc_TpPbhQmKH8YoLF6aqLw77CiEoQxFL6SaM", data, constant.ZeroPayload, 0, 0, DefaultFixedRetryStrategy())
		assert.NoError(t, err)
		t.Log(hash.String())
		r, err := json.Marshal(receipt)
		assert.NoError(t, err)
		t.Log(string(r))
	})

	t.Run("Pre-contract approve", func(t *testing.T) {
		function, err := abi.NewAbi(counterAbi).GetLatticeFunction("getCount")
		assert.NoError(t, err)
		data, err := function.Encode()
		assert.NoError(t, err)
		receipt, err := latticeApi.PreCallContract(context.Background(), chainId, credentials.AccountAddress, "zltc_ebuyF9qei2hoESDzpFVM2cVFC9ViJXyNn", data, constant.ZeroPayload)
		assert.Nil(t, err)
		receiptBytes, _ := json.Marshal(receipt)
		t.Log(string(receiptBytes))
		assert.NotNil(t, receipt)
	})
}

func TestLattice_JsonRpc(t *testing.T) {
	t.Run("Get subchain brief info", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		subchainBriefInfos, err := latticeApi.HttpApi().GetSubchainBriefInfo(ctx, "-1")
		assert.NoError(t, err)
		t.Log(subchainBriefInfos[0])
	})

	t.Run("Get contract information", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		info, err := latticeApi.HttpApi().GetContractInformation(ctx, "17", "zltc_n7jX7sTsPFpqEjuVUWD8sBskEB384o5cp")
		assert.NoError(t, err)
		states := info.GetContractStates()
		fmt.Println(states)
		t.Log(info)
	})

	t.Run("Get contract lifecycle proposal", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		proposal, err := latticeApi.HttpApi().GetContractLifecycleProposal(ctx, "17", "zltc_n7jX7sTsPFpqEjuVUWD8sBskEB384o5cp", types.ProposalStateINITIAL, "20250703", "20250703")
		assert.NoError(t, err)
		t.Log(proposal)
	})

	t.Run("Get recent daemon blocks", func(t *testing.T) {
		var b byte = 2
		a := fmt.Sprintf("%08b", b)
		fmt.Println(a)
	})
}

func TestLattice_CanDial(t *testing.T) {
	err := latticeApi.HttpApi().CanDial(2 * time.Second)
	assert.NoError(t, err)
}
