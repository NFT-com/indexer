package api

import (
	"crypto/subtle"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func Server(authUsername, authPassword, graphDbString, logLevel string) error {

	// Start HTTP Server
	e := echo.New()

	e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte(authUsername)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(authPassword)) == 1 {
			return true, nil
		}
		return false, nil
	}))

	// Ping - Health Check
	e.GET("/ping", pingPong)
	// Single endpoint to enqueue a job for a given NFT
	e.POST("/metadata/refresh", refreshTokenMetaData)

	e.Logger.Fatal(e.Start(":8080"))

	return nil
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
