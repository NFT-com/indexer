package ethereum_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NFT-com/indexer/networks/ethereum"
)

func TestTopicHash(t *testing.T) {
	tests := []struct {
		name  string
		topic string
		hash  string
	}{
		{
			name:  "should return transfer event hash",
			topic: "Transfer",
			hash:  "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
		},
		{
			name:  "should return transfer single event hash",
			topic: "TransferSingle",
			hash:  "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62",
		},
		{
			name:  "should return transfer batch event hash",
			topic: "TransferBatch",
			hash:  "0x4a39dc06d4c0dbc64b70af90fd698a233a518aa5d07e595d983b8c0526c8f7fb",
		},
		{
			name:  "should return uri event hash",
			topic: "URI",
			hash:  "0x6bb7ff708619ba0610cba295a58592e0451dee2622938c8755667688daf3529b",
		},
		{
			name:  "should return zero event hash",
			topic: "",
			hash:  "0x0000000000000000000000000000000000000000000000000000000000000000",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			hash := ethereum.TopicHash(tt.topic)
			assert.Equal(t, tt.hash, hash.Hex())
		})
	}
}
