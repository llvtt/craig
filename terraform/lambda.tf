resource "aws_lambda_function" "scrape_craig_lambda" {
  function_name = "scrape-craiglist"
  package_type  = "Image"
  timeout       = 300
  image_config {
    entry_point = ["/cloudwatch-events"]
  }

  role = aws_iam_role.iam_for_lambda.arn

  image_uri = "${var.image_uri}:${var.tag_name}"

  environment {
    variables = {
      "SLACK_ACCESS_TOKEN" : var.slack_access_token,
      "SLACK_SIGNING_SECRET" : var.slack_signing_secret,
    }
  }

  depends_on = [aws_iam_role_policy_attachment.lambda_logs, aws_cloudwatch_log_group.allow_cloudwatch]
}

resource "aws_lambda_function" "slack-events" {
  function_name = "slack-events"
  package_type  = "Image"
  timeout       = 3
  image_config {
    entry_point = ["/slack-events"]
  }

  role = aws_iam_role.iam_for_lambda.arn

  image_uri = "${var.image_uri}:${var.tag_name}"

  environment {
    variables = {
      "SLACK_ACCESS_TOKEN" : var.slack_access_token,
      "SLACK_SIGNING_SECRET" : var.slack_signing_secret,
    }
  }
}

