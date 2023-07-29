module "s3_lambda_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = var.code_bucket_name
  acl    = "private"

  // https://github.com/terraform-aws-modules/terraform-aws-s3-bucket/issues/223#issuecomment-1545649581
  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = true
  }
}

module "s3_data_bucket" {
  source = "terraform-aws-modules/s3-bucket/aws"
  bucket = var.bucket_name
  acl    = "private"

  // https://github.com/terraform-aws-modules/terraform-aws-s3-bucket/issues/223#issuecomment-1545649581
  control_object_ownership = true
  object_ownership         = "ObjectWriter"

  versioning = {
    enabled = false
  }
}

resource "aws_s3_bucket_policy" "allow_access" {
  bucket = module.s3_data_bucket.s3_bucket_id
  policy = data.aws_iam_policy_document.allow_access.json
}

# https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html#Conditions_IPAddress
# https://discuss.hashicorp.com/t/true-false-condition-not-working-in-for-each-for-iam-role-condition/40541
# https://developer.hashicorp.com/terraform/language/expressions/dynamic-blocks
data "aws_iam_policy_document" "allow_access" {
  dynamic "statement" {
    for_each = var.allowed_bucket_ingress_sources

    content {
      sid = "AccessPolicy${statement.value["name"]}"

      principals {
        type        = "AWS"
        identifiers = ["*"]
      }

      actions = [
        "s3:GetObject",
        "s3:GetObjectVersion",
      ]

      resources = [
        module.s3_data_bucket.s3_bucket_arn,
        "${module.s3_data_bucket.s3_bucket_arn}/*",
      ]

      condition {
        test     = "IpAddress"
        values   = statement.value["addresses"]
        variable = "aws:SourceIp"
      }
    }
  }
}
