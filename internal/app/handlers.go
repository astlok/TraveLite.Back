package app

import (
	routemapH "TraveLite/internal/app/route/delivery"
	routemapR "TraveLite/internal/app/route/repository"
	routemapU "TraveLite/internal/app/route/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func HandlersInit(db *sqlx.DB) *echo.Echo {
	e := echo.New()

	routeRepo := routemapR.NewMapRepo(db)
	routeUseCase := routemapU.NewMapUseCase(routeRepo)
	routeHandlers := routemapH.NewMapHandler(routeUseCase)

	api := e.Group("/api/v1")
	api.POST("/route", routeHandlers.CreateRoute)
	api.GET("/route/:id", routeHandlers.GetRoute)
	api.GET("/route", routeHandlers.GetRoutesWithFilters)

	return e
}
