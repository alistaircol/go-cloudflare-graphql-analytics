package main

import (
	"context"
	"errors"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/analytics"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/s3"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) error {
	// TODO: get the granularity from the event
	var interval string

	period, err := MakeAnalyticsPeriod(interval)
	if err != nil {
		//
	}

	res, err := period.GetCloudflareAnalytics()
	if err != nil {
		//
	}

	svc, err := s3.MakeS3Service()
	b, _ := io.ReadAll(res.Body)

	period.UploadPrimary(svc, b)
	period.UploadSecondary(svc, b)

	return nil
}

func MakeAnalyticsPeriod(interval string) (analytics.Period, error) {
	switch interval {
	case "d":
		return analytics.Daily{}, nil
	case "w":
		return analytics.Weekly{}, nil
	case "m":
		return analytics.Monthly{}, nil
	}

	return nil, errors.New("invalid period given")
}
