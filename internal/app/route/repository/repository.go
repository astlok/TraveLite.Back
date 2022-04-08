package route

import (
	"TraveLite/internal/models"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

const InsertTrek = `INSERT INTO travelite.trek (
	name,
	difficult,
	distance,
	days,
	description,
	best_time_to_go,
	type,
	climb,
	region,
	creator_id,
	mod_status,
	route,
	start
) VALUES (
	:name,
	:difficult,
    :distance,
	:days,
	:description,
	:best_time_to_go,
	:type,
	:climb,
	:region,
	:creator_id,
	:mod_status,
	ST_GeogFromText(:route),
	ST_GeogFromText(:start)
) RETURNING id;`

const InsertMarks = ` INSERT INTO travelite.marks (
	trek_id, 
	point, 
	title, 
	description, 
	image
) VALUES (
	:trek_id,
	ST_GeogFromText(:point),
	:title,
	:description,
	:image
);`

const SelectMarksByRouteID = `select trek_id, st_astext(point) as point, title, description, image from travelite.marks where trek_id=$1;`

const SelectRouteByID = `SELECT id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, is_moderate, ST_AsText(route) AS ROUTE, ST_AsText(start) AS START from travelite.trek WHERE id = $1;`

const SelectAllRouteWithoutRouteLine = "SELECT id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, mod_status, ST_AsText(start) AS START from travelite.trek"

const SelectAllRoutesInPolygon = `SELECT id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, mod_status, ST_AsText(start) AS START from travelite.trek 
	WHERE ST_Intersects(
		ST_GeogFromText($1),
		START);`

type RouteRepo struct {
	db *sqlx.DB
}

func NewMapRepo(conn *sqlx.DB) *RouteRepo {
	return &RouteRepo{
		db: conn,
	}
}

func (m *RouteRepo) CreateRoute(route models.DBRoute) (uint64, error) {
	//TODO: использовать транзакции
	rows, err := m.db.NamedQuery(InsertTrek, &route)
	if err != nil {
		return 0, err
	}

	var id uint64
	if rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (m *RouteRepo) CreateMarks(marks []models.DBMark) error {
	//TODO: Сделать нормальную вставку сразу всех марок
	for _, mark := range marks {
		_, err := m.db.NamedExec(InsertMarks, &mark)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *RouteRepo) SelectRouteByID(id uint64) (models.DBRoute, error) {
	var route models.DBRoute

	err := m.db.Get(&route, SelectRouteByID, id)
	if err != nil {
		return models.DBRoute{}, err
	}

	return route, nil
}

func (r *RouteRepo) SelectMarksByRouteID(id uint64) ([]models.DBMark, error) {
	var marks []models.DBMark

	err := r.db.Select(&marks, SelectMarksByRouteID, id)
	if err != nil {
		return nil, err
	}

	return marks, nil
}

func (r *RouteRepo) SelectAllRoutesWithCond(conditions models.RouteConditions) ([]models.DBRoute, error) {
	var routes []models.DBRoute

	query, args, err := buildConditionsString(SelectAllRouteWithoutRouteLine, conditions, models.AllowedRouteFilters)
	if err != nil {
		return nil, err
	}

	err = r.db.Select(&routes, query, args...)

	if err != nil {
		return nil, err
	}

	return routes, nil
}

func buildConditionsString(query string, conditions models.RouteConditions, allowedFilters map[string]string) (string, []interface{}, error) {
	filterString := "WHERE 1=1"
	var inputArgs []interface{}

	// filter key is the url name of the filter used as the lookup for the allowed filters list
	for filterKey, filterValList := range conditions.FiltersVal {
		if realFilterName, ok := allowedFilters[filterKey]; ok {
			if len(filterValList) == 0 {
				continue
			}

			filterString = fmt.Sprintf("%s AND %s IN (?)", filterString, realFilterName)
			inputArgs = append(inputArgs, filterValList)
		}
	}

	for filterKey, filterCol := range conditions.FiltersCol {
		if realFilterName, ok := allowedFilters[filterKey]; ok {
			if len(filterCol) == 0 {
				continue
			}

			filterString = fmt.Sprintf("%s AND %s = ?", filterString, realFilterName)
			inputArgs = append(inputArgs, filterCol)
		}
	}

	if conditions.RoutesInPolygon.NorthEast != "" && conditions.RoutesInPolygon.SouthWest != "" {
		ne := strings.Split(conditions.RoutesInPolygon.NorthEast, " ")
		sw := strings.Split(conditions.RoutesInPolygon.SouthWest, " ")

		neLat := ne[0]
		neLong := ne[1]
		swLat := sw[0]
		swLong := sw[1]

		filterString = fmt.Sprintf("%s AND  ST_Intersects(ST_GeogFromText(?), START)", filterString)

		// get string from ne and sw like SRID=4326;POLYGON((55.745359 37.658375, 55.745526 37.705746, 55.724144 37.709792, 55.723866 37.627189, 55.745359 37.658375))
		polygon := "SRID=4326;POLYGON((" + neLat + " " + neLong + ", " + neLat + " " + swLong + ", " + swLat + " " + swLong + ", " + swLat + " " + neLong + ", " + neLat + " " + neLong + "))"
		inputArgs = append(inputArgs, polygon)
	}

	//TODO: сделат ордер бай рабочим
	if conditions.Desc {
		filterString = fmt.Sprintf("%s ORDER BY %s DESC", filterString, conditions.OrderBy.String())
	} else {
		filterString = fmt.Sprintf("%s ORDER BY %s", filterString, conditions.OrderBy.String())
	}

	if conditions.Limit != 0 {
		filterString = fmt.Sprintf("%s LIMIT ?", filterString)
		inputArgs = append(inputArgs, conditions.Limit)
	}

	if conditions.Offset != 0 {
		filterString = fmt.Sprintf("%s OFFSET ?", filterString)
		inputArgs = append(inputArgs, conditions.Offset)
	}

	query, args, err := sqlx.In(fmt.Sprintf("%s %s", query, filterString), inputArgs...)
	if err != nil {
		return "", nil, err
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	return query, args, nil
}
