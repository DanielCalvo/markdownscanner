data "terraform_remote_state" "dcalvo-dev-zone" {
  backend = "s3"
  config = {
    bucket = "dani-terraform-states"
    key    = "dcalvo.dev/dns-domain.tfstate"
    region = "eu-west-1"
  }
}