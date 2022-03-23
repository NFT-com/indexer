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
	"gopkg.in/olahol/melody.v1"

	"github.com/NFT-com/indexer/service/api"
	"github.com/NFT-com/indexer/service/handler"
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

// run has the server startup code.
func run() error {
	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagBind             string
		flagDBConnectionInfo string
		flagLogLevel         string
	)

	pflag.StringVarP(&flagBind, "bind", "b", "8081", "jobs api binding port")
	pflag.StringVarP(&flagDBConnectionInfo, "database", "d", "", "data source name for database connection")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Msg("could not parse log level")
		return err
	}
	log = log.Level(level)
	eLog := lecho.From(log)

	// Initialize echo webserver.
	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Logger = eLog
	server.Use(lecho.Middleware(lecho.Config{Logger: eLog}))

	// Open database connection.
	db, err := sql.Open(databaseDriver, flagDBConnectionInfo)
	if err != nil {
		log.Error().Err(err).Msg("could not open SQL connection")
		return err
	}

	// Create the database store.
	store, err := postgres.NewStore(db)
	if err != nil {
		log.Error().Err(err).Msg("could not create store")
		return err
	}

	// Create the broadcaster.
	broadcaster := melody.New()

	// Business logic handler.
	handler := handler.New(store, broadcaster)

	// Request validator.
	validator := validator.New()

	// REST API Handler.
	apiHandler := api.NewHandler(broadcaster, handler, validator)
	apiHandler.ApplyRoutes(server)

	failed := make(chan error)

	go func() {
		log.Info().Msg("jobs api server starting")

		err = server.Start(fmt.Sprint(":", flagBind))
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn().Err(err).Msg("jobs api server failed")
			failed <- err
			return
		}

		log.Info().Msg("jobs api server done")
	}()

	select {
	case <-sig:
		log.Info().Msg("jobs api server stopping")
	case err = <-failed:
		log.Error().Err(err).Msg("jobs api server aborted")
		return err
	}
	go func() {
		<-sig
		log.Warn().Msg("forcing exit")
		os.Exit(1)
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server.
	err = server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not shut down jobs api")
		return err
	}

	return nil
}
