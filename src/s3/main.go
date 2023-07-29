package s3

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"strconv"
)

func MakeS3Service() (s3.Client, error) {
	bucket := os.Getenv("PAYLOAD_BUCKET_NAME")
	if bucket == "" {
		return s3.Client{}, errors.New("no PAYLOAD_BUCKET_NAME given")
	}

	var cfg aws.Config

	cfg, _ = config.LoadDefaultConfig(context.TODO())

	uri := os.Getenv("AWS_ENDPOINT")
	if uri != "" {
		log.Printf("Non-empty AWS_ENDPOINT given, setting in config: %s", uri)

		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID {
				return aws.Endpoint{
					URL: uri,
				}, nil
			}

			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		cfg, _ = config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver))
	}

	var svc s3.Client

	svc = *s3.NewFromConfig(cfg)

	l := os.Getenv("LOCAL_DEVELOPMENT")
	localDevelopment, _ := strconv.ParseBool(l)

	// We need to set up our client slightly differently when using localstack
	if localDevelopment {
		log.Print("LOCAL_DEVELOPMENT environment set, disabling SSL verification and using S3ForcePathStyle")

		//config.DisableSSL = aws.Bool(true)
		//config.S3ForcePathStyle = aws.Bool(true)

		svc = *s3.NewFromConfig(cfg, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}

	return svc, nil
}
