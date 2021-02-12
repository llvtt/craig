// slack-events responds to events driven from Slack through API Gateway
resource "aws_lambda_function" "slack-events" {
  function_name = "slack-events"
  timeout       = 5
  runtime       = "go1.x"
  handler       = "slack-events"

  role = aws_iam_role.iam_for_lambda.arn

  image_uri = "${var.image_uri}:${var.tag_name}"

  environment {
    variables = {
      "SLACK_ACCESS_TOKEN" : var.slack_access_token,
      "SLACK_SIGNING_SECRET" : var.slack_signing_secret,
    }
  }
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
