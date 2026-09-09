output "ecr_repository_url" {
  description = "Amazon ECR Repository URI"
  value       = aws_ecr_repository.wallet_service.repository_url
}

output "dynamodb_table_arn" {
  description = "DynamoDB Idempotency Table ARN"
  value       = aws_dynamodb_table.idempotency.arn
}

output "ecs_cluster_name" {
  description = "AWS ECS Fargate Cluster Name"
  value       = aws_ecs_cluster.wallet_cluster.name
}

output "ecs_service_name" {
  description = "AWS ECS Service Name"
  value       = aws_ecs_service.wallet_service.name
}

output "cloudflare_ingress_url" {
  description = "Public Production Endpoint (Cloudflare Proxied HTTPS)"
  value       = "https://${var.subdomain_name}/healthz"
}

output "alb_dns_name" {
  description = "AWS Application Load Balancer Public Regional DNS Name"
  value       = aws_lb.wallet_alb.dns_name
}

output "alb_arn" {
  description = "AWS Application Load Balancer ARN"
  value       = aws_lb.wallet_alb.arn
}

output "target_group_arn" {
  description = "AWS ALB Target Group ARN (IP Target Type for Fargate awsvpc)"
  value       = aws_lb_target_group.wallet_tg.arn
}

output "architecture_summary" {
  description = "Summary of provisioned FinTech infrastructure"
  value = {
    cloud_tier      = "AWS ECS Fargate Hybrid (Base On-Demand + Burst Spot)"
    load_balancer   = "AWS Application Load Balancer (Dual-AZ Decoupled Ingress)"
    database_tier   = "DynamoDB On-Demand (ACID Idempotency Engine)"
    edge_tier       = "Cloudflare Anycast WAF + Static CNAME Decoupling"
    compliance_spec = "PCI-DSS v4.0 Network Micro-Segmentation & Non-Root Scratch"
  }
}

