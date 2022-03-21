package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/NFT-com/indexer/service/validator"
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
		os.Exit(1)
	}

	os.Exit(0)
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
		log.Error().Err(err).Msg("could not parse log level")
		return err
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
		log.Error().Err(err).Msg("could not open SQL connection")
		return err
	}

	store, err := postgres.NewStore(db)
	if err != nil {
		log.Error().Err(err).Msg("could not create store")
		return err
	}

	broadcaster := melody.New()
	businessController := handler.New(store, broadcaster)

	validator := validator.New()

	apiHandler := api.NewHandler(broadcaster, businessController, validator)
	apiHandler.ApplyRoutes(server)

	failed := make(chan error)

	go func() {
		log.Info().Msg("jobs api server starting")

		err = server.Start(fmt.Sprint(":", flagPort))
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

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error().Err(err).Msg("could not shut down jobs api")
		return err
	}

	return nil
}
