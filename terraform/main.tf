provider "aws" {
  region = var.aws_region
}

data "aws_iam_policy_document" "policy" {
  statement {
    sid    = ""
    effect = "Allow"

    principals {
      identifiers = ["lambda.amazonaws.com"]
      type        = "Service"
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.policy.json
}

data "aws_ecr_image" "image" {
  repository_name = "craig"
  image_tag       = var.tag_name
}

// TODO: extract lambda, apigateway integration into its own module
// Consider doing away with docker to reduce deploy time
resource "aws_cloudwatch_log_group" "allow_cloudwatch" {
  name              = "/aws/lambda/${aws_lambda_function.slack-events.function_name}"
  retention_in_days = 1
}

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole
resource "aws_iam_policy" "lambda_logging" {
  name        = "lambda_logging"
  path        = "/"
  description = "IAM policy for logging from a lambda"

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:*:*:*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy_attachment" "lambda_logs" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.lambda_logging.arn
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.scrape_craig_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.scrape_craigslist_trigger_rule.arn
}

resource "aws_cloudwatch_event_rule" "scrape_craigslist_trigger_rule" {
  name                = "ScrapeCraigslistTriggerRule"
  description         = "Cron schedule to make craig scrape craigslist"
  schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "scrape_craigslist_trigger" {
  arn  = aws_lambda_function.scrape_craig_lambda.arn
  rule = aws_cloudwatch_event_rule.scrape_craigslist_trigger_rule.name
}

