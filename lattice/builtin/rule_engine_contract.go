package builtin

import (
	"github.com/LatticeBCLab/go-lattice/abi"
	"github.com/LatticeBCLab/go-lattice/common/convert"
	myabi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/samber/lo"
)

func NewRuleEngineContract() RuleEngineContract {
	return &ruleEngineContract{
		abi: abi.NewAbi(RuleEngineBuiltinContract.AbiString),
	}
}

type Rule struct {
	Name           string `json:"name"`
	GRule          string `json:"grule"`
	Type           uint8  `json:"type"`
	FactJSONString string `json:"factJsonString"`
}

type AccessParams struct {
	ContractAddr string `json:"contractAddr"`
	ResourceID   string `json:"resourceId"`
	Operation    string `json:"operation"`
	Rules        []Rule `json:"rules"`
}

type CreateConnectInfo struct {
	ConnectID    string `json:"connectId"`
	AccessType   string `json:"accessType"`
	AccessConfig string `json:"accessConfig"`
	Entity       string `json:"entity"`
}

type UpgradeConnectInfo struct {
	ConnectID    string `json:"connectId"`
	AccessType   string `json:"accessType"`
	AccessConfig string `json:"accessConfig"`
}

type StrategyNode struct {
	NodeID       string `json:"nodeId"`
	NodeName     string `json:"nodeName"`
	NodePeerID   string `json:"nodePeerId"`
	StrategyType string `json:"strategyType"`
}

type Strategy struct {
	ResourceID    string              `json:"resourceId"`
	Connects      []CreateConnectInfo `json:"connects"`
	ResourceName  string              `json:"resourceName"`
	Abstract      string              `json:"abstract"`
	Operation     string              `json:"operation"`
	StrategyNodes []StrategyNode      `json:"strategyNodes"`
	Rules         []Rule              `json:"rules"`
}

type UpgradeStrategy struct {
	ResourceID    string               `json:"resourceId"`
	Connects      []UpgradeConnectInfo `json:"connects"`
	ResourceName  string               `json:"resourceName"`
	Abstract      string               `json:"abstract"`
	Operation     string               `json:"operation"`
	StrategyNodes []StrategyNode       `json:"strategyNodes"`
	Rules         []Rule               `json:"rules"`
}

type Signatory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Sign string `json:"sign"`
}

type ContractParams struct {
	ContractID        string      `json:"contractId"`
	ContractName      string      `json:"contractName"`
	ContractAbstract  string      `json:"contractAbstract"`
	SignMode          string      `json:"signMode"`
	HasPrivacyCompute bool        `json:"hasPrivacyCompute"`
	ActivationTime    uint64      `json:"activationTime"`
	EndTime           uint64      `json:"endTime"`
	Strategies        []Strategy  `json:"strategies"`
	Signatories       []Signatory `json:"signatories"`
	Code              string      `json:"code"`
}

type UpgradeContractParams struct {
	ContractID        string            `json:"contractId"`
	ContractName      string            `json:"contractName"`
	ContractAbstract  string            `json:"contractAbstract"`
	SignMode          string            `json:"signMode"`
	HasPrivacyCompute bool              `json:"hasPrivacyCompute"`
	ActivationTime    uint64            `json:"activationTime"`
	EndTime           uint64            `json:"endTime"`
	Strategies        []UpgradeStrategy `json:"strategies"`
	Signatories       []Signatory       `json:"signatories"`
	Code              string            `json:"code"`
}

type DataResourceColumn struct {
	ColName    string `json:"colName"`
	ColType    string `json:"colType"`
	ColComment string `json:"colComment"`
}

type DataResourceTable struct {
	TableID      string               `json:"tableId"`
	DatasourceID string               `json:"datasourceId"`
	TableName    string               `json:"tableName"`
	Columns      []DataResourceColumn `json:"columns"`
}

type ResourceInfoParams struct {
	ResourceID        string              `json:"resourceId"`
	ResourceName      string              `json:"resourceName"`
	ResourceStage     string              `json:"resourceStage"`
	Owner             string              `json:"owner"`
	Desc              string              `json:"desc"`
	ValidityStartTime string              `json:"validityStartTime"`
	ValidityEndTime   string              `json:"validityEndTime"`
	Tables            []DataResourceTable `json:"tables"`
}

