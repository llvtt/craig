package main

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3ClientWrapper struct {
	sess   *session.Session
	s3     *s3.S3
	bucket string
}

func InitS3Client() *S3ClientWrapper {
	bucket := os.Getenv("CRAIG_S3_BUCKET")
	if len(bucket) == 0 {
		panic("CRAIG_S3_BUCKET is undefined")
	}
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	return &S3ClientWrapper{
		sess,
		s3.New(sess),
		bucket,
	}
}

func (self *S3ClientWrapper) GetObject(key string) io.Reader {
	result, err := self.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(self.bucket),
		Key:    aws.String(key),
	})
	panicOnErr(err)
	return result.Body
}

func (self *S3ClientWrapper) PutObject(key string, reader io.Reader) {
	uploader := s3manager.NewUploader(self.sess)
	panicOnErr(uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(self.bucket),
		Key:    aws.String(key),
		Body:   reader,
	}))
}
