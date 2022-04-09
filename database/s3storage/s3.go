package s3storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Storage struct {
	session *s3.S3
}

func NewS3Storage() *S3Storage {
	return &S3Storage{
		session: s3.New(session.Must(session.NewSession()), &aws.Config{
			Region:      aws.String("ru-msk"),
			Credentials: credentials.NewStaticCredentials("9iecEu1sazvpBUmFsXb7jv", "3Qu3U52eGKkWxYAXAsA7qZbezfVQdDySjPyfsaWgGRDH", ""),
			Endpoint:    aws.String("hb.bizmrg.com"),
		}),
	}
}

func (s *S3Storage) Get() *s3.S3 {
	return s.session
}
