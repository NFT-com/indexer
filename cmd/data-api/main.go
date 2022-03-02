package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/ziflex/lecho/v2"

	api "github.com/NFT-com/indexer/service/api/data"
	handler "github.com/NFT-com/indexer/service/handler/data"
	"github.com/NFT-com/indexer/service/postgres"
	"github.com/NFT-com/indexer/service/validator"
)

const (
	databaseDriver = "postgres"
)

func main() {
	err := run()
	if err != nil {
		// TODO: Improve this mixing logging
		// https://github.com/NFT-com/indexer/issues/32
		log.Fatal(err)
	}
}

func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagPort             string
		flagDBConnectionInfo string
		flagLogLevel         string
	)

	pflag.StringVarP(&flagPort, "port", "p", "8081", "server port")
	pflag.StringVarP(&flagDBConnectionInfo, "db", "d", "", "data source name for database connection")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}
	log = log.Level(level)
	eLog := lecho.From(log)

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Logger = eLog
	server.Use(lecho.Middleware(lecho.Config{Logger: eLog}))

	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		return fmt.Errorf("could not open SQL connection: %w", err)
	}

	postgresStore, err := postgres.NewStore(db)
	if err != nil {
		return fmt.Errorf("could not create store: %w", err)
	}

	handler := handler.NewHandler(postgresStore)

	// Request validator.
	validator := validator.New()

	apiHandler := api.NewHandler(handler, validator)
	apiHandler.RegisterEndpoints(server)

	failed := make(chan error)

	go func() {
		log.Info().Msg("data api server starting")

		err = server.Start(fmt.Sprint(":", flagPort))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			failed <- err
			return
		}

		log.Info().Msg("data api server done")
	}()

	select {
	case <-sig:
		log.Info().Msg("data api server stopping")
	case err = <-failed:
		log.Error().Err(err).Msg("data api server aborted")
		return err
	}
	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not gracefully shutdown data api")
		return err
	}

	return nil
}
