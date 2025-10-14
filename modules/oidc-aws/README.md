# Terraform Module: AWS OIDC Role for GitHub Actions

This module creates the necessary AWS IAM resources to allow GitHub Actions to securely authenticate with AWS using OpenID Connect (OIDC). It establishes a trust relationship between your AWS account and a specific GitHub repository, enabling passwordless, short-lived credential generation for your CI/CD workflows.

This pattern is a foundational element of the "no long-lived infrastructure" philosophy, as it eliminates the need to store static `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` credentials in GitHub Secrets.

## Core Components

1.  **IAM OIDC Identity Provider:** Configures AWS to trust GitHub's OIDC provider (`token.actions.githubusercontent.com`).
2.  **IAM Role:** An assumable role that GitHub Actions workflows can use.
3.  **IAM Trust Policy:** A policy attached to the role that strictly scopes which GitHub repository and branch are allowed to assume it.

## Security

The trust policy is the most critical security component. It uses OIDC claims to ensure that only workflows from the specified repository and branch can assume the role.

-   `token.actions.githubusercontent.com:aud`: Verifies the token was intended for AWS (`sts.amazonaws.com`).
-   `token.actions.githubusercontent.com:sub`: Verifies the subject (the source of the request), locking it down to a specific repository and branch (`repo:owner/repo-name:ref:refs/heads/branch-name`).

This module only configures the **trust policy** (who can assume the role). It is your responsibility to attach a separate, minimal **permissions policy** that defines what the role is allowed to do (e.g., `s3:GetObject`, `ec2:DescribeInstances`).

## Usage

Add the following to your Terraform configuration to create the role. For deterministic builds, always pin the module `source` to a specific version tag.

```terraform
module "github_oidc_role" {
  source = "github.com/pwjohnston2002/secrets-workflow//modules/oidc-aws?ref=v0.2.0" # Pin to a specific version

  github_owner = "my-github-org"
  github_repo  = "my-awesome-repo"
  branch_name  = "main"
}

output "github_oidc_role_arn" {
  description = "The ARN of the IAM role for GitHub Actions to assume."
  value       = module.github_oidc_role.role_arn
}
```

After applying this Terraform, store the output `github_oidc_role_arn` as a secret in your GitHub repository (e.g., `AWS_ROLE_TO_ASSUME`) to be used by your workflows.

### Validation Example

A complete, working configuration that uses this module can be found in the `examples/oidc-aws-setup` directory. This example is used for automated testing in CI.

## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> terraform | >= 1.0 |
| <a name="requirement_aws"></a> aws | ~> 5.0 |
| <a name="requirement_tls"></a> tls | ~> 4.0 |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_branch_name"></a> branch\_name | The branch name to scope the OIDC trust policy to. | `string` | `"main"` | no |
| <a name="input_github_owner"></a> github\_owner | The GitHub organization or user that owns the repository. | `string` | n/a | yes |
| <a name="input_github_repo"></a> github\_repo | The name of the GitHub repository. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_role_arn"></a> role\_arn | The ARN of the created IAM role, which can be assumed by GitHub Actions using OIDC. |
