package analytics

import (
	"bytes"
	"fmt"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/cloudflare"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
	"os"
	"time"
)

const GraphqlQueryDay string = `query {
  viewer {
    analytics: zones(filter: {zoneTag: $zone}) {
      httpRequests1hGroups(filter: {datetime_gt: $datetime_gt, datetime_leq: $datetime_leq}, limit: 24, orderBy: [datetime_DESC]) {
        sum {
          requests
        }
        dimensions {
          datetime
        }
        uniq {
          uniques
        }
      }
    }
  }
}
`

func (q Daily) GetCloudflareAnalytics() (*http.Response, error) {
	zone := os.Getenv("CLOUDFLARE_ZONE")
	req := q.GetCloudflareQueryRequest(zone)

	return cloudflare.GetAnalytics(req)
}

func (q Daily) GetCloudflareQueryRequest(zone string) cloudflare.GraphqlQueryRequest {
	var query = GraphqlQueryDay

	type GraphqlQueryVariablesDay struct {
		Zone        string `json:"zone"`
		DatetimeGt  string `json:"datetime_gt"`
		DatetimeLeq string `json:"datetime_leq"`
	}

	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	from := until.Add(-(time.Hour * 24))

	variables := GraphqlQueryVariablesDay{
		Zone:        zone,
		DatetimeGt:  from.Format(time.RFC3339),
		DatetimeLeq: until.Format(time.RFC3339),
	}

	return cloudflare.GraphqlQueryRequest{
		Query:     query,
		Variables: variables,
	}
}

func (q Daily) UploadPrimary(client s3.S3, body []byte) {
	_, _ = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String("1d.json"),
		Body:   bytes.NewReader(body),
	})
}

func (q Daily) UploadSecondary(client s3.S3, body []byte) {
	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	object := aws.String(fmt.Sprintf("d-%s.json", until.Format("2006-01-02T15:04")))

	_, _ = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    object,
		Body:   bytes.NewReader(body),
	})
}
