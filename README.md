# A Bank FinTech Mobile Wallet: SRE & Cloud-Native Architecture Showcase

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Container Standard](https://img.shields.io/badge/Container-Distroless%20%2F%20Non--Root-2496ED?style=flat-square&logo=docker)](Dockerfile)
[![Kubernetes](https://img.shields.io/badge/Orchestration-Helm%20v3-326CE5?style=flat-square&logo=helm)](charts/wallet-service/)
[![AWS Compute](https://img.shields.io/badge/Compute-AWS%20ECS%20Fargate%20(Firecracker)-FF9900?style=flat-square&logo=amazon-aws)](aws/)
[![Compliance](https://img.shields.io/badge/Compliance-PCI--DSS%20%7C%20ISO%2027001-success?style=flat-square)](ManagedServices.md)

---

## 1. Executive Summary

This repository presents the **Production Reference Architecture** for **A Bank's Mobile Wallet Transaction Platform**. It resolves the core tension in modern banking infrastructure: balancing **sovereign data compliance and cost predictability** (on-premise / Kubernetes) with **hyper-elastic bursting and reduced compliance blast radius** (cloud serverless).

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│              DUAL-RUN FINTECH ARCHITECTURE: HYBRID CLOUD & ON-PREM               │
├──────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│   [ Client Mobile App ] ───(HTTPS / mTLS / X-Idempotency-Key)                    │
│                                      │                                           │
│                                      ▼                                           │
│                     [ AWS Application Load Balancer / API Gateway ]              │
│                                      │                                           │
│         ┌────────────────────────────┴─────────────────────────────┐             │
│         ▼                                                          ▼             │
│  [ AWS ECS Fargate ]                                       [ Kubernetes / K3s ]  │
│  (Cloud Serverless Tier)                                   (On-Prem Parity Tier) │
│  • MicroVM Jails (Firecracker)                             • Production Helm     │
│  • Golang Wallet Microservice                              • HPA Autoscale       │
│  • Zero Host OS Patching                                   • Non-Root Security   │
│         │                                                          │             │
│         ├────────────────────────────┬─────────────────────────────┘             │
│         ▼                            ▼                                           │
│  [ Amazon DynamoDB ]         [ AWS Lambda Engine ]                               │
│  (Idempotency & Session)     (Asynchronous Webhooks)                             │
│  • X-Idempotency-Key Locks   • MPU Switch Callbacks                              │
│  • 24-Hour Atomic TTL        • SMS/OTP Notification Triggers                     │
│  • Pay-Per-Request ($0 idle) • Zero Persistent Server Overhead                   │
│                                                                                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. The Architectural Dilemma: Bare Metal vs AWS Fargate

| Architectural Vector | Self-Hosted Bare Metal / K3s | AWS ECS Fargate (Serverless) |
| :--- | :--- | :--- |
| **PCI-DSS Compliance** | SRE team owns host Linux kernel CVE patching, SSH key rotation, and hypervisor audits. | Host OS compliance offloaded to AWS's PCI Level 1 Attestation of Compliance (AoC). |
| **Process Isolation** | Containers share host Linux kernel (vulnerable to kernel zero-day escapes). | **Firecracker MicroVMs**: Every pod runs in a dedicated hardware KVM jail. |
| **Payday Spikes (25k TPS)**| Sized for peak; **85% capacity sits idle** 26 days a month. | Autoscales in seconds; scales to zero/baseline after promotional surges. |
| **Steady-State Cost** | 35% - 50% cheaper for 24/7/365 sustained compute. | Higher unit cost per vCPU-hour if kept static indefinitely. |
| **Synthesis** | **Core General Ledger** runs on sovereign, private infrastructure. | **Mobile APIs & Promo Bursts** run on AWS Fargate + Lambda. |

---

## 3. Core Engine: Golang Double-Entry Ledger (`wallet-service`)

The transactional core (`cmd/wallet-service/main.go`) is engineered in **Golang** for sub-millisecond execution and deterministic garbage collection.

### 3.1 Mathematical Balance Conservation
Every fund movement enforces strict atomic double-entry bookkeeping:
$$
\sum \Delta \text{Balance}_{\text{Debit}} + \sum \Delta \text{Balance}_{\text{Credit}} = 0
$$

### 3.2 Network Retry Defense (`X-Idempotency-Key`)
In mobile banking environments with intermittent cellular connectivity, clients frequently retry requests. Without idempotency, users suffer duplicate debiting.
- Every mutating request requires an `X-Idempotency-Key` HTTP header.
- Cached results are replayed instantly (`idempotent_replay: true`) without re-executing ledger writes.
- Payload divergence returns `409 Conflict`.

### 3.3 Endpoints
- `POST /api/v1/wallets/transfer` — Executes atomic transfer between accounts.
- `GET /api/v1/wallets/{account_id}/balance` — Retrieves real-time balance.
- `GET /healthz` — Liveness probe (HTTP 200).
- `GET /readyz` — Readiness probe (validates database connection pool).
- `GET /metrics` — Exposes transaction volume and internal SRE telemetry.

---

## 4. Container Hardening (Distroless Non-Root)

```dockerfile
# Hardened multi-stage scratch container (<15MB)
FROM scratch
COPY --from=builder /wallet-service /wallet-service
USER 10001:10001
EXPOSE 8080
ENTRYPOINT ["/wallet-service"]
```
- **Zero Shell**: No `/bin/sh` or `/bin/bash` in the runtime container (eliminates remote shell breakouts).
- **Read-Only Root Filesystem**: Prevents runtime binary tampering.
- **Unprivileged Execution**: Runs as UID `10001`.

---

## 5. Enterprise Kubernetes Packaging (`charts/wallet-service`)

Eliminates flat YAML configuration sprawl through parameterized Helm templates:
- `values.yaml` — Configurable replica counts, image tags, and resource quotas.
- `templates/deployment.yaml` — Restricts pod security capabilities (`drop: [ALL]`).
- `templates/hpa.yaml` — Horizontal Pod Autoscaler targeting 70% CPU / 80% Memory thresholds.
- `templates/service.yaml` — ClusterIP internal service abstraction.

---

## 6. Cloud-Native Serverless Tier (`aws/`)

- `fargate-task-definition.json` — ECS task definition utilizing AWS Fargate with CloudWatch log shipping.
- `dynamodb-idempotency-table.json` — Serverless DynamoDB table configured for Pay-Per-Request (zero cost when idle) with a 24-hour TTL attribute.
- `lambda-webhook-handler.py` — Serverless event-driven receiver validating MPU (Myanmar Payment Union) payment callback signatures.
- `deploy.sh` & `teardown.sh` — Automated deployment and 100% clean teardown scripts.

---

## 7. Comparative Case Study: How KBZPay Operates on Huawei Cloud

| Banking Function | KBZPay Architecture (Huawei Cloud) | A Bank Blueprint (AWS Hybrid) |
| :--- | :--- | :--- |
| **Core Microservices** | Huawei CCE (Cloud Container Engine) | Kubernetes (On-Prem) / EKS |
| **Serverless Bursting** | Huawei CCI (Cloud Container Instance) | **AWS ECS Fargate** |
| **Asynchronous Webhooks**| Huawei FunctionGraph | **AWS Lambda** |
| **Primary Ledger Database**| Huawei GaussDB (Distributed Multi-AZ) | Amazon Aurora / Distributed SQL |
| **Idempotency Cache** | Distributed In-Memory Cache | **Amazon DynamoDB (On-Demand)** |

---

## 8. Author & Audit Trail
- **Author**: Thaw Zin @ Chris (`pretamane`)
- **Target Enterprise**: A Bank (Mobile Wallet Division)
- **Master Documentation Bible**: [`ManagedServices.md`](../ManagedServices.md)
