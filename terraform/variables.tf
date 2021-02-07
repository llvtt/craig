variable "aws_region" {
  description = "The AWS region"
  default = "us-east-1"
}

variable "image_uri" {
  description = "URI for the image in ECR to deploy to lambda"
}

variable "tag_name" {
  description = "The docker image tag"
}
