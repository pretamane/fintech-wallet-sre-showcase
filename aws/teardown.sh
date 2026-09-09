#!/usr/bin/env bash
# ==============================================================================
# A Bank Mobile Wallet - Complete Resource Dismantle Script
# Purpose: Eradicates all live AWS cloud resources to enforce $0.00 zero-budget.
# ==============================================================================
set -euo pipefail

AWS_REGION="us-east-1"
REPO_NAME="a-bank-wallet-service"
TABLE_NAME="a-bank-wallet-idempotency"

echo "===> [TEARDOWN 1/3]: Deleting DynamoDB Idempotency Table..."
if aws dynamodb describe-table --table-name "${TABLE_NAME}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    aws dynamodb delete-table --table-name "${TABLE_NAME}" --region "${AWS_REGION}"
    echo "DynamoDB table '${TABLE_NAME}' deleted."
else
    echo "DynamoDB table '${TABLE_NAME}' does not exist."
fi

echo "===> [TEARDOWN 2/3]: Purging and Deleting ECR Repository..."
if aws ecr describe-repositories --repository-names "${REPO_NAME}" --region "${AWS_REGION}" >/dev/null 2>&1; then
    aws ecr delete-repository --repository-name "${REPO_NAME}" --force --region "${AWS_REGION}"
    echo "ECR repository '${REPO_NAME}' deleted."
else
    echo "ECR repository '${REPO_NAME}' does not exist."
fi

echo "===> [TEARDOWN 3/3]: Deregistering ECS Task Definitions..."
TASK_ARNS=$(aws ecs list-task-definitions --family-prefix "${REPO_NAME}" --region "${AWS_REGION}" --query 'taskDefinitionArns[]' --output text || true)
for ARN in $TASK_ARNS; do
    if [ -n "$ARN" ] && [ "$ARN" != "None" ]; then
        echo "Deregistering task definition: $ARN"
        aws ecs deregister-task-definition --task-definition "$ARN" --region "${AWS_REGION}" >/dev/null
    fi
done

echo "=============================================================================="
echo "TEARDOWN COMPLETE: Zero lingering cloud footprint. $0.00 spend preserved."
echo "=============================================================================="
