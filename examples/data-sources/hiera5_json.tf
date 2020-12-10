data "hiera5_json" "aws_tags" {
  key = "aws_tags"
}

locals {
  aws_tags = jsondecode(data.hiera5_json.aws_tags.value)
}