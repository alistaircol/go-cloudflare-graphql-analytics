resource "aws_cloudwatch_event_rule" "weekly_on_daily_schedule" {
  name                = "ac93uk-cloudflare-analytics-weekly-on-daily-schedule"
  description         = "Dispatch event to run weekly cloudflare analytics every day"
  schedule_expression = "cron(0 * * * ? *)"
}

resource "aws_cloudwatch_event_target" "weekly_on_daily_schedule" {
  rule = aws_cloudwatch_event_rule.weekly_on_daily_schedule.name
  arn  = aws_lambda_function.analytics.arn

  input = jsonencode({
    period : "w"
  })
}
