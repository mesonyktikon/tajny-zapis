package storage

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mesonyktikon/tajny-zapis/config"
)

var s3Client *s3.S3

func init() {
	s3Client = s3.New(session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(config.AwsRegion),
		},
	})))
}

func GeneratePresignedPutUrl(s3Key string, fileSize int64) (string, error) {
	req, _ := s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:        aws.String(config.BucketName),
		Key:           aws.String(s3Key),
		ContentLength: aws.Int64(fileSize),
	})
	return req.Presign(1 * time.Minute)
}

func GeneratePresignedGetUrl(s3Key string) (string, error) {
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.BucketName),
		Key:    aws.String(s3Key),
	})
	return req.Presign(1 * time.Minute)
}
