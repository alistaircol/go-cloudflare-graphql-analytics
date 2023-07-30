> :warning: **Work in progress** but mostly complete and hopefully correct.

![Chart](https://raw.githubusercontent.com/alistaircol/go-cloudflare-graphql-analytics/main/.github/chart.png)

# Overview

A `go` app which runs in a [AWS Lambda](https://aws.amazon.com/lambda/) on a schedule with [Amazon EventBridge](https://aws.amazon.com/eventbridge/) to invoke the lambda which will query Cloudflare analytics for a domain zone and save the payloads in an [S3](https://aws.amazon.com/s3/) bucket for use in [chartjs](https://www.chartjs.org/).

## Local Development

```bash
cp .env.example .env
```

Set `CLOUDFLARE_*` and `AWS_S3_BUCKET` for local development. These values are read by `task` and are used to build local infrastructure with `localstack`.

```bash
docker compose up -d
```

Start localstack.

```bash
task go:build
```

Builds the go binary, and creates an archive for creating lambda function later.

```bash
task localstack:setup
```

Will create the bucket, lambda, event bus, copy resources, etc. to the buckets. Run `task` to see all available options.

```bash
task event:day
```

Will dispatch an event with the `Detail` to pass to the lambda, to tell it to query analytics for the previous day. See `task` for the other event types.

```bash
task s3:list
```

You should then see files have been created by the lambda.

You can then view the chart at, e.g. `http://analytics-bucket.s3.localhost.localstack.cloud:4566/d.html`

## Production Deployment

See `terraform/README.md`.
