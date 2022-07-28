package main

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/pflag"
	"net/http"
	"os"
	"os/signal"
)

const (
	success = 0
	failure = 1
)

func main() {
	os.Exit(run())
}

func pingPong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func refreshTokenMetaData(c echo.Context) error {
	// Get contractAddress and tokenId
	contractAddress := c.FormValue("contractAddress")
	tokenId := c.FormValue("tokenId")

	// TODO: Enqueue a job to refresh metadata for token/tokenId

	return c.String(http.StatusOK, "contractAddress:"+contractAddress+", tokenId:"+tokenId)
}

func run() int {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	var (
		flagUsername string
		flagPassword string

		flagLogLevel string

		flagGraphDB string
	)

	pflag.StringVarP(&flagUsername, "username", "u", "admin", "Basic HTTP Auth Username")
	pflag.StringVarP(&flagPassword, "password", "p", "admin", "Basic HTTP Auth Password")
	pflag.StringVarP(&flagLogLevel, "log-level", "l", "info", "severity level for log output")
	pflag.StringVarP(&flagGraphDB, "graph-database", "g", "host=127.0.0.1 port=5432 user=postgres password=postgres dbname=graph sslmode=disable", "Postgres connection details for graph database")
	pflag.Parse()

	// Start HTTP Server
	e := echo.New()

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(flagUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(flagPassword)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Ping - Health Check
	e.GET("/ping", pingPong)
	// Single endpoint to enqueue a job for a given NFT
	e.POST("/metadata/refresh", refreshTokenMetaData)

	e.Logger.Fatal(e.Start(":8080"))

	return success
}
