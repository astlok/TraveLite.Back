package route

import (
	route "TraveLite/internal/app/route/repository"
	"TraveLite/internal/models"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"strings"
)

type RouteUseCase struct {
	routeRepo *route.RouteRepo
}

func NewMapUseCase(routeRepo *route.RouteRepo) *RouteUseCase {
	return &RouteUseCase{
		routeRepo: routeRepo,
	}
}

func (r *RouteUseCase) CreateRoute(route models.Route) (models.Route, error) {
	dbRoute, err := routeToDBRoute(route)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't transform route to DB route")
	}

	id, err := r.routeRepo.CreateRoute(dbRoute)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't insert route to DB")
	}

	route.ID = id

	marks := genDBMarks(route)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't generate DB marks from route")
	}

	err = r.routeRepo.CreateMarks(marks)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't insert marks to DB")
	}

	return route, nil
}

func (r *RouteUseCase) GetRouteByID(id uint64) (models.Route, error) {
	dbRoute, err := r.routeRepo.SelectRouteByID(id)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't select route by id from db")
	}

	var uRoute models.Route

	uRoute, err = dbRouteToRoute(dbRoute)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't copy from db route to route")
	}

	return uRoute, nil
}

func routeToDBRoute(route models.Route) (models.DBRoute, error) {
	var dbRoute models.DBRoute

	err := copier.Copy(&dbRoute, &route)
	if err != nil {
		return models.DBRoute{}, errors.Wrap(err, "can't copy fields from route to DB route")
	}

	dbRoute.Route += "LINESTRING Z("
	for i, coor := range route.Route {
		dbRoute.Route += coor.Latitude + " " + coor.Longitude + " " + coor.Height
		if i != len(route.Route)-1 {
			dbRoute.Route += ", "
		}
	}
	dbRoute.Route += ")"

	dbRoute.Start += "POINT Z(" + route.Start.Latitude + " " + route.Start.Longitude + " " + route.Start.Height + ")"
	//for i, coor := range route.Route {
	//	dbRoute.Route += coor.Latitude + " " + coor.Longitude + " " + coor.Height
	//	if i != len(route.Route) {
	//		dbRoute.Route += ","
	//	}
	//}
	//dbRoute.Route += ")"

	return dbRoute, nil
}

func genDBMarks(route models.Route) []models.DBMark {
	var dbMarks []models.DBMark

	for _, mark := range route.Marks {
		var m models.DBMark
		m.TrekId = route.ID
		m.Point = "POINT Z(" + mark.Point.Latitude + " " + mark.Point.Longitude + " " + mark.Point.Height + ")"
		m.Title = mark.Title
		//TODO: add photo service
		//m.Photo = mark.Photo
		m.Description = mark.Description

		dbMarks = append(dbMarks, m)
	}

	return dbMarks
}

func dbRouteToRoute(dbRoute models.DBRoute) (models.Route, error) {
	var r models.Route
	err := copier.Copy(r, dbRoute)

	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't copy fields from DB route to route")
	}

	strings.Replace(dbRoute.Route, "LINESTRING Z (", "", -1)
	strings.Replace(dbRoute.Route, ")", "", -1)

	return r, nil
}
