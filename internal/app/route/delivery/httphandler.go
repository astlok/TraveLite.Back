package route

import (
	route "TraveLite/internal/app/route/usecase"
	"TraveLite/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RouteHandler struct {
	routeUseCase *route.RouteUseCase
}

func NewMapHandler(routeUseCase *route.RouteUseCase) *RouteHandler {
	return &RouteHandler{
		routeUseCase: routeUseCase,
	}
}

// CreateRoute godoc
// @Summary Create a route
// @Description Create a new route
// @Tags Route
// @Accept json
// @Produce json
// @Param route body models.Route true "new route"
// @Success 201 {object} models.Route
// @Failure 500 {object} echo.HTTPError
// @Router /route [post]
func (m *RouteHandler) CreateRoute(c echo.Context) error {
	routeCreate := &models.Route{}

	if err := c.Bind(routeCreate); err != nil {
		return err
	}

	routeGet, err := m.routeUseCase.CreateRoute(*routeCreate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, routeGet)
}
