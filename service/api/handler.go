package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	DeliveryJobIDParamKey = "delivery_id"
	ParsingJobIDParamKey  = "parsing_id"
)

type Handler struct {
}

func NewHandler() (*Handler, error) {
	s := Handler{}

	return &s, nil
}

func (a *Handler) ApplyRoutes(server *echo.Echo) {
	deliveryJobGroup := server.Group("/deliveries")
	{
		deliveryJobGroup.POST("", a.Noop)
		deliveryJobGroup.GET("", a.Noop)
		deliveryJobGroup.GET("/"+DeliveryJobIDParamKey, a.Noop)
		deliveryJobGroup.PATCH("/"+DeliveryJobIDParamKey, a.Noop)
		deliveryJobGroup.DELETE("/"+DeliveryJobIDParamKey, a.Noop)
		deliveryJobGroup.POST("/"+DeliveryJobIDParamKey+"/requeue", a.Noop)
	}

	parsingsJobGroup := server.Group("/parsings")
	{
		parsingsJobGroup.POST("", a.Noop)
		parsingsJobGroup.GET("", a.Noop)
		parsingsJobGroup.GET("/"+ParsingJobIDParamKey, a.Noop)
		parsingsJobGroup.PATCH("/"+ParsingJobIDParamKey, a.Noop)
		parsingsJobGroup.DELETE("/"+ParsingJobIDParamKey, a.Noop)
		parsingsJobGroup.POST("/"+ParsingJobIDParamKey+"/requeue", a.Noop)
	}
}

func (a *Handler) Noop(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, nil)
}
