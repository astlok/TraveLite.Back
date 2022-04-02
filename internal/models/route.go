package models

type Coordinates struct {
	Latitude  string `json:"latitude" example:"56.2348"`
	Longitude string `json:"longitude" example:"64.1352"`
	Height    string `json:"height" example:"5642"`
}

type Mark struct {
	Point       Coordinates `json:"point"`
	Title       string      `json:"title" example:"Pohod228"`
	Description string      `json:"description,omitempty" example:"Lexa zachem mi eto delaem"`
	Photo       string      `json:"photo,omitempty" example:"Tut mogla bit vasha reklama"`
}

type Route struct {
	ID           uint64        `json:"id,omitempty" example:"2"`
	Name         string        `json:"name" example:"Lexa"`
	Difficult    uint8         `json:"difficult" example:"3"`
	Days         uint8         `json:"days" example:"5"`
	Description  string        `json:"description,omitempty" example:"Lexa dava verstai skoree"`
	BestTimeToGo string        `json:"best_time_to_go,omitempty" example:"Лето"`
	Type         string        `json:"type" example:"Пеший"`
	Climb        int           `json:"climb" example:"3800"`
	Region       string        `json:"region" example:"Хабаровский край"`
	CreatorID    int           `json:"creator_id" example:"5"`
	IsModerate   bool          `json:"is_moderate" example:"true"`
	Marks        []Mark        `json:"marks,omitempty"`
	Route        []Coordinates `json:"route,omitempty"`
	Start        Coordinates   `json:"start"`
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
