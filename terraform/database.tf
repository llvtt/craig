resource "aws_dynamodb_table" "searches" {
  name           = "searches"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "query"

  attribute {
    name = "query"
    type = "S"
  }
}

resource "aws_dynamodb_table" "items" {
  name           = "items"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "url"

  attribute {
    name = "url"
    type = "S"
  }
}

resource "aws_dynamodb_table" "price_logs" {
  name           = "price_logs"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "item_url"

  attribute {
    name = "item_url"
    type = "S"
  }
}

data "aws_iam_policy_document" "dynamodb_policy_document" {
  statement {
    actions = [
      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]

    resources = [
      aws_dynamodb_table.searches.arn,
      aws_dynamodb_table.items.arn,
      aws_dynamodb_table.price_logs.arn
    ]
  }
}

resource "aws_iam_policy" "interact_with_dynamodb" {
  name = "dynamodb-policy"
  description = "grant privileges to interact with Craig dynamodb tables"
  policy = data.aws_iam_policy_document.dynamodb_policy_document.json
}

resource "aws_iam_role_policy_attachment" "dynamodb_attachment" {
  role = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.interact_with_dynamodb.arn
}
