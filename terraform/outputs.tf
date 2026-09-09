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

output "architecture_summary" {
  description = "Summary of provisioned FinTech infrastructure"
  value = {
    cloud_tier      = "AWS ECS Fargate Serverless"
    database_tier   = "DynamoDB On-Demand (ACID Idempotency)"
    edge_tier       = "Cloudflare Anycast WAF + Origin Port 8080 Rewrite"
    compliance_spec = "PCI-DSS v4.0 Non-Root Distroless Runtime"
  }
}
