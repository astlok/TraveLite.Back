package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(dsn string) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Postgres{
		db: db,
	}, nil
}

func NewDSN(user string, dbName string, pswd string, host string, port string) string {
	return fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", user, dbName, pswd, host, port)
}

func (p *Postgres) GetPostgres() *sqlx.DB {
	return p.db
}

func (p *Postgres) Close() error {
	err := p.db.Close()
	return err
}
