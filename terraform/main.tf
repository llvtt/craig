provider "aws" {
  region = "${var.aws_region}"
}

provider "archive" {}

data "archive_file" "zip" {
  type        = "zip"
  source_file = "craig"
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
  assume_role_policy = "${data.aws_iam_policy_document.policy.json}"
}

resource "aws_lapnmbda_function" "lambda" {
  function_name = "craig"

  filename         = "${data.archive_file.zip.output_path}"
  source_code_hash = "${data.archive_file.zip.output_base64sha256}"

  role = "${aws_iam_role.iam_for_lambda.arn}"

  # TODO: ensure the name is correct
  # rename craig package to "core" or "craig-core"
  handler = "craig"
  runtime = "go1.x"

  environment {
    variables = {
      CRAIG_SLACK_ENDPOINT = "${var.slack_endpoint}"
    }
  }
}
