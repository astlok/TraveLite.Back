package app

import (
	filesH "TraveLite/internal/app/files/delivery"
	filesR "TraveLite/internal/app/files/repository"
	filesU "TraveLite/internal/app/files/usecase"
	routemapH "TraveLite/internal/app/route/delivery"
	routemapR "TraveLite/internal/app/route/repository"
	routemapU "TraveLite/internal/app/route/usecase"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func HandlersInit(db *sqlx.DB, s3Store *s3.S3) *echo.Echo {
	e := echo.New()

	routeRepo := routemapR.NewMapRepo(db)
	routeUseCase := routemapU.NewMapUseCase(routeRepo)
	routeHandlers := routemapH.NewMapHandler(routeUseCase)

	filesRepo := filesR.NewFilesRepo(db, s3Store)
	filesUseCase := filesU.NewFilesUseCase(filesRepo)
	filesHandlers := filesH.NewFileHandler(filesUseCase)

	api := e.Group("/api/v1")

	route := api.Group("/route")
	route.POST("", routeHandlers.CreateRoute)
	route.GET("/:id", routeHandlers.GetRoute)
	route.GET("", routeHandlers.GetRoutesWithFilters)

	files := api.Group("/files")

	files.POST("", filesHandlers.Create)
	//files.PUT("/:id", filesHandlers.Upload)
	//files.GET()

	return e
}
