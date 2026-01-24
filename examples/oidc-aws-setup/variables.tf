variable "github_owner" { 
  type = string 
}

variable "github_repo"  { 
  type = string 
}

variable "branch_name"  { 
  type = string
  default = "main" 
}

variable "region"  { 
  type = string
  default = "us-east-1" 
}