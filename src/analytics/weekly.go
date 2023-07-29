package analytics

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alistaircol/go-cloudflare-graphql-analytics/cloudflare"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/http"
	"os"
	"time"
)

const GraphqlQueryWeek string = `query {
  viewer {
    analytics: zones(filter: {zoneTag: $zone}) {
      httpRequests1dGroups(filter: {date_gt: $date_gt, date_leq: $date_leq}, limit: 7, orderBy: [date_DESC]) {
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

func (q Weekly) GetCloudflareAnalytics() (*http.Response, error) {
	zone := os.Getenv("CLOUDFLARE_ZONE")
	req := q.GetCloudflareQueryRequest(zone)

	return cloudflare.GetAnalytics(req)
}

func (q Weekly) GetCloudflareQueryRequest(zone string) cloudflare.GraphqlQueryRequest {
	var query = GraphqlQueryDay

	type GraphqlQueryVariablesWeek struct {
		Zone    string `json:"zone"`
		DateGt  string `json:"date_gt"`
		DateLeq string `json:"date_leq"`
	}

	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	from := until.Add(-(time.Hour * 24 * 7))

	variables := GraphqlQueryVariablesWeek{
		Zone:    zone,
		DateGt:  from.Format("2006-01-02"),
		DateLeq: until.Format("2006-01-02"),
	}

	return cloudflare.GraphqlQueryRequest{
		Query:     query,
		Variables: variables,
	}
}

func (q Weekly) UploadPrimary(client s3.Client, body []byte) {
	_, _ = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    aws.String("1w.json"),
		Body:   bytes.NewReader(body),
	})
}

func (q Weekly) UploadSecondary(client s3.Client, body []byte) {
	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	object := aws.String(fmt.Sprintf("w-%s.json", until.Format("2006-01-02")))

	_, _ = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET")),
		Key:    object,
		Body:   bytes.NewReader(body),
	})
}
