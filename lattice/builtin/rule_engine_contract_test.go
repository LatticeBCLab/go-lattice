package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRuleEngineContract(t *testing.T) {
	contract := NewRuleEngineContract()
	assert.NotNil(t, contract)
	assert.NotNil(t, contract.MyAbi())
	assert.Equal(t, "zltc_aVwgkXjLrjc1TQaXAbD6WwgMFJ1mVuEGF", contract.ContractAddress())
	assert.Contains(t, contract.MyAbi().Methods, "accessContract")
	assert.Contains(t, contract.MyAbi().Methods, "createContract")
	assert.Contains(t, contract.MyAbi().Methods, "createResource")
	assert.Contains(t, contract.MyAbi().Methods, "upgradeContract")
	assert.Contains(t, contract.MyAbi().Methods, "signContract")
	assert.Contains(t, contract.MyAbi().Methods, "createSource")
	assert.Contains(t, contract.MyAbi().Methods, "createProduct")
}

func TestRuleEngineContractEncode(t *testing.T) {
	contract := NewRuleEngineContract()

	t.Run("AccessContract", func(t *testing.T) {
		actual, err := contract.AccessContract(AccessParams{
			ContractAddr: "zltc_gLpU8MFUgdECP5wJdhoZXVr4kmYH5xuir",
			ResourceID:   "resource-1",
			Operation:    "read",
			Rules: []Rule{
				{
					Name:           "rule-1",
					GRule:          "a > 1",
					Type:           1,
					FactJSONString: "{\"a\":2}",
				},
			},
		})
		assert.NoError(t, err)
		expected := "0xd9c478130000000000000000000000000000000000000000000000000000000000000020000000000000000000000000af8c6ed17254f26715099823c1e30293a717bcef000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000a7265736f757263652d31000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004726561640000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000672756c652d310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000561203e203100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000077b2261223a327d00000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})

	t.Run("CreateContract", func(t *testing.T) {
		actual, err := contract.CreateContract(ContractParams{
			ContractID:        "contract-1",
			ContractName:      "name-1",
			ContractAbstract:  "contract-abstract",
			SignMode:          "all",
			HasPrivacyCompute: false,
			ActivationTime:    1,
			EndTime:           2,
			Strategies: []Strategy{
				{
					ResourceID: "resource-1",
					Connects: []CreateConnectInfo{
						{
							ConnectID:    "connect-1",
							AccessType:   "api",
							AccessConfig: "{}",
							Entity:       "node-1",
						},
					},
					ResourceName: "resource-name",
					Abstract:     "resource-abstract",
					Operation:    "read",
					StrategyNodes: []StrategyNode{
						{
							NodeID:       "n1",
							NodeName:     "node1",
							NodePeerID:   "peer1",
							StrategyType: "or",
						},
					},
					Rules: []Rule{
						{
							Name:           "rule-1",
							GRule:          "a > 1",
							Type:           1,
							FactJSONString: "{\"a\":2}",
						},
					},
				},
			},
			Signatories: []Signatory{
				{
					ID:   "s1",
					Name: "signer-1",
					Sign: "0x0102",
				},
			},
			Code: "0x01",
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Contains(t, actual, "0x")
	})

	t.Run("CreateResourceAndProduct", func(t *testing.T) {
		args := ResourceInfoParams{
			ResourceID:        "res-8f9d2a1b",
			ResourceName:      "user_behavior_logs",
			ResourceStage:     "production",
			SafeLevel:         "L3",
			Industry:          "internet",
			Privacy:           true,
			Feature:           true,
			Owner:             "data-team-admin",
			Desc:              "Daily user click and browsing behavior logs",
			ValidityStartTime: "2024-01-01 00:00:00",
			ValidityEndTime:   "2099-12-31 23:59:59",
			Tables: []DataResourceTable{
				{
					TableID:      "tbl-log-01",
					DatasourceID: "ds-mysql-core",
					TableName:    "user_clicks",
					Columns: []DataResourceColumn{
						{
							ColName:    "user_id",
							ColType:    "varchar(64)",
							ColComment: "Unique identifier for the user",
						},
					},
				},
			},
			MetadataConfig: []string{"pii:user_id", "retention:365d"},
		}

		actual1, err := contract.CreateResource(args)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual1)
		assert.Contains(t, actual1, "0x")

		actual2, err := contract.CreateProduct(args)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual2)
		assert.Contains(t, actual2, "0x")
	})

	t.Run("UpgradeContract", func(t *testing.T) {
		actual, err := contract.UpgradeContract(UpgradeContractParams{
			ContractID:        "contract-1",
			ContractName:      "name-1",
			ContractAbstract:  "contract-abstract",
			SignMode:          "all",
			HasPrivacyCompute: false,
			ActivationTime:    1,
			EndTime:           2,
			Strategies: []UpgradeStrategy{
				{
					ResourceID: "resource-1",
					Connects: []UpgradeConnectInfo{
						{
							ConnectID:    "connect-1",
							AccessType:   "api",
							AccessConfig: "{}",
						},
					},
					ResourceName: "resource-name",
					Abstract:     "resource-abstract",
					Operation:    "read",
					StrategyNodes: []StrategyNode{
						{
							NodeID:       "n1",
							NodeName:     "node1",
							NodePeerID:   "peer1",
							StrategyType: "or",
						},
					},
					Rules: []Rule{
						{
							Name:           "rule-1",
							GRule:          "a > 1",
							Type:           1,
							FactJSONString: "{\"a\":2}",
						},
					},
				},
			},
			Signatories: []Signatory{
				{
					ID:   "s1",
					Name: "signer-1",
					Sign: "0x0102",
				},
			},
			Code: "0x01",
		}, "0x9293c604c644bfac34f498998cc3402f203d4d6b")
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Contains(t, actual, "0x")
	})

	t.Run("SignContract", func(t *testing.T) {
		actual, err := contract.SignContract(
			"0x9293c604c644bfac34f498998cc3402f203d4d6b",
			[]Signatory{
				{
					ID:   "s1",
					Name: "signer-1",
					Sign: "0x0102",
				},
			},
		)
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)
		assert.Contains(t, actual, "0x")
	})

	t.Run("CreateSource", func(t *testing.T) {
		actual, err := contract.CreateSource(SourceInfoParams{
			ID:         "source-1",
			Name:       "name-1",
			Owner:      "owner-1",
			Abstract:   "source-abstract",
			SourceType: "mysql",
			Config:     "{}",
		})
		assert.NoError(t, err)
		expected := "0x1b321beb000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001c000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000008736f757263652d3100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000066e616d652d31000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000076f776e65722d3100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f736f757263652d6162737472616374000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000056d7973716c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000027b7d000000000000000000000000000000000000000000000000000000000000"
		assert.Equal(t, expected, actual)
	})
}