type SourceInfoParams struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Abstract   string `json:"abstract"`
	SourceType string `json:"sourceType"`
	Config     string `json:"config"`
}

type RuleEngineContract interface {
	MyAbi() *myabi.ABI
	ContractAddress() string
	AccessContract(args AccessParams) (string, error)
	CreateContract(args ContractParams) (string, error)
	CreateResource(args ResourceInfoParams) (string, error)
	UpgradeContract(args UpgradeContractParams, contractAddr string) (string, error)
	SignContract(contractAddr string, signatories []Signatory) (string, error)
	CreateSource(args SourceInfoParams) (string, error)
	CreateProduct(args ResourceInfoParams) (string, error)
}

type ruleEngineContract struct {
	abi abi.LatticeAbi
}

func (c *ruleEngineContract) MyAbi() *myabi.ABI {
	return c.abi.RawAbi()
}

func (c *ruleEngineContract) ContractAddress() string {
	return RuleEngineBuiltinContract.Address
}

func (c *ruleEngineContract) AccessContract(args AccessParams) (string, error) {
	param, err := toAccessContractParam(args)
	if err != nil {
		return "", err
	}

	code, err := c.abi.RawAbi().Pack("accessContract", param)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) CreateContract(args ContractParams) (string, error) {
	param, err := toCreateContractParam(args)
	if err != nil {
		return "", err
	}

	code, err := c.abi.RawAbi().Pack("createContract", param)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) CreateResource(args ResourceInfoParams) (string, error) {
	param := toResourceInfoParam(args)
	code, err := c.abi.RawAbi().Pack("createResource", param)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) UpgradeContract(args UpgradeContractParams, contractAddr string) (string, error) {
	param, err := toUpgradeContractParam(args)
	if err != nil {
		return "", err
	}

	parsedAddr, err := convert.ZltcToAddress(contractAddr)
	if err != nil {
		return "", err
	}

	code, err := c.abi.RawAbi().Pack("upgradeContract", param, parsedAddr)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) SignContract(contractAddr string, signatories []Signatory) (string, error) {
	parsedAddr, err := convert.ZltcToAddress(contractAddr)
	if err != nil {
		return "", err
	}

	params, err := toSignatoryParams(signatories)
	if err != nil {
		return "", err
	}

	code, err := c.abi.RawAbi().Pack("signContract", parsedAddr, params)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) CreateSource(args SourceInfoParams) (string, error) {
	param := sourceInfoParam{
		ID:         args.ID,
		Name:       args.Name,
		Owner:      args.Owner,
		Abstract:   args.Abstract,
		SourceType: args.SourceType,
		Config:     args.Config,
	}

	code, err := c.abi.RawAbi().Pack("createSource", param)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

func (c *ruleEngineContract) CreateProduct(args ResourceInfoParams) (string, error) {
	param := toResourceInfoParam(args)
	code, err := c.abi.RawAbi().Pack("createProduct", param)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(code), nil
}

type ruleParam struct {
	Name           string `abi:"name"`
	GRule          string `abi:"grule"`
	Type           uint8  `abi:"Type"`
	FactJSONString string `abi:"factJsonString"`
}

type accessContractParam struct {
	ContractAddr common.Address `abi:"contractAddr"`
	ResourceID   string         `abi:"resourceId"`
	Operation    string         `abi:"operation"`
	Rules        []ruleParam    `abi:"rules"`
}

type createConnectInfoParam struct {
	ConnectID    string `abi:"connectId"`
	AccessType   string `abi:"accessType"`
	AccessConfig string `abi:"accessConfig"`
	Entity       string `abi:"entity"`
}

type upgradeConnectInfoParam struct {
	ConnectID    string `abi:"connectId"`
	AccessType   string `abi:"accessType"`
	AccessConfig string `abi:"accessConfig"`
}

type strategyNodeParam struct {
	NodeID       string `abi:"nodeId"`
	NodeName     string `abi:"nodeName"`
	NodePeerID   string `abi:"nodePeerId"`
	StrategyType string `abi:"strategyType"`
}

type createStrategyParam struct {
	ResourceID    string                   `abi:"resourceId"`
	Connects      []createConnectInfoParam `abi:"connects"`
	ResourceName  string                   `abi:"resourceName"`
	Abstract      string                   `abi:"abstract"`
	Operation     string                   `abi:"operation"`
	StrategyNodes []strategyNodeParam      `abi:"strategyNodes"`
	Rules         []ruleParam              `abi:"rules"`
}

type upgradeStrategyParam struct {
	ResourceID    string                    `abi:"resourceId"`
	Connects      []upgradeConnectInfoParam `abi:"connects"`
	ResourceName  string                    `abi:"resourceName"`
	Abstract      string                    `abi:"abstract"`
	Operation     string                    `abi:"operation"`
	StrategyNodes []strategyNodeParam       `abi:"strategyNodes"`
	Rules         []ruleParam               `abi:"rules"`
}

type signatoryParam struct {
	ID   string `abi:"id"`
	Name string `abi:"name"`
	Sign []byte `abi:"sign"`
}

type createContractParam struct {
	ContractID        string                `abi:"contractId"`
	ContractName      string                `abi:"contractName"`
	ContractAbstract  string                `abi:"contractAbstract"`
	SignMode          string                `abi:"signMode"`
	HasPrivacyCompute bool                  `abi:"hasPrivacyCompute"`
	ActivationTime    uint64                `abi:"activationTime"`
	EndTime           uint64                `abi:"endTime"`
	Strategies        []createStrategyParam `abi:"strategies"`
	Signatories       []signatoryParam      `abi:"signatories"`
	Code              []byte                `abi:"code"`
}

type upgradeContractParam struct {
	ContractID        string                 `abi:"contractId"`
	ContractName      string                 `abi:"contractName"`
	ContractAbstract  string                 `abi:"contractAbstract"`
	SignMode          string                 `abi:"signMode"`
	HasPrivacyCompute bool                   `abi:"hasPrivacyCompute"`
	ActivationTime    uint64                 `abi:"activationTime"`
	EndTime           uint64                 `abi:"endTime"`
	Strategies        []upgradeStrategyParam `abi:"strategies"`
	Signatories       []signatoryParam       `abi:"signatories"`
	Code              []byte                 `abi:"code"`
}

type dataResourceColumnParam struct {
	ColName    string `abi:"colName"`
	ColType    string `abi:"colType"`
	ColComment string `abi:"colComment"`
}

type dataResourceTableParam struct {
	TableID      string                    `abi:"tableId"`
	DatasourceID string                    `abi:"datasourceId"`
	TableName    string                    `abi:"tableName"`
	Columns      []dataResourceColumnParam `abi:"columns"`
}

type resourceInfoParam struct {
	ResourceID        string                   `abi:"resourceId"`
	ResourceName      string                   `abi:"resourceName"`
	ResourceStage     string                   `abi:"resourceStage"`
	Owner             string                   `abi:"owner"`
	Desc              string                   `abi:"desc"`
	ValidityStartTime string                   `abi:"validityStartTime"`
	ValidityEndTime   string                   `abi:"validityEndTime"`
	Tables            []dataResourceTableParam `abi:"tables"`
}

type sourceInfoParam struct {
	ID         string `abi:"id"`
	Name       string `abi:"name"`
	Owner      string `abi:"owner"`
	Abstract   string `abi:"abstract"`
	SourceType string `abi:"sourceType"`
	Config     string `abi:"config"`
}

func toAccessContractParam(args AccessParams) (accessContractParam, error) {
	addr, err := convert.ZltcToAddress(args.ContractAddr)
	if err != nil {
		return accessContractParam{}, err
	}

	return accessContractParam{
		ContractAddr: addr,
		ResourceID:   args.ResourceID,
		Operation:    args.Operation,
		Rules:        toRuleParams(args.Rules),
	}, nil
}

func toCreateContractParam(args ContractParams) (createContractParam, error) {
	signatories, err := toSignatoryParams(args.Signatories)
	if err != nil {
		return createContractParam{}, err
	}

	code, err := hexutil.Decode(args.Code)
	if err != nil {
		return createContractParam{}, err
	}

	return createContractParam{
		ContractID:        args.ContractID,
		ContractName:      args.ContractName,
		ContractAbstract:  args.ContractAbstract,
		SignMode:          args.SignMode,
		HasPrivacyCompute: args.HasPrivacyCompute,
		ActivationTime:    args.ActivationTime,
		EndTime:           args.EndTime,
		Strategies:        toCreateStrategyParams(args.Strategies),
		Signatories:       signatories,
		Code:              code,
	}, nil
}

func toUpgradeContractParam(args UpgradeContractParams) (upgradeContractParam, error) {
	signatories, err := toSignatoryParams(args.Signatories)
	if err != nil {
		return upgradeContractParam{}, err
	}

	code, err := hexutil.Decode(args.Code)
	if err != nil {
		return upgradeContractParam{}, err
	}

	return upgradeContractParam{
		ContractID:        args.ContractID,
		ContractName:      args.ContractName,
		ContractAbstract:  args.ContractAbstract,
		SignMode:          args.SignMode,
		HasPrivacyCompute: args.HasPrivacyCompute,
		ActivationTime:    args.ActivationTime,
		EndTime:           args.EndTime,
		Strategies:        toUpgradeStrategyParams(args.Strategies),
		Signatories:       signatories,
		Code:              code,
	}, nil
}

func toResourceInfoParam(args ResourceInfoParams) resourceInfoParam {
	return resourceInfoParam{
		ResourceID:        args.ResourceID,
		ResourceName:      args.ResourceName,
		ResourceStage:     args.ResourceStage,
		Owner:             args.Owner,
		Desc:              args.Desc,
		ValidityStartTime: args.ValidityStartTime,
		ValidityEndTime:   args.ValidityEndTime,
		Tables: lo.Map(args.Tables, func(table DataResourceTable, _ int) dataResourceTableParam {
			return dataResourceTableParam{
				TableID:      table.TableID,
				DatasourceID: table.DatasourceID,
				TableName:    table.TableName,
				Columns: lo.Map(table.Columns, func(column DataResourceColumn, _ int) dataResourceColumnParam {
					return dataResourceColumnParam{
						ColName:    column.ColName,
						ColType:    column.ColType,
						ColComment: column.ColComment,
					}
				}),
			}
		}),
	}
}

func toRuleParams(rules []Rule) []ruleParam {
	return lo.Map(rules, func(rule Rule, _ int) ruleParam {
		return ruleParam{
			Name:           rule.Name,
			GRule:          rule.GRule,
			Type:           rule.Type,
			FactJSONString: rule.FactJSONString,
		}
	})
}

func toCreateStrategyParams(strategies []Strategy) []createStrategyParam {
	return lo.Map(strategies, func(strategy Strategy, _ int) createStrategyParam {
		return createStrategyParam{
			ResourceID:   strategy.ResourceID,
			ResourceName: strategy.ResourceName,
			Abstract:     strategy.Abstract,
			Operation:    strategy.Operation,
			Connects: lo.Map(strategy.Connects, func(connect CreateConnectInfo, _ int) createConnectInfoParam {
				return createConnectInfoParam{
					ConnectID:    connect.ConnectID,
					AccessType:   connect.AccessType,
					AccessConfig: connect.AccessConfig,
					Entity:       connect.Entity,
				}
			}),
			StrategyNodes: lo.Map(strategy.StrategyNodes, func(node StrategyNode, _ int) strategyNodeParam {
				return strategyNodeParam{
					NodeID:       node.NodeID,
					NodeName:     node.NodeName,
					NodePeerID:   node.NodePeerID,
					StrategyType: node.StrategyType,
				}
			}),
			Rules: toRuleParams(strategy.Rules),
		}
	})
}

func toUpgradeStrategyParams(strategies []UpgradeStrategy) []upgradeStrategyParam {
	return lo.Map(strategies, func(strategy UpgradeStrategy, _ int) upgradeStrategyParam {
		return upgradeStrategyParam{
			ResourceID:   strategy.ResourceID,
			ResourceName: strategy.ResourceName,
			Abstract:     strategy.Abstract,
			Operation:    strategy.Operation,
			Connects: lo.Map(strategy.Connects, func(connect UpgradeConnectInfo, _ int) upgradeConnectInfoParam {
				return upgradeConnectInfoParam{
					ConnectID:    connect.ConnectID,
					AccessType:   connect.AccessType,
					AccessConfig: connect.AccessConfig,
				}
			}),
			StrategyNodes: lo.Map(strategy.StrategyNodes, func(node StrategyNode, _ int) strategyNodeParam {
				return strategyNodeParam{
					NodeID:       node.NodeID,
					NodeName:     node.NodeName,
					NodePeerID:   node.NodePeerID,
					StrategyType: node.StrategyType,
				}
			}),
			Rules: toRuleParams(strategy.Rules),
		}
	})
}

func toSignatoryParams(signatories []Signatory) ([]signatoryParam, error) {
	params := make([]signatoryParam, 0, len(signatories))
	for _, signatory := range signatories {
		sign, err := hexutil.Decode(signatory.Sign)
		if err != nil {
			return nil, err
		}
		params = append(params, signatoryParam{
			ID:   signatory.ID,
			Name: signatory.Name,
			Sign: sign,
		})
	}
	return params, nil
}

var RuleEngineBuiltinContract = Contract{
	Description: "规则引擎合约",
	Address:     "zltc_aVwgkXjLrjc1TQaXAbD6WwgMFJ1mVuEGF",
	AbiString: `[
		{
			"type": "function",
			"name": "accessContract",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "contractAddr",
							"type": "address"
						},
						{
							"name": "resourceId",
							"type": "string"
						},
						{
							"name": "operation",
							"type": "string"
						},
						{
							"name": "rules",
							"type": "tuple[]",
							"components": [
								{
									"name": "name",
									"type": "string"
								},
								{
									"name": "grule",
									"type": "string"
								},
								{
									"name": "Type",
									"type": "uint8"
								},
								{
									"name": "factJsonString",
									"type": "string"
								}
							]
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "view"
		},
		{
			"type": "function",
			"name": "createContract",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "contractId",
							"type": "string"
						},
						{
							"name": "contractName",
							"type": "string"
						},
						{
							"name": "contractAbstract",
							"type": "string"
						},
						{
							"name": "signMode",
							"type": "string"
						},
						{
							"name": "hasPrivacyCompute",
							"type": "bool"
						},
						{
							"name": "activationTime",
							"type": "uint64"
						},
						{
							"name": "endTime",
							"type": "uint64"
						},
						{
							"name": "strategies",
							"type": "tuple[]",
							"components": [
								{
									"name": "resourceId",
									"type": "string"
								},
								{
									"name": "connects",
									"type": "tuple[]",
									"components": [
										{
											"name": "connectId",
											"type": "string"
										},
										{
											"name": "accessType",
											"type": "string"
										},
										{
											"name": "accessConfig",
											"type": "string"
										},
										{
											"name": "entity",
											"type": "string"
										}
									]
								},
								{
									"name": "resourceName",
									"type": "string"
								},
								{
									"name": "abstract",
									"type": "string"
								},
								{
									"name": "operation",
									"type": "string"
								},
								{
									"name": "strategyNodes",
									"type": "tuple[]",
									"components": [
										{
											"name": "nodeId",
											"type": "string"
										},
										{
											"name": "nodeName",
											"type": "string"
										},
										{
											"name": "nodePeerId",
											"type": "string"
										},
										{
											"name": "strategyType",
											"type": "string"
										}
									]
								},
								{
									"name": "rules",
									"type": "tuple[]",
									"components": [
										{
											"name": "name",
											"type": "string"
										},
										{
											"name": "grule",
											"type": "string"
										},
										{
											"name": "Type",
											"type": "uint8"
										},
										{
											"name": "factJsonString",
											"type": "string"
										}
									]
								}
							]
						},
						{
							"name": "signatories",
							"type": "tuple[]",
							"components": [
								{
									"name": "id",
									"type": "string"
								},
								{
									"name": "name",
									"type": "string"
								},
								{
									"name": "sign",
									"type": "bytes"
								}
							]
						},
						{
							"name": "code",
							"type": "bytes"
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "view"
		},
		{
			"type": "function",
			"name": "createResource",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "resourceId",
							"type": "string"
						},
						{
							"name": "resourceName",
							"type": "string"
						},
						{
							"name": "resourceStage",
							"type": "string"
						},
						{
							"name": "owner",
							"type": "string"
						},
						{
							"name": "desc",
							"type": "string"
						},
						{
							"name": "validityStartTime",
							"type": "string"
						},
						{
							"name": "validityEndTime",
							"type": "string"
						},
						{
							"name": "tables",
							"type": "tuple[]",
							"components": [
								{
									"name": "tableId",
									"type": "string"
								},
								{
									"name": "datasourceId",
									"type": "string"
								},
								{
									"name": "tableName",
									"type": "string"
								},
								{
									"name": "columns",
									"type": "tuple[]",
									"components": [
										{
											"name": "colName",
											"type": "string"
										},
										{
											"name": "colType",
											"type": "string"
										},
										{
											"name": "colComment",
											"type": "string"
										}
									]
								}
							]
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		},
		{
			"type": "function",
			"name": "upgradeContract",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "contractId",
							"type": "string"
						},
						{
							"name": "contractName",
							"type": "string"
						},
						{
							"name": "contractAbstract",
							"type": "string"
						},
						{
							"name": "signMode",
							"type": "string"
						},
						{
							"name": "hasPrivacyCompute",
							"type": "bool"
						},
						{
							"name": "activationTime",
							"type": "uint64"
						},
						{
							"name": "endTime",
							"type": "uint64"
						},
						{
							"name": "strategies",
							"type": "tuple[]",
							"components": [
								{
									"name": "resourceId",
									"type": "string"
								},
								{
									"name": "connects",
									"type": "tuple[]",
									"components": [
										{
											"name": "connectId",
											"type": "string"
										},
										{
											"name": "accessType",
											"type": "string"
										},
										{
											"name": "accessConfig",
											"type": "string"
										}
									]
								},
								{
									"name": "resourceName",
									"type": "string"
								},
								{
									"name": "abstract",
									"type": "string"
								},
								{
									"name": "operation",
									"type": "string"
								},
								{
									"name": "strategyNodes",
									"type": "tuple[]",
									"components": [
										{
											"name": "nodeId",
											"type": "string"
										},
										{
											"name": "nodeName",
											"type": "string"
										},
										{
											"name": "nodePeerId",
											"type": "string"
										},
										{
											"name": "strategyType",
											"type": "string"
										}
									]
								},
								{
									"name": "rules",
									"type": "tuple[]",
									"components": [
										{
											"name": "name",
											"type": "string"
										},
										{
											"name": "grule",
											"type": "string"
										},
										{
											"name": "Type",
											"type": "uint8"
										},
										{
											"name": "factJsonString",
											"type": "string"
										}
									]
								}
							]
						},
						{
							"name": "signatories",
							"type": "tuple[]",
							"components": [
								{
									"name": "id",
									"type": "string"
								},
								{
									"name": "name",
									"type": "string"
								},
								{
									"name": "sign",
									"type": "bytes"
								}
							]
						},
						{
							"name": "code",
							"type": "bytes"
						}
					]
				},
				{
					"name": "contractAddr",
					"type": "address"
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		},
		{
			"type": "function",
			"name": "signContract",
			"inputs": [
				{
					"name": "contractAddr",
					"type": "address"
				},
				{
					"name": "signatories",
					"type": "tuple[]",
					"components": [
						{
							"name": "id",
							"type": "string"
						},
						{
							"name": "name",
							"type": "string"
						},
						{
							"name": "sign",
							"type": "bytes"
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		},
		{
			"type": "function",
			"name": "createSource",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "id",
							"type": "string"
						},
						{
							"name": "name",
							"type": "string"
						},
						{
							"name": "owner",
							"type": "string"
						},
						{
							"name": "abstract",
							"type": "string"
						},
						{
							"name": "sourceType",
							"type": "string"
						},
						{
							"name": "config",
							"type": "string"
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		},
		{
			"type": "function",
			"name": "createProduct",
			"inputs": [
				{
					"name": "args",
					"type": "tuple",
					"components": [
						{
							"name": "resourceId",
							"type": "string"
						},
						{
							"name": "resourceName",
							"type": "string"
						},
						{
							"name": "resourceStage",
							"type": "string"
						},
						{
							"name": "owner",
							"type": "string"
						},
						{
							"name": "desc",
							"type": "string"
						},
						{
							"name": "validityStartTime",
							"type": "string"
						},
						{
							"name": "validityEndTime",
							"type": "string"
						},
						{
							"name": "tables",
							"type": "tuple[]",
							"components": [
								{
									"name": "tableId",
									"type": "string"
								},
								{
									"name": "datasourceId",
									"type": "string"
								},
								{
									"name": "tableName",
									"type": "string"
								},
								{
									"name": "columns",
									"type": "tuple[]",
									"components": [
										{
											"name": "colName",
											"type": "string"
										},
										{
											"name": "colType",
											"type": "string"
										},
										{
											"name": "colComment",
											"type": "string"
										}
									]
								}
							]
						}
					]
				}
			],
			"outputs": [],
			"stateMutability": "nonpayable"
		}
	]`,
}
