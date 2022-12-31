provider "aws" { //Pin version!
  region  = "eu-west-1"
  profile = "default"
}

provider "aws" {
  alias  = "us-east-1"
  region = "us-east-1"
}

terraform {
  required_version = "1.3.6"
  backend "s3" {
    bucket = "dani-terraform-states"
    key    = "mdscanner.dev/mdscanner.dev.tfstate"
    region = "eu-west-1"
  }
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.48.0"
    }
  }
}
