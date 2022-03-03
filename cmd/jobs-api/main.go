package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/NFT-com/indexer/service/controller"
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
	"github.com/NFT-com/indexer/service/postgres"
)

func main() {
	if err := run(); err != nil {
		fmt.Printf("failure: %v\n", err)
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
		flagDBDriver         string
		flagDBConnectionInfo string
		flagLogLevel         string
	)

	pflag.StringVarP(&flagPort, "port", "p", "8081", "server port")
	pflag.StringVar(&flagDBDriver, "driver", "postgres", "postgres connection info")
	pflag.StringVarP(&flagDBConnectionInfo, "db", "d", "", "postgres connection info")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "log level")
	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		return err
	}
	log = log.Level(level)
	eLog := lecho.From(log)

	failed := make(chan error)
	done := make(chan struct{})

	server := echo.New()
	server.HideBanner = true
	server.HidePort = true
	server.Logger = eLog
	server.Use(lecho.Middleware(lecho.Config{Logger: eLog}))

	db, err := sql.Open(flagDBDriver, flagDBConnectionInfo)
	if err != nil {
		return err
	}

	postgresStore, err := postgres.NewStore(db)
	if err != nil {
		return err
	}

	businessController := controller.NewController(postgresStore, postgresStore)
	apiHandler := api.NewHandler(businessController, businessController)

	apiHandler.ApplyRoutes(server)

	go func() {
		log.Info().Msg("dispatcher server starting")

		err := server.Start(fmt.Sprint(":", flagPort))

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Warn().Err(err).Msg("dispatcher server failed")
			failed <- err
			return
		}

		log.Info().Msg("dispatcher server stopped")
		close(done)
	}()

	select {
	case <-sig:
		log.Info().Msg("dispatcher server stopping")
	case <-done:
		log.Info().Msg("dispatcher server done")
	case err := <-failed:
		log.Error().Err(err).Msg("dispatcher server aborted")
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
		log.Error().Err(err).Msg("could not shut down dispatcher API")
		return err
	}

	return nil
}
