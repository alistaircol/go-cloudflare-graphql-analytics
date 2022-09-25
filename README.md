> :warning: **Work in progress**

# Overview

Some `go` apps which will eventually run in [AWS Lambda](https://aws.amazon.com/lambda/) on a schedule with [Amazon EventBridge](https://aws.amazon.com/eventbridge/) and put analytics payloads in an [S3](https://aws.amazon.com/s3/) bucket for use in [chartjs](https://www.chartjs.org/).

There are three executables:

* `bin/analytics1d` - from the last UTC hour, get the previous 24 hours `requests` and `uniques`
* `bin/analytics1w` - from the last UTC day, get the previous 7 days `requests` and `uniques`
* `bin/analytics1m` - from the last UTC day, get the previous 7 days `requests` and `uniques`

## MVP

<details>
<summary>A bash script to get day/week/month analytics from Coudflare's graphql analytics api and example output</summary>

`query.sh`:

```bash
#!/usr/bin/env bash
ZONE="YOUR_CLOUDFLARE_ANALYTICS_ZONE"
GRAPHQL_QUERY=$(cat <<"EOF"
query {
  viewer {
    analytics: zones(filter: {zoneTag: $zone}) {
      day: httpRequests1hGroups(filter: $filter1d, limit: 24, orderBy: [datetime_DESC]) {
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
      week: httpRequests1dGroups(filter: $filter1w, limit: 7, orderBy: [date_DESC]) {
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
      month: httpRequests1dGroups(filter: $filter1m, limit: 30, orderBy: [date_DESC]) {
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
EOF
)

# using python datetime because `date` isn't very cross-platform friendly 
jq  -n \
  --arg zone "$ZONE" \
  --arg graphql_query "$GRAPHQL_QUERY" \
  --arg filter_1d_from "$(python3 -c "from datetime import datetime, timedelta; print((datetime.utcnow() - timedelta(days=1)).strftime(\"%Y-%m-%dT%H:00:00Z\"))")" \
  --arg filter_1d_until "$(python3 -c "from datetime import datetime; print(datetime.utcnow().strftime(\"%Y-%m-%dT%H:00:00Z\"))")" \
  --arg filter_1w_from "$(python3 -c "from datetime import datetime, timedelta; print((datetime.utcnow() - timedelta(days=7)).strftime(\"%Y-%m-%d\"))")" \
  --arg filter_1w_until "$(python3 -c "from datetime import datetime; print(datetime.utcnow().strftime(\"%Y-%m-%d\"))")" \
  --arg filter_1m_from "$(python3 -c "from datetime import datetime, timedelta; print((datetime.utcnow() - timedelta(days=30)).strftime(\"%Y-%m-%d\"))")" \
  --arg filter_1m_until "$(python3 -c "from datetime import datetime; print(datetime.utcnow().strftime(\"%Y-%m-%d\"))")" \
  "$(cat <<"EOF"
{
    "query": $graphql_query,
    "variables": {
        "zone": $zone,
        "filter1d": {
            "datetime_gt": $filter_1d_from,
            "datetime_leq": $filter_1d_until
        },
        "filter1w": {
            "date_gt": $filter_1w_from,
            "date_leq": $filter_1w_until
        },
        "filter1m": {
            "date_gt": $filter_1m_from,
            "date_leq": $filter_1m_until
        }
    }
}
EOF
)" > request.json

curl --silent \
  --request POST \
  --header "Content-Type: application/json" \
  --header "X-Auth-Email: your+account@cloudflare.com" \
  --header "Authorization: Bearer your_analytics_api_token" \
  --data @request.json \
  https://api.cloudflare.com/client/v4/graphql/ | jq . > stats.json

rm request.json
```

</details>

<details>
<summary>stats.json</summary>

```json
{
  "data": {
    "viewer": {
      "analytics": [
        {
          "day": [
            {
              "dimensions": {
                "datetime": "2022-09-25T00:00:00Z"
              },
              "sum": {
                "requests": 6
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T23:00:00Z"
              },
              "sum": {
                "requests": 20
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T22:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T21:00:00Z"
              },
              "sum": {
                "requests": 34
              },
              "uniq": {
                "uniques": 14
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T20:00:00Z"
              },
              "sum": {
                "requests": 46
              },
              "uniq": {
                "uniques": 21
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T19:00:00Z"
              },
              "sum": {
                "requests": 219
              },
              "uniq": {
                "uniques": 12
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T18:00:00Z"
              },
              "sum": {
                "requests": 63
              },
              "uniq": {
                "uniques": 15
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T17:00:00Z"
              },
              "sum": {
                "requests": 91
              },
              "uniq": {
                "uniques": 27
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T16:00:00Z"
              },
              "sum": {
                "requests": 32
              },
              "uniq": {
                "uniques": 12
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T15:00:00Z"
              },
              "sum": {
                "requests": 31
              },
              "uniq": {
                "uniques": 13
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T14:00:00Z"
              },
              "sum": {
                "requests": 27
              },
              "uniq": {
                "uniques": 6
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T13:00:00Z"
              },
              "sum": {
                "requests": 13
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T12:00:00Z"
              },
              "sum": {
                "requests": 33
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T11:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T10:00:00Z"
              },
              "sum": {
                "requests": 24
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T09:00:00Z"
              },
              "sum": {
                "requests": 38
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T08:00:00Z"
              },
              "sum": {
                "requests": 25
              },
              "uniq": {
                "uniques": 10
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T07:00:00Z"
              },
              "sum": {
                "requests": 41
              },
              "uniq": {
                "uniques": 17
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T06:00:00Z"
              },
              "sum": {
                "requests": 34
              },
              "uniq": {
                "uniques": 23
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T05:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 6
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T04:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T03:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T02:00:00Z"
              },
              "sum": {
                "requests": 23
              },
              "uniq": {
                "uniques": 10
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T01:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            }
          ],
          "month": [
            {
              "dimensions": {
                "date": "2022-09-25"
              },
              "sum": {
                "requests": 9
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "date": "2022-09-24"
              },
              "sum": {
                "requests": 911
              },
              "uniq": {
                "uniques": 170
              }
            },
            {
              "dimensions": {
                "date": "2022-09-23"
              },
              "sum": {
                "requests": 810
              },
              "uniq": {
                "uniques": 137
              }
            },
            {
              "dimensions": {
                "date": "2022-09-22"
              },
              "sum": {
                "requests": 888
              },
              "uniq": {
                "uniques": 134
              }
            },
            {
              "dimensions": {
                "date": "2022-09-21"
              },
              "sum": {
                "requests": 742
              },
              "uniq": {
                "uniques": 89
              }
            },
            {
              "dimensions": {
                "date": "2022-09-20"
              },
              "sum": {
                "requests": 750
              },
              "uniq": {
                "uniques": 99
              }
            },
            {
              "dimensions": {
                "date": "2022-09-19"
              },
              "sum": {
                "requests": 767
              },
              "uniq": {
                "uniques": 114
              }
            },
            {
              "dimensions": {
                "date": "2022-09-18"
              },
              "sum": {
                "requests": 914
              },
              "uniq": {
                "uniques": 129
              }
            },
            {
              "dimensions": {
                "date": "2022-09-17"
              },
              "sum": {
                "requests": 768
              },
              "uniq": {
                "uniques": 130
              }
            },
            {
              "dimensions": {
                "date": "2022-09-16"
              },
              "sum": {
                "requests": 717
              },
              "uniq": {
                "uniques": 129
              }
            },
            {
              "dimensions": {
                "date": "2022-09-15"
              },
              "sum": {
                "requests": 936
              },
              "uniq": {
                "uniques": 148
              }
            },
            {
              "dimensions": {
                "date": "2022-09-14"
              },
              "sum": {
                "requests": 705
              },
              "uniq": {
                "uniques": 162
              }
            },
            {
              "dimensions": {
                "date": "2022-09-13"
              },
              "sum": {
                "requests": 957
              },
              "uniq": {
                "uniques": 183
              }
            },
            {
              "dimensions": {
                "date": "2022-09-12"
              },
              "sum": {
                "requests": 747
              },
              "uniq": {
                "uniques": 138
              }
            },
            {
              "dimensions": {
                "date": "2022-09-11"
              },
              "sum": {
                "requests": 703
              },
              "uniq": {
                "uniques": 127
              }
            },
            {
              "dimensions": {
                "date": "2022-09-10"
              },
              "sum": {
                "requests": 721
              },
              "uniq": {
                "uniques": 174
              }
            },
            {
              "dimensions": {
                "date": "2022-09-09"
              },
              "sum": {
                "requests": 1237
              },
              "uniq": {
                "uniques": 189
              }
            },
            {
              "dimensions": {
                "date": "2022-09-08"
              },
              "sum": {
                "requests": 1137
              },
              "uniq": {
                "uniques": 127
              }
            },
            {
              "dimensions": {
                "date": "2022-09-07"
              },
              "sum": {
                "requests": 2031
              },
              "uniq": {
                "uniques": 132
              }
            },
            {
              "dimensions": {
                "date": "2022-09-06"
              },
              "sum": {
                "requests": 2562
              },
              "uniq": {
                "uniques": 159
              }
            },
            {
              "dimensions": {
                "date": "2022-09-05"
              },
              "sum": {
                "requests": 2013
              },
              "uniq": {
                "uniques": 159
              }
            },
            {
              "dimensions": {
                "date": "2022-09-04"
              },
              "sum": {
                "requests": 599
              },
              "uniq": {
                "uniques": 115
              }
            },
            {
              "dimensions": {
                "date": "2022-09-03"
              },
              "sum": {
                "requests": 661
              },
              "uniq": {
                "uniques": 107
              }
            },
            {
              "dimensions": {
                "date": "2022-09-02"
              },
              "sum": {
                "requests": 1139
              },
              "uniq": {
                "uniques": 130
              }
            },
            {
              "dimensions": {
                "date": "2022-09-01"
              },
              "sum": {
                "requests": 987
              },
              "uniq": {
                "uniques": 167
              }
            },
            {
              "dimensions": {
                "date": "2022-08-31"
              },
              "sum": {
                "requests": 840
              },
              "uniq": {
                "uniques": 119
              }
            },
            {
              "dimensions": {
                "date": "2022-08-30"
              },
              "sum": {
                "requests": 1084
              },
              "uniq": {
                "uniques": 128
              }
            },
            {
              "dimensions": {
                "date": "2022-08-29"
              },
              "sum": {
                "requests": 730
              },
              "uniq": {
                "uniques": 95
              }
            },
            {
              "dimensions": {
                "date": "2022-08-28"
              },
              "sum": {
                "requests": 492
              },
              "uniq": {
                "uniques": 102
              }
            },
            {
              "dimensions": {
                "date": "2022-08-27"
              },
              "sum": {
                "requests": 909
              },
              "uniq": {
                "uniques": 102
              }
            }
          ],
          "week": [
            {
              "dimensions": {
                "date": "2022-09-25"
              },
              "sum": {
                "requests": 9
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "date": "2022-09-24"
              },
              "sum": {
                "requests": 911
              },
              "uniq": {
                "uniques": 170
              }
            },
            {
              "dimensions": {
                "date": "2022-09-23"
              },
              "sum": {
                "requests": 810
              },
              "uniq": {
                "uniques": 137
              }
            },
            {
              "dimensions": {
                "date": "2022-09-22"
              },
              "sum": {
                "requests": 888
              },
              "uniq": {
                "uniques": 134
              }
            },
            {
              "dimensions": {
                "date": "2022-09-21"
              },
              "sum": {
                "requests": 742
              },
              "uniq": {
                "uniques": 89
              }
            },
            {
              "dimensions": {
                "date": "2022-09-20"
              },
              "sum": {
                "requests": 750
              },
              "uniq": {
                "uniques": 99
              }
            },
            {
              "dimensions": {
                "date": "2022-09-19"
              },
              "sum": {
                "requests": 767
              },
              "uniq": {
                "uniques": 114
              }
            }
          ]
        }
      ]
    }
  },
  "errors": null
}
```

</details>

## Go

Use [`taskfile`](https://taskfile.dev/) to build (and eventually upload binaries to lambda executables):

```
task
```

The following executables will require the following environment variables (and possibly some AWS credentials to upload the JSON from graphql API to S3):

```dotenv
CLOUDFLARE_EMAIL=
CLOUDFLARE_TOKEN=
CLOUDFLARE_ZONE=
```

For testing purposes you could use:

`test.sh`:

```bash
export CLOUDFLARE_ZONE="zone"
export CLOUDFLARE_EMAIL="email"
export CLOUDFLARE_TOKEN="token"

# iteration one: test.sh > day.json
./bin/analytics1d | jq .

# iteration two: test.sh > week.json
./bin/analytics1w | jq .

# iteration three: test.sh > month.json
./bin/analytics1m | jq .
```

<details>
<summary>day.json</summary>

```json
{
  "data": {
    "viewer": {
      "analytics": [
        {
          "day": [
            {
              "dimensions": {
                "datetime": "2022-09-25T20:00:00Z"
              },
              "sum": {
                "requests": 6
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T19:00:00Z"
              },
              "sum": {
                "requests": 58
              },
              "uniq": {
                "uniques": 16
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T18:00:00Z"
              },
              "sum": {
                "requests": 15
              },
              "uniq": {
                "uniques": 4
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T17:00:00Z"
              },
              "sum": {
                "requests": 47
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T16:00:00Z"
              },
              "sum": {
                "requests": 16
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T15:00:00Z"
              },
              "sum": {
                "requests": 46
              },
              "uniq": {
                "uniques": 9
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T14:00:00Z"
              },
              "sum": {
                "requests": 36
              },
              "uniq": {
                "uniques": 6
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T13:00:00Z"
              },
              "sum": {
                "requests": 107
              },
              "uniq": {
                "uniques": 6
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T12:00:00Z"
              },
              "sum": {
                "requests": 13
              },
              "uniq": {
                "uniques": 2
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T11:00:00Z"
              },
              "sum": {
                "requests": 36
              },
              "uniq": {
                "uniques": 9
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T10:00:00Z"
              },
              "sum": {
                "requests": 51
              },
              "uniq": {
                "uniques": 11
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T09:00:00Z"
              },
              "sum": {
                "requests": 20
              },
              "uniq": {
                "uniques": 8
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T08:00:00Z"
              },
              "sum": {
                "requests": 56
              },
              "uniq": {
                "uniques": 10
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T07:00:00Z"
              },
              "sum": {
                "requests": 19
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T06:00:00Z"
              },
              "sum": {
                "requests": 27
              },
              "uniq": {
                "uniques": 9
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T05:00:00Z"
              },
              "sum": {
                "requests": 23
              },
              "uniq": {
                "uniques": 11
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T04:00:00Z"
              },
              "sum": {
                "requests": 19
              },
              "uniq": {
                "uniques": 6
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T03:00:00Z"
              },
              "sum": {
                "requests": 25
              },
              "uniq": {
                "uniques": 11
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T02:00:00Z"
              },
              "sum": {
                "requests": 69
              },
              "uniq": {
                "uniques": 12
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T01:00:00Z"
              },
              "sum": {
                "requests": 44
              },
              "uniq": {
                "uniques": 9
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-25T00:00:00Z"
              },
              "sum": {
                "requests": 59
              },
              "uniq": {
                "uniques": 12
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T23:00:00Z"
              },
              "sum": {
                "requests": 20
              },
              "uniq": {
                "uniques": 7
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T22:00:00Z"
              },
              "sum": {
                "requests": 17
              },
              "uniq": {
                "uniques": 5
              }
            },
            {
              "dimensions": {
                "datetime": "2022-09-24T21:00:00Z"
              },
              "sum": {
                "requests": 34
              },
              "uniq": {
                "uniques": 14
              }
            }
          ]
        }
      ]
    }
  },
  "errors": null
}
```

</details>

<details>
<summary>week.json</summary>

```json
{
  "data": {
    "viewer": {
      "analytics": [
        {
          "week": [
            {
              "dimensions": {
                "date": "2022-09-25"
              },
              "sum": {
                "requests": 828
              },
              "uniq": {
                "uniques": 109
              }
            },
            {
              "dimensions": {
                "date": "2022-09-24"
              },
              "sum": {
                "requests": 911
              },
              "uniq": {
                "uniques": 170
              }
            },
            {
              "dimensions": {
                "date": "2022-09-23"
              },
              "sum": {
                "requests": 810
              },
              "uniq": {
                "uniques": 137
              }
            },
            {
              "dimensions": {
                "date": "2022-09-22"
              },
              "sum": {
                "requests": 888
              },
              "uniq": {
                "uniques": 134
              }
            },
            {
              "dimensions": {
                "date": "2022-09-21"
              },
              "sum": {
                "requests": 742
              },
              "uniq": {
                "uniques": 89
              }
            },
            {
              "dimensions": {
                "date": "2022-09-20"
              },
              "sum": {
                "requests": 750
              },
              "uniq": {
                "uniques": 99
              }
            },
            {
              "dimensions": {
                "date": "2022-09-19"
              },
              "sum": {
                "requests": 767
              },
              "uniq": {
                "uniques": 114
              }
            }
          ]
        }
      ]
    }
  },
  "errors": null
}
```

</details>

<details>
<summary>month.json</summary>

```json
{
  "data": {
    "viewer": {
      "analytics": [
        {
          "month": [
            {
              "dimensions": {
                "date": "2022-09-25"
              },
              "sum": {
                "requests": 828
              },
              "uniq": {
                "uniques": 109
              }
            },
            {
              "dimensions": {
                "date": "2022-09-24"
              },
              "sum": {
                "requests": 911
              },
              "uniq": {
                "uniques": 170
              }
            },
            {
              "dimensions": {
                "date": "2022-09-23"
              },
              "sum": {
                "requests": 810
              },
              "uniq": {
                "uniques": 137
              }
            },
            {
              "dimensions": {
                "date": "2022-09-22"
              },
              "sum": {
                "requests": 888
              },
              "uniq": {
                "uniques": 134
              }
            },
            {
              "dimensions": {
                "date": "2022-09-21"
              },
              "sum": {
                "requests": 742
              },
              "uniq": {
                "uniques": 89
              }
            },
            {
              "dimensions": {
                "date": "2022-09-20"
              },
              "sum": {
                "requests": 750
              },
              "uniq": {
                "uniques": 99
              }
            },
            {
              "dimensions": {
                "date": "2022-09-19"
              },
              "sum": {
                "requests": 767
              },
              "uniq": {
                "uniques": 114
              }
            },
            {
              "dimensions": {
                "date": "2022-09-18"
              },
              "sum": {
                "requests": 914
              },
              "uniq": {
                "uniques": 129
              }
            },
            {
              "dimensions": {
                "date": "2022-09-17"
              },
              "sum": {
                "requests": 768
              },
              "uniq": {
                "uniques": 130
              }
            },
            {
              "dimensions": {
                "date": "2022-09-16"
              },
              "sum": {
                "requests": 717
              },
              "uniq": {
                "uniques": 129
              }
            },
            {
              "dimensions": {
                "date": "2022-09-15"
              },
              "sum": {
                "requests": 936
              },
              "uniq": {
                "uniques": 148
              }
            },
            {
              "dimensions": {
                "date": "2022-09-14"
              },
              "sum": {
                "requests": 705
              },
              "uniq": {
                "uniques": 162
              }
            },
            {
              "dimensions": {
                "date": "2022-09-13"
              },
              "sum": {
                "requests": 957
              },
              "uniq": {
                "uniques": 183
              }
            },
            {
              "dimensions": {
                "date": "2022-09-12"
              },
              "sum": {
                "requests": 747
              },
              "uniq": {
                "uniques": 138
              }
            },
            {
              "dimensions": {
                "date": "2022-09-11"
              },
              "sum": {
                "requests": 703
              },
              "uniq": {
                "uniques": 127
              }
            },
            {
              "dimensions": {
                "date": "2022-09-10"
              },
              "sum": {
                "requests": 721
              },
              "uniq": {
                "uniques": 174
              }
            },
            {
              "dimensions": {
                "date": "2022-09-09"
              },
              "sum": {
                "requests": 1237
              },
              "uniq": {
                "uniques": 189
              }
            },
            {
              "dimensions": {
                "date": "2022-09-08"
              },
              "sum": {
                "requests": 1137
              },
              "uniq": {
                "uniques": 127
              }
            },
            {
              "dimensions": {
                "date": "2022-09-07"
              },
              "sum": {
                "requests": 2031
              },
              "uniq": {
                "uniques": 132
              }
            },
            {
              "dimensions": {
                "date": "2022-09-06"
              },
              "sum": {
                "requests": 2562
              },
              "uniq": {
                "uniques": 159
              }
            },
            {
              "dimensions": {
                "date": "2022-09-05"
              },
              "sum": {
                "requests": 2013
              },
              "uniq": {
                "uniques": 159
              }
            },
            {
              "dimensions": {
                "date": "2022-09-04"
              },
              "sum": {
                "requests": 599
              },
              "uniq": {
                "uniques": 115
              }
            },
            {
              "dimensions": {
                "date": "2022-09-03"
              },
              "sum": {
                "requests": 661
              },
              "uniq": {
                "uniques": 107
              }
            },
            {
              "dimensions": {
                "date": "2022-09-02"
              },
              "sum": {
                "requests": 1139
              },
              "uniq": {
                "uniques": 130
              }
            },
            {
              "dimensions": {
                "date": "2022-09-01"
              },
              "sum": {
                "requests": 987
              },
              "uniq": {
                "uniques": 167
              }
            },
            {
              "dimensions": {
                "date": "2022-08-31"
              },
              "sum": {
                "requests": 840
              },
              "uniq": {
                "uniques": 119
              }
            },
            {
              "dimensions": {
                "date": "2022-08-30"
              },
              "sum": {
                "requests": 1084
              },
              "uniq": {
                "uniques": 128
              }
            },
            {
              "dimensions": {
                "date": "2022-08-29"
              },
              "sum": {
                "requests": 730
              },
              "uniq": {
                "uniques": 95
              }
            },
            {
              "dimensions": {
                "date": "2022-08-28"
              },
              "sum": {
                "requests": 492
              },
              "uniq": {
                "uniques": 102
              }
            },
            {
              "dimensions": {
                "date": "2022-08-27"
              },
              "sum": {
                "requests": 909
              },
              "uniq": {
                "uniques": 102
              }
            }
          ]
        }
      ]
    }
  },
  "errors": null
}
```

</details>

TODO:

* go script to upload json to s3
* chartjs pages