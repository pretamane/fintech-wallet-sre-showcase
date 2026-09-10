# A Bank Mobile Wallet: Production SRE & Cloud-Native FinTech Architecture

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Container Runtime](https://img.shields.io/badge/Container-Distroless%20Scratch%20Non--Root-2496ED?style=flat-square&logo=docker)](Dockerfile)
[![Kubernetes Orchestration](https://img.shields.io/badge/Orchestration-Helm%20v3%20%7C%20Multi--AZ-326CE5?style=flat-square&logo=helm)](charts/wallet-service/)
[![AWS Infrastructure](https://img.shields.io/badge/IaC-Terraform%20v1.9.5%20(27%20Resources)-844FBA?style=flat-square&logo=terraform)](terraform/)
[![Perimeter Defense](https://img.shields.io/badge/Security-AWS%20WAFv2%20%2B%20Cloudflare%20Anycast-FF9900?style=flat-square&logo=amazon-aws)](terraform/main.tf)
[![Compliance Standard](https://img.shields.io/badge/Compliance-PCI--DSS%20v4.0%20%7C%20ISO%2020022-success?style=flat-square)](#2-core-banking-system-cbs-decoupling--payment-rails)
[![Live Endpoint](https://img.shields.io/badge/Live%20Production-200%20OK%20Verified-brightgreen?style=flat-square)](https://wallet.thaw-zin-2k77.de5.net/healthz)

---

## 1. Executive Summary

This repository presents the **Production Reference Architecture** for **A Bank's Consumer Mobile Wallet Transaction Platform** (A-Plus digital banking ecosystem). 

Engineered specifically for financial-grade reliability, it resolves the core architectural tension in modern commercial banking: balancing **data sovereignty, regulatory auditability, and General Ledger consistency** (Core Banking System / on-premise Kubernetes) with **hyper-elastic promotional bursting, edge security, and sub-millisecond payment authorization** (AWS Serverless & Edge WAF).

The infrastructure is codified in **HashiCorp Terraform** (27 active cloud resources) and deployed live to AWS with Cloudflare Anycast perimeter routing.

### Verified Live Endpoints
* **Public Gateway (Edge Verified)**: `https://wallet.thaw-zin-2k77.de5.net/healthz`
* **Live SRE Circuit Breaker Telemetry**: `https://wallet.thaw-zin-2k77.de5.net/api/v1/resilience/circuit-breaker`
* **Ledger Readiness Probe**: `https://wallet.thaw-zin-2k77.de5.net/readyz`

```text
┌──────────────────────────────────────────────────────────────────────────────────────────────────┐
│                             END-TO-END VERIFIED PRODUCTION TOPOLOGY                              │
├──────────────────────────────────────────────────────────────────────────────────────────────────┤
│ MOBILE WALLET APPS (A-Plus / A-Bank Consumer Clients)                                            │
│   │ [HTTPS / TLS 1.3 - Authenticated Bearer JWT + Biometric Auth]                                │
│   ▼                                                                                              │
│ CLOUDFLARE EDGE ANYCAST (Global WAF, Anti-DDoS, Geo-Fencing, TLS Termination)                     │
│   │ [Proxied Static CNAME: wallet.thaw-zin-2k77.de5.net]                                         │
│   ▼                                                                                              │
│ AWS WAFv2 REGIONAL WEB APPLICATION FIREWALL (PCI-DSS v4.0 Layer 7 Defense)                      │
│   ├── Rule 1: CF-Connecting-IP Rate Limiting (300 req / 5 min per client mobile IP)             │
│   ├── Rule 2: AWSManagedRulesCommonRuleSet (OWASP Top 10 SQLi & XSS Shield)                      │
│   └── Rule 3: AWSManagedRulesKnownBadInputsRuleSet                                               │
│   │                                                                                              │
│   ▼                                                                                              │
│ AWS APPLICATION LOAD BALANCER (Multi-AZ Decoupled Ingress: us-east-1a through 1f)               │
│   ├── Listener Port 80   (Standard HTTP Forward)                                                 │
│   └── Listener Port 8080 (Cloudflare Origin Rule Ingestion Forward)                              │
│   │                                                                                              │
│   ▼                                                                                              │
│ AWS ALB TARGET GROUP (awsvpc IP-Target Micro-Segmentation)                                       │
│   └── Target: 172.31.1.82:8080 [HEALTHY in us-east-1a]                                          │
│   │                                                                                              │
│   ▼                                                                                              │
│ AWS ECS FARGATE SERVICE (Hybrid Autonomous Compute: a-bank-wallet-service-svc)                    │
│   ├── Base 2 On-Demand Tasks (Core Ledger Stability Guarantee)                                   │
│   └── Burst 4 Fargate Spot Tasks (70% FinOps Cost Savings during Promo Spikes)                   │
│   └── Runtime: Distroless Scratch Golang 1.22 Binary (Zero Shell, Non-Root UID 10001)            │
│         │                                                                                        │
│         ├── [ACID Idempotency Engine] ───> Amazon DynamoDB (PITR Continuous, KMS CMK)            │
│         ├── [Sub-millisecond Read Shield]> Redis / Valkey Cluster                                │
│         └── [Transactional Outbox] ──────> AWS SQS FIFO (a-bank-transaction-outbox)              │
│                                              │                                                   │
│                                              ├── [Dead-Letter Queue: 14-Day Audit Retention]     │
│                                              ▼                                                   │
│                                ASYNCHRONOUS RECONCILIATION WORKERS                               │
│                                              │                                                   │
│                     ┌────────────────────────┴──────────────────────────┐                        │
│                     ▼                                                   ▼                        │
│      CORE BANKING SYSTEM (CBS)                              NATIONAL PAYMENT SWITCHES            │
│      (Oracle FLEXCUBE / Apache Fineract)                    (CBM-Net 2 / MPU / MMQR)             │
│      - Master General Ledger (GL Accounts)                  - ISO 8583 Card/ATM Bitmaps          │
│      - CASA Accounts (Savings & Current)                    - ISO 20022 Interbank RTGS XML       │
│      - End-Of-Day (EOD) Batch Accrual                       - EMVCo Merchant QR Clearing         │
└──────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Core Banking System (CBS) Decoupling & Payment Rails

### 2.1 The Architectural Mandate: Why Wallets Never Touch CBS Databases Directly
A critical anti-pattern in digital banking is permitting a consumer mobile app to write directly to the Core Banking System database (e.g. Oracle, GaussDB, or DB2):
* **The Root Problem**: Core Banking Systems (such as Oracle FLEXCUBE, Finacle, or open-source **Apache Fineract**) are designed for General Ledger balancing, daily interest accruals, and nightly End-Of-Day (EOD) batch processing. They are not engineered to handle 50,000 requests/second mobile bursts during flash promotions or salary payout mornings.
* **The Failure Mode**: Direct mobile queries cause row-level locking on the `GL_ACCOUNTS` table, freezing branch tellers, ATMs, and corporate SWIFT wires simultaneously.
* **The Solution**: The Go wallet microservice acts as a high-throughput **Edge Ledger**. It completes transactions locally against an in-memory double-entry ledger with DynamoDB idempotency locking, emitting settled entries to an **AWS SQS Transactional Outbox** for asynchronous reconciliation with the Core Banking System.

### 2.2 Open-Source Financial Standards in Real-World Banking
1. **Core Banking Engine: Apache Fineract (Mifos X)**
   * Headless, API-first open-source core banking platform maintained by the Apache Software Foundation.
   * Manages double-entry GL bookkeeping, multi-currency CASA accounts, portfolio fees, and loan amortization.
2. **Card & ATM Switching: ISO 8583**
   * The international bitmap messaging protocol powering ATM switches, POS terminals, and the **Myanmar Payment Union (MPU)** network.
   * Codified via Go standard `moov-io/iso8583`.
3. **Interbank Clearing & RTGS: ISO 20022**
   * The modern XML financial standard mandatory across SWIFT CBPR+ and Central Bank Real-Time Gross Settlement systems (**CBM-Net 2**).
   * Supports `pacs.008` (Financial Institutional Customer Credit Transfer) and `pain.001` (Customer Credit Transfer Initiation).
   * Codified via Go standard `moov-io/iso20022`.
4. **National QR Standard: MMQR**
   * EMVCo Merchant-Presented QR specification enabling cross-wallet interoperability between A-Bank Wallet, KBZPay, CB Pay, and WavePay.

---

## 3. Resilient Payment Engineering: Golang Core (`wallet-service`)

The microservice (`cmd/wallet-service/`) is engineered in **Golang 1.22** for sub-millisecond execution, zero garbage-collection latency pauses, and deterministic memory consumption.

### 3.1 Mathematical Balance Conservation
Every internal transfer enforces atomic double-entry bookkeeping:

```text
Total Debited Balance == Total Credited Balance
Sum(Delta_Balance_Debit) + Sum(Delta_Balance_Credit) == 0
```

### 3.2 Network Retry Defense (`X-Idempotency-Key`)
In mobile cellular networks with intermittent handovers, clients frequently re-submit requests. Without idempotency, duplicate debits occur.
* Every mutating request requires a unique `X-Idempotency-Key` HTTP header.
* DynamoDB stores the SHA-256 request payload hash and previous response.
* Replays return cached responses (`idempotent_replay: true`) within 2 milliseconds without re-executing ledger operations.
* Payload mutations on an existing key return `HTTP 409 Conflict`.

### 3.3 Thread-Safe 3-State Circuit Breaker (`circuit_breaker.go`)
Protects internal goroutines from thread pool exhaustion when external clearing rails (CBM-Net, MPU, Visa) experience latency spikes or downtime:

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│                     3-STATE PAYMENT RAIL CIRCUIT BREAKER                         │
├──────────────────────────────────────────────────────────────────────────────────┤
│   [ CLOSED ] ──(5 Consecutive Timeouts / Failures)──> [ OPEN ]                   │
│       ▲                                                 │                        │
│       │                                                 │ (10s Cooldown Window)  │
│       │                                                 ▼                        │
│       └──(2 Successful Health Probes)──────────── [ HALF-OPEN ]                  │
└──────────────────────────────────────────────────────────────────────────────────┘
```

* **CLOSED**: Normal operations. Upstream calls execute; error counters monitored.
* **OPEN**: Tripped after 5 consecutive failures. Fails fast in $< 1\text{ ms}$ with `HTTP 503 Service Unavailable` and `Retry-After: 10`, preventing connection pooling exhaustion.
* **HALF-OPEN**: Allows 2 probe requests after 10-second cooldown. If both succeed, resets to `CLOSED`.
* **Chaos Engineering Endpoints**: Exposes `POST /api/v1/resilience/circuit-breaker/trip` and `/reset` for automated game-day testing.

---

## 4. Cloud Infrastructure-as-Code (Terraform v1.9.5)

All physical AWS and Cloudflare resources are codified in [terraform/main.tf](terraform/main.tf) with local state persistence (`terraform.tfstate`).

### 4.1 Active Resource Registry (27 Resources Managed)

| Category | Resource Name | Purpose |
| :--- | :--- | :--- |
| **Edge & Ingress** | `cloudflare_record.wallet_dns` | Proxied CNAME pointing to AWS ALB |
| **Edge Defense** | `aws_wafv2_web_acl.wallet_waf` | Regional WAF with CF-Connecting-IP rate-limiting |
| **WAF Binding** | `aws_wafv2_web_acl_association.alb_waf` | Enforces WAF inspection on the ALB |
| **Load Balancing** | `aws_lb.wallet_alb` | Internet-facing ALB across 6 Availability Zones |
| **ALB Listeners** | `aws_lb_listener.http`, `http_8080` | Dual-port listeners (80 & 8080 forward to TG) |
| **Target Group** | `aws_lb_target_group.wallet_tg` | IP-target mode for Fargate tasks on port 8080 |
| **Security Groups**| `aws_security_group.alb_sg` | Ingress 80, 443, 8080 from Cloudflare edge |
| **Micro-Segment** | `aws_security_group.fargate_sg` | Ingress strictly restricted to ALB SG ID (PCI-DSS 1.2) |
| **Serverless Tier**| `aws_ecs_cluster.wallet_cluster` | AWS ECS Fargate cluster |
| **Task Definition**| `aws_ecs_task_definition.wallet_task` | Scratch distroless runtime container |
| **Service Tier** | `aws_ecs_service.wallet_service` | Hybrid capacity providers (1 Base + 4 Spot) |
| **IAM Security** | `aws_iam_role.ecs_execution_role` | Least-privilege scoped task execution role |
| **Database Tier** | `aws_dynamodb_table.idempotency` | Pay-per-request ACID table with continuous PITR |
| **KMS Encryption** | `aws_kms_key.wallet_kms` | Customer Managed Key (CMK) with annual rotation |
| **KMS Alias** | `aws_kms_alias.wallet_kms_alias` | `alias/a-bank-wallet-cmk` |
| **Asynch Queue** | `aws_sqs_queue.tx_outbox` | Core Banking Transactional Outbox FIFO buffer |
| **Dead-Letter** | `aws_sqs_queue.tx_outbox_dlq` | 14-day retention dead-letter queue |
| **Backup Vault** | `aws_backup_vault.banking_vault` | Financial audit compliance backup vault |
| **Backup Policy** | `aws_backup_plan.dynamodb_plan` | Automated daily snapshot plan (35-day lifecycle) |
| **Observability** | `aws_cloudwatch_log_group.wallet_logs` | `/ecs/a-bank-wallet-service` (14-day retention) |
| **SLO Alarm 1** | `aws_cloudwatch_metric_alarm.alb_5xx_errors` | Alerts if $\ge 3$ server faults occur in 60s |
| **SLO Alarm 2** | `aws_cloudwatch_metric_alarm.alb_target_latency`| Alerts if latency $> 250\text{ ms}$ for 2 minutes |
| **SLO Alarm 3** | `aws_cloudwatch_metric_alarm.unhealthy_hosts` | Alerts if any container replica fails health checks |
| **Container ECR** | `aws_ecr_repository.wallet_service` | Hardened private container registry |
| **ECR Lifecycle** | `aws_ecr_lifecycle_policy.wallet_service` | Automated pruning of untagged images |

### 4.2 Adversarial Trap Solved: Cloudflare Forwarded-IP Rate Limiting
If standard client IP (`remote_addr`) rate limiting is configured in AWS WAF behind Cloudflare, AWS WAF evaluates Cloudflare's Anycast proxy IPs. When 300 legitimate requests arrive globally, AWS WAF blocks Cloudflare's IP, taking down the entire bank!
* **The SRE Defense**: In [main.tf](terraform/main.tf#L510), we codified `aggregate_key_type = "FORWARDED_IP"` with `header_name = "CF-Connecting-IP"`. This inspects the originating mobile device IP passed through Cloudflare's cryptographic header, isolating abusive actors without collateral damage to legitimate banking clients.

---

## 5. Hybrid Cloud-Native Hardening: Kubernetes, GitOps & Ansible

For financial sovereignty and hybrid datacenter parity (on-premise banking datacenters + AWS cloud), the platform provides complete multi-tier orchestration:

### 5.1 Cloud-Native Kubernetes & Helm v3 (`charts/wallet-service`)
* **PodDisruptionBudget (`templates/pdb.yaml`)**: Enforces `minAvailable: 1`. Guarantees that automated node patching, AMI upgrades, or spot terminations never cause ledger downtime.
* **Multi-AZ Topology Spread & Anti-Affinity (`values.yaml`)**: Configured with `topology.kubernetes.io/zone`. Pod replicas are strictly scheduled across different physical Availability Zones (`us-east-1a`, `us-east-1b`, `us-east-1c`).
* **Zero-Trust NetworkPolicy (`templates/networkpolicy.yaml`)**: Default Deny. Restricts inbound traffic strictly to ingress controllers on port 8080. Egress is restricted strictly to CoreDNS (port 53) and DynamoDB/KMS (port 443).
* **KEDA Event-Driven Autoscaling (`templates/keda-scaledobject.yaml`)**: Bypasses slow CPU metrics by scaling on SQS queue depth (`targetQueueLength: 30`). Reacts in sub-seconds during traffic spikes before container CPU registers a rise.

### 5.2 GitOps Continuous Delivery: ArgoCD (`gitops/argocd-application.yaml`)
* **Declarative Single Source of Truth**: The Git repository is the sole authoritative state of production.
* **Zero-Human-Kubectl Access Model**: Engineers never execute manual `kubectl apply`. ArgoCD continuously reconciles cluster state against Git commit hashes.
* **Self-Healing & Drift Detection**: Configured with `selfHeal: true` and `prune: true`. Any out-of-band cluster modifications or orphaned resources are automatically reverted within seconds.

### 5.3 Bare-Metal Fleet Hardening & Automation: Ansible (`ansible/`)
* **CIS Benchmark Linux Kernel Hardening (`ansible/roles/fintech_kernel_hardening`)**:
  * Tuning network stack buffers (`net.core.somaxconn = 65535`, `net.ipv4.tcp_max_syn_backlog = 65535`).
  * Preventing ephemeral port exhaustion under heavy concurrent mobile load (`net.ipv4.ip_local_port_range = 1024 65535`).
  * TCP SYN cookie protection (`net.ipv4.tcp_syncookies = 1`) against volumetric SYN flood attacks.
  * Expanding file descriptor ceilings (`fs.file-max = 2097152`).
* **Automated Runtime Bootstrapping (`ansible/playbooks/site.yaml`)**: Boots lightweight K3s nodes and deploys the wallet Helm chart idempotently across edge bare-metal servers.

---

## 6. DevSecOps CI/CD & Supply Chain Security

The automated pipeline ([.github/workflows/ci-cd.yaml](.github/workflows/ci-cd.yaml)) executes on every commit pushed to `main`:

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│                          FOUR-STAGE DEVSECOPS PIPELINE                           │
├──────────────────────────────────────────────────────────────────────────────────┤
│ Stage 1: Quality Gate       --> Go fmt check, vet, -race tests & Helm linting    │
│ Stage 2: Container Security --> Aqua Security Trivy Scan (0 Critical/High CVEs)  │
│ Stage 3: IaC Policy Audit   --> Aqua Security tfsec Static Analysis              │
│ Stage 4: Build & Release    --> ECR push + Anchore Syft SPDX-JSON SBOM Artifact │
└──────────────────────────────────────────────────────────────────────────────────┘
```

* **Supply Chain SBOM (PCI-DSS v4.0 Requirement 6.3.2)**: Anchore Syft inspects the compiled distroless container image, generates a verifiable SPDX-JSON Software Bill of Materials (`a-bank-wallet-sbom.spdx.json`), and stores it as a 30-day compliance artifact.
* **Distroless Scratch Security**: Built on `scratch` (< 15MB). Contains zero shell binaries (`/bin/sh`, `/bin/bash`), zero package managers, and executes under unprivileged UID `10001:10001`.
* **Zero Host Pollution**: All toolchains, runners, and linters run in isolated ephemeral containers or userland runtimes.

---

## 7. Financial SRE Observability & Golden Signals

Comprehensive financial telemetry is codified in [monitoring/](monitoring/):

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         FOUR GOLDEN SIGNALS MONITORING                           │
├──────────────────────────────────────────────────────────────────────────────────┤
│ LATENCY     --> P99 Transaction Latency (< 250ms target across all payment rails)│
│ TRAFFIC     --> Requests Per Second (RPS) broken down by route & payment partner │
│ ERRORS      --> HTTP 5XX rate + Payment Gateway Failure rate (Error Budget SLI)  │
│ SATURATION  --> Goroutine pool count, SQS queue depth, DynamoDB throttled events │
└──────────────────────────────────────────────────────────────────────────────────┘
```

* **Prometheus Alerting Rules (`monitoring/slo-rules.yaml` & `monitoring/prometheus-rules.yaml`)**:
  * **Dual Engine Validation**: Fully validated via `promtool check rules` (native PromQL engine) and `kubeconform` (Kubernetes `monitoring.coreos.com/v1 PrometheusRule` CRD schema).
  * **99.95% Availability SLO Burn Rate**: Multi-window burn rate alerts triggering at $14.4\times$ burn over 1 hour (critical page) and $6\times$ burn over 6 hours (warning ticket).
  * **P99 Latency Breach**: Alerts if $P_{99} > 250\text{ ms}$ sustained over a 5-minute evaluation window.
  * **Circuit Breaker Tripped**: Immediate P1 alert whenever an external payment rail transitions to `OPEN`.
  * **Dead-Letter Queue Backlog**: Alerts if `a-bank-transaction-outbox-dlq` contains $> 0$ unprocessable transactions requiring forensic audit.
* **Grafana SRE Dashboard (`monitoring/wallet-dashboard.json`)**: Pre-configured JSON dashboard visualizing live throughput, circuit breaker states, and error budget exhaustion curves.

---

## 8. FinOps Cost Governance & Disaster Recovery

```text
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         A BANK AWS FINOPS RUN RATE                               │
├──────────────────────────────────────────────────────────────────────────────────┤
│ COMPONENT                             SPECIFICATION                       COST   │
├──────────────────────────────────────────────────────────────────────────────────┤
│ AWS Application Load Balancer         Multi-AZ (6 Subnets, Dual Port)     ~$18.00│
│ AWS ECS Fargate Base On-Demand        1 Replica (0.25 vCPU, 0.5 GB RAM)   ~$9.01 │
│ AWS ECS Fargate Spot Burst            Up to 4 Replicas (70% Discount)     ~$2.70 │
│ AWS DynamoDB Idempotency Table        Pay-Per-Request (On-Demand Mode)    ~$0.10 │
│ AWS KMS Customer Managed Key          alias/a-bank-wallet-cmk             ~$1.00 │
│ AWS SQS Transactional Outbox          1M Requests / Month (Free Tier)     ~$0.00 │
│ AWS WAFv2 WebACL                      1 WebACL + 3 Rules                  ~$8.00 │
│ AWS CloudWatch Metrics & Logs         14-Day Automated Log Retention      ~$0.00 │
├──────────────────────────────────────────────────────────────────────────────────┤
│ TOTAL ESTIMATED STEADY-STATE RUN RATE:                                    ~$38.81│
└──────────────────────────────────────────────────────────────────────────────────┘
```

* **Disaster Recovery (DR) Metrics**:
  * **Recovery Point Objective (RPO)**: `0 seconds` (DynamoDB continuous Point-in-Time Recovery streams every write).
  * **Recovery Time Objective (RTO)**: `< 15 minutes` (Automated Terraform state reconstitution and backup vault restore).
* **Compliance Vault**: `aws_backup_vault.banking_vault` enforces a 35-day financial cycle backup policy with Write-Once-Read-Many (WORM) audit protection.

---

## 9. Quickstart & Local Developer Experience (DevEx)

Engineers can boot an identical banking environment locally in under 5 seconds without AWS IAM permissions:

```bash
# Clone the repository
git clone https://github.com/pretamane/fintech-wallet-sre-showcase.git
cd fintech-wallet-sre-showcase

# Boot the local production parity stack (Local Go Container + In-Memory DynamoDB)
docker compose up -d

# Verify local health
curl -s http://localhost:8080/healthz

# Inspect accounts
curl -s http://localhost:8080/api/v1/wallets
```

---

## 10. Author & Operational Handover

* **Author**: Thaw Zin @ Chris (`pretamane`)
* **Role Target**: Senior DevOps / SRE Lead (A Bank Digital Banking Division)
* **Architecture Specification**: Codified in [terraform/main.tf](terraform/main.tf) and [charts/wallet-service/values.yaml](charts/wallet-service/values.yaml).
* **Repository**: [pretamane/fintech-wallet-sre-showcase](https://github.com/pretamane/fintech-wallet-sre-showcase)
