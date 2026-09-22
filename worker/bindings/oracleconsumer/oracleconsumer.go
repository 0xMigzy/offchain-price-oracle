// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracleconsumer

import (
	"context"
	"errors"
	"math/big"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
	_ = time.Tick
	_ = context.Background
)

// OracleconsumerMetaData contains all meta data concerning the Oracleconsumer contract.
var OracleconsumerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addSubmitter\",\"inputs\":[{\"name\":\"_submitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getLastPrice\",\"inputs\":[{\"name\":\"pair\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPriceHistory\",\"inputs\":[{\"name\":\"pair\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"lastPrices\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"price\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"submitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"paused\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"pendingOwner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"priceHistory\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"price\",\"type\":\"int256\",\"internalType\":\"int256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"submitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeSubmitter\",\"inputs\":[{\"name\":\"_submitter\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitPrice\",\"inputs\":[{\"name\":\"pair\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"price\",\"type\":\"int256\",\"internalType\":\"int256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"submitters\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unpause\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferStarted\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Paused\",\"inputs\":[{\"name\":\"by\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PriceSubmitted\",\"inputs\":[{\"name\":\"pair\",\"type\":\"string\",\"indexed\":true,\"internalType\":\"string\"},{\"name\":\"price\",\"type\":\"int256\",\"indexed\":false,\"internalType\":\"int256\"},{\"name\":\"timestamp\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"submitter\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Unpaused\",\"inputs\":[{\"name\":\"by\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"ContractPaused\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"EmptyPair\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPrice\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NoChange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotAuthorizedSubmitter\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotPendingOwner\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]}]",
	Bin: "0x608060405234801561000f575f5ffd5b503360035f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600160025f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff165f73ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a36114c18061010c5f395ff3fe608060405234801561000f575f5ffd5b50600436106100f3575f3560e01c80638456cb5911610095578063c36f41e211610064578063c36f41e21461024c578063d5eaf6bf14610268578063e30c39781461029a578063f2fde38b146102b8576100f3565b80638456cb59146101d65780638da5cb5b146101e057806395316ed6146101fe578063a8c608c81461021a576100f3565b80633dbc43b5116100d15780633dbc43b5146101745780633f4ba83a146101a45780635c975abb146101ae57806379ba5097146101cc576100f3565b806306560435146100f7578063072900f91461012757806312f8acb514610143575b5f5ffd5b610111600480360381019061010c9190610fe4565b6102d4565b60405161011e9190611029565b60405180910390f35b610141600480360381019061013c9190610fe4565b6102f1565b005b61015d600480360381019061015891906110a3565b6103cf565b60405161016b92919061111e565b60405180910390f35b61018e600480360381019061018991906110a3565b61047b565b60405161019b9190611145565b60405180910390f35b6101ac6104a8565b005b6101b66105d3565b6040516101c39190611029565b60405180910390f35b6101d46105e6565b005b6101de6107b0565b005b6101e86108dd565b6040516101f5919061116d565b60405180910390f35b610218600480360381019061021391906111b0565b610902565b005b610234600480360381019061022f919061136f565b610bfb565b604051610243939291906113c9565b60405180910390f35b61026660048036038101906102619190610fe4565b610c72565b005b610282600480360381019061027d91906113fe565b610d4f565b604051610291939291906113c9565b60405180910390f35b6102a2610dab565b6040516102af919061116d565b60405180910390f35b6102d260048036038101906102cd9190610fe4565b610dd0565b005b6002602052805f5260405f205f915054906101000a900460ff1681565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610377576040517f30cd747100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600160025f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555050565b5f5f5f5f85856040516103e3929190611473565b90815260200160405180910390206040518060600160405290815f820154815260200160018201548152602001600282015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815250509050805f0151816020015192509250509250929050565b5f6001838360405161048e929190611473565b908152602001604051809103902080549050905092915050565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461052e576040517f30cd747100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460149054906101000a900460ff16610574576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600460146101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa60405160405180910390a2565b600460149054906101000a900460ff1681565b60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461066c576040517f1853971c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1660035f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f60045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a350565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610836576040517f30cd747100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460149054906101000a900460ff161561087d576040517fa88ee57700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001600460146101000a81548160ff0219169083151502179055503373ffffffffffffffffffffffffffffffffffffffff167f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a25860405160405180910390a2565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60025f3373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f9054906101000a900460ff16610982576040517f34a8dd7f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600460149054906101000a900460ff16156109c9576040517fab35696f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8383905003610a05576040517f9272a79800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8113610a3d576040517ebfc92100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60405180606001604052808381526020014281526020013373ffffffffffffffffffffffffffffffffffffffff168152509050805f8585604051610a83929190611473565b90815260200160405180910390205f820151815f0155602082015181600101556040820151816002015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555090505060018484604051610afe929190611473565b908152602001604051809103902081908060018154018082558091505060019003905f5260205f2090600302015f909190919091505f820151815f0155602082015181600101556040820151816002015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050503373ffffffffffffffffffffffffffffffffffffffff168484604051610bb4929190611473565b60405180910390207ff7c154a9aaa42edfde55542ee1732563a130bad00c142426902340112fb139ef8442604051610bed92919061111e565b60405180910390a350505050565b6001828051602081018201805184825260208301602085012081835280955050505050508181548110610c2c575f80fd5b905f5260205f2090600302015f9150915050805f015490806001015490806002015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610cf8576040517f30cd747100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60025f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548160ff02191690831515021790555050565b5f818051602081018201805184825260208301602085012081835280955050505050505f91509050805f015490806001015490806002015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905083565b60045f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610e56576040517f30cd747100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610ebb576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8060045f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff1660035f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e2270060405160405180910390a350565b5f604051905090565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f610fb382610f8a565b9050919050565b610fc381610fa9565b8114610fcd575f5ffd5b50565b5f81359050610fde81610fba565b92915050565b5f60208284031215610ff957610ff8610f82565b5b5f61100684828501610fd0565b91505092915050565b5f8115159050919050565b6110238161100f565b82525050565b5f60208201905061103c5f83018461101a565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f84011261106357611062611042565b5b8235905067ffffffffffffffff8111156110805761107f611046565b5b60208301915083600182028301111561109c5761109b61104a565b5b9250929050565b5f5f602083850312156110b9576110b8610f82565b5b5f83013567ffffffffffffffff8111156110d6576110d5610f86565b5b6110e28582860161104e565b92509250509250929050565b5f819050919050565b611100816110ee565b82525050565b5f819050919050565b61111881611106565b82525050565b5f6040820190506111315f8301856110f7565b61113e602083018461110f565b9392505050565b5f6020820190506111585f83018461110f565b92915050565b61116781610fa9565b82525050565b5f6020820190506111805f83018461115e565b92915050565b61118f816110ee565b8114611199575f5ffd5b50565b5f813590506111aa81611186565b92915050565b5f5f5f604084860312156111c7576111c6610f82565b5b5f84013567ffffffffffffffff8111156111e4576111e3610f86565b5b6111f08682870161104e565b935093505060206112038682870161119c565b9150509250925092565b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b61125782611211565b810181811067ffffffffffffffff8211171561127657611275611221565b5b80604052505050565b5f611288610f79565b9050611294828261124e565b919050565b5f67ffffffffffffffff8211156112b3576112b2611221565b5b6112bc82611211565b9050602081019050919050565b828183375f83830152505050565b5f6112e96112e484611299565b61127f565b9050828152602081018484840111156113055761130461120d565b5b6113108482856112c9565b509392505050565b5f82601f83011261132c5761132b611042565b5b813561133c8482602086016112d7565b91505092915050565b61134e81611106565b8114611358575f5ffd5b50565b5f8135905061136981611345565b92915050565b5f5f6040838503121561138557611384610f82565b5b5f83013567ffffffffffffffff8111156113a2576113a1610f86565b5b6113ae85828601611318565b92505060206113bf8582860161135b565b9150509250929050565b5f6060820190506113dc5f8301866110f7565b6113e9602083018561110f565b6113f6604083018461115e565b949350505050565b5f6020828403121561141357611412610f82565b5b5f82013567ffffffffffffffff8111156114305761142f610f86565b5b61143c84828501611318565b91505092915050565b5f81905092915050565b5f61145a8385611445565b93506114678385846112c9565b82840190509392505050565b5f61147f82848661144f565b9150819050939250505056fea26469706673582212209304f9e601922a9c9d02ff513be0892d600eb435f33e70659d70dfe751bde8be64736f6c63430008230033",
}

