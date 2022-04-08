package files

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
)

type FilesRepo struct {
	db *sqlx.DB
	s3 *s3.S3
}

func NewFilesRepo(db *sqlx.DB, s3 *s3.S3) *FilesRepo {
	return &FilesRepo{
		db,
		s3,
	}
}
