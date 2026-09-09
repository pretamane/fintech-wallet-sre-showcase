variable "aws_region" {
  type        = string
  description = "AWS Primary Production Region"
  default     = "us-east-1"
}

variable "environment" {
  type        = string
  description = "Deployment environment tier (production, staging, uat)"
  default     = "production"
}

variable "service_name" {
  type        = string
  description = "Name of the microservice"
  default     = "a-bank-wallet-service"
}

variable "container_image" {
  type        = string
  description = "Full ECR container image URI"
  default     = "464868388812.dkr.ecr.us-east-1.amazonaws.com/a-bank-wallet-service:1.0.0"
}

variable "container_port" {
  type        = number
  description = "TCP port exposed by the Go ledger service"
  default     = 8080
}

variable "cpu_units" {
  type        = string
  description = "ECS Fargate CPU units (256 = 0.25 vCPU)"
  default     = "256"
}

variable "memory_units" {
  type        = string
  description = "ECS Fargate memory allocation in MB"
  default     = "512"
}

variable "fargate_desired_count" {
  type        = number
  description = "Number of desired active Fargate task instances"
  default     = 1
}

variable "cloudflare_email" {
  type        = string
  description = "Cloudflare administrative account email"
  default     = "thawzin252467@gmail.com"
}

variable "cloudflare_api_key" {
  type        = string
  description = "Cloudflare Global API Key"
  sensitive   = true
}

variable "cloudflare_zone_id" {
  type        = string
  description = "Cloudflare DNS Zone ID for thaw-zin-2k77.de5.net"
  default     = "7e72740785ff0c8f5484443c8025de82"
}

variable "subdomain_name" {
  type        = string
  description = "Target public subdomain name"
  default     = "wallet.thaw-zin-2k77.de5.net"
}