// OracleconsumerABI is the input ABI used to generate the binding from.
// Deprecated: Use OracleconsumerMetaData.ABI instead.
var OracleconsumerABI = OracleconsumerMetaData.ABI

// OracleconsumerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use OracleconsumerMetaData.Bin instead.
var OracleconsumerBin = OracleconsumerMetaData.Bin

// DeployOracleconsumer deploys a new Ethereum contract, binding an instance of Oracleconsumer to it.
func DeployOracleconsumer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Oracleconsumer, error) {
	parsed, err := OracleconsumerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(OracleconsumerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Oracleconsumer{OracleconsumerCaller: OracleconsumerCaller{contract: contract}, OracleconsumerTransactor: OracleconsumerTransactor{contract: contract}, OracleconsumerFilterer: OracleconsumerFilterer{contract: contract}}, nil
}

// Oracleconsumer is an auto generated Go binding around an Ethereum contract.
type Oracleconsumer struct {
	OracleconsumerCaller     // Read-only binding to the contract
	OracleconsumerTransactor // Write-only binding to the contract
	OracleconsumerFilterer   // Log filterer for contract events
}

// OracleconsumerCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleconsumerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleconsumerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleconsumerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleconsumerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleconsumerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleconsumerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleconsumerSession struct {
	Contract     *Oracleconsumer   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleconsumerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleconsumerCallerSession struct {
	Contract *OracleconsumerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// OracleconsumerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleconsumerTransactorSession struct {
	Contract     *OracleconsumerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// OracleconsumerRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleconsumerRaw struct {
	Contract *Oracleconsumer // Generic contract binding to access the raw methods on
}

// OracleconsumerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleconsumerCallerRaw struct {
	Contract *OracleconsumerCaller // Generic read-only contract binding to access the raw methods on
}

// OracleconsumerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleconsumerTransactorRaw struct {
	Contract *OracleconsumerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleconsumer creates a new instance of Oracleconsumer, bound to a specific deployed contract.
func NewOracleconsumer(address common.Address, backend bind.ContractBackend) (*Oracleconsumer, error) {
	contract, err := bindOracleconsumer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Oracleconsumer{OracleconsumerCaller: OracleconsumerCaller{contract: contract}, OracleconsumerTransactor: OracleconsumerTransactor{contract: contract}, OracleconsumerFilterer: OracleconsumerFilterer{contract: contract}}, nil
}

// NewOracleconsumerCaller creates a new read-only instance of Oracleconsumer, bound to a specific deployed contract.
func NewOracleconsumerCaller(address common.Address, caller bind.ContractCaller) (*OracleconsumerCaller, error) {
	contract, err := bindOracleconsumer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerCaller{contract: contract}, nil
}

// NewOracleconsumerTransactor creates a new write-only instance of Oracleconsumer, bound to a specific deployed contract.
func NewOracleconsumerTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleconsumerTransactor, error) {
	contract, err := bindOracleconsumer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerTransactor{contract: contract}, nil
}

// NewOracleconsumerFilterer creates a new log filterer instance of Oracleconsumer, bound to a specific deployed contract.
func NewOracleconsumerFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleconsumerFilterer, error) {
	contract, err := bindOracleconsumer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerFilterer{contract: contract}, nil
}

