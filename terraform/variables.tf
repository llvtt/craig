variable "aws_region" {
  description = "The AWS region"
  default     = "us-east-1"
}

variable "image_uri" {
  description = "URI for the image in ECR to deploy to lambda"
}

variable "slack_access_token" {
  description = "Slack API token"
}

variable "slack_signing_secret" {
  description = "Slack signing secret"
}
