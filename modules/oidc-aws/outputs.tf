output "role_arn" {
  description = "The ARN of the created IAM role for GitHub OIDC."
  value       = aws_iam_role.github_oidc_role.arn
}
