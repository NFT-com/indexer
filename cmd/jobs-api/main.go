package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/spf13/pflag"
	"github.com/ziflex/lecho/v2"

	"github.com/NFT-com/indexer/api/controller"
	"github.com/NFT-com/indexer/config/params"
	"github.com/NFT-com/indexer/storage/jobs"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

// run has the server startup code.
func run() int {

	// Signal catching for clean shutdown.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Command line parameter initialization.
	var (
		flagLogLevel string

		flagJobsDB  string
		flagBindAPI string

		flagOpenConnections uint
		flagIdleConnections uint
	)

	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "output level for logging")

	pflag.StringVarP(&flagJobsDB, "jobs-database", "j", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=jobs sslmode=disable", "Postgres connection details for jobs database")
	pflag.StringVarP(&flagBindAPI, "bind-api", "a", "127.0.0.1:8080", "host and port to bind to jobs API endpoint")

	pflag.UintVar(&flagOpenConnections, "open-connections", 16, "limit for open database connections")
	pflag.UintVar(&flagIdleConnections, "idle-connections", 4, "limit for idle database connections")

	pflag.Parse()

	// Logger initialization.
	zerolog.TimestampFunc = func() time.Time { return time.Now().UTC() }
	log := zerolog.New(os.Stderr).With().Timestamp().Logger().Level(zerolog.DebugLevel)
	level, err := zerolog.ParseLevel(flagLogLevel)
	if err != nil {
		log.Error().Err(err).Msg("could not parse log level")
		return failure
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
	jobsDB, err := sql.Open(params.DialectPostgres, flagJobsDB)
	if err != nil {
		log.Error().Err(err).Msg("could not open SQL connection")
		return failure
	}
	jobsDB.SetMaxOpenConns(int(flagOpenConnections))
	jobsDB.SetMaxIdleConns(int(flagIdleConnections))

	// Create the database store.
	parsingRepo := jobs.NewParsingRepository(jobsDB)
	actionRepo := jobs.NewActionRepository(jobsDB)

	// Initialize the REST API controllers.
	parsings := controller.NewParsings(parsingRepo)
	actions := controller.NewActions(actionRepo)

	// Declare the parsing jobs REST API routes.
	server.POST("/parsings/", parsings.Create)
	server.GET("/parsings/:parsing_id", parsings.Read)
	server.PATCH("/parsings/:parsing_id", parsings.Update)
	server.GET("/parsings/", parsings.Index)

	// Declare the action jobs REST API routes.
	server.POST("/actions/", actions.Create)
	server.GET("/actions/:action_id", actions.Read)
	server.PATCH("/actions/:action_id", actions.Update)
	server.GET("/actions/", actions.Index)

	failed := make(chan error)

	go func() {
		log.Info().Msg("jobs api server starting")

		err = server.Start(flagBindAPI)
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
		return failure
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
		return failure
	}

	return success
}
