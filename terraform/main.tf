provider "aws" {
  region = var.aws_region
}

locals {
  function_name = "dynamo-test"
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

resource "aws_lambda_function" "dynamo-test" {
  function_name = local.function_name
  package_type = "Image"
  timeout = 3

  role = aws_iam_role.iam_for_lambda.arn

  image_uri = "${var.image_uri}:${var.tag_name}"
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
