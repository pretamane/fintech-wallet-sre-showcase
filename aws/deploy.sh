#!/usr/bin/env bash
# ==============================================================================
# A Bank Mobile Wallet - Automated Cloud-Native Provisioning
# Architecture: AWS ECR + DynamoDB + ECS Fargate (Serverless Tier)
# Cost Profile: 100% On-Demand / Free-Tier (< $0.05 total run)
# ==============================================================================
set -euo pipefail

AWS_REGION="us-east-1"
REPO_NAME="a-bank-wallet-service"
TABLE_NAME="a-bank-wallet-idempotency"

echo "===> [STEP 1/3]: Provisioning Private AWS ECR Repository..."
if ! aws ecr describe-repositories --repository-names "${REPO_NAME}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    aws ecr create-repository \
        --repository-name "${REPO_NAME}" \
        --image-scanning-configuration scanOnPush=true \
        --region "${AWS_REGION}"
    echo "ECR Repository '${REPO_NAME}' created successfully."
else
    echo "ECR Repository '${REPO_NAME}' already exists."
fi

echo "===> [STEP 2/3]: Provisioning DynamoDB Idempotency Table (On-Demand)..."
if ! aws dynamodb describe-table --table-name "${TABLE_NAME}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    aws dynamodb create-table \
        --cli-input-json file://$(dirname "$0")/dynamodb-idempotency-table.json \
        --region "${AWS_REGION}"
    echo "DynamoDB table '${TABLE_NAME}' created. Awaiting ACTIVE state..."
    aws dynamodb wait table-exists --table-name "${TABLE_NAME}" --region "${AWS_REGION}"
else
    echo "DynamoDB table '${TABLE_NAME}' already exists."
fi

echo "===> [STEP 3/3]: Registering ECS Fargate Task Definition..."
aws ecs register-task-definition \
    --cli-input-json file://$(dirname "$0")/fargate-task-definition.json \
    --region "${AWS_REGION}"

echo "=============================================================================="
echo "DEPLOYMENT COMPLETED SUCCESSFULLY"
echo "ECR URI: 464868388812.dkr.ecr.${AWS_REGION}.amazonaws.com/${REPO_NAME}"
echo "=============================================================================="
