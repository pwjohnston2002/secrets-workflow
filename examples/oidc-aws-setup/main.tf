terraform {
  required_version = ">= 1.6.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.region
}

module "github_oidc_role" {
  source = "../../modules/oidc-aws" # Adjust if your module path differs

  github_owner = var.github_owner
  github_repo  = var.github_repo
  branch_name  = var.branch_name
}

output "role_arn" {
  description = "IAM role ARN assumed by GitHub Actions via OIDC."
  value       = module.github_oidc_role.role_arn
}