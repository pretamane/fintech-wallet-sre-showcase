#!/usr/bin/env bash
# ==============================================================================
# A BANK SRE PLATFORM • FINOPS COST GOVERNANCE & LIFECYCLE CONTROLLER
# ==============================================================================
# Provides idempotent commands to pause compute billing (scale to 0) or resume
# (scale to 1) for live interview demonstrations:
# Usage:
#   ./scripts/finops-lifecycle.sh status
#   ./scripts/finops-lifecycle.sh pause   # Scales ECS tasks to 0 (halts compute billing)
#   ./scripts/finops-lifecycle.sh resume  # Scales ECS tasks to 1 (restores in ~30s)
# ==============================================================================
set -euo pipefail

CLUSTER_NAME="a-bank-wallet-cluster"
SERVICE_NAME="a-bank-wallet-service-svc"
REGION="us-east-1"
TG_ARN="arn:aws:elasticloadbalancing:us-east-1:464868388812:targetgroup/a-bank-wallet-tg/dc1ef5a5b515c80f"

action="${1:-status}"

echo "┌──────────────────────────────────────────────────────────────────────────┐"
echo "│ A BANK FINOPS CONTROLLER • SRE LIFECYCLE & COST GOVERNANCE               │"
echo "├──────────────────────────────────────────────────────────────────────────┤"
echo "│ Cluster:   $CLUSTER_NAME (AWS $REGION)"
echo "│ Service:   $SERVICE_NAME"
echo "│ Action:    $(echo "$action" | tr '[:lower:]' '[:upper:]')"
echo "└──────────────────────────────────────────────────────────────────────────┘"
echo ""

case "$action" in
  status)
    echo "--- [1] ECS Service Compute State ---"
    aws ecs describe-services \
      --cluster "$CLUSTER_NAME" \
      --services "$SERVICE_NAME" \
      --region "$REGION" \
      --query "services[0].{Desired:desiredCount,Running:runningCount,Pending:pendingCount,TaskDef:taskDefinition,Status:status}" \
      --output table

    echo "--- [2] ALB Target Group Health ---"
    aws elbv2 describe-target-health \
      --target-group-arn "$TG_ARN" \
      --region "$REGION" \
      --query "TargetHealthDescriptions[*].{Target:Target.Id,Port:Target.Port,Health:TargetHealth.State,Reason:TargetHealth.Reason}" \
      --output table
    ;;

  pause)
    echo "Scaling ECS Fargate service tasks to 0 (Halting compute cost)..."
    aws ecs update-service \
      --cluster "$CLUSTER_NAME" \
      --service "$SERVICE_NAME" \
      --desired-count 0 \
      --region "$REGION" \
      --query "service.{Desired:desiredCount,Running:runningCount,Status:status}" \
      --output table
    echo "Compute paused. Configuration, ALB, KMS, WAF, DynamoDB, and SQS preserved."
    ;;

  resume)
    echo "Scaling ECS Fargate service tasks to 1 (Restoring live compute)..."
    aws ecs update-service \
      --cluster "$CLUSTER_NAME" \
      --service "$SERVICE_NAME" \
      --desired-count 1 \
      --region "$REGION" \
      --query "service.{Desired:desiredCount,Running:runningCount,Status:status}" \
      --output table
    echo "Waiting for task deployment to reach running state..."
    aws ecs wait services-stable --cluster "$CLUSTER_NAME" --services "$SERVICE_NAME" --region "$REGION"
    echo "Service is stable. Validating live health probe..."
    curl -sI "https://wallet.thaw-zin-2k77.de5.net/healthz" | head -n 1
    ;;

  *)
    echo "Unknown action: $action. Valid options: status | pause | resume"
    exit 1
    ;;
esac
