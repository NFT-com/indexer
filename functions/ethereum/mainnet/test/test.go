package main

import (
	"context"
	"crypto/sha256"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/block"
	"github.com/NFT-com/indexer/event"
)

const (
	EnvVarNodeURL  = "NODE_URL"
	EnvVarLogLevel = "LOG_LEVEL"
)

func main() {
	nodeURL, ok := os.LookupEnv(EnvVarNodeURL)
	if !ok {
		log.Fatalln("missing environment variable")
	}

	logLevel, ok := os.LookupEnv(EnvVarLogLevel)
	if !ok {
		logLevel = "info"
	}

	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		log.Fatalln("failed to parse log level", err)
	}
	logger = logger.Level(level)

	client, err := ethclient.Dial(nodeURL)
	if err != nil {
		log.Fatalln(err)
	}

	handler := New(logger, client)

	lambda.Start(handler.Handle)
}

type Handler struct {
	log    zerolog.Logger
	client *ethclient.Client
}

func New(log zerolog.Logger, client *ethclient.Client) *Handler {
	h := Handler{
		log:    log,
		client: client,
	}

	return &h
}

func (h *Handler) Handle(ctx context.Context, b *block.Block) error {
	blockHash := common.HexToHash(b.Hash)
	query := ethereum.FilterQuery{
		BlockHash: &blockHash,
	}

	logs, err := h.client.FilterLogs(ctx, query)
	if err != nil {
		return err
	}

	for _, l := range logs {
		eventJson, err := l.MarshalJSON()
		if err != nil {
			h.log.Error().Err(err).Str("block_hash", b.Hash).Msg("failed to marshal event")
			continue
		}

		hash := sha256.Sum256(eventJson)
		_ = event.Event{
			ID:              common.Bytes2Hex(hash[:]),
			Network:         b.NetworkID,
			Chain:           b.ChainID,
			Block:           l.BlockNumber,
			TransactionHash: l.TxHash,
			Address:         l.Address,
			Topic:           l.Topics[0],
			IndexedData:     l.Topics[1:],
			Data:            l.Data,
		}
	}

	return nil
}
