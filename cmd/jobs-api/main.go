package main

import (
	"context"
	"database/sql"
	"errors"
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
		flagBind            string
		flagJobsDB          string
		flagLogLevel        string
		flagOpenConnections uint
		flagIdleConnections uint
	)

	pflag.StringVarP(&flagBind, "bind", "b", "127.0.0.1:8081", "host and port for jobs API endpoint")
	pflag.StringVarP(&flagJobsDB, "jobs-database", "d", "", "server details for Postgres database")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "output level for logging")
	pflag.UintVar(&flagOpenConnections, "db-connection-limit", 16, "maximum number of database connections, -1 for unlimited")
	pflag.UintVar(&flagIdleConnections, "db-idle-connection-limit", 4, "maximum number of idle connections")

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
	jobsDB, err := sql.Open(databaseDriver, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not open SQL connection")
		return err
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	// Create the database store.
	jobsStore, err := postgres.NewStore(jobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not create store")
		return err
	}

	// Business logic handler.
	handler := handler.New(jobsStore)

	// Request validator.
	validator := validator.New()

	// REST API Handler.
	apiHandler := api.NewHandler(handler, validator)
	apiHandler.ApplyRoutes(server)

	failed := make(chan error)

	go func() {
		log.Info().Msg("jobs api server starting")

		err = server.Start(flagBind)
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
