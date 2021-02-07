resource "aws_apigatewayv2_api" "api" {
  name          = "craig-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_route" "list_searches" {
  api_id = aws_apigatewayv2_api.api.id
  route_key = "POST /slack/events"
  target = "integrations/${aws_apigatewayv2_integration.slack-events.id}"
}

resource "aws_apigatewayv2_integration" "slack-events" {
  api_id = aws_apigatewayv2_api.api.id
  integration_type = "AWS_PROXY"
  integration_method = "POST"
  integration_uri = aws_lambda_function.slack-events.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_lambda_permission" "main" {
  action = "lambda:InvokeFunction"
  function_name = aws_lambda_function.slack-events.function_name
  principal = "apigateway.amazonaws.com"
  source_arn = "${aws_apigatewayv2_api.api.execution_arn}/*/*"
}

resource "aws_apigatewayv2_stage" "development" {
  api_id = aws_apigatewayv2_api.api.id
  name   = "development"
  auto_deploy = true

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.access_logs.arn
    format = jsonencode({ "requestId":"$context.requestId", "ip": "$context.identity.sourceIp", "caller":"$context.identity.caller", "user":"$context.identity.user", "requestTime":"$context.requestTime", "httpMethod":"$context.httpMethod", "status":"$context.status", "protocol":"$context.protocol", "responseLength":"$context.responseLength"})
  }
}

resource "aws_cloudwatch_log_group" "access_logs" {
  name              = "${aws_apigatewayv2_api.api.name}-access-logs"
  retention_in_days = 1
}
