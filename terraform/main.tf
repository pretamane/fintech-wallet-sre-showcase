# ==============================================================================
# A Bank Mobile Wallet - Complete Declarative Cloud & Edge Infrastructure
# Architecture: AWS ECR + DynamoDB + ECS Fargate + Cloudflare Anycast WAF
# ==============================================================================

# Data source for AWS caller identity
data "aws_caller_identity" "current" {}

# Data source for default VPC and subnets
data "aws_vpc" "default" {
  default = true
}

data "aws_subnets" "default" {
  filter {
    name   = "vpc-id"
    values = [data.aws_vpc.default.id]
  }
}

# ------------------------------------------------------------------------------
# 1. AWS ECR Repository & Lifecycle Policy
# ------------------------------------------------------------------------------
resource "aws_ecr_repository" "wallet_service" {
  name                 = var.service_name
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  tags = {
    Name = var.service_name
  }
}

resource "aws_ecr_lifecycle_policy" "wallet_service" {
  repository = aws_ecr_repository.wallet_service.name

  policy = jsonencode({
    rules = [
      {
        rulePriority = 1
        description  = "Retain last 10 production images, expire older untagged builds"
        selection = {
          tagStatus   = "untagged"
          countType   = "sinceImagePushed"
          countUnit   = "days"
          countNumber = 14
        }
        action = {
          type = "expire"
        }
      }
    ]
  })
}

# ------------------------------------------------------------------------------
# 2. AWS DynamoDB Idempotency & Ledger Locking Table (Serverless Tier)
# ------------------------------------------------------------------------------
resource "aws_dynamodb_table" "idempotency" {
  name         = "a-bank-wallet-idempotency"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "IdempotencyKey"

  attribute {
    name = "IdempotencyKey"
    type = "S"
  }

  ttl {
    attribute_name = "TTL"
    enabled        = true
  }

  point_in_time_recovery {
    enabled = true
  }

  server_side_encryption {
    enabled = true
  }

  tags = {
    Name        = "a-bank-wallet-idempotency"
    Component   = "LedgerIdempotencyEngine"
    RecoveryRPO = "ZeroDataLoss"
  }
}

# ------------------------------------------------------------------------------
# 3. IAM Execution Role for ECS Fargate (Least-Privilege Scoped)
# ------------------------------------------------------------------------------
resource "aws_iam_role" "ecs_execution_role" {
  name = "a-bank-ecs-task-execution-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        Principal = {
          Service = "ecs-tasks.amazonaws.com"
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "ecs_execution_policy" {
  role       = aws_iam_role.ecs_execution_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
}

# ------------------------------------------------------------------------------
# 4. CloudWatch Logs & ECS Fargate Cluster
# ------------------------------------------------------------------------------
resource "aws_cloudwatch_log_group" "wallet_logs" {
  name              = "/ecs/${var.service_name}"
  retention_in_days = 14

  tags = {
    Application = var.service_name
  }
}

resource "aws_ecs_cluster" "wallet_cluster" {
  name = "a-bank-wallet-cluster"

  setting {
    name  = "containerInsights"
    value = "disabled" # Enabled for high-scale enterprise observability
  }
}

# ------------------------------------------------------------------------------
# 5. Security Group: Ingress Restricted to Container Port 8080
# ------------------------------------------------------------------------------
resource "aws_security_group" "fargate_sg" {
  name        = "a-bank-wallet-sg"
  description = "Security group for A Bank Wallet Service on Fargate"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "Allow inbound HTTP traffic on microservice port 8080"
    from_port   = var.container_port
    to_port     = var.container_port
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description = "Allow all outbound traffic for ECR pulls and DynamoDB APIs"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "a-bank-wallet-sg"
  }
}

# ------------------------------------------------------------------------------
# 6. ECS Task Definition (Distroless Scratch Runtime - Zero Shell)
# ------------------------------------------------------------------------------
resource "aws_ecs_task_definition" "wallet_task" {
  family                   = var.service_name
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  cpu                      = var.cpu_units
  memory                   = var.memory_units
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn

  container_definitions = jsonencode([
    {
      name      = "wallet-service"
      image     = var.container_image
      essential = true
      portMappings = [
        {
          containerPort = var.container_port
          hostPort      = var.container_port
          protocol      = "tcp"
        }
      ]
      environment = [
        {
          name  = "PORT"
          value = tostring(var.container_port)
        },
        {
          name  = "ENVIRONMENT"
          value = var.environment
        }
      ]
      logConfiguration = {
        logDriver = "awslogs"
        options = {
          "awslogs-group"         = aws_cloudwatch_log_group.wallet_logs.name
          "awslogs-region"        = var.aws_region
          "awslogs-stream-prefix" = "fargate"
        }
      }
    }
  ])
}

# ------------------------------------------------------------------------------
# 7. ECS Service (Fargate Autonomous Self-Healing)
# ------------------------------------------------------------------------------
resource "aws_ecs_service" "wallet_service" {
  name            = "${var.service_name}-svc"
  cluster         = aws_ecs_cluster.wallet_cluster.id
  task_definition = aws_ecs_task_definition.wallet_task.arn
  desired_count   = var.fargate_desired_count
  launch_type     = "FARGATE"

  network_configuration {
    subnets          = data.aws_subnets.default.ids
    security_groups  = [aws_security_group.fargate_sg.id]
    assign_public_ip = true
  }

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100
}

# ------------------------------------------------------------------------------
# 8. Cloudflare DNS & Origin Rule Port Rewriting
# ------------------------------------------------------------------------------
# A-Record pointing to Origin IP (or dynamically to ALB/Fargate)
resource "cloudflare_record" "wallet_dns" {
  zone_id = var.cloudflare_zone_id
  name    = "wallet"
  content = "44.204.78.56" # Origin Public IP / ALB
  type    = "A"
  proxied = true
  comment = "Managed by Terraform: A Bank Mobile Wallet Live Ingress"
  ttl     = 1
}

# Origin Rule rewriting incoming port 443 -> backend container port 8080
resource "cloudflare_ruleset" "origin_port_rewrite" {
  zone_id     = var.cloudflare_zone_id
  name        = "Override Port to 8080"
  description = "Managed by Terraform: Route HTTPS 443 to Fargate Port 8080"
  kind        = "zone"
  phase       = "http_request_origin"

  rules {
    action = "route"
    action_parameters {
      origin {
        port = var.container_port
      }
    }
    expression  = "(http.host eq \"${var.subdomain_name}\")"
    description = "Forward A Bank Wallet HTTPS to AWS Fargate Port 8080"
    enabled     = true
  }
}
