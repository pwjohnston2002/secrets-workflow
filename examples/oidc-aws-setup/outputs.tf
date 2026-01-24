output "example_github_oidc_role_arn" {
  description = "The ARN of the IAM role for GitHub Actions to assume."
  value       = module.github_oidc_role.role_arn
}
