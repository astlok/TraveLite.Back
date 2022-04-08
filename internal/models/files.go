package models

type Entity int

const (
	ToRoute Entity = iota
	ToUser
	ToComment
)

func NewEntity(e string) Entity {
	switch e {
	case "route":
		return Entity(0)
	case "user":
		return Entity(1)
	case "comment":
		return Entity(2)
	default:
		return Entity(0)
	}
}

func (e Entity) String() string {
	return [...]string{"route", "user", "comment"}[e]
}

type File struct {
	ID    uint64
	Link  string
	Whose Entity
}
