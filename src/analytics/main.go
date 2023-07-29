package analytics

import (
	"github.com/alistaircol/go-cloudflare-graphql-analytics/cloudflare"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"net/http"
)

type Period interface {
	GetCloudflareAnalytics() (*http.Response, error)
	GetCloudflareQueryRequest(zone string) cloudflare.GraphqlQueryRequest
	UploadPrimary(client s3.Client, body []byte)
	UploadSecondary(client s3.Client, body []byte)
}

type Daily struct {
	//
}

type Weekly struct {
	//
}

type Monthly struct {
	//
}
