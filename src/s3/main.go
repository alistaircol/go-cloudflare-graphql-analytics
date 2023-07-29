package s3

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"strconv"
)

func MakeS3Service() (s3.S3, error) {
	bucket := os.Getenv("AWS_S3_BUCKET")
	if bucket == "" {
		return s3.S3{}, errors.New("no AWS_S3_BUCKET given")
	}

	log.Print("Creating AWS config for S3 service session")
	config := &aws.Config{}

	uri := os.Getenv("AWS_ENDPOINT")
	if uri != "" {
		log.Printf("Non-empty AWS_ENDPOINT given, setting in config: %s", uri)
		config.Endpoint = aws.String(uri)
	}

	l := os.Getenv("LOCAL_DEVELOPMENT")
	localDevelopment, _ := strconv.ParseBool(l)

	// We need to set up our client slightly differently when using localstack
	if localDevelopment {
		log.Print("LOCAL_DEVELOPMENT environment set, disabling SSL verification and using S3ForcePathStyle")

		config.DisableSSL = aws.Bool(true)
		config.S3ForcePathStyle = aws.Bool(true)
	}

	log.Print("Creating AWS session from config")
	sess, _ := session.NewSession(config)

	log.Print("Creating S3 service from session")
	svc := s3.New(sess)

	return *svc, nil
}
