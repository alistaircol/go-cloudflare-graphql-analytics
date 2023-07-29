resource "aws_iam_role" "lambda_role" {
  name_prefix        = "lambda-ac93uk-cloudflare-analytics-"
  path               = "/"
  assume_role_policy = data.aws_iam_policy_document.lambda_policy.json
}

resource "aws_iam_role_policy_attachment" "lambda_execution" {
  role       = aws_iam_role.lambda_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

data "aws_iam_policy_document" "lambda_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type = "Service"
      identifiers = [
        "lambda.amazonaws.com",
        "s3.amazonaws.com",
      ]
    }
  }
}

data "aws_iam_policy_document" "lambda_data_bucket" {
  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
      "s3:PutObject",
      "s3:ListBucket",
    ]

    resources = [
      "arn:aws:s3:::${var.bucket_name}",
      "arn:aws:s3:::${var.bucket_name}/*",
    ]
  }
}

resource "aws_iam_policy" "lambda_data_bucket" {
  name   = "ac93uk_cloudflare_analytics_lambda_allow_access_to_data_bucket"
  path   = "/"
  policy = data.aws_iam_policy_document.lambda_data_bucket.json
}

resource "aws_iam_role_policy_attachment" "lambda_data_bucket" {
  policy_arn = aws_iam_policy.lambda_data_bucket.arn
  role       = aws_iam_role.lambda_role.name
}
