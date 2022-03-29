package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// badRequest returns a bad request error warped with the initial error.
func badRequest(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

// internalError returns an internal server error warped with the initial error.
func internalError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
