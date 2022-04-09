package route

import (
	route "TraveLite/internal/app/route/usecase"
	"TraveLite/internal/models"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"

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

//TODO: fix doc region param

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
// @Param limit query int false "limit" default(0)
// @Param offset query int false "offset" default(0)
// @Param type query string false "Route type" Enums:(Пеший, Горный, Водный, Альпинизм, Велотуризм, Бег, Мото, Авто, Скитур, Лыжный, Горный велотуризм, Бездорожье, Ски-альпинизм, Снегоступы)
// @Param region query string false "Region, ТОЛЬКО РАБОТАЕТ С ТЕМИ, ЧТО УЖЕ ЗАБИТЫ В БАЗЕ"
// @Param sort query string false "Sort by this, default rate" Enums:(rate, radius)
// @Param desc query boolean false "sort desc or asc(by default)" default(false)
// @Param difficult query []int false "difficult=1,4  FIRST is from SECOND is to"
// @Param days query []int false "days=2,6  FIRST is from SECOND is to"
// @Param distance query []int false "distance=20,50  FIRST is from SECOND is to"
func (h *RouteHandler) GetRoutesWithFilters(c echo.Context) error {
	conditions := models.NewRouteConditions()

	errs := echo.QueryParamsBinder(c).
		String("ne", &conditions.RoutesInPolygon.NorthEast).
		String("sw", &conditions.RoutesInPolygon.SouthWest).
		Int("offset", &conditions.Offset).
		Int("limit", &conditions.Limit).
		Bool("desc", &conditions.Desc).
		CustomFunc("sort", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("sort", values, "more than 1 value in sort param", nil))
				return errs
			}

			conditions.OrderBy = models.NewSort(values[0])

			return nil
		}).
		CustomFunc("type", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("type", values, "more than 1 value in type param", nil))
				return errs
			}

			conditions.FiltersCol["type"] = values[0]

			return nil
		}).
		CustomFunc("region", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("region", values, "more than 1 value in type param", nil))
				return errs
			}

			conditions.FiltersCol["region"] = values[0]

			return nil
		}).
		CustomFunc("difficult", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("difficult", values, "more than 1 value in difficult param, must be like difficult=1,3", nil))
				return errs
			}

			values[0] = strings.ReplaceAll(values[0], " ", "")
			conditions.FiltersVal["difficult"] = strings.Split(values[0], ",")

			return nil
		}).
		CustomFunc("days", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("days", values, "more than 1 value in difficult param, must be like days=1,3", nil))
				return errs
			}

			values[0] = strings.ReplaceAll(values[0], " ", "")
			conditions.FiltersVal["days"] = strings.Split(values[0], ",")

			return nil
		}).
		CustomFunc("distance", func(values []string) []error {
			var errs []error
			if len(values) > 1 {
				errs = append(errs, echo.NewBindingError("distance", values, "more than 1 value in distance param, must be like days=1,3", nil))
				return errs
			}

			values[0] = strings.ReplaceAll(values[0], " ", "")
			conditions.FiltersVal["distance"] = strings.Split(values[0], ",")

			return nil
		}).
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

	routes, err := h.routeUseCase.GetAllRoutesWithConditions(conditions)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, routes)
}
