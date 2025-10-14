output "role_arn" {
  description = "The ARN of the created IAM role, which can be assumed by GitHub Actions using OIDC."
  value       = aws_iam_role.github_oidc_role.arn
}
