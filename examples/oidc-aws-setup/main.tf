terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

module "oidc_aws_role" {
  source = "../../modules/oidc-aws"

  github_owner = var.github_owner
  github_repo  = var.github_repo
  branch_name  = var.branch_name
}