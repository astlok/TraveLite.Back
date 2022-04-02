package route

import (
	route "TraveLite/internal/app/route/repository"
	"TraveLite/internal/models"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"strings"
)

var ErrIncorrectPoint = errors.New("Point don't have 3 coordinates")

type RouteUseCase struct {
	routeRepo *route.RouteRepo
}

func NewMapUseCase(routeRepo *route.RouteRepo) *RouteUseCase {
	return &RouteUseCase{
		routeRepo: routeRepo,
	}
}

func (u *RouteUseCase) CreateRoute(route models.Route) (models.Route, error) {
	dbRoute, err := routeToDBRoute(route)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't transform route to DB route")
	}

	id, err := u.routeRepo.CreateRoute(dbRoute)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't insert route to DB")
	}

	route.ID = id

	marks := genDBMarks(route)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't generate DB marks from route")
	}

	err = u.routeRepo.CreateMarks(marks)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't insert marks to DB")
	}

	return route, nil
}

func (u *RouteUseCase) GetRouteByID(id uint64) (models.Route, error) {
	dbRoute, err := u.routeRepo.SelectRouteByID(id)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't select route by id from db")
	}

	dbMarks, err := u.routeRepo.SelectMarksByRouteID(id)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't select marks by id from db")
	}

	uMarks, err := dbMarksToMarks(dbMarks)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't copy from db marks to marks")
	}

	var uRoute models.Route

	uRoute, err = dbRouteToRoute(dbRoute)
	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't copy from db route to route")
	}

	uRoute.Marks = uMarks

	return uRoute, nil
}

func (u *RouteUseCase) GetAllRoutesWithoutRouteLing() ([]models.Route, error) {
	var routes []models.Route

	dbRoutes, err := u.routeRepo.SelectAllRoutesWithoutRouteLine()
	if err != nil {
		return nil, errors.Wrap(err, "can't select all routes without route line from db")
	}

	err = copier.Copy(&routes, &dbRoutes)
	if err != nil {
		return nil, errors.Wrap(err, "can't copy dbRoutes to routes")
	}

	return routes, nil
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
	err := copier.Copy(&r, &dbRoute)

	if err != nil {
		return models.Route{}, errors.Wrap(err, "can't copy fields from DB route to route")
	}

	dbRoute.Route = strings.Replace(dbRoute.Route, "LINESTRING Z (", "", -1)
	dbRoute.Route = strings.Replace(dbRoute.Route, ")", "", -1)

	sPoints := strings.Split(dbRoute.Route, ",")
	for _, sPoint := range sPoints {
		sCoors := strings.Split(sPoint, " ")

		if len(sCoors) != 3 {
			//TODO: Refactor error
			return models.Route{}, ErrIncorrectPoint
		}
		var point models.Coordinates

		point.Latitude = sCoors[0]
		point.Longitude = sCoors[1]
		point.Height = sCoors[2]

		r.Route = append(r.Route, point)
	}

	dbRoute.Start = strings.Replace(dbRoute.Start, "POINT Z (", "", -1)
	dbRoute.Start = strings.Replace(dbRoute.Start, ")", "", -1)

	sStart := strings.Split(dbRoute.Start, " ")
	if len(sStart) != 3 {
		//TODO: Refactor error
		return models.Route{}, ErrIncorrectPoint
	}

	r.Start.Latitude = sStart[0]
	r.Start.Longitude = sStart[1]
	r.Start.Height = sStart[2]

	return r, nil
}

func dbMarksToMarks(dbMarks []models.DBMark) ([]models.Mark, error) {
	marks := make([]models.Mark, len(dbMarks))

	for i, _ := range dbMarks {
		dbMarks[i].Point = strings.Replace(dbMarks[i].Point, "POINT Z (", "", -1)
		dbMarks[i].Point = strings.Replace(dbMarks[i].Point, ")", "", -1)
	}

	err := copier.Copy(&marks, &dbMarks)
	if err != nil || len(marks) != len(dbMarks) {
		return nil, errors.Wrap(err, "can't copy fields from db marks to marks")
	}

	for i, dbMark := range dbMarks {
		coor := strings.Split(dbMark.Point, " ")
		if len(coor) != 3 {
			//TODO: Refactor error
			return nil, ErrIncorrectPoint
		}

		marks[i].Point.Latitude = coor[0]
		marks[i].Point.Longitude = coor[1]
		marks[i].Point.Height = coor[2]
	}

	return marks, nil
}
