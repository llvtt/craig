provider "aws" {
  region = var.aws_region
}

provider "archive" {}

locals {
  craig_binary = "../../../main/lambda/main"
  function_name = "craig"
}

data "archive_file" "zip" {
  type        = "zip"
  source_file = local.craig_binary
  output_path = "craig.zip"
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

resource "aws_lambda_function" "craig_lambda" {
  function_name = local.function_name
  handler       = "main"

  filename         = data.archive_file.zip.output_path
  source_code_hash = data.archive_file.zip.output_base64sha256

  role = aws_iam_role.iam_for_lambda.arn

  runtime = "go1.x"

  environment {
    variables = {
      CRAIG_SLACK_ENDPOINT = var.slack_endpoint
    }
  }

  depends_on = [aws_iam_role_policy_attachment.lambda_logs, aws_cloudwatch_log_group.allow_cloudwatch]
}

resource "aws_lambda_permission" "allow_cloudwatch" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.craig_lambda.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.scrape_craigslist_trigger_rule.arn
}

resource "aws_cloudwatch_event_rule" "scrape_craigslist_trigger_rule" {
  name                = "ScrapeCraigslistTriggerRule"
  description         = "Cron schedule to make craig scrape craigslist"
  schedule_expression = "rate(1 hour)"
}

resource "aws_cloudwatch_event_target" "scrape_craigslist_trigger" {
  arn  = aws_lambda_function.craig_lambda.arn
  rule = aws_cloudwatch_event_rule.scrape_craigslist_trigger_rule.name
}

resource "aws_cloudwatch_log_group" "allow_cloudwatch" {
  name              = "/aws/lambda/${local.function_name}"
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