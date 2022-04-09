package files

import (
	"TraveLite/internal/models"
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jmoiron/sqlx"
)

const InsertFile = `INSERT INTO travelite.files(
    id, filename, owner, owner_id
) VALUES (
          :id,
          :filename,
          :owner,
          :owner_id
);`

type FilesRepo struct {
	db      *sqlx.DB
	s3store *s3.S3
	bucket  string
}

func NewFilesRepo(db *sqlx.DB, s3 *s3.S3, bucket string) *FilesRepo {
	return &FilesRepo{
		db,
		s3,
		bucket,
	}
}

func (f *FilesRepo) Create(file models.FileMeta) error {
	dbFile, err := file.ToDBFile()
	if err != nil {
		return err
	}

	_, err = f.db.NamedQuery(InsertFile, &dbFile)
	if err != nil {
		return err
	}
	return nil
}

func (f *FilesRepo) Upload(file []byte, fileID string) error {
	buf := bytes.NewReader(file)

	_, err := f.s3store.PutObjectWithContext(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(f.bucket),
		Key:         aws.String(fileID),
		Body:        buf,
		ContentType: aws.String("image"),
	})

	if err != nil {
		return err
	}

	return nil
}
