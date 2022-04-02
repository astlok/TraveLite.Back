package route

import (
	"TraveLite/internal/models"
	"github.com/jmoiron/sqlx"
)

const InsertTrek = `INSERT INTO travelite.trek (
	name,
	difficult,
	days,
	description,
	best_time_to_go,
	type,
	climb,
	region,
	creator_id,
	is_moderate,
	route,
	start
) VALUES (
	:name,
	:difficult,
	:days,
	:description,
	:best_time_to_go,
	:type,
	:climb,
	:region,
	:creator_id,
	:is_moderate,
	ST_GeomFromText(:route, 4326),
	ST_GeomFromText(:start, 4326)
) RETURNING id;`

const InsertMarks = ` INSERT INTO travelite.marks (
	trek_id, 
	point, 
	title, 
	description, 
	image
) VALUES (
	:trek_id,
	ST_GeomFromText(:point, 4326),
	:title,
	:description,
	:image
);`

const SelectMarksByRouteID = `select trek_id, st_astext(point) as point, title, description, image from travelite.marks where trek_id=$1;`

const SelectRouteByID = `SELECT id, name, difficult, days, description, best_time_to_go, type, climb, region, creator_id, is_moderate, ST_AsText(route) AS ROUTE, ST_AsText(start) AS START from travelite.trek WHERE id = $1;`

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
