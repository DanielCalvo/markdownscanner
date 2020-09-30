data "terraform_remote_state" "dcalvo-dev-zone" {
  backend = "s3"
  config = {
    bucket = "dani-terraform-states"
    key    = "dcalvo.dev/dns-domain.tfstate"
    region = "eu-west-1"
  }
}

module "s3_website" {
  source = "github.com/DanielCalvo/studies/Projects/dcalvo.dev/terraform_modules/s3_website"
  domain-name = var.domain-name
  bucket-name = var.bucket-name
}