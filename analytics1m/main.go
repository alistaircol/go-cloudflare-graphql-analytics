package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type graphqlQueryVariables struct {
	Zone    string `json:"zone"`
	DateGt  string `json:"date_gt"`
	DateLeq string `json:"date_leq"`
}

type graphqlQueryRequest struct {
	Query     string                `json:"query"`
	Variables graphqlQueryVariables `json:"variables"`
}

func handler(ctx context.Context) error {
	analytics := getAnalytics()
	body, _ := ioutil.ReadAll(analytics.Body)

	uploadAnalyticsPrimary(body)
	uploadAnalyticsSecondary(body)

	return nil
}

func main() {
	lambda.Start(handler)
}

func getEnvironmentVariable(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		color.Red("%s has not been set", key)
		os.Exit(1)
	}

	return value
}

func getAnalytics() *http.Response {
	cfZone := getEnvironmentVariable("CLOUDFLARE_ZONE")
	cfEmail := getEnvironmentVariable("CLOUDFLARE_EMAIL")
	cfToken := getEnvironmentVariable("CLOUDFLARE_TOKEN")

	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	from := until.Add(-(time.Hour * 24 * 30))

	graphqlQuery := `query {
  viewer {
    analytics: zones(filter: {zoneTag: $zone}) {
      httpRequests1dGroups(filter: {date_gt: $date_gt, date_leq: $date_leq}, limit: 30, orderBy: [date_DESC]) {
        sum {
          requests
        }
        dimensions {
          date
        }
        uniq {
          uniques
        }
      }
    }
  }
}
`

	// https://stackoverflow.com/a/62479701/5873008
	graphqlVariables := &graphqlQueryVariables{
		Zone:    cfZone,
		DateGt:  from.Format("2006-01-02"),
		DateLeq: until.Format("2006-01-02"),
	}

	graphqlQueryRequest := &graphqlQueryRequest{
		Query:     graphqlQuery,
		Variables: *graphqlVariables,
	}

	graphqlQueryRequestJson, _ := json.Marshal(graphqlQueryRequest)
	httpRequestBody := bytes.NewBuffer(graphqlQueryRequestJson)

	// https://stackoverflow.com/a/41034588/5873008
	httpClient := &http.Client{}
	httpRequest, _ := http.NewRequest(http.MethodPost, "https://api.cloudflare.com/client/v4/graphql/", httpRequestBody)
	httpRequest.Header.Set("Content-Type", "application/res")
	httpRequest.Header.Set("X-Auth-Email", cfEmail)
	httpRequest.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cfToken))

	// https://stackoverflow.com/a/24512194/5873008
	httpResponse, _ := httpClient.Do(httpRequest)

	if httpResponse.StatusCode != 200 {
		fmt.Println("no successful http response")
		fmt.Printf("%+v", httpResponse)
		os.Exit(1)
	}

	return httpResponse
}

func uploadAnalyticsPrimary(body []byte) {
	// https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	bucket := getEnvironmentVariable("AWS_S3_BUCKET")
	awsS3Client := s3.NewFromConfig(awsConfig)
	awsS3Uploader := manager.NewUploader(awsS3Client)

	object := aws.String("1m.json")
	result, err := awsS3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    object,
		Body:   bytes.NewReader(body),
	})

	log.Printf("primary result: %v", result)
}

func uploadAnalyticsSecondary(body []byte) {
	// https://aws.github.io/aws-sdk-go-v2/docs/sdk-utilities/s3/
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	bucket := getEnvironmentVariable("AWS_S3_BUCKET")
	awsS3Client := s3.NewFromConfig(awsConfig)
	awsS3Uploader := manager.NewUploader(awsS3Client)

	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	object := aws.String(fmt.Sprintf("m-%s.json", until.Format("2006-01-02T15:04")))

	result, err := awsS3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    object,
		Body:   bytes.NewReader(body),
	})

	log.Printf("secondary result: %v", result)
}
