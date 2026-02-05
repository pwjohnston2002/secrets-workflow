variable "region" {
  type    = string
  default = "us-east-1"
}

variable "run_id" {
  type        = string
  description = "Unique identifier for this test run to avoid collisions"
}