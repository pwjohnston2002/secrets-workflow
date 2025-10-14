output "oidc_role_arn" {
  description = "The ARN of the IAM role created by the oidc-aws module."
  value       = module.oidc_aws_role.role_arn
}
