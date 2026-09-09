terraform {
  required_version = ">= 1.5.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.40"
    }
    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 4.25"
    }
  }

  # Production SRE Remote State Backend (S3 + DynamoDB Locking)
  # Uncomment and configure with your organization's backend bucket:
  # backend "s3" {
  #   bucket         = "a-bank-sre-terraform-state-464868388812"
  #   key            = "mobile-wallet/production/terraform.tfstate"
  #   region         = "us-east-1"
  #   dynamodb_table = "a-bank-terraform-locks"
  #   encrypt        = true
  # }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = "a-bank-mobile-wallet"
      Environment = var.environment
      ManagedBy   = "Terraform"
      Compliance  = "PCI-DSS-v4.0"
      CostCenter  = "Digital-Banking-SRE"
    }
  }
}

provider "cloudflare" {
  email   = var.cloudflare_email
  api_key = var.cloudflare_api_key
}
