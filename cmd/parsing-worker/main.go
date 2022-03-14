package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"github.com/NFT-com/indexer/service/postgres"
	"github.com/NFT-com/indexer/workers/parsing"
)

const (
	EnvVarLogLevel     = "LOG_LEVEL"
	EnvVarDBDriver     = "DATABASE_DRIVER"
	EnvVarDBConnection = "DATABASE_CONNECTION"

	DefaultLogLevel = "info"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}

func run() error {
	logLevel, ok := os.LookupEnv(EnvVarLogLevel)
	if !ok {
		logLevel = DefaultLogLevel
	}
	dbDriver, ok := os.LookupEnv(EnvVarDBDriver)
	if !ok {
		return errors.New("failed to get database driver")
	}
	dbConnection, ok := os.LookupEnv(EnvVarDBConnection)
	if !ok {
		return errors.New("failed to get database connection")
	}

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %v", err)
	}
	log = log.Level(level)

	db, err := sql.Open(dbDriver, dbConnection)
	if err != nil {
		return fmt.Errorf("failed to open slq connection: %v", err)
	}

	store, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("failed to create store: %v", err)
	}

	handler := parsing.NewHandler(log, store)
	lambda.Start(handler.Handle)

	return nil
}
