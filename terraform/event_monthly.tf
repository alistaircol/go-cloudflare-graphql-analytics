resource "aws_cloudwatch_event_rule" "monthly_on_daily_schedule" {
  name                = "ac93uk-cloudflare-analytics-monthly-on-daily-schedule"
  description         = "Dispatch event to run monthly cloudflare analytics every day"
  schedule_expression = "cron(0 * * * ? *)"
}

resource "aws_cloudwatch_event_target" "monthly_on_daily_schedule" {
  rule = aws_cloudwatch_event_rule.monthly_on_daily_schedule.name
  arn  = aws_lambda_function.analytics.arn

  input = jsonencode({
    period : "m"
  })
}
