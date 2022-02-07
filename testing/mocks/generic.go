package mocks

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
)

var (
	GenericErrChannel = make(chan error)

	GenericContractType = "erc721"
	GenericContractABI  = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"contractURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"receivers\",\"type\":\"address[]\"}],\"name\":\"gift\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reveal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"extension\",\"type\":\"string\"}],\"name\":\"setBaseExtension\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"URI\",\"type\":\"string\"}],\"name\":\"setBaseURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"URI\",\"type\":\"string\"}],\"name\":\"setContractURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"URI\",\"type\":\"string\"}],\"name\":\"setNotRevealURI\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reserve\",\"type\":\"uint256\"}],\"name\":\"setReserve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"setVaultAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"switchSale\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

	GenericEthereumBlockHeader = &types.Header{
		ParentHash:  common.HexToHash("0x120fc9108e799e0cd10987dab1aa10ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
		UncleHash:   common.HexToHash("0x120fc9108e799e0cd10987dab1aa20ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
		Coinbase:    common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc41"),
		Root:        common.HexToHash("0x120fc9108e799e0cd10987dab1aa30ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
		TxHash:      common.HexToHash("0x120fc9108e799e0cd10987dab1aa40ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
		ReceiptHash: common.HexToHash("0x120fc9108e799e0cd10987dab1aa50ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
		Difficulty:  big.NewInt(1),
		Number:      big.NewInt(2),
	}
	GenericEthereumLogs = []types.Log{
		{
			Address: common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc43"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"),
				common.HexToHash("0x000000000000000000000000283af0b28c62c092c9727f1ee09c02ca627eb7f5"),
				common.HexToHash("0xc2a774a71c5c0b82d7dc3ddc44182605d50feeda3a3cef90dd117adcaa317be3"),
			},
			BlockNumber: 1,
			TxHash:      common.HexToHash("0x120fc9108e799e0cd10987dab1aa60ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			TxIndex:     1,
			BlockHash:   common.HexToHash("0x120fc9108e799e0cd10987dab1aa70ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			Index:       1,
		},
		{
			Address: common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc43"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x0000000000000000000000005a654ba2d2262cfd870a5ddbe0d5c672d0c88568"),
				common.HexToHash("0x000000000000000000000000568b47030b2bff3e49fbfc7f1ec5e270de0a28a6"),
				common.HexToHash("0x90d9b112c23dc6a63cf1391750c536959bb9699fdb4f6dba5b594f041ae6fba4"),
			},
			BlockNumber: 2,
			TxHash:      common.HexToHash("0x120fc9108e799e0cd10987dab1aa80ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			TxIndex:     1,
			BlockHash:   common.HexToHash("0x120fc9108e799e0cd10987dab1aa90ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			Index:       1,
		},
		{
			Address: common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc43"),
			Topics: []common.Hash{
				common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
				common.HexToHash("0x000000000000000000000000283af0b28c62c092c9727f1ee09c02ca627eb7f5"),
				common.HexToHash("0x000000000000000000000000e5383c637515dbf520c95ccd79ee657e5471bb6e"),
				common.HexToHash("0xc2a774a71c5c0b82d7dc3ddc44182605d50feeda3a3cef90dd117adcaa317be3"),
			},
			BlockNumber: 3,
			TxHash:      common.HexToHash("0x120fc9108e799e0cd10987dab1aa11ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			TxIndex:     1,
			BlockHash:   common.HexToHash("0x120fc9108e799e0cd10987dab1aa12ad33cc6a9e2924a16ad6b0e7762c9c1d22"),
			Index:       1,
		},
	}

	GenericBlock  = block.Block("0x120fc9108e799e0cd10987dab1aa30ad33cc6a9e2924a16ad6b0e7762c9c1d22")
	GenericEvents = []*event.Event{
		{
			ID:              "0x120fc9108e799e0cd10987dab1aa30ad33cc6a9e2924a16ad6b0e7762c9c1d21",
			Network:         "ethereum",
			Chain:           "mainnet",
			Block:           1,
			TransactionHash: common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3e2"),
			Address:         common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc43"),
			Topic:           common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			IndexedData: []common.Hash{
				common.HexToHash("0x0000000000000000000000005a654ba2d2262cfd870a5ddbe0d5c672d0c88568"),
				common.HexToHash("0x000000000000000000000000568b47030b2bff3e49fbfc7f1ec5e270de0a28a6"),
				common.HexToHash("0x90d9b112c23dc6a63cf1391750c536959bb9699fdb4f6dba5b594f041ae6fba4"),
			},
		},
		{
			ID:              "0x120fc9108e799e0cd10987dab1aa30ad33cc6a9e2924a16ad6b0e7762c9c1d20",
			Network:         "ethereum",
			Chain:           "mainnet",
			Block:           2,
			TransactionHash: common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3e1"),
			Address:         common.HexToAddress("0x68b3465833fb72a70ecdf485e0e4c7bd8665fc43"),
			Topic:           common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"),
			IndexedData: []common.Hash{
				common.HexToHash("0x000000000000000000000000283af0b28c62c092c9727f1ee09c02ca627eb7f5"),
				common.HexToHash("0x000000000000000000000000e5383c637515dbf520c95ccd79ee657e5471bb6e"),
				common.HexToHash("0xc2a774a71c5c0b82d7dc3ddc44182605d50feeda3a3cef90dd117adcaa317be3"),
			},
		},
	}
)
