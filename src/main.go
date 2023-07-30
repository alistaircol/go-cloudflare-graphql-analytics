package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/analytics"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/s3"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-xray-sdk-go/xray"
	"io"
	"log"
)

type Detail struct {
	Period string `json:"period"`
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.CloudWatchEvent) error {
	ctx, seg := xray.BeginSegment(context.Background(), "go-cloudflare-graphql-analytics")

	log.Print("main.go: Going to get details from event")
	var detail Detail
	if err := json.Unmarshal(event.Detail, &detail); err != nil {
		return err
	}

	log.Printf("main.go: event detail is: %+v", detail)
	log.Print("main.go: Going to get details from event")
	period, err := MakeAnalyticsPeriod(detail.Period)
	if err != nil {
		log.Printf("main.go: error from MakeAnalyticsPeriod: %+v", err)
		return err
	}

	log.Printf("main.go: period: %+v", period)
	res, err := period.GetCloudflareAnalytics()
	if err != nil {
		log.Printf("main.go: error from period.GetCloudflareAnalytics: %+v", err)
		return err
	}

	log.Print("main.go: Going to make S3 service")
	svc, err := s3.MakeS3Service()
	b, _ := io.ReadAll(res.Body)

	log.Printf("main.go: Going to upload: %s", string(b))

	log.Print("main.go: Going upload primary")
	period.UploadPrimary(svc, b)
	log.Print("main.go: Going upload secondary")
	period.UploadSecondary(svc, b)

	seg.Close(nil)

	return nil
}

func MakeAnalyticsPeriod(interval string) (analytics.Period, error) {
	switch interval {
	case "d":
		log.Print("main.go: returning daily from event detail")
		return analytics.Daily{}, nil
	case "w":
		log.Print("main.go: returning weekly from event detail")
		return analytics.Weekly{}, nil
	case "m":
		log.Print("main.go: returning monthly from event detail")
		return analytics.Monthly{}, nil
	}

	return nil, errors.New("invalid period given")
}
