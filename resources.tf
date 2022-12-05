resource "aws_cloudwatch_event_rule" "analytics_hourly" {
  name                = "blog-analytics-hourly"
  description         = "Blog analytics: hourly"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_cloudwatch_event_rule" "analytics_daily" {
  name                = "blog-analytics-daily"
  description         = "Blog analytics: daily"
  schedule_expression = "cron(0 * * * ? *)"
}

# create a terraform script to make three go AWS lambda functions
# create a terraform script to invoke an AWS lambda function using event bridge on a schedule
resource "aws_lambda_function" "analytics_1d" {
  function_name = "analytics_1d"
  runtime       = var.aws_lambda_runtime
  handler       = "main"
  role          = aws_iam_role.lambda_role.arn
}

resource "aws_lambda_function" "analytics_1w" {
  function_name = "analytics_1w"
  runtime       = var.aws_lambda_runtime
  handler       = "main"
  role          = aws_iam_role.lambda_role.arn
}

resource "aws_lambda_function" "analytics_1m" {
  function_name = "analytics_1m"
  runtime       = var.aws_lambda_runtime
  handler       = "main"
  role          = aws_iam_role.lambda_role.arn
}

# Create an IAM role for the Lambda functions
resource "aws_iam_role" "lambda_role" {
  name               = "lambda_role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

resource "aws_s3_bucket" "analytics" {
  bucket = "analytics"
}

# Attach the AWSLambdaBasicExecutionRole policy to the IAM role
resource "aws_iam_policy" "lambda_policy" {
  name   = "lambda_policy"
  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:*"
      ],
      "Effect": "Allow",
      "Resource": "arn:aws:logs:*:*:*"
    },
    {
      "Sid": "WritePermission",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:PutObjectAcl"
      ],
      "Resource": [
        "${aws_s3_bucket.analytics.arn}/*"
      ]
    }
  ]
}
EOF
}

resource "aws_iam_policy_attachment" "lambda_policy_attachment" {
  name       = "lambda_policy_attachment"
  roles      = [aws_iam_role.lambda_role.name]
  policy_arn = aws_iam_policy.lambda_policy.arn
}

# Create a target for the EventBridge rule
resource "aws_cloudwatch_event_target" "analytics_1d_target_hourly" {
  rule      = aws_cloudwatch_event_rule.analytics_hourly
  target_id = "lambda_target"
  arn       = aws_lambda_function.analytics_1d.arn
}

resource "aws_cloudwatch_event_target" "analytics_1d_target_daily" {
  rule      = aws_cloudwatch_event_rule.analytics_daily
  target_id = "lambda_target"
  arn       = aws_lambda_function.analytics_1d.arn
}

resource "aws_cloudwatch_event_target" "analytics_1w_target_daily" {
  rule      = aws_cloudwatch_event_rule.analytics_daily
  target_id = "lambda_target"
  arn       = aws_lambda_function.analytics_1w.arn
}

resource "aws_cloudwatch_event_target" "analytics_1m_target_daily" {
  rule      = aws_cloudwatch_event_rule.analytics_daily
  target_id = "lambda_target"
  arn       = aws_lambda_function.analytics_1m.arn
}