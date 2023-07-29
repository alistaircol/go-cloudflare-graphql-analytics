data "aws_s3_object" "lambda_executable" {
  bucket = module.s3_lambda_bucket.s3_bucket_id
  key    = "build.zip"
}

resource "aws_lambda_function" "analytics" {
  function_name = "ac93uk-cloudflare-analytics"
  role          = aws_iam_role.lambda_role.arn
  handler       = "main"
  runtime       = "go1.x"

  s3_bucket         = module.s3_lambda_bucket.s3_bucket_id
  s3_key            = data.aws_s3_object.lambda_executable.key
  s3_object_version = data.aws_s3_object.lambda_executable.version_id

  tracing_config {
    mode = "Active"
  }

  environment {
    variables = {
      CLOUDFLARE_ZONE  = var.cloudflare_zone
      CLOUDFLARE_EMAIL = var.cloudflare_email
      CLOUDFLARE_TOKEN = var.cloudflare_token
      AWS_S3_BUCKET    = module.s3_data_bucket.s3_bucket_id
    }
  }

  lifecycle {
    ignore_changes = [s3_object_version]
  }

  depends_on = [
    data.aws_s3_object.lambda_executable
  ]
}
