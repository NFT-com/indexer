// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package cryptokitties

import (
	"errors"
	"math/big"
	"strings"

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
)

// MainMetaData contains all meta data concerning the Main contract.
var MainMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[{\"name\":\"_interfaceID\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"},{\"name\":\"_preferredTransport\",\"type\":\"string\"}],\"name\":\"tokenMetadata\",\"outputs\":[{\"name\":\"infoUrl\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"promoCreatedCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GEN0_STARTING_PRICE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSiringAuctionAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"pregnantKitties\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_kittyId\",\"type\":\"uint256\"}],\"name\":\"isPregnant\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GEN0_AUCTION_DURATION\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"siringAuction\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setGeneScienceAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_kittyId\",\"type\":\"uint256\"},{\"name\":\"_startingPrice\",\"type\":\"uint256\"},{\"name\":\"_endingPrice\",\"type\":\"uint256\"},{\"name\":\"_duration\",\"type\":\"uint256\"}],\"name\":\"createSaleAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"sireAllowedToAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_matronId\",\"type\":\"uint256\"},{\"name\":\"_sireId\",\"type\":\"uint256\"}],\"name\":\"canBreedWith\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"kittyIndexToApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_kittyId\",\"type\":\"uint256\"},{\"name\":\"_startingPrice\",\"type\":\"uint256\"},{\"name\":\"_endingPrice\",\"type\":\"uint256\"},{\"name\":\"_duration\",\"type\":\"uint256\"}],\"name\":\"createSiringAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"val\",\"type\":\"uint256\"}],\"name\":\"setAutoBirthFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"},{\"name\":\"_sireId\",\"type\":\"uint256\"}],\"name\":\"approveSiring\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_genes\",\"type\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"createPromoKitty\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"secs\",\"type\":\"uint256\"}],\"name\":\"setSecondsPerBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawBalance\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"GEN0_CREATION_LIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"newContractAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"setSaleAuctionAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"count\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_v2Address\",\"type\":\"address\"}],\"name\":\"setNewAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"secondsPerBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"tokensOfOwner\",\"outputs\":[{\"name\":\"ownerTokens\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_matronId\",\"type\":\"uint256\"}],\"name\":\"giveBirth\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdrawAuctionBalances\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"cooldowns\",\"outputs\":[{\"name\":\"\",\"type\":\"uint32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"kittyIndexToOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"autoBirthFee\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"erc721Metadata\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_genes\",\"type\":\"uint256\"}],\"name\":\"createGen0Auction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_kittyId\",\"type\":\"uint256\"}],\"name\":\"isReadyToBreed\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"PROMO_CREATION_LIMIT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_contractAddress\",\"type\":\"address\"}],\"name\":\"setMetadataAddress\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"saleAuction\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_id\",\"type\":\"uint256\"}],\"name\":\"getKitty\",\"outputs\":[{\"name\":\"isGestating\",\"type\":\"bool\"},{\"name\":\"isReady\",\"type\":\"bool\"},{\"name\":\"cooldownIndex\",\"type\":\"uint256\"},{\"name\":\"nextActionAt\",\"type\":\"uint256\"},{\"name\":\"siringWithId\",\"type\":\"uint256\"},{\"name\":\"birthTime\",\"type\":\"uint256\"},{\"name\":\"matronId\",\"type\":\"uint256\"},{\"name\":\"sireId\",\"type\":\"uint256\"},{\"name\":\"generation\",\"type\":\"uint256\"},{\"name\":\"genes\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sireId\",\"type\":\"uint256\"},{\"name\":\"_matronId\",\"type\":\"uint256\"}],\"name\":\"bidOnSiringAuction\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"gen0CreatedCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"geneScience\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_matronId\",\"type\":\"uint256\"},{\"name\":\"_sireId\",\"type\":\"uint256\"}],\"name\":\"breedWithAuto\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"matronId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sireId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"cooldownEndBlock\",\"type\":\"uint256\"}],\"name\":\"Pregnant\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"kittyId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"matronId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sireId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"genes\",\"type\":\"uint256\"}],\"name\":\"Birth\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newContract\",\"type\":\"address\"}],\"name\":\"ContractUpgrade\",\"type\":\"event\"}]",
}

// MainABI is the input ABI used to generate the binding from.
// Deprecated: Use MainMetaData.ABI instead.
var MainABI = MainMetaData.ABI

// Main is an auto generated Go binding around an Ethereum contract.
type Main struct {
	MainCaller     // Read-only binding to the contract
	MainTransactor // Write-only binding to the contract
	MainFilterer   // Log filterer for contract events
}

// MainCaller is an auto generated read-only Go binding around an Ethereum contract.
type MainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MainSession struct {
	Contract     *Main             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MainCallerSession struct {
	Contract *MainCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MainTransactorSession struct {
	Contract     *MainTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MainRaw is an auto generated low-level Go binding around an Ethereum contract.
type MainRaw struct {
	Contract *Main // Generic contract binding to access the raw methods on
}

// MainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MainCallerRaw struct {
	Contract *MainCaller // Generic read-only contract binding to access the raw methods on
}

// MainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MainTransactorRaw struct {
	Contract *MainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMain creates a new instance of Main, bound to a specific deployed contract.
func NewMain(address common.Address, backend bind.ContractBackend) (*Main, error) {
	contract, err := bindMain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Main{MainCaller: MainCaller{contract: contract}, MainTransactor: MainTransactor{contract: contract}, MainFilterer: MainFilterer{contract: contract}}, nil
}

// NewMainCaller creates a new read-only instance of Main, bound to a specific deployed contract.
func NewMainCaller(address common.Address, caller bind.ContractCaller) (*MainCaller, error) {
	contract, err := bindMain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MainCaller{contract: contract}, nil
}

// NewMainTransactor creates a new write-only instance of Main, bound to a specific deployed contract.
func NewMainTransactor(address common.Address, transactor bind.ContractTransactor) (*MainTransactor, error) {
	contract, err := bindMain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MainTransactor{contract: contract}, nil
}

// NewMainFilterer creates a new log filterer instance of Main, bound to a specific deployed contract.
func NewMainFilterer(address common.Address, filterer bind.ContractFilterer) (*MainFilterer, error) {
	contract, err := bindMain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MainFilterer{contract: contract}, nil
}

// bindMain binds a generic wrapper to an already deployed contract.
func bindMain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Main *MainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Main.Contract.MainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Main *MainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.Contract.MainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Main *MainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Main.Contract.MainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Main *MainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Main.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Main *MainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Main *MainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Main.Contract.contract.Transact(opts, method, params...)
}

// GEN0AUCTIONDURATION is a free data retrieval call binding the contract method 0x19c2f201.
//
// Solidity: function GEN0_AUCTION_DURATION() view returns(uint256)
func (_Main *MainCaller) GEN0AUCTIONDURATION(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "GEN0_AUCTION_DURATION")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GEN0AUCTIONDURATION is a free data retrieval call binding the contract method 0x19c2f201.
//
// Solidity: function GEN0_AUCTION_DURATION() view returns(uint256)
func (_Main *MainSession) GEN0AUCTIONDURATION() (*big.Int, error) {
	return _Main.Contract.GEN0AUCTIONDURATION(&_Main.CallOpts)
}

// GEN0AUCTIONDURATION is a free data retrieval call binding the contract method 0x19c2f201.
//
// Solidity: function GEN0_AUCTION_DURATION() view returns(uint256)
func (_Main *MainCallerSession) GEN0AUCTIONDURATION() (*big.Int, error) {
	return _Main.Contract.GEN0AUCTIONDURATION(&_Main.CallOpts)
}

// GEN0CREATIONLIMIT is a free data retrieval call binding the contract method 0x680eba27.
//
// Solidity: function GEN0_CREATION_LIMIT() view returns(uint256)
func (_Main *MainCaller) GEN0CREATIONLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "GEN0_CREATION_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GEN0CREATIONLIMIT is a free data retrieval call binding the contract method 0x680eba27.
//
// Solidity: function GEN0_CREATION_LIMIT() view returns(uint256)
func (_Main *MainSession) GEN0CREATIONLIMIT() (*big.Int, error) {
	return _Main.Contract.GEN0CREATIONLIMIT(&_Main.CallOpts)
}

// GEN0CREATIONLIMIT is a free data retrieval call binding the contract method 0x680eba27.
//
// Solidity: function GEN0_CREATION_LIMIT() view returns(uint256)
func (_Main *MainCallerSession) GEN0CREATIONLIMIT() (*big.Int, error) {
	return _Main.Contract.GEN0CREATIONLIMIT(&_Main.CallOpts)
}

// GEN0STARTINGPRICE is a free data retrieval call binding the contract method 0x0e583df0.
//
// Solidity: function GEN0_STARTING_PRICE() view returns(uint256)
func (_Main *MainCaller) GEN0STARTINGPRICE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "GEN0_STARTING_PRICE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GEN0STARTINGPRICE is a free data retrieval call binding the contract method 0x0e583df0.
//
// Solidity: function GEN0_STARTING_PRICE() view returns(uint256)
func (_Main *MainSession) GEN0STARTINGPRICE() (*big.Int, error) {
	return _Main.Contract.GEN0STARTINGPRICE(&_Main.CallOpts)
}

// GEN0STARTINGPRICE is a free data retrieval call binding the contract method 0x0e583df0.
//
// Solidity: function GEN0_STARTING_PRICE() view returns(uint256)
func (_Main *MainCallerSession) GEN0STARTINGPRICE() (*big.Int, error) {
	return _Main.Contract.GEN0STARTINGPRICE(&_Main.CallOpts)
}

// PROMOCREATIONLIMIT is a free data retrieval call binding the contract method 0xdefb9584.
//
// Solidity: function PROMO_CREATION_LIMIT() view returns(uint256)
func (_Main *MainCaller) PROMOCREATIONLIMIT(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "PROMO_CREATION_LIMIT")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PROMOCREATIONLIMIT is a free data retrieval call binding the contract method 0xdefb9584.
//
// Solidity: function PROMO_CREATION_LIMIT() view returns(uint256)
func (_Main *MainSession) PROMOCREATIONLIMIT() (*big.Int, error) {
	return _Main.Contract.PROMOCREATIONLIMIT(&_Main.CallOpts)
}

// PROMOCREATIONLIMIT is a free data retrieval call binding the contract method 0xdefb9584.
//
// Solidity: function PROMO_CREATION_LIMIT() view returns(uint256)
func (_Main *MainCallerSession) PROMOCREATIONLIMIT() (*big.Int, error) {
	return _Main.Contract.PROMOCREATIONLIMIT(&_Main.CallOpts)
}

// AutoBirthFee is a free data retrieval call binding the contract method 0xb0c35c05.
//
// Solidity: function autoBirthFee() view returns(uint256)
func (_Main *MainCaller) AutoBirthFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "autoBirthFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AutoBirthFee is a free data retrieval call binding the contract method 0xb0c35c05.
//
// Solidity: function autoBirthFee() view returns(uint256)
func (_Main *MainSession) AutoBirthFee() (*big.Int, error) {
	return _Main.Contract.AutoBirthFee(&_Main.CallOpts)
}

// AutoBirthFee is a free data retrieval call binding the contract method 0xb0c35c05.
//
// Solidity: function autoBirthFee() view returns(uint256)
func (_Main *MainCallerSession) AutoBirthFee() (*big.Int, error) {
	return _Main.Contract.AutoBirthFee(&_Main.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 count)
func (_Main *MainCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 count)
func (_Main *MainSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Main.Contract.BalanceOf(&_Main.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) view returns(uint256 count)
func (_Main *MainCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _Main.Contract.BalanceOf(&_Main.CallOpts, _owner)
}

// CanBreedWith is a free data retrieval call binding the contract method 0x46d22c70.
//
// Solidity: function canBreedWith(uint256 _matronId, uint256 _sireId) view returns(bool)
func (_Main *MainCaller) CanBreedWith(opts *bind.CallOpts, _matronId *big.Int, _sireId *big.Int) (bool, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "canBreedWith", _matronId, _sireId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanBreedWith is a free data retrieval call binding the contract method 0x46d22c70.
//
// Solidity: function canBreedWith(uint256 _matronId, uint256 _sireId) view returns(bool)
func (_Main *MainSession) CanBreedWith(_matronId *big.Int, _sireId *big.Int) (bool, error) {
	return _Main.Contract.CanBreedWith(&_Main.CallOpts, _matronId, _sireId)
}

// CanBreedWith is a free data retrieval call binding the contract method 0x46d22c70.
//
// Solidity: function canBreedWith(uint256 _matronId, uint256 _sireId) view returns(bool)
func (_Main *MainCallerSession) CanBreedWith(_matronId *big.Int, _sireId *big.Int) (bool, error) {
	return _Main.Contract.CanBreedWith(&_Main.CallOpts, _matronId, _sireId)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() view returns(address)
func (_Main *MainCaller) CeoAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "ceoAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() view returns(address)
func (_Main *MainSession) CeoAddress() (common.Address, error) {
	return _Main.Contract.CeoAddress(&_Main.CallOpts)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() view returns(address)
func (_Main *MainCallerSession) CeoAddress() (common.Address, error) {
	return _Main.Contract.CeoAddress(&_Main.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() view returns(address)
func (_Main *MainCaller) CfoAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "cfoAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() view returns(address)
func (_Main *MainSession) CfoAddress() (common.Address, error) {
	return _Main.Contract.CfoAddress(&_Main.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() view returns(address)
func (_Main *MainCallerSession) CfoAddress() (common.Address, error) {
	return _Main.Contract.CfoAddress(&_Main.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() view returns(address)
func (_Main *MainCaller) CooAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "cooAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() view returns(address)
func (_Main *MainSession) CooAddress() (common.Address, error) {
	return _Main.Contract.CooAddress(&_Main.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() view returns(address)
func (_Main *MainCallerSession) CooAddress() (common.Address, error) {
	return _Main.Contract.CooAddress(&_Main.CallOpts)
}

// Cooldowns is a free data retrieval call binding the contract method 0x9d6fac6f.
//
// Solidity: function cooldowns(uint256 ) view returns(uint32)
func (_Main *MainCaller) Cooldowns(opts *bind.CallOpts, arg0 *big.Int) (uint32, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "cooldowns", arg0)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// Cooldowns is a free data retrieval call binding the contract method 0x9d6fac6f.
//
// Solidity: function cooldowns(uint256 ) view returns(uint32)
func (_Main *MainSession) Cooldowns(arg0 *big.Int) (uint32, error) {
	return _Main.Contract.Cooldowns(&_Main.CallOpts, arg0)
}

// Cooldowns is a free data retrieval call binding the contract method 0x9d6fac6f.
//
// Solidity: function cooldowns(uint256 ) view returns(uint32)
func (_Main *MainCallerSession) Cooldowns(arg0 *big.Int) (uint32, error) {
	return _Main.Contract.Cooldowns(&_Main.CallOpts, arg0)
}

// Erc721Metadata is a free data retrieval call binding the contract method 0xbc4006f5.
//
// Solidity: function erc721Metadata() view returns(address)
func (_Main *MainCaller) Erc721Metadata(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "erc721Metadata")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Erc721Metadata is a free data retrieval call binding the contract method 0xbc4006f5.
//
// Solidity: function erc721Metadata() view returns(address)
func (_Main *MainSession) Erc721Metadata() (common.Address, error) {
	return _Main.Contract.Erc721Metadata(&_Main.CallOpts)
}

// Erc721Metadata is a free data retrieval call binding the contract method 0xbc4006f5.
//
// Solidity: function erc721Metadata() view returns(address)
func (_Main *MainCallerSession) Erc721Metadata() (common.Address, error) {
	return _Main.Contract.Erc721Metadata(&_Main.CallOpts)
}

// Gen0CreatedCount is a free data retrieval call binding the contract method 0xf1ca9410.
//
// Solidity: function gen0CreatedCount() view returns(uint256)
func (_Main *MainCaller) Gen0CreatedCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "gen0CreatedCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Gen0CreatedCount is a free data retrieval call binding the contract method 0xf1ca9410.
//
// Solidity: function gen0CreatedCount() view returns(uint256)
func (_Main *MainSession) Gen0CreatedCount() (*big.Int, error) {
	return _Main.Contract.Gen0CreatedCount(&_Main.CallOpts)
}

// Gen0CreatedCount is a free data retrieval call binding the contract method 0xf1ca9410.
//
// Solidity: function gen0CreatedCount() view returns(uint256)
func (_Main *MainCallerSession) Gen0CreatedCount() (*big.Int, error) {
	return _Main.Contract.Gen0CreatedCount(&_Main.CallOpts)
}

// GeneScience is a free data retrieval call binding the contract method 0xf2b47d52.
//
// Solidity: function geneScience() view returns(address)
func (_Main *MainCaller) GeneScience(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "geneScience")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GeneScience is a free data retrieval call binding the contract method 0xf2b47d52.
//
// Solidity: function geneScience() view returns(address)
func (_Main *MainSession) GeneScience() (common.Address, error) {
	return _Main.Contract.GeneScience(&_Main.CallOpts)
}

// GeneScience is a free data retrieval call binding the contract method 0xf2b47d52.
//
// Solidity: function geneScience() view returns(address)
func (_Main *MainCallerSession) GeneScience() (common.Address, error) {
	return _Main.Contract.GeneScience(&_Main.CallOpts)
}

// GetKitty is a free data retrieval call binding the contract method 0xe98b7f4d.
//
// Solidity: function getKitty(uint256 _id) view returns(bool isGestating, bool isReady, uint256 cooldownIndex, uint256 nextActionAt, uint256 siringWithId, uint256 birthTime, uint256 matronId, uint256 sireId, uint256 generation, uint256 genes)
func (_Main *MainCaller) GetKitty(opts *bind.CallOpts, _id *big.Int) (struct {
	IsGestating   bool
	IsReady       bool
	CooldownIndex *big.Int
	NextActionAt  *big.Int
	SiringWithId  *big.Int
	BirthTime     *big.Int
	MatronId      *big.Int
	SireId        *big.Int
	Generation    *big.Int
	Genes         *big.Int
}, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "getKitty", _id)

	outstruct := new(struct {
		IsGestating   bool
		IsReady       bool
		CooldownIndex *big.Int
		NextActionAt  *big.Int
		SiringWithId  *big.Int
		BirthTime     *big.Int
		MatronId      *big.Int
		SireId        *big.Int
		Generation    *big.Int
		Genes         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.IsGestating = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.IsReady = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.CooldownIndex = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.NextActionAt = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.SiringWithId = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.BirthTime = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.MatronId = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.SireId = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.Generation = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.Genes = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetKitty is a free data retrieval call binding the contract method 0xe98b7f4d.
//
// Solidity: function getKitty(uint256 _id) view returns(bool isGestating, bool isReady, uint256 cooldownIndex, uint256 nextActionAt, uint256 siringWithId, uint256 birthTime, uint256 matronId, uint256 sireId, uint256 generation, uint256 genes)
func (_Main *MainSession) GetKitty(_id *big.Int) (struct {
	IsGestating   bool
	IsReady       bool
	CooldownIndex *big.Int
	NextActionAt  *big.Int
	SiringWithId  *big.Int
	BirthTime     *big.Int
	MatronId      *big.Int
	SireId        *big.Int
	Generation    *big.Int
	Genes         *big.Int
}, error) {
	return _Main.Contract.GetKitty(&_Main.CallOpts, _id)
}

// GetKitty is a free data retrieval call binding the contract method 0xe98b7f4d.
//
// Solidity: function getKitty(uint256 _id) view returns(bool isGestating, bool isReady, uint256 cooldownIndex, uint256 nextActionAt, uint256 siringWithId, uint256 birthTime, uint256 matronId, uint256 sireId, uint256 generation, uint256 genes)
func (_Main *MainCallerSession) GetKitty(_id *big.Int) (struct {
	IsGestating   bool
	IsReady       bool
	CooldownIndex *big.Int
	NextActionAt  *big.Int
	SiringWithId  *big.Int
	BirthTime     *big.Int
	MatronId      *big.Int
	SireId        *big.Int
	Generation    *big.Int
	Genes         *big.Int
}, error) {
	return _Main.Contract.GetKitty(&_Main.CallOpts, _id)
}

// IsPregnant is a free data retrieval call binding the contract method 0x1940a936.
//
// Solidity: function isPregnant(uint256 _kittyId) view returns(bool)
func (_Main *MainCaller) IsPregnant(opts *bind.CallOpts, _kittyId *big.Int) (bool, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "isPregnant", _kittyId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPregnant is a free data retrieval call binding the contract method 0x1940a936.
//
// Solidity: function isPregnant(uint256 _kittyId) view returns(bool)
func (_Main *MainSession) IsPregnant(_kittyId *big.Int) (bool, error) {
	return _Main.Contract.IsPregnant(&_Main.CallOpts, _kittyId)
}

// IsPregnant is a free data retrieval call binding the contract method 0x1940a936.
//
// Solidity: function isPregnant(uint256 _kittyId) view returns(bool)
func (_Main *MainCallerSession) IsPregnant(_kittyId *big.Int) (bool, error) {
	return _Main.Contract.IsPregnant(&_Main.CallOpts, _kittyId)
}

// IsReadyToBreed is a free data retrieval call binding the contract method 0xd3e6f49f.
//
// Solidity: function isReadyToBreed(uint256 _kittyId) view returns(bool)
func (_Main *MainCaller) IsReadyToBreed(opts *bind.CallOpts, _kittyId *big.Int) (bool, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "isReadyToBreed", _kittyId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsReadyToBreed is a free data retrieval call binding the contract method 0xd3e6f49f.
//
// Solidity: function isReadyToBreed(uint256 _kittyId) view returns(bool)
func (_Main *MainSession) IsReadyToBreed(_kittyId *big.Int) (bool, error) {
	return _Main.Contract.IsReadyToBreed(&_Main.CallOpts, _kittyId)
}

// IsReadyToBreed is a free data retrieval call binding the contract method 0xd3e6f49f.
//
// Solidity: function isReadyToBreed(uint256 _kittyId) view returns(bool)
func (_Main *MainCallerSession) IsReadyToBreed(_kittyId *big.Int) (bool, error) {
	return _Main.Contract.IsReadyToBreed(&_Main.CallOpts, _kittyId)
}

// KittyIndexToApproved is a free data retrieval call binding the contract method 0x481af3d3.
//
// Solidity: function kittyIndexToApproved(uint256 ) view returns(address)
func (_Main *MainCaller) KittyIndexToApproved(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "kittyIndexToApproved", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KittyIndexToApproved is a free data retrieval call binding the contract method 0x481af3d3.
//
// Solidity: function kittyIndexToApproved(uint256 ) view returns(address)
func (_Main *MainSession) KittyIndexToApproved(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.KittyIndexToApproved(&_Main.CallOpts, arg0)
}

// KittyIndexToApproved is a free data retrieval call binding the contract method 0x481af3d3.
//
// Solidity: function kittyIndexToApproved(uint256 ) view returns(address)
func (_Main *MainCallerSession) KittyIndexToApproved(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.KittyIndexToApproved(&_Main.CallOpts, arg0)
}

// KittyIndexToOwner is a free data retrieval call binding the contract method 0xa45f4bfc.
//
// Solidity: function kittyIndexToOwner(uint256 ) view returns(address)
func (_Main *MainCaller) KittyIndexToOwner(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "kittyIndexToOwner", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// KittyIndexToOwner is a free data retrieval call binding the contract method 0xa45f4bfc.
//
// Solidity: function kittyIndexToOwner(uint256 ) view returns(address)
func (_Main *MainSession) KittyIndexToOwner(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.KittyIndexToOwner(&_Main.CallOpts, arg0)
}

// KittyIndexToOwner is a free data retrieval call binding the contract method 0xa45f4bfc.
//
// Solidity: function kittyIndexToOwner(uint256 ) view returns(address)
func (_Main *MainCallerSession) KittyIndexToOwner(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.KittyIndexToOwner(&_Main.CallOpts, arg0)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Main *MainCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Main *MainSession) Name() (string, error) {
	return _Main.Contract.Name(&_Main.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Main *MainCallerSession) Name() (string, error) {
	return _Main.Contract.Name(&_Main.CallOpts)
}

// NewContractAddress is a free data retrieval call binding the contract method 0x6af04a57.
//
// Solidity: function newContractAddress() view returns(address)
func (_Main *MainCaller) NewContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "newContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NewContractAddress is a free data retrieval call binding the contract method 0x6af04a57.
//
// Solidity: function newContractAddress() view returns(address)
func (_Main *MainSession) NewContractAddress() (common.Address, error) {
	return _Main.Contract.NewContractAddress(&_Main.CallOpts)
}

// NewContractAddress is a free data retrieval call binding the contract method 0x6af04a57.
//
// Solidity: function newContractAddress() view returns(address)
func (_Main *MainCallerSession) NewContractAddress() (common.Address, error) {
	return _Main.Contract.NewContractAddress(&_Main.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address owner)
func (_Main *MainCaller) OwnerOf(opts *bind.CallOpts, _tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "ownerOf", _tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address owner)
func (_Main *MainSession) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _Main.Contract.OwnerOf(&_Main.CallOpts, _tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 _tokenId) view returns(address owner)
func (_Main *MainCallerSession) OwnerOf(_tokenId *big.Int) (common.Address, error) {
	return _Main.Contract.OwnerOf(&_Main.CallOpts, _tokenId)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Main *MainCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Main *MainSession) Paused() (bool, error) {
	return _Main.Contract.Paused(&_Main.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool)
func (_Main *MainCallerSession) Paused() (bool, error) {
	return _Main.Contract.Paused(&_Main.CallOpts)
}

// PregnantKitties is a free data retrieval call binding the contract method 0x183a7947.
//
// Solidity: function pregnantKitties() view returns(uint256)
func (_Main *MainCaller) PregnantKitties(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "pregnantKitties")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PregnantKitties is a free data retrieval call binding the contract method 0x183a7947.
//
// Solidity: function pregnantKitties() view returns(uint256)
func (_Main *MainSession) PregnantKitties() (*big.Int, error) {
	return _Main.Contract.PregnantKitties(&_Main.CallOpts)
}

// PregnantKitties is a free data retrieval call binding the contract method 0x183a7947.
//
// Solidity: function pregnantKitties() view returns(uint256)
func (_Main *MainCallerSession) PregnantKitties() (*big.Int, error) {
	return _Main.Contract.PregnantKitties(&_Main.CallOpts)
}

// PromoCreatedCount is a free data retrieval call binding the contract method 0x05e45546.
//
// Solidity: function promoCreatedCount() view returns(uint256)
func (_Main *MainCaller) PromoCreatedCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "promoCreatedCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PromoCreatedCount is a free data retrieval call binding the contract method 0x05e45546.
//
// Solidity: function promoCreatedCount() view returns(uint256)
func (_Main *MainSession) PromoCreatedCount() (*big.Int, error) {
	return _Main.Contract.PromoCreatedCount(&_Main.CallOpts)
}

// PromoCreatedCount is a free data retrieval call binding the contract method 0x05e45546.
//
// Solidity: function promoCreatedCount() view returns(uint256)
func (_Main *MainCallerSession) PromoCreatedCount() (*big.Int, error) {
	return _Main.Contract.PromoCreatedCount(&_Main.CallOpts)
}

// SaleAuction is a free data retrieval call binding the contract method 0xe6cbe351.
//
// Solidity: function saleAuction() view returns(address)
func (_Main *MainCaller) SaleAuction(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "saleAuction")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SaleAuction is a free data retrieval call binding the contract method 0xe6cbe351.
//
// Solidity: function saleAuction() view returns(address)
func (_Main *MainSession) SaleAuction() (common.Address, error) {
	return _Main.Contract.SaleAuction(&_Main.CallOpts)
}

// SaleAuction is a free data retrieval call binding the contract method 0xe6cbe351.
//
// Solidity: function saleAuction() view returns(address)
func (_Main *MainCallerSession) SaleAuction() (common.Address, error) {
	return _Main.Contract.SaleAuction(&_Main.CallOpts)
}

// SecondsPerBlock is a free data retrieval call binding the contract method 0x7a7d4937.
//
// Solidity: function secondsPerBlock() view returns(uint256)
func (_Main *MainCaller) SecondsPerBlock(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "secondsPerBlock")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SecondsPerBlock is a free data retrieval call binding the contract method 0x7a7d4937.
//
// Solidity: function secondsPerBlock() view returns(uint256)
func (_Main *MainSession) SecondsPerBlock() (*big.Int, error) {
	return _Main.Contract.SecondsPerBlock(&_Main.CallOpts)
}

// SecondsPerBlock is a free data retrieval call binding the contract method 0x7a7d4937.
//
// Solidity: function secondsPerBlock() view returns(uint256)
func (_Main *MainCallerSession) SecondsPerBlock() (*big.Int, error) {
	return _Main.Contract.SecondsPerBlock(&_Main.CallOpts)
}

// SireAllowedToAddress is a free data retrieval call binding the contract method 0x46116e6f.
//
// Solidity: function sireAllowedToAddress(uint256 ) view returns(address)
func (_Main *MainCaller) SireAllowedToAddress(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "sireAllowedToAddress", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SireAllowedToAddress is a free data retrieval call binding the contract method 0x46116e6f.
//
// Solidity: function sireAllowedToAddress(uint256 ) view returns(address)
func (_Main *MainSession) SireAllowedToAddress(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.SireAllowedToAddress(&_Main.CallOpts, arg0)
}

// SireAllowedToAddress is a free data retrieval call binding the contract method 0x46116e6f.
//
// Solidity: function sireAllowedToAddress(uint256 ) view returns(address)
func (_Main *MainCallerSession) SireAllowedToAddress(arg0 *big.Int) (common.Address, error) {
	return _Main.Contract.SireAllowedToAddress(&_Main.CallOpts, arg0)
}

// SiringAuction is a free data retrieval call binding the contract method 0x21717ebf.
//
// Solidity: function siringAuction() view returns(address)
func (_Main *MainCaller) SiringAuction(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "siringAuction")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// SiringAuction is a free data retrieval call binding the contract method 0x21717ebf.
//
// Solidity: function siringAuction() view returns(address)
func (_Main *MainSession) SiringAuction() (common.Address, error) {
	return _Main.Contract.SiringAuction(&_Main.CallOpts)
}

// SiringAuction is a free data retrieval call binding the contract method 0x21717ebf.
//
// Solidity: function siringAuction() view returns(address)
func (_Main *MainCallerSession) SiringAuction() (common.Address, error) {
	return _Main.Contract.SiringAuction(&_Main.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceID) view returns(bool)
func (_Main *MainCaller) SupportsInterface(opts *bind.CallOpts, _interfaceID [4]byte) (bool, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "supportsInterface", _interfaceID)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceID) view returns(bool)
func (_Main *MainSession) SupportsInterface(_interfaceID [4]byte) (bool, error) {
	return _Main.Contract.SupportsInterface(&_Main.CallOpts, _interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 _interfaceID) view returns(bool)
func (_Main *MainCallerSession) SupportsInterface(_interfaceID [4]byte) (bool, error) {
	return _Main.Contract.SupportsInterface(&_Main.CallOpts, _interfaceID)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Main *MainCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Main *MainSession) Symbol() (string, error) {
	return _Main.Contract.Symbol(&_Main.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Main *MainCallerSession) Symbol() (string, error) {
	return _Main.Contract.Symbol(&_Main.CallOpts)
}

// TokenMetadata is a free data retrieval call binding the contract method 0x0560ff44.
//
// Solidity: function tokenMetadata(uint256 _tokenId, string _preferredTransport) view returns(string infoUrl)
func (_Main *MainCaller) TokenMetadata(opts *bind.CallOpts, _tokenId *big.Int, _preferredTransport string) (string, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "tokenMetadata", _tokenId, _preferredTransport)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenMetadata is a free data retrieval call binding the contract method 0x0560ff44.
//
// Solidity: function tokenMetadata(uint256 _tokenId, string _preferredTransport) view returns(string infoUrl)
func (_Main *MainSession) TokenMetadata(_tokenId *big.Int, _preferredTransport string) (string, error) {
	return _Main.Contract.TokenMetadata(&_Main.CallOpts, _tokenId, _preferredTransport)
}

// TokenMetadata is a free data retrieval call binding the contract method 0x0560ff44.
//
// Solidity: function tokenMetadata(uint256 _tokenId, string _preferredTransport) view returns(string infoUrl)
func (_Main *MainCallerSession) TokenMetadata(_tokenId *big.Int, _preferredTransport string) (string, error) {
	return _Main.Contract.TokenMetadata(&_Main.CallOpts, _tokenId, _preferredTransport)
}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address _owner) view returns(uint256[] ownerTokens)
func (_Main *MainCaller) TokensOfOwner(opts *bind.CallOpts, _owner common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "tokensOfOwner", _owner)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address _owner) view returns(uint256[] ownerTokens)
func (_Main *MainSession) TokensOfOwner(_owner common.Address) ([]*big.Int, error) {
	return _Main.Contract.TokensOfOwner(&_Main.CallOpts, _owner)
}

// TokensOfOwner is a free data retrieval call binding the contract method 0x8462151c.
//
// Solidity: function tokensOfOwner(address _owner) view returns(uint256[] ownerTokens)
func (_Main *MainCallerSession) TokensOfOwner(_owner common.Address) ([]*big.Int, error) {
	return _Main.Contract.TokensOfOwner(&_Main.CallOpts, _owner)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Main *MainCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Main.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Main *MainSession) TotalSupply() (*big.Int, error) {
	return _Main.Contract.TotalSupply(&_Main.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Main *MainCallerSession) TotalSupply() (*big.Int, error) {
	return _Main.Contract.TotalSupply(&_Main.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _to, uint256 _tokenId) returns()
func (_Main *MainTransactor) Approve(opts *bind.TransactOpts, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "approve", _to, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _to, uint256 _tokenId) returns()
func (_Main *MainSession) Approve(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.Approve(&_Main.TransactOpts, _to, _tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address _to, uint256 _tokenId) returns()
func (_Main *MainTransactorSession) Approve(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.Approve(&_Main.TransactOpts, _to, _tokenId)
}

// ApproveSiring is a paid mutator transaction binding the contract method 0x4dfff04f.
//
// Solidity: function approveSiring(address _addr, uint256 _sireId) returns()
func (_Main *MainTransactor) ApproveSiring(opts *bind.TransactOpts, _addr common.Address, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "approveSiring", _addr, _sireId)
}

// ApproveSiring is a paid mutator transaction binding the contract method 0x4dfff04f.
//
// Solidity: function approveSiring(address _addr, uint256 _sireId) returns()
func (_Main *MainSession) ApproveSiring(_addr common.Address, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.ApproveSiring(&_Main.TransactOpts, _addr, _sireId)
}

// ApproveSiring is a paid mutator transaction binding the contract method 0x4dfff04f.
//
// Solidity: function approveSiring(address _addr, uint256 _sireId) returns()
func (_Main *MainTransactorSession) ApproveSiring(_addr common.Address, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.ApproveSiring(&_Main.TransactOpts, _addr, _sireId)
}

// BidOnSiringAuction is a paid mutator transaction binding the contract method 0xed60ade6.
//
// Solidity: function bidOnSiringAuction(uint256 _sireId, uint256 _matronId) payable returns()
func (_Main *MainTransactor) BidOnSiringAuction(opts *bind.TransactOpts, _sireId *big.Int, _matronId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "bidOnSiringAuction", _sireId, _matronId)
}

// BidOnSiringAuction is a paid mutator transaction binding the contract method 0xed60ade6.
//
// Solidity: function bidOnSiringAuction(uint256 _sireId, uint256 _matronId) payable returns()
func (_Main *MainSession) BidOnSiringAuction(_sireId *big.Int, _matronId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.BidOnSiringAuction(&_Main.TransactOpts, _sireId, _matronId)
}

// BidOnSiringAuction is a paid mutator transaction binding the contract method 0xed60ade6.
//
// Solidity: function bidOnSiringAuction(uint256 _sireId, uint256 _matronId) payable returns()
func (_Main *MainTransactorSession) BidOnSiringAuction(_sireId *big.Int, _matronId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.BidOnSiringAuction(&_Main.TransactOpts, _sireId, _matronId)
}

// BreedWithAuto is a paid mutator transaction binding the contract method 0xf7d8c883.
//
// Solidity: function breedWithAuto(uint256 _matronId, uint256 _sireId) payable returns()
func (_Main *MainTransactor) BreedWithAuto(opts *bind.TransactOpts, _matronId *big.Int, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "breedWithAuto", _matronId, _sireId)
}

// BreedWithAuto is a paid mutator transaction binding the contract method 0xf7d8c883.
//
// Solidity: function breedWithAuto(uint256 _matronId, uint256 _sireId) payable returns()
func (_Main *MainSession) BreedWithAuto(_matronId *big.Int, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.BreedWithAuto(&_Main.TransactOpts, _matronId, _sireId)
}

// BreedWithAuto is a paid mutator transaction binding the contract method 0xf7d8c883.
//
// Solidity: function breedWithAuto(uint256 _matronId, uint256 _sireId) payable returns()
func (_Main *MainTransactorSession) BreedWithAuto(_matronId *big.Int, _sireId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.BreedWithAuto(&_Main.TransactOpts, _matronId, _sireId)
}

// CreateGen0Auction is a paid mutator transaction binding the contract method 0xc3bea9af.
//
// Solidity: function createGen0Auction(uint256 _genes) returns()
func (_Main *MainTransactor) CreateGen0Auction(opts *bind.TransactOpts, _genes *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "createGen0Auction", _genes)
}

// CreateGen0Auction is a paid mutator transaction binding the contract method 0xc3bea9af.
//
// Solidity: function createGen0Auction(uint256 _genes) returns()
func (_Main *MainSession) CreateGen0Auction(_genes *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateGen0Auction(&_Main.TransactOpts, _genes)
}

// CreateGen0Auction is a paid mutator transaction binding the contract method 0xc3bea9af.
//
// Solidity: function createGen0Auction(uint256 _genes) returns()
func (_Main *MainTransactorSession) CreateGen0Auction(_genes *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateGen0Auction(&_Main.TransactOpts, _genes)
}

// CreatePromoKitty is a paid mutator transaction binding the contract method 0x56129134.
//
// Solidity: function createPromoKitty(uint256 _genes, address _owner) returns()
func (_Main *MainTransactor) CreatePromoKitty(opts *bind.TransactOpts, _genes *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "createPromoKitty", _genes, _owner)
}

// CreatePromoKitty is a paid mutator transaction binding the contract method 0x56129134.
//
// Solidity: function createPromoKitty(uint256 _genes, address _owner) returns()
func (_Main *MainSession) CreatePromoKitty(_genes *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Main.Contract.CreatePromoKitty(&_Main.TransactOpts, _genes, _owner)
}

// CreatePromoKitty is a paid mutator transaction binding the contract method 0x56129134.
//
// Solidity: function createPromoKitty(uint256 _genes, address _owner) returns()
func (_Main *MainTransactorSession) CreatePromoKitty(_genes *big.Int, _owner common.Address) (*types.Transaction, error) {
	return _Main.Contract.CreatePromoKitty(&_Main.TransactOpts, _genes, _owner)
}

// CreateSaleAuction is a paid mutator transaction binding the contract method 0x3d7d3f5a.
//
// Solidity: function createSaleAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainTransactor) CreateSaleAuction(opts *bind.TransactOpts, _kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "createSaleAuction", _kittyId, _startingPrice, _endingPrice, _duration)
}

// CreateSaleAuction is a paid mutator transaction binding the contract method 0x3d7d3f5a.
//
// Solidity: function createSaleAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainSession) CreateSaleAuction(_kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateSaleAuction(&_Main.TransactOpts, _kittyId, _startingPrice, _endingPrice, _duration)
}

// CreateSaleAuction is a paid mutator transaction binding the contract method 0x3d7d3f5a.
//
// Solidity: function createSaleAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainTransactorSession) CreateSaleAuction(_kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateSaleAuction(&_Main.TransactOpts, _kittyId, _startingPrice, _endingPrice, _duration)
}

// CreateSiringAuction is a paid mutator transaction binding the contract method 0x4ad8c938.
//
// Solidity: function createSiringAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainTransactor) CreateSiringAuction(opts *bind.TransactOpts, _kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "createSiringAuction", _kittyId, _startingPrice, _endingPrice, _duration)
}

// CreateSiringAuction is a paid mutator transaction binding the contract method 0x4ad8c938.
//
// Solidity: function createSiringAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainSession) CreateSiringAuction(_kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateSiringAuction(&_Main.TransactOpts, _kittyId, _startingPrice, _endingPrice, _duration)
}

// CreateSiringAuction is a paid mutator transaction binding the contract method 0x4ad8c938.
//
// Solidity: function createSiringAuction(uint256 _kittyId, uint256 _startingPrice, uint256 _endingPrice, uint256 _duration) returns()
func (_Main *MainTransactorSession) CreateSiringAuction(_kittyId *big.Int, _startingPrice *big.Int, _endingPrice *big.Int, _duration *big.Int) (*types.Transaction, error) {
	return _Main.Contract.CreateSiringAuction(&_Main.TransactOpts, _kittyId, _startingPrice, _endingPrice, _duration)
}

// GiveBirth is a paid mutator transaction binding the contract method 0x88c2a0bf.
//
// Solidity: function giveBirth(uint256 _matronId) returns(uint256)
func (_Main *MainTransactor) GiveBirth(opts *bind.TransactOpts, _matronId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "giveBirth", _matronId)
}

// GiveBirth is a paid mutator transaction binding the contract method 0x88c2a0bf.
//
// Solidity: function giveBirth(uint256 _matronId) returns(uint256)
func (_Main *MainSession) GiveBirth(_matronId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.GiveBirth(&_Main.TransactOpts, _matronId)
}

// GiveBirth is a paid mutator transaction binding the contract method 0x88c2a0bf.
//
// Solidity: function giveBirth(uint256 _matronId) returns(uint256)
func (_Main *MainTransactorSession) GiveBirth(_matronId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.GiveBirth(&_Main.TransactOpts, _matronId)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Main *MainTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Main *MainSession) Pause() (*types.Transaction, error) {
	return _Main.Contract.Pause(&_Main.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_Main *MainTransactorSession) Pause() (*types.Transaction, error) {
	return _Main.Contract.Pause(&_Main.TransactOpts)
}

// SetAutoBirthFee is a paid mutator transaction binding the contract method 0x4b85fd55.
//
// Solidity: function setAutoBirthFee(uint256 val) returns()
func (_Main *MainTransactor) SetAutoBirthFee(opts *bind.TransactOpts, val *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setAutoBirthFee", val)
}

// SetAutoBirthFee is a paid mutator transaction binding the contract method 0x4b85fd55.
//
// Solidity: function setAutoBirthFee(uint256 val) returns()
func (_Main *MainSession) SetAutoBirthFee(val *big.Int) (*types.Transaction, error) {
	return _Main.Contract.SetAutoBirthFee(&_Main.TransactOpts, val)
}

// SetAutoBirthFee is a paid mutator transaction binding the contract method 0x4b85fd55.
//
// Solidity: function setAutoBirthFee(uint256 val) returns()
func (_Main *MainTransactorSession) SetAutoBirthFee(val *big.Int) (*types.Transaction, error) {
	return _Main.Contract.SetAutoBirthFee(&_Main.TransactOpts, val)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(address _newCEO) returns()
func (_Main *MainTransactor) SetCEO(opts *bind.TransactOpts, _newCEO common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setCEO", _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(address _newCEO) returns()
func (_Main *MainSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCEO(&_Main.TransactOpts, _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(address _newCEO) returns()
func (_Main *MainTransactorSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCEO(&_Main.TransactOpts, _newCEO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(address _newCFO) returns()
func (_Main *MainTransactor) SetCFO(opts *bind.TransactOpts, _newCFO common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setCFO", _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(address _newCFO) returns()
func (_Main *MainSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCFO(&_Main.TransactOpts, _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(address _newCFO) returns()
func (_Main *MainTransactorSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCFO(&_Main.TransactOpts, _newCFO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(address _newCOO) returns()
func (_Main *MainTransactor) SetCOO(opts *bind.TransactOpts, _newCOO common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setCOO", _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(address _newCOO) returns()
func (_Main *MainSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCOO(&_Main.TransactOpts, _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(address _newCOO) returns()
func (_Main *MainTransactorSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetCOO(&_Main.TransactOpts, _newCOO)
}

// SetGeneScienceAddress is a paid mutator transaction binding the contract method 0x24e7a38a.
//
// Solidity: function setGeneScienceAddress(address _address) returns()
func (_Main *MainTransactor) SetGeneScienceAddress(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setGeneScienceAddress", _address)
}

// SetGeneScienceAddress is a paid mutator transaction binding the contract method 0x24e7a38a.
//
// Solidity: function setGeneScienceAddress(address _address) returns()
func (_Main *MainSession) SetGeneScienceAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetGeneScienceAddress(&_Main.TransactOpts, _address)
}

// SetGeneScienceAddress is a paid mutator transaction binding the contract method 0x24e7a38a.
//
// Solidity: function setGeneScienceAddress(address _address) returns()
func (_Main *MainTransactorSession) SetGeneScienceAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetGeneScienceAddress(&_Main.TransactOpts, _address)
}

// SetMetadataAddress is a paid mutator transaction binding the contract method 0xe17b25af.
//
// Solidity: function setMetadataAddress(address _contractAddress) returns()
func (_Main *MainTransactor) SetMetadataAddress(opts *bind.TransactOpts, _contractAddress common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setMetadataAddress", _contractAddress)
}

// SetMetadataAddress is a paid mutator transaction binding the contract method 0xe17b25af.
//
// Solidity: function setMetadataAddress(address _contractAddress) returns()
func (_Main *MainSession) SetMetadataAddress(_contractAddress common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetMetadataAddress(&_Main.TransactOpts, _contractAddress)
}

// SetMetadataAddress is a paid mutator transaction binding the contract method 0xe17b25af.
//
// Solidity: function setMetadataAddress(address _contractAddress) returns()
func (_Main *MainTransactorSession) SetMetadataAddress(_contractAddress common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetMetadataAddress(&_Main.TransactOpts, _contractAddress)
}

// SetNewAddress is a paid mutator transaction binding the contract method 0x71587988.
//
// Solidity: function setNewAddress(address _v2Address) returns()
func (_Main *MainTransactor) SetNewAddress(opts *bind.TransactOpts, _v2Address common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setNewAddress", _v2Address)
}

// SetNewAddress is a paid mutator transaction binding the contract method 0x71587988.
//
// Solidity: function setNewAddress(address _v2Address) returns()
func (_Main *MainSession) SetNewAddress(_v2Address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetNewAddress(&_Main.TransactOpts, _v2Address)
}

// SetNewAddress is a paid mutator transaction binding the contract method 0x71587988.
//
// Solidity: function setNewAddress(address _v2Address) returns()
func (_Main *MainTransactorSession) SetNewAddress(_v2Address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetNewAddress(&_Main.TransactOpts, _v2Address)
}

// SetSaleAuctionAddress is a paid mutator transaction binding the contract method 0x6fbde40d.
//
// Solidity: function setSaleAuctionAddress(address _address) returns()
func (_Main *MainTransactor) SetSaleAuctionAddress(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setSaleAuctionAddress", _address)
}

// SetSaleAuctionAddress is a paid mutator transaction binding the contract method 0x6fbde40d.
//
// Solidity: function setSaleAuctionAddress(address _address) returns()
func (_Main *MainSession) SetSaleAuctionAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetSaleAuctionAddress(&_Main.TransactOpts, _address)
}

// SetSaleAuctionAddress is a paid mutator transaction binding the contract method 0x6fbde40d.
//
// Solidity: function setSaleAuctionAddress(address _address) returns()
func (_Main *MainTransactorSession) SetSaleAuctionAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetSaleAuctionAddress(&_Main.TransactOpts, _address)
}

// SetSecondsPerBlock is a paid mutator transaction binding the contract method 0x5663896e.
//
// Solidity: function setSecondsPerBlock(uint256 secs) returns()
func (_Main *MainTransactor) SetSecondsPerBlock(opts *bind.TransactOpts, secs *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setSecondsPerBlock", secs)
}

// SetSecondsPerBlock is a paid mutator transaction binding the contract method 0x5663896e.
//
// Solidity: function setSecondsPerBlock(uint256 secs) returns()
func (_Main *MainSession) SetSecondsPerBlock(secs *big.Int) (*types.Transaction, error) {
	return _Main.Contract.SetSecondsPerBlock(&_Main.TransactOpts, secs)
}

// SetSecondsPerBlock is a paid mutator transaction binding the contract method 0x5663896e.
//
// Solidity: function setSecondsPerBlock(uint256 secs) returns()
func (_Main *MainTransactorSession) SetSecondsPerBlock(secs *big.Int) (*types.Transaction, error) {
	return _Main.Contract.SetSecondsPerBlock(&_Main.TransactOpts, secs)
}

// SetSiringAuctionAddress is a paid mutator transaction binding the contract method 0x14001f4c.
//
// Solidity: function setSiringAuctionAddress(address _address) returns()
func (_Main *MainTransactor) SetSiringAuctionAddress(opts *bind.TransactOpts, _address common.Address) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "setSiringAuctionAddress", _address)
}

// SetSiringAuctionAddress is a paid mutator transaction binding the contract method 0x14001f4c.
//
// Solidity: function setSiringAuctionAddress(address _address) returns()
func (_Main *MainSession) SetSiringAuctionAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetSiringAuctionAddress(&_Main.TransactOpts, _address)
}

// SetSiringAuctionAddress is a paid mutator transaction binding the contract method 0x14001f4c.
//
// Solidity: function setSiringAuctionAddress(address _address) returns()
func (_Main *MainTransactorSession) SetSiringAuctionAddress(_address common.Address) (*types.Transaction, error) {
	return _Main.Contract.SetSiringAuctionAddress(&_Main.TransactOpts, _address)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _tokenId) returns()
func (_Main *MainTransactor) Transfer(opts *bind.TransactOpts, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "transfer", _to, _tokenId)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _tokenId) returns()
func (_Main *MainSession) Transfer(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.Transfer(&_Main.TransactOpts, _to, _tokenId)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address _to, uint256 _tokenId) returns()
func (_Main *MainTransactorSession) Transfer(_to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.Transfer(&_Main.TransactOpts, _to, _tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _tokenId) returns()
func (_Main *MainTransactor) TransferFrom(opts *bind.TransactOpts, _from common.Address, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "transferFrom", _from, _to, _tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _tokenId) returns()
func (_Main *MainSession) TransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.TransferFrom(&_Main.TransactOpts, _from, _to, _tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address _from, address _to, uint256 _tokenId) returns()
func (_Main *MainTransactorSession) TransferFrom(_from common.Address, _to common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Main.Contract.TransferFrom(&_Main.TransactOpts, _from, _to, _tokenId)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Main *MainTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Main *MainSession) Unpause() (*types.Transaction, error) {
	return _Main.Contract.Unpause(&_Main.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_Main *MainTransactorSession) Unpause() (*types.Transaction, error) {
	return _Main.Contract.Unpause(&_Main.TransactOpts)
}

// WithdrawAuctionBalances is a paid mutator transaction binding the contract method 0x91876e57.
//
// Solidity: function withdrawAuctionBalances() returns()
func (_Main *MainTransactor) WithdrawAuctionBalances(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "withdrawAuctionBalances")
}

// WithdrawAuctionBalances is a paid mutator transaction binding the contract method 0x91876e57.
//
// Solidity: function withdrawAuctionBalances() returns()
func (_Main *MainSession) WithdrawAuctionBalances() (*types.Transaction, error) {
	return _Main.Contract.WithdrawAuctionBalances(&_Main.TransactOpts)
}

// WithdrawAuctionBalances is a paid mutator transaction binding the contract method 0x91876e57.
//
// Solidity: function withdrawAuctionBalances() returns()
func (_Main *MainTransactorSession) WithdrawAuctionBalances() (*types.Transaction, error) {
	return _Main.Contract.WithdrawAuctionBalances(&_Main.TransactOpts)
}

// WithdrawBalance is a paid mutator transaction binding the contract method 0x5fd8c710.
//
// Solidity: function withdrawBalance() returns()
func (_Main *MainTransactor) WithdrawBalance(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Main.contract.Transact(opts, "withdrawBalance")
}

// WithdrawBalance is a paid mutator transaction binding the contract method 0x5fd8c710.
//
// Solidity: function withdrawBalance() returns()
func (_Main *MainSession) WithdrawBalance() (*types.Transaction, error) {
	return _Main.Contract.WithdrawBalance(&_Main.TransactOpts)
}

// WithdrawBalance is a paid mutator transaction binding the contract method 0x5fd8c710.
//
// Solidity: function withdrawBalance() returns()
func (_Main *MainTransactorSession) WithdrawBalance() (*types.Transaction, error) {
	return _Main.Contract.WithdrawBalance(&_Main.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Main *MainTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Main.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Main *MainSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Main.Contract.Fallback(&_Main.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Main *MainTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Main.Contract.Fallback(&_Main.TransactOpts, calldata)
}

// MainApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Main contract.
type MainApprovalIterator struct {
	Event *MainApproval // Event containing the contract specifics and raw log

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
func (it *MainApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MainApproval)
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
		it.Event = new(MainApproval)
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
func (it *MainApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MainApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MainApproval represents a Approval event raised by the Main contract.
type MainApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address owner, address approved, uint256 tokenId)
func (_Main *MainFilterer) FilterApproval(opts *bind.FilterOpts) (*MainApprovalIterator, error) {

	logs, sub, err := _Main.contract.FilterLogs(opts, "Approval")
	if err != nil {
		return nil, err
	}
	return &MainApprovalIterator{contract: _Main.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address owner, address approved, uint256 tokenId)
func (_Main *MainFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *MainApproval) (event.Subscription, error) {

	logs, sub, err := _Main.contract.WatchLogs(opts, "Approval")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MainApproval)
				if err := _Main.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address owner, address approved, uint256 tokenId)
func (_Main *MainFilterer) ParseApproval(log types.Log) (*MainApproval, error) {
	event := new(MainApproval)
	if err := _Main.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MainBirthIterator is returned from FilterBirth and is used to iterate over the raw logs and unpacked data for Birth events raised by the Main contract.
type MainBirthIterator struct {
	Event *MainBirth // Event containing the contract specifics and raw log

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
func (it *MainBirthIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MainBirth)
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
		it.Event = new(MainBirth)
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
func (it *MainBirthIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MainBirthIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MainBirth represents a Birth event raised by the Main contract.
type MainBirth struct {
	Owner    common.Address
	KittyId  *big.Int
	MatronId *big.Int
	SireId   *big.Int
	Genes    *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBirth is a free log retrieval operation binding the contract event 0x0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5.
//
// Solidity: event Birth(address owner, uint256 kittyId, uint256 matronId, uint256 sireId, uint256 genes)
func (_Main *MainFilterer) FilterBirth(opts *bind.FilterOpts) (*MainBirthIterator, error) {

	logs, sub, err := _Main.contract.FilterLogs(opts, "Birth")
	if err != nil {
		return nil, err
	}
	return &MainBirthIterator{contract: _Main.contract, event: "Birth", logs: logs, sub: sub}, nil
}

// WatchBirth is a free log subscription operation binding the contract event 0x0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5.
//
// Solidity: event Birth(address owner, uint256 kittyId, uint256 matronId, uint256 sireId, uint256 genes)
func (_Main *MainFilterer) WatchBirth(opts *bind.WatchOpts, sink chan<- *MainBirth) (event.Subscription, error) {

	logs, sub, err := _Main.contract.WatchLogs(opts, "Birth")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MainBirth)
				if err := _Main.contract.UnpackLog(event, "Birth", log); err != nil {
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

// ParseBirth is a log parse operation binding the contract event 0x0a5311bd2a6608f08a180df2ee7c5946819a649b204b554bb8e39825b2c50ad5.
//
// Solidity: event Birth(address owner, uint256 kittyId, uint256 matronId, uint256 sireId, uint256 genes)
func (_Main *MainFilterer) ParseBirth(log types.Log) (*MainBirth, error) {
	event := new(MainBirth)
	if err := _Main.contract.UnpackLog(event, "Birth", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MainContractUpgradeIterator is returned from FilterContractUpgrade and is used to iterate over the raw logs and unpacked data for ContractUpgrade events raised by the Main contract.
type MainContractUpgradeIterator struct {
	Event *MainContractUpgrade // Event containing the contract specifics and raw log

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
func (it *MainContractUpgradeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MainContractUpgrade)
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
		it.Event = new(MainContractUpgrade)
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
func (it *MainContractUpgradeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MainContractUpgradeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MainContractUpgrade represents a ContractUpgrade event raised by the Main contract.
type MainContractUpgrade struct {
	NewContract common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterContractUpgrade is a free log retrieval operation binding the contract event 0x450db8da6efbe9c22f2347f7c2021231df1fc58d3ae9a2fa75d39fa446199305.
//
// Solidity: event ContractUpgrade(address newContract)
func (_Main *MainFilterer) FilterContractUpgrade(opts *bind.FilterOpts) (*MainContractUpgradeIterator, error) {

	logs, sub, err := _Main.contract.FilterLogs(opts, "ContractUpgrade")
	if err != nil {
		return nil, err
	}
	return &MainContractUpgradeIterator{contract: _Main.contract, event: "ContractUpgrade", logs: logs, sub: sub}, nil
}

// WatchContractUpgrade is a free log subscription operation binding the contract event 0x450db8da6efbe9c22f2347f7c2021231df1fc58d3ae9a2fa75d39fa446199305.
//
// Solidity: event ContractUpgrade(address newContract)
func (_Main *MainFilterer) WatchContractUpgrade(opts *bind.WatchOpts, sink chan<- *MainContractUpgrade) (event.Subscription, error) {

	logs, sub, err := _Main.contract.WatchLogs(opts, "ContractUpgrade")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MainContractUpgrade)
				if err := _Main.contract.UnpackLog(event, "ContractUpgrade", log); err != nil {
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

// ParseContractUpgrade is a log parse operation binding the contract event 0x450db8da6efbe9c22f2347f7c2021231df1fc58d3ae9a2fa75d39fa446199305.
//
// Solidity: event ContractUpgrade(address newContract)
func (_Main *MainFilterer) ParseContractUpgrade(log types.Log) (*MainContractUpgrade, error) {
	event := new(MainContractUpgrade)
	if err := _Main.contract.UnpackLog(event, "ContractUpgrade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MainPregnantIterator is returned from FilterPregnant and is used to iterate over the raw logs and unpacked data for Pregnant events raised by the Main contract.
type MainPregnantIterator struct {
	Event *MainPregnant // Event containing the contract specifics and raw log

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
func (it *MainPregnantIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MainPregnant)
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
		it.Event = new(MainPregnant)
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
func (it *MainPregnantIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MainPregnantIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MainPregnant represents a Pregnant event raised by the Main contract.
type MainPregnant struct {
	Owner            common.Address
	MatronId         *big.Int
	SireId           *big.Int
	CooldownEndBlock *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterPregnant is a free log retrieval operation binding the contract event 0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80.
//
// Solidity: event Pregnant(address owner, uint256 matronId, uint256 sireId, uint256 cooldownEndBlock)
func (_Main *MainFilterer) FilterPregnant(opts *bind.FilterOpts) (*MainPregnantIterator, error) {

	logs, sub, err := _Main.contract.FilterLogs(opts, "Pregnant")
	if err != nil {
		return nil, err
	}
	return &MainPregnantIterator{contract: _Main.contract, event: "Pregnant", logs: logs, sub: sub}, nil
}

// WatchPregnant is a free log subscription operation binding the contract event 0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80.
//
// Solidity: event Pregnant(address owner, uint256 matronId, uint256 sireId, uint256 cooldownEndBlock)
func (_Main *MainFilterer) WatchPregnant(opts *bind.WatchOpts, sink chan<- *MainPregnant) (event.Subscription, error) {

	logs, sub, err := _Main.contract.WatchLogs(opts, "Pregnant")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MainPregnant)
				if err := _Main.contract.UnpackLog(event, "Pregnant", log); err != nil {
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

// ParsePregnant is a log parse operation binding the contract event 0x241ea03ca20251805084d27d4440371c34a0b85ff108f6bb5611248f73818b80.
//
// Solidity: event Pregnant(address owner, uint256 matronId, uint256 sireId, uint256 cooldownEndBlock)
func (_Main *MainFilterer) ParsePregnant(log types.Log) (*MainPregnant, error) {
	event := new(MainPregnant)
	if err := _Main.contract.UnpackLog(event, "Pregnant", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MainTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Main contract.
type MainTransferIterator struct {
	Event *MainTransfer // Event containing the contract specifics and raw log

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
func (it *MainTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MainTransfer)
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
		it.Event = new(MainTransfer)
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
func (it *MainTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MainTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MainTransfer represents a Transfer event raised by the Main contract.
type MainTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address from, address to, uint256 tokenId)
func (_Main *MainFilterer) FilterTransfer(opts *bind.FilterOpts) (*MainTransferIterator, error) {

	logs, sub, err := _Main.contract.FilterLogs(opts, "Transfer")
	if err != nil {
		return nil, err
	}
	return &MainTransferIterator{contract: _Main.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address from, address to, uint256 tokenId)
func (_Main *MainFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *MainTransfer) (event.Subscription, error) {

	logs, sub, err := _Main.contract.WatchLogs(opts, "Transfer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MainTransfer)
				if err := _Main.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address from, address to, uint256 tokenId)
func (_Main *MainFilterer) ParseTransfer(log types.Log) (*MainTransfer, error) {
	event := new(MainTransfer)
	if err := _Main.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
