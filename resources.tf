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