// bindOracleconsumer binds a generic wrapper to an already deployed contract.
func bindOracleconsumer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := OracleconsumerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracleconsumer *OracleconsumerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracleconsumer.Contract.OracleconsumerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracleconsumer *OracleconsumerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.OracleconsumerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracleconsumer *OracleconsumerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.OracleconsumerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Oracleconsumer *OracleconsumerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Oracleconsumer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Oracleconsumer *OracleconsumerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Oracleconsumer *OracleconsumerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.contract.Transact(opts, method, params...)
}

// GetLastPrice is a free data retrieval call binding the contract method 0x12f8acb5.
//
// Solidity: function getLastPrice(string pair) view returns(int256, uint256)
func (_Oracleconsumer *OracleconsumerCaller) GetLastPrice(opts *bind.CallOpts, pair string) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "getLastPrice", pair)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetLastPrice is a free data retrieval call binding the contract method 0x12f8acb5.
//
// Solidity: function getLastPrice(string pair) view returns(int256, uint256)
func (_Oracleconsumer *OracleconsumerSession) GetLastPrice(pair string) (*big.Int, *big.Int, error) {
	return _Oracleconsumer.Contract.GetLastPrice(&_Oracleconsumer.CallOpts, pair)
}

// GetLastPrice is a free data retrieval call binding the contract method 0x12f8acb5.
//
// Solidity: function getLastPrice(string pair) view returns(int256, uint256)
func (_Oracleconsumer *OracleconsumerCallerSession) GetLastPrice(pair string) (*big.Int, *big.Int, error) {
	return _Oracleconsumer.Contract.GetLastPrice(&_Oracleconsumer.CallOpts, pair)
}

// GetPriceHistory is a free data retrieval call binding the contract method 0x3dbc43b5.
//
// Solidity: function getPriceHistory(string pair) view returns(uint256)
func (_Oracleconsumer *OracleconsumerCaller) GetPriceHistory(opts *bind.CallOpts, pair string) (*big.Int, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "getPriceHistory", pair)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPriceHistory is a free data retrieval call binding the contract method 0x3dbc43b5.
//
// Solidity: function getPriceHistory(string pair) view returns(uint256)
func (_Oracleconsumer *OracleconsumerSession) GetPriceHistory(pair string) (*big.Int, error) {
	return _Oracleconsumer.Contract.GetPriceHistory(&_Oracleconsumer.CallOpts, pair)
}

// GetPriceHistory is a free data retrieval call binding the contract method 0x3dbc43b5.
//
// Solidity: function getPriceHistory(string pair) view returns(uint256)
func (_Oracleconsumer *OracleconsumerCallerSession) GetPriceHistory(pair string) (*big.Int, error) {
	return _Oracleconsumer.Contract.GetPriceHistory(&_Oracleconsumer.CallOpts, pair)
}

