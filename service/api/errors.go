package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// BadRequest returns a bad request error warped with the initial error.
func BadRequest(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusBadRequest, err.Error())
}

// InternalError returns an internal server error warped with the initial error.
func InternalError(err error) *echo.HTTPError {
	return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
}
