package route

import (
	route "TraveLite/internal/app/route/usecase"
	"TraveLite/internal/models"
	"github.com/pkg/errors"
	"net/http"
	"strconv"

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
func (h *RouteHandler) CreateRoute(c echo.Context) error {
	routeCreate := &models.Route{}

	if err := c.Bind(routeCreate); err != nil {
		return err
	}

	routeGet, err := h.routeUseCase.CreateRoute(*routeCreate)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, routeGet)
}

// GetRoute godoc
// @Summary Get route by id
// @Description Get one route by id
// @Tags Route
// @Produce json
// @Param id path uint64 true "Route id"
// @Success 201 {object} models.Route
// @Failure 500 {object} echo.HTTPError
// @Router /route/{id} [get]
func (h *RouteHandler) GetRoute(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return errors.Wrap(err, "Can't parse id from string to uint64")
	}

	r, err := h.routeUseCase.GetRouteByID(id)
	if err != nil {
		return errors.Wrap(err, "Can't get route by id")
	}

	return c.JSON(http.StatusOK, r)
}

// GetRoutesWithFilters godoc
// @Summary Get all routes
// @Description  Get all routes without marks and route line
// @Tags Route
// @Produce json
// @Success 200 {object} []models.Route
// @Failure 500 {object} echo.HTTPError
// @Router /route [get]
// @Param ne query string false "55.745359 37.658375"
// @Param sw query string false "55.971152 63.507595"
func (h *RouteHandler) GetRoutesWithFilters(c echo.Context) error {
	var filters models.RouteFilters

	errs := echo.QueryParamsBinder(c).
		String("ne", &filters.RoutesInPolygon.NorthEast).
		String("sw", &filters.RoutesInPolygon.SouthWest).
		BindErrors()
	if errs != nil {
		errMess := "Can't binding query param: "
		for _, err := range errs {
			bErr := err.(*echo.BindingError)
			c.Logger().Errorf("in case you want to access what field: %s values: %v failed", bErr.Field, bErr.Values)
			errMess += bErr.Field + ", "
		}
		return echo.NewHTTPError(http.StatusBadRequest, errMess)
	}

	routes, err := h.routeUseCase.GetAllRoutesWithFilters(filters)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, routes)
}