// LastPrices is a free data retrieval call binding the contract method 0xd5eaf6bf.
//
// Solidity: function lastPrices(string ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerCaller) LastPrices(opts *bind.CallOpts, arg0 string) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "lastPrices", arg0)

	outstruct := new(struct {
		Price     *big.Int
		Timestamp *big.Int
		Submitter common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Price = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Submitter = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// LastPrices is a free data retrieval call binding the contract method 0xd5eaf6bf.
//
// Solidity: function lastPrices(string ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerSession) LastPrices(arg0 string) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	return _Oracleconsumer.Contract.LastPrices(&_Oracleconsumer.CallOpts, arg0)
}

// LastPrices is a free data retrieval call binding the contract method 0xd5eaf6bf.
//
// Solidity: function lastPrices(string ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerCallerSession) LastPrices(arg0 string) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	return _Oracleconsumer.Contract.LastPrices(&_Oracleconsumer.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracleconsumer *OracleconsumerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracleconsumer *OracleconsumerSession) Owner() (common.Address, error) {
	return _Oracleconsumer.Contract.Owner(&_Oracleconsumer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Oracleconsumer *OracleconsumerCallerSession) Owner() (common.Address, error) {
	return _Oracleconsumer.Contract.Owner(&_Oracleconsumer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracleconsumer *OracleconsumerCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracleconsumer *OracleconsumerSession) Paused() (bool, error) {
	return _Oracleconsumer.Contract.Paused(&_Oracleconsumer.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Oracleconsumer *OracleconsumerCallerSession) Paused() (bool, error) {
	return _Oracleconsumer.Contract.Paused(&_Oracleconsumer.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracleconsumer *OracleconsumerCaller) PendingOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "pendingOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracleconsumer *OracleconsumerSession) PendingOwner() (common.Address, error) {
	return _Oracleconsumer.Contract.PendingOwner(&_Oracleconsumer.CallOpts)
}

// PendingOwner is a free data retrieval call binding the contract method 0xe30c3978.
//
// Solidity: function pendingOwner() view returns(address)
func (_Oracleconsumer *OracleconsumerCallerSession) PendingOwner() (common.Address, error) {
	return _Oracleconsumer.Contract.PendingOwner(&_Oracleconsumer.CallOpts)
}

// PriceHistory is a free data retrieval call binding the contract method 0xa8c608c8.
//
// Solidity: function priceHistory(string , uint256 ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerCaller) PriceHistory(opts *bind.CallOpts, arg0 string, arg1 *big.Int) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "priceHistory", arg0, arg1)

	outstruct := new(struct {
		Price     *big.Int
		Timestamp *big.Int
		Submitter common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Price = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Submitter = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// PriceHistory is a free data retrieval call binding the contract method 0xa8c608c8.
//
// Solidity: function priceHistory(string , uint256 ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerSession) PriceHistory(arg0 string, arg1 *big.Int) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	return _Oracleconsumer.Contract.PriceHistory(&_Oracleconsumer.CallOpts, arg0, arg1)
}

// PriceHistory is a free data retrieval call binding the contract method 0xa8c608c8.
//
// Solidity: function priceHistory(string , uint256 ) view returns(int256 price, uint256 timestamp, address submitter)
func (_Oracleconsumer *OracleconsumerCallerSession) PriceHistory(arg0 string, arg1 *big.Int) (struct {
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
}, error) {
	return _Oracleconsumer.Contract.PriceHistory(&_Oracleconsumer.CallOpts, arg0, arg1)
}

// Submitters is a free data retrieval call binding the contract method 0x06560435.
//
// Solidity: function submitters(address ) view returns(bool)
func (_Oracleconsumer *OracleconsumerCaller) Submitters(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Oracleconsumer.contract.Call(opts, &out, "submitters", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Submitters is a free data retrieval call binding the contract method 0x06560435.
//
// Solidity: function submitters(address ) view returns(bool)
func (_Oracleconsumer *OracleconsumerSession) Submitters(arg0 common.Address) (bool, error) {
	return _Oracleconsumer.Contract.Submitters(&_Oracleconsumer.CallOpts, arg0)
}

// Submitters is a free data retrieval call binding the contract method 0x06560435.
//
// Solidity: function submitters(address ) view returns(bool)
func (_Oracleconsumer *OracleconsumerCallerSession) Submitters(arg0 common.Address) (bool, error) {
	return _Oracleconsumer.Contract.Submitters(&_Oracleconsumer.CallOpts, arg0)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracleconsumer *OracleconsumerTransactor) AcceptOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "acceptOwnership")
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracleconsumer *OracleconsumerSession) AcceptOwnership() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.AcceptOwnership(&_Oracleconsumer.TransactOpts)
}

// AcceptOwnership is a paid mutator transaction binding the contract method 0x79ba5097.
//
// Solidity: function acceptOwnership() returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) AcceptOwnership() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.AcceptOwnership(&_Oracleconsumer.TransactOpts)
}

// AddSubmitter is a paid mutator transaction binding the contract method 0x072900f9.
//
// Solidity: function addSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerTransactor) AddSubmitter(opts *bind.TransactOpts, _submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "addSubmitter", _submitter)
}

// AddSubmitter is a paid mutator transaction binding the contract method 0x072900f9.
//
// Solidity: function addSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerSession) AddSubmitter(_submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.AddSubmitter(&_Oracleconsumer.TransactOpts, _submitter)
}

