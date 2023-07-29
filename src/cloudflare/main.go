package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type GraphqlQueryRequest struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables"`
}

func GetAnalytics(req GraphqlQueryRequest) (*http.Response, error) {
	d, _ := json.Marshal(req)

	log.Printf("cloudflare/main.go: request data: %s", d)

	// https://stackoverflow.com/a/41034588/5873008
	c := &http.Client{}

	email := os.Getenv("CLOUDFLARE_EMAIL")
	token := os.Getenv("CLOUDFLARE_TOKEN")

	request, _ := http.NewRequest(http.MethodPost, "https://api.cloudflare.com/client/v4/graphql/", bytes.NewBuffer(d))
	request.Header.Set("Content-Type", "application/res")
	request.Header.Set("X-Auth-Email", email)
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := c.Do(request)

	if err != nil || res.StatusCode != 200 {
		fmt.Println("error occurred, or no successful http response")
		return &http.Response{}, err
	}

	return res, nil
}
