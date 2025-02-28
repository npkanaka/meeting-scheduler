# variables.tf - Essential variables with no defaults for sensitive data

variable "aws_region" {
  description = "The AWS region to deploy to"
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "The name of the project"
  type        = string
  default     = "meeting-scheduler"
}

variable "db_username" {
  description = "The username for the database"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "The password for the database"
  type        = string
  sensitive   = true
  # No default - must be provided by user
}

variable "db_name" {
  description = "The name of the database"
  type        = string
  default     = "scheduler"
}

variable "container_image" {
  description = "The container image to deploy"
  type        = string
  # No default - must be provided by user
}