package models

import (
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
)

type Entity int

//const (
//	ToRoute Entity = iota
//	ToUser
//	ToComment
//)
//
//func NewEntity(e string) Entity {
//	switch e {
//	case "route":
//		return Entity(0)
//	case "user":
//		return Entity(1)
//	case "comment":
//		return Entity(2)
//	default:
//		return Entity(0)
//	}
//}
//
//func (e Entity) String() string {
//	return [...]string{"route", "user", "comment"}[e]
//}

type FileMeta struct {
	ID       uuid.UUID `json:"id,omitempty" example:"9b7b28f3-966e-4fbe-b0a2-a7c5f9a90f0f.jbg"`
	Filename string    `json:"filename" validate:"required" example:"kek.txt"`
	Owner    string    `json:"owner" validate:"required" example:"route"`
	OwnerID  uint64    `json:"owner_id" validate:"required" example:"2"`
}

func (f FileMeta) ToDBFile() (DBFile, error) {
	var dbFile DBFile

	err := copier.Copy(&dbFile, &f)
	if err != nil {
		return DBFile{}, errors.Wrap(err, "can't convert file to db file")
	}

	return dbFile, nil
}

type DBFile struct {
	ID       uuid.UUID `db:"id"`
	Filename string    `db:"filename"`
	Owner    string    `db:"owner"`
	OwnerID  uint64    `db:"owner_id"`
}

type FileInfo struct {
	Link string `json:"link" exapmle:"https://trailite-1.hb.bizmrg.com/000050270006_2.jpg"`
}
