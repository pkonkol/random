variable "AWS_REGION" {
  default = "eu-west-1"
}

variable "AWS_DEFAULT_ZONE" {
  default = "eu-west-1b"
}

variable "AWS_SECONDARY_ZONE" {
  default = "eu-west-1c"
}

variable "DEFAULT_TAGS" {
  type = map(string)
  default = {
    env = "dev"
    app = "scraper"
  }
}