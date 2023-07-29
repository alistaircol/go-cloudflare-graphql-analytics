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
  default = "go1.x"
}

variable "bucket_name" {}
variable "code_bucket_name" {}

variable "cloudflare_zone" {}
variable "cloudflare_email" {}
variable "cloudflare_token" {}


variable "allowed_bucket_ingress_sources" {
  type = list(object({
    name      = string
    addresses = set(string)
  }))
}