// AddSubmitter is a paid mutator transaction binding the contract method 0x072900f9.
//
// Solidity: function addSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) AddSubmitter(_submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.AddSubmitter(&_Oracleconsumer.TransactOpts, _submitter)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracleconsumer *OracleconsumerTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracleconsumer *OracleconsumerSession) Pause() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.Pause(&_Oracleconsumer.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) Pause() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.Pause(&_Oracleconsumer.TransactOpts)
}

// RemoveSubmitter is a paid mutator transaction binding the contract method 0xc36f41e2.
//
// Solidity: function removeSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerTransactor) RemoveSubmitter(opts *bind.TransactOpts, _submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "removeSubmitter", _submitter)
}

// RemoveSubmitter is a paid mutator transaction binding the contract method 0xc36f41e2.
//
// Solidity: function removeSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerSession) RemoveSubmitter(_submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.RemoveSubmitter(&_Oracleconsumer.TransactOpts, _submitter)
}

// RemoveSubmitter is a paid mutator transaction binding the contract method 0xc36f41e2.
//
// Solidity: function removeSubmitter(address _submitter) returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) RemoveSubmitter(_submitter common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.RemoveSubmitter(&_Oracleconsumer.TransactOpts, _submitter)
}

// SubmitPrice is a paid mutator transaction binding the contract method 0x95316ed6.
//
// Solidity: function submitPrice(string pair, int256 price) returns()
func (_Oracleconsumer *OracleconsumerTransactor) SubmitPrice(opts *bind.TransactOpts, pair string, price *big.Int) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "submitPrice", pair, price)
}

// SubmitPrice is a paid mutator transaction binding the contract method 0x95316ed6.
//
// Solidity: function submitPrice(string pair, int256 price) returns()
func (_Oracleconsumer *OracleconsumerSession) SubmitPrice(pair string, price *big.Int) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.SubmitPrice(&_Oracleconsumer.TransactOpts, pair, price)
}

// SubmitPrice is a paid mutator transaction binding the contract method 0x95316ed6.
//
// Solidity: function submitPrice(string pair, int256 price) returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) SubmitPrice(pair string, price *big.Int) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.SubmitPrice(&_Oracleconsumer.TransactOpts, pair, price)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracleconsumer *OracleconsumerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracleconsumer *OracleconsumerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.TransferOwnership(&_Oracleconsumer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Oracleconsumer.Contract.TransferOwnership(&_Oracleconsumer.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracleconsumer *OracleconsumerTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Oracleconsumer.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracleconsumer *OracleconsumerSession) Unpause() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.Unpause(&_Oracleconsumer.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Oracleconsumer *OracleconsumerTransactorSession) Unpause() (*types.Transaction, error) {
	return _Oracleconsumer.Contract.Unpause(&_Oracleconsumer.TransactOpts)
}

