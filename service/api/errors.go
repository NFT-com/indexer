package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func unpackError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func validateError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func parsingError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

func apiError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
