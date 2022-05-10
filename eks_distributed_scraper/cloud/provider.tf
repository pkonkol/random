provider "aws" {
  region = var.AWS_REGION

  default_tags {
    tags = var.DEFAULT_TAGS
  }
}