// OracleconsumerOwnershipTransferStartedIterator is returned from FilterOwnershipTransferStarted and is used to iterate over the raw logs and unpacked data for OwnershipTransferStarted events raised by the Oracleconsumer contract.
type OracleconsumerOwnershipTransferStartedIterator struct {
	Event *OracleconsumerOwnershipTransferStarted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleconsumerOwnershipTransferStartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleconsumerOwnershipTransferStarted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleconsumerOwnershipTransferStarted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleconsumerOwnershipTransferStartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleconsumerOwnershipTransferStartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleconsumerOwnershipTransferStarted represents a OwnershipTransferStarted event raised by the Oracleconsumer contract.
type OracleconsumerOwnershipTransferStarted struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferStarted is a free log retrieval operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) FilterOwnershipTransferStarted(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OracleconsumerOwnershipTransferStartedIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oracleconsumer.contract.FilterLogs(opts, "OwnershipTransferStarted", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerOwnershipTransferStartedIterator{contract: _Oracleconsumer.contract, event: "OwnershipTransferStarted", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferStarted is a free log subscription operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) WatchOwnershipTransferStarted(opts *bind.WatchOpts, sink chan<- *OracleconsumerOwnershipTransferStarted, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oracleconsumer.contract.WatchLogs(opts, "OwnershipTransferStarted", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleconsumerOwnershipTransferStarted)
				if err := _Oracleconsumer.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferStarted is a log parse operation binding the contract event 0x38d16b8cac22d99fc7c124b9cd0de2d3fa1faef420bfe791d8c362d765e22700.
//
// Solidity: event OwnershipTransferStarted(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) ParseOwnershipTransferStarted(log types.Log) (*OracleconsumerOwnershipTransferStarted, error) {
	event := new(OracleconsumerOwnershipTransferStarted)
	if err := _Oracleconsumer.contract.UnpackLog(event, "OwnershipTransferStarted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleconsumerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Oracleconsumer contract.
type OracleconsumerOwnershipTransferredIterator struct {
	Event *OracleconsumerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleconsumerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleconsumerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleconsumerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleconsumerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleconsumerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleconsumerOwnershipTransferred represents a OwnershipTransferred event raised by the Oracleconsumer contract.
type OracleconsumerOwnershipTransferred struct {
	From common.Address
	To   common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*OracleconsumerOwnershipTransferredIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oracleconsumer.contract.FilterLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerOwnershipTransferredIterator{contract: _Oracleconsumer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OracleconsumerOwnershipTransferred, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Oracleconsumer.contract.WatchLogs(opts, "OwnershipTransferred", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleconsumerOwnershipTransferred)
				if err := _Oracleconsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed from, address indexed to)
func (_Oracleconsumer *OracleconsumerFilterer) ParseOwnershipTransferred(log types.Log) (*OracleconsumerOwnershipTransferred, error) {
	event := new(OracleconsumerOwnershipTransferred)
	if err := _Oracleconsumer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleconsumerPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Oracleconsumer contract.
type OracleconsumerPausedIterator struct {
	Event *OracleconsumerPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleconsumerPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleconsumerPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleconsumerPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleconsumerPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleconsumerPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleconsumerPaused represents a Paused event raised by the Oracleconsumer contract.
type OracleconsumerPaused struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) FilterPaused(opts *bind.FilterOpts, by []common.Address) (*OracleconsumerPausedIterator, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _Oracleconsumer.contract.FilterLogs(opts, "Paused", byRule)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerPausedIterator{contract: _Oracleconsumer.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *OracleconsumerPaused, by []common.Address) (event.Subscription, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _Oracleconsumer.contract.WatchLogs(opts, "Paused", byRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleconsumerPaused)
				if err := _Oracleconsumer.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) ParsePaused(log types.Log) (*OracleconsumerPaused, error) {
	event := new(OracleconsumerPaused)
	if err := _Oracleconsumer.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleconsumerPriceSubmittedIterator is returned from FilterPriceSubmitted and is used to iterate over the raw logs and unpacked data for PriceSubmitted events raised by the Oracleconsumer contract.
type OracleconsumerPriceSubmittedIterator struct {
	Event *OracleconsumerPriceSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleconsumerPriceSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleconsumerPriceSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleconsumerPriceSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleconsumerPriceSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleconsumerPriceSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleconsumerPriceSubmitted represents a PriceSubmitted event raised by the Oracleconsumer contract.
type OracleconsumerPriceSubmitted struct {
	Pair      common.Hash
	Price     *big.Int
	Timestamp *big.Int
	Submitter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPriceSubmitted is a free log retrieval operation binding the contract event 0xf7c154a9aaa42edfde55542ee1732563a130bad00c142426902340112fb139ef.
//
// Solidity: event PriceSubmitted(string indexed pair, int256 price, uint256 timestamp, address indexed submitter)
func (_Oracleconsumer *OracleconsumerFilterer) FilterPriceSubmitted(opts *bind.FilterOpts, pair []string, submitter []common.Address) (*OracleconsumerPriceSubmittedIterator, error) {

	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	var submitterRule []interface{}
	for _, submitterItem := range submitter {
		submitterRule = append(submitterRule, submitterItem)
	}

	logs, sub, err := _Oracleconsumer.contract.FilterLogs(opts, "PriceSubmitted", pairRule, submitterRule)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerPriceSubmittedIterator{contract: _Oracleconsumer.contract, event: "PriceSubmitted", logs: logs, sub: sub}, nil
}

// WatchPriceSubmitted is a free log subscription operation binding the contract event 0xf7c154a9aaa42edfde55542ee1732563a130bad00c142426902340112fb139ef.
//
// Solidity: event PriceSubmitted(string indexed pair, int256 price, uint256 timestamp, address indexed submitter)
func (_Oracleconsumer *OracleconsumerFilterer) WatchPriceSubmitted(opts *bind.WatchOpts, sink chan<- *OracleconsumerPriceSubmitted, pair []string, submitter []common.Address) (event.Subscription, error) {

	var pairRule []interface{}
	for _, pairItem := range pair {
		pairRule = append(pairRule, pairItem)
	}

	var submitterRule []interface{}
	for _, submitterItem := range submitter {
		submitterRule = append(submitterRule, submitterItem)
	}

	logs, sub, err := _Oracleconsumer.contract.WatchLogs(opts, "PriceSubmitted", pairRule, submitterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleconsumerPriceSubmitted)
				if err := _Oracleconsumer.contract.UnpackLog(event, "PriceSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePriceSubmitted is a log parse operation binding the contract event 0xf7c154a9aaa42edfde55542ee1732563a130bad00c142426902340112fb139ef.
//
// Solidity: event PriceSubmitted(string indexed pair, int256 price, uint256 timestamp, address indexed submitter)
func (_Oracleconsumer *OracleconsumerFilterer) ParsePriceSubmitted(log types.Log) (*OracleconsumerPriceSubmitted, error) {
	event := new(OracleconsumerPriceSubmitted)
	if err := _Oracleconsumer.contract.UnpackLog(event, "PriceSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// OracleconsumerUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the Oracleconsumer contract.
type OracleconsumerUnpausedIterator struct {
	Event *OracleconsumerUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleconsumerUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleconsumerUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleconsumerUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleconsumerUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleconsumerUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleconsumerUnpaused represents a Unpaused event raised by the Oracleconsumer contract.
type OracleconsumerUnpaused struct {
	By  common.Address
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) FilterUnpaused(opts *bind.FilterOpts, by []common.Address) (*OracleconsumerUnpausedIterator, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _Oracleconsumer.contract.FilterLogs(opts, "Unpaused", byRule)
	if err != nil {
		return nil, err
	}
	return &OracleconsumerUnpausedIterator{contract: _Oracleconsumer.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *OracleconsumerUnpaused, by []common.Address) (event.Subscription, error) {

	var byRule []interface{}
	for _, byItem := range by {
		byRule = append(byRule, byItem)
	}

	logs, sub, err := _Oracleconsumer.contract.WatchLogs(opts, "Unpaused", byRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleconsumerUnpaused)
				if err := _Oracleconsumer.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address indexed by)
func (_Oracleconsumer *OracleconsumerFilterer) ParseUnpaused(log types.Log) (*OracleconsumerUnpaused, error) {
	event := new(OracleconsumerUnpaused)
	if err := _Oracleconsumer.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
