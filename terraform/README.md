# Terraform

## Preamble

Will create all relevant resources to invoke lambda on various schedules to query Cloudflare analytics and dump the payload in a bucket.

## Starting

```bash
cp terraform.tfvars.example terraform.tfvars
```

Populate the `terraform.tfvars` as you see fit.

## AWS Configuration

Will need to have configuration details for the profile set in `terraform.tfvars`.

To set up this profile:

```text
$ aws configure --profile="analytics"
AWS Access Key ID [None]: abc123
AWS Secret Access Key [None]: xyx789
Default region name [None]: eu-west-2
Default output format [None]: json
```

You should create a programmatic user with the relevant roles to create buckets, policies, event bridge rules/targets, etc.

## Creating Infrastructure

> **Note**
> The `terraform apply` will need to be run atleast twice.
> It will fail when creating the lambda function.
> You will need to upload the zipped binary to the code bucket, and then re-run the `terraform apply`.
> It will not be able to create the lambda function because the object will not yet exist in the bucket.
> e.g.:
> `Error: getting S3 Bucket (lambda-ac93uk-cloudflare-analytics) Object (build.zip): NotFound: Not Found`

From the root of the repository:

```bash
task terraform:init
task terraform:plan
task terraform:apply
```

Upon first `task terraform:apply`, it will fail. Upload the zipped binary to the bucket.

```bash
# ensure you are in the root of this repository
task go:build
task terraform:upload_to_code_bucket
```

Then continue with the rest of the infrastructure:

```bash
task terraform:apply
```
