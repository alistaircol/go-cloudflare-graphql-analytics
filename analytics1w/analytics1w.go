package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
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

func main() {
	cfZone := get("CLOUDFLARE_ZONE")
	cfEmail := get("CLOUDFLARE_EMAIL")
	cfToken := get("CLOUDFLARE_TOKEN")

	now := time.Now().UTC()
	until := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.UTC)
	from := until.Add(-(time.Hour * 24 * 7))

	graphqlQuery := `query {
  viewer {
    analytics: zones(filter: {zoneTag: $zone}) {
      week: httpRequests1dGroups(filter: {date_gt: $date_gt, date_leq: $date_leq}, limit: 7, orderBy: [date_DESC]) {
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

	// fmt.Println(httpResponse.StatusCode)
	// TODO: check successful response
	res, _ := ioutil.ReadAll(httpResponse.Body)
	fmt.Printf("%s", res)
}

func get(key string) string {
	value, ok := os.LookupEnv(key)

	if !ok {
		color.Red("%s has not been set", key)
		os.Exit(1)
	}

	return value
}
