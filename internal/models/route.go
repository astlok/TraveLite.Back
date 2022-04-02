package models

type Coordinates struct {
	Latitude  string `json:"latitude" example:"56.2348" validate:"required"`
	Longitude string `json:"longitude" example:"64.1352" validate:"required"`
	Height    string `json:"height" example:"5642" validate:"required"`
}

type Mark struct {
	Point       Coordinates `json:"point" validate:"required"`
	Title       string      `json:"title" example:"Pohod228" validate:"required"`
	Description string      `json:"description,omitempty" example:"Lexa zachem mi eto delaem"`
	Photo       string      `json:"photo,omitempty" example:"Tut mogla bit vasha reklama"`
}

type Route struct {
	ID           uint64        `json:"id,omitempty" example:"2"`
	Name         string        `json:"name" example:"Lexa" validate:"required"`
	Difficult    uint8         `json:"difficult" example:"3" validate:"required"`
	Days         uint8         `json:"days" example:"5" validate:"required"`
	Description  string        `json:"description,omitempty" example:"Lexa dava verstai skoree"`
	BestTimeToGo string        `json:"best_time_to_go,omitempty" example:"Лето" validate:"best_time_to_go"`
	Type         string        `json:"type" example:"Пеший" validate:"required"`
	Climb        int           `json:"climb" example:"3800" validate:"required"`
	Region       string        `json:"region" example:"Хабаровский край" validate:"required"`
	CreatorID    int           `json:"creator_id" example:"5" validate:"required"`
	IsModerate   bool          `json:"is_moderate" example:"true" validate:"required"`
	Marks        []Mark        `json:"marks,omitempty" validate:"required"`
	Route        []Coordinates `json:"route,omitempty" validate:"required"`
	Start        Coordinates   `json:"start" validate:"required"`
}

type DBRoute struct {
	ID           uint64 `db:"id"`
	Name         string `db:"name"`
	Difficult    uint8  `db:"difficult"`
	Days         uint8  `db:"days"`
	Description  string `db:"description"`
	BestTimeToGo string `db:"best_time_to_go"`
	Type         string `db:"type"`
	Climb        int    `db:"climb"`
	Region       string `db:"region"`
	CreatorID    int    `db:"creator_id"`
	IsModerate   bool   `db:"is_moderate"`
	Route        string `db:"route"`
	Start        string `db:"start"`
}

type DBMark struct {
	TrekId      uint64 `db:"trek_id"`
	Point       string `db:"point"`
	Title       string `db:"title"`
	Description string `db:"description"`
	Image       string `db:"image"`
}

type SearchPolygon struct {
	NorthEast string
	SouthWest string
}

type RouteFilters struct {
	RoutesInPolygon SearchPolygon
}
