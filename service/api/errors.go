package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func unpackError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err)
}

func parsingError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err)
}

func apiError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
