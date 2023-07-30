resource "aws_cloudwatch_event_rule" "daily_on_hourly_schedule" {
  name                = "ac93uk-cloudflare-analytics-daily-on-hourly-schedule"
  description         = "Dispatch event to run daily cloudflare analytics every hour"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_cloudwatch_event_rule" "daily_on_daily_schedule" {
  name                = "ac93uk-cloudflare-analytics-daily-on-daily-schedule"
  description         = "Dispatch event to run daily cloudflare analytics every day"
  schedule_expression = "cron(0 * * * ? *)"
}

resource "aws_cloudwatch_event_target" "daily_on_hourly_schedule" {
  rule  = aws_cloudwatch_event_rule.daily_on_hourly_schedule.name
  arn   = aws_lambda_function.analytics.arn
  input = jsonencode({
    Detail: {
      period : "d"
    }
  })
}

resource "aws_cloudwatch_event_target" "daily_on_daily_schedule" {
  rule  = aws_cloudwatch_event_rule.daily_on_daily_schedule.name
  arn   = aws_lambda_function.analytics.arn
  input = jsonencode({
    Detail: {
      period : "d"
    }
  })
}
