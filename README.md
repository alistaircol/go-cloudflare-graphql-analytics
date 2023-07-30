> :warning: **Work in progress** but mostly complete and hopefully correct.

![Chart](https://raw.githubusercontent.com/alistaircol/go-cloudflare-graphql-analytics/main/.github/chart.png)

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [Overview](#overview)
  - [Local Development](#local-development)
  - [Production Deployment](#production-deployment)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

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

---

<!-- task-start -->
<!-- DO NOT EDIT THIS SECTION - IT IS UPDATED BY A GITHUB WORKFLOW -->
```
task: Available tasks for this project:
* go:lint:                                         Format go files
* go:build:                                        Build lambda binary
* go:test:                                         Run go tests
* go:coverage:                                     Run go tests and generate coverage report
* localstack:shell:                                Shell inside localstack container
* localstack:setup:                                Setup local infrastructure (creates bucket and lambda function)
* localstack:create_bucket:                        Create a bucket which lambda downloads files to
* localstack:copy_resources_to_bucket:             Copy chartjs resources to bucket
* localstack:lambda_create_function:               Create lambda function
* localstack:update_lambda_function:               Update lambda code
* event:day:                                       Dispatch Eventbridge notification to invoke lambda to gather daily stats
* event:week:                                      Dispatch Eventbridge notification to invoke lambda to gather weekly stats
* event:month:                                     Dispatch Eventbridge notification to invoke lambda to gather monthly stats
* s3:list:                                         List objects in the bucket
* s3:get:                                          Get the payload file from the bucket (e.g. `task s3:get -- filename.json`)
* s3:delete:                                       Delete all objects in the bucket
* terraform:init:                                  Run `terraform init`
* terraform:fmt:                                   Run `terraform fmt`
* terraform:plan:                                  Run `terraform plan`
* terraform:apply:                                 Run `terraform apply`
* terraform:console:                               Run `terraform console`
* terraform:upload_to_code_bucket:                 Uploads zipped binary to the code bucket
* terraform:update_lambda_code:                    Upload new binary and re-fresh lambda
* terraform:upload_resources_to_data_bucket:       Upload chart resources to data bucket
```
<!-- task-end -->

## Production Deployment

See `terraform/README.md`.
