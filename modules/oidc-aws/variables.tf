variable "github_owner" {
  description = "The GitHub organization or user that owns the repository."
  type        = string
}

variable "github_repo" {
  description = "The name of the GitHub repository."
  type        = string
}

variable "branch_name" {
  description = "The branch name to scope the OIDC trust policy to."
  type        = string
  default     = "main"
}
