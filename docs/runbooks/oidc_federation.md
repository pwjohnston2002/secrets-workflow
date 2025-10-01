# Runbook: OIDC Federation with AWS

This runbook explains the concept, setup, and security principles of using OpenID Connect (OIDC) to allow GitHub Actions to securely authenticate with AWS without long-lived credentials.

## 1. The Concept: "Sign in with GitHub" for AWS

The simplest way to understand OIDC federation is to think of a "Sign in with Google" button on a website. You aren't creating a new password for that site; you are trusting Google to vouch for your identity.

In our case:
- **GitHub Actions** is the user trying to sign in.
- **AWS** is the service that requires authentication.
- **OIDC** is the "Sign in with..." protocol that connects them.

Instead of storing a static `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY` in GitHub Secrets, the workflow asks GitHub for a special, short-lived OIDC token. It then presents this token to AWS, which verifies its authenticity and exchanges it for temporary AWS credentials.

This eliminates the risk of long-lived keys being compromised.

## 2. How It Works

The authentication flow follows these steps:

1.  A GitHub Actions workflow job starts.
2.  The job requests a unique OIDC token from GitHub's identity provider. This token contains claims about the job, such as the repository and branch it's running from (e.g., `repo:your-org/your-repo:ref:refs/heads/main`).
3.  The workflow presents this token to the AWS Security Token Service (STS).
4.  AWS STS validates the token against the OIDC provider it has been configured to trust (GitHub's).
5.  AWS STS checks if the claims in the token (like the repository name) match the conditions defined in a special IAM Role's "Trust Policy."
6.  If the trust policy is satisfied, AWS STS grants the workflow temporary, short-lived credentials by allowing it to assume the IAM Role.
7.  The workflow uses these temporary credentials to interact with AWS resources. The credentials automatically expire after a set duration.

## 3. Implementation Steps

There are two parts to implementing this pattern: creating the trusted role in AWS and configuring the workflow to use it.

### Step 1: Create the IAM Role in AWS

This repository provides a reusable Terraform module to create the necessary AWS resources. Use the `modules/oidc-aws` module in your own Terraform configuration.

**Example Usage:**
```terraform
# In your project's Terraform code

module "github_oidc_role" {
  source = "github.com/your-org/secrets-workflow//modules/oidc-aws?ref=v1.0.0"

  github_owner = "your-org"
  github_repo  = "your-repo-name"
  branch_name  = "main" # Or another branch like 'develop'
}

output "github_oidc_role_arn" {
  value = module.github_oidc_role.role_arn
}
```

After running `terraform apply`, the output will give you the ARN of the role you need for the next step.

### Step 2: Configure the GitHub Actions Workflow

In your workflow file, you must grant it permission to request an OIDC token and configure the `aws-actions/configure-aws-credentials` action.

```yaml
jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      id-token: write   # Required to request the OIDC token
      contents: read    # Required to checkout the repository
    steps:
      - name: Configure AWS Credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE_TO_ASSUME }} # The ARN from the Terraform output
          aws-region: ${{ secrets.AWS_REGION }}
```

Store the role ARN from the Terraform output as a repository secret (e.g., `AWS_ROLE_TO_ASSUME`).

## 4. Troubleshooting Common Errors

*   **Error: `Not authorized to perform sts:AssumeRoleWithWebIdentity`**: This is the most common error. It means the claims in the OIDC token from GitHub did not match the trust policy on your IAM role. Double-check the `repo:` and `ref:` values in your IAM role's trust policy to ensure they exactly match the repository and branch you are running the workflow from.
*   **Error: `Error: The id-token could not be retrieved`**: You are missing the `permissions: id-token: write` block in your workflow file. The job does not have permission to request a token from GitHub.
