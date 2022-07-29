package main

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
	"os"
	"os/signal"

	"github.com/NFT-com/indexer/service/api"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func run() int {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagUsername string
		flagPassword string

		flagLogLevel string
		flagGraphDB  string
	)

	pflag.StringVarP(&flagUsername, "username", "u", "admin", "Basic HTTP Auth Username")
	pflag.StringVarP(&flagPassword, "password", "p", "admin", "Basic HTTP Auth Password")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")
	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.Parse()

	err := api.Server(flagUsername, flagPassword, flagGraphDB, flagLogLevel)

	if err != nil {
		log.Error().Err(err).Str("api", "server").Msg("error running HTTP server")
		return failure
	}
	return success
}
