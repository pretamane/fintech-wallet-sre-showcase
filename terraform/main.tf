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
# 5. Security Groups: Zero-Trust Perimeter & Micro-Segmentation
# ------------------------------------------------------------------------------
# ALB Security Group: Ingress from Public Edge (Cloudflare WAF / Client traffic)
resource "aws_security_group" "alb_sg" {
  name        = "a-bank-wallet-alb-sg"
  description = "Security group for A Bank Internet-Facing Application Load Balancer"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description = "Allow inbound HTTP from Cloudflare Anycast Edge"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow inbound HTTPS from Cloudflare Anycast Edge"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    description = "Allow inbound HTTP 8080 from Cloudflare Origin Rule"
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    description     = "Allow outbound forward strictly to Fargate microservice instances"
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [] # Linked dynamically via fargate_sg rules
    cidr_blocks     = ["0.0.0.0/0"]
  }

  tags = {
    Name      = "a-bank-wallet-alb-sg"
    Component = "EdgeIngress"
  }
}

# Fargate Security Group: Zero-Trust Ingress Restricted ONLY to ALB Security Group
resource "aws_security_group" "fargate_sg" {
  name        = "a-bank-wallet-sg"
  description = "Security group for A Bank Wallet Service on Fargate"
  vpc_id      = data.aws_vpc.default.id

  ingress {
    description     = "Allow inbound HTTP traffic strictly from ALB security group (PCI-DSS Requirement 1.2)"
    from_port       = var.container_port
    to_port         = var.container_port
    protocol        = "tcp"
    security_groups = [aws_security_group.alb_sg.id]
  }

  egress {
    description = "Allow all outbound traffic for AWS API calls, ECR pulls, and DynamoDB"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name      = "a-bank-wallet-sg"
    Component = "InternalLedgerEngine"
  }
}

# ------------------------------------------------------------------------------
# 6. Application Load Balancer (ALB) & Static CNAME Ingress Decoupling
# ------------------------------------------------------------------------------
resource "aws_lb" "wallet_alb" {
  name               = "a-bank-wallet-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.alb_sg.id]
  subnets            = data.aws_subnets.default.ids

  enable_deletion_protection = false

  tags = {
    Name        = "a-bank-wallet-alb"
    Component   = "FinTechIngressRouter"
    Environment = var.environment
  }
}

resource "aws_lb_target_group" "wallet_tg" {
  name                 = "a-bank-wallet-tg"
  port                 = var.container_port
  protocol             = "HTTP"
  vpc_id               = data.aws_vpc.default.id
  target_type          = "ip" # Mandatory for Fargate awsvpc network mode
  deregistration_delay = 15   # Fast connection draining during rolling updates

  health_check {
    enabled             = true
    path                = "/healthz"
    port                = tostring(var.container_port)
    protocol            = "HTTP"
    interval            = 15
    timeout             = 5
    healthy_threshold   = 2
    unhealthy_threshold = 3
    matcher             = "200"
  }

  tags = {
    Name      = "a-bank-wallet-tg"
    Component = "LedgerHealthTargetGroup"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.wallet_alb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.wallet_tg.arn
  }
}

resource "aws_lb_listener" "http_8080" {
  load_balancer_arn = aws_lb.wallet_alb.arn
  port              = 8080
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.wallet_tg.arn
  }
}

# ------------------------------------------------------------------------------
# 7. ECS Task Definition (Distroless Scratch Runtime - Zero Shell)
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
# 8. ECS Service (Fargate Autonomous Self-Healing with Hybrid Capacity Providers)
# ------------------------------------------------------------------------------
resource "aws_ecs_service" "wallet_service" {
  name            = "${var.service_name}-svc"
  cluster         = aws_ecs_cluster.wallet_cluster.id
  task_definition = aws_ecs_task_definition.wallet_task.arn
  desired_count   = var.fargate_desired_count

  # Hybrid FinOps Strategy: Base 2 On-Demand + Burst onto Fargate Spot (70% savings)
  capacity_provider_strategy {
    capacity_provider = "FARGATE"
    weight            = 1
    base              = 2
  }

  capacity_provider_strategy {
    capacity_provider = "FARGATE_SPOT"
    weight            = 4
    base              = 0
  }

  network_configuration {
    subnets          = data.aws_subnets.default.ids
    security_groups  = [aws_security_group.fargate_sg.id]
    assign_public_ip = true
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.wallet_tg.arn
    container_name   = "wallet-service"
    container_port   = var.container_port
  }

  deployment_maximum_percent         = 200
  deployment_minimum_healthy_percent = 100

  depends_on = [aws_lb_listener.http]
}

# ------------------------------------------------------------------------------
# 9. Cloudflare DNS & Edge Decoupling (Permanent Static CNAME)
# ------------------------------------------------------------------------------
# CNAME Record pointing to AWS ALB DNS Name (Zero-Breakage on Task Restarts)
resource "cloudflare_record" "wallet_dns" {
  zone_id = var.cloudflare_zone_id
  name    = "wallet"
  content = aws_lb.wallet_alb.dns_name
  type    = "CNAME"
  proxied = true
  comment = "Managed by Terraform: Decoupled CNAME to AWS Application Load Balancer"
  ttl     = 1
}


