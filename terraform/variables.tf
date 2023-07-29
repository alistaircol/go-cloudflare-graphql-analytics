variable "aws_profile" {
  type    = string
  default = "analytics"
}

variable "aws_region" {
  type    = string
  default = "eu-west-2"
}

variable "aws_lambda_runtime" {
  type    = string
  default = "go2.x"
}