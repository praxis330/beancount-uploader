package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"os"
)

func InitAWSUploader(title string) AWSUploader { // Get credentials with environment variables.
	return AWSUploader{
		S3Config{
			AWSConfig: AWSConfig{
				ID:     os.Getenv("AWS_ACCESS_KEY_ID"),
				Secret: os.Getenv("AWS_SECRET_ACCESS_KEY"),
				Region: os.Getenv("AWS_REGION"),
			},
			Bucket: os.Getenv("BEANCOUNT_BUCKET"),
			Key:    title,
		},
	}
}

type Uploader interface {
	Upload(b *BeancountItem) error
}

type AWSConfig struct {
	ID     string
	Secret string
	Region string
}

type S3Config struct {
	AWSConfig
	Bucket string
	Key    string
}

type AWSUploader struct {
	S3Config
}

func (a *AWSUploader) Upload(b *BeancountItem) error {
	json, _ := json.Marshal(b)

	err := a.uploadToS3(json)

	return err
}

func (a *AWSUploader) uploadToS3(json []byte) error {

	creds := credentials.NewStaticCredentials(a.ID, a.Secret, "")

	config := aws.NewConfig().WithCredentials(creds).WithRegion(a.Region)
	session := session.New(config)
	svc := s3.New(session)

	params := &s3.PutObjectInput{
		Bucket:      aws.String(a.Bucket),
		Key:         aws.String(a.Key),
		ACL:         aws.String("public-read"),
		ContentType: aws.String("text/plain"),
		Body:        bytes.NewReader(json),
	}

	_, err := svc.PutObject(params)

	return err
}
