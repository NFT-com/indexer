package ethereum

import "github.com/ethereum/go-ethereum/common"

const (
	TopicTransfer       = "Transfer"
	TopicTransferSingle = "TransferSingle"
	TopicTransferBatch  = "TransferBatch"
	TopicURI            = "URI"
)

func TopicHash(topic string) common.Hash {
	switch topic {

	case TopicTransfer:
		return common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")

	case TopicTransferSingle:
		return common.HexToHash("0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62")

	case TopicTransferBatch:
		return common.HexToHash("0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb")

	case TopicURI:
		return common.HexToHash("0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b")

	default:
		return common.Hash{}
	}
}
