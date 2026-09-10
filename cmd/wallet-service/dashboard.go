package main

const rootDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>A Bank Payment Services • SRE & FinTech Operations Platform</title>
  <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 32 32'%3E%3Crect width='32' height='32' rx='6' fill='%231e3a8a'/%3E%3Ctext x='16' y='23' font-family='-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,sans-serif' font-size='20' font-weight='800' fill='%23ffffff' text-anchor='middle'%3EA%3C/text%3E%3C/svg%3E">
  <style>
    :root {
      --bg-body: #f8fafc;
      --bg-card: #ffffff;
      --bg-subtle: #f1f5f9;
      --border-light: #e2e8f0;
      --border-dark: #cbd5e1;
      --text-heading: #0f172a;
      --text-body: #334155;
      --text-muted: #64748b;
      --primary-navy: #1e3a8a;
      --primary-hover: #1e40af;
      --accent-blue: #0284c7;
      --success-green: #059669;
      --success-bg: #ecfdf5;
      --warning-amber: #d97706;
      --warning-bg: #fffbeb;
      --danger-red: #dc2626;
      --danger-bg: #fef2f2;
      --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
      --font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-body);
      color: var(--text-body);
      font-family: var(--font-sans);
      line-height: 1.5;
      padding: 1.5rem 1rem;
    }
    .container { max-width: 1260px; margin: 0 auto; }
    
    /* Header */
    header {
      background-color: var(--bg-card);
      border: 1px solid var(--border-light);
      border-radius: 8px;
      padding: 1.25rem 1.75rem;
      margin-bottom: 1.25rem;
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      gap: 1rem;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
    }
    .brand-section { display: flex; align-items: center; gap: 1rem; }
    .brand-badge {
      background-color: var(--primary-navy);
      color: #ffffff;
      font-weight: 800;
      font-size: 0.85rem;
      padding: 0.45rem 0.85rem;
      border-radius: 4px;
      letter-spacing: 0.08em;
    }
    h1 {
      font-size: 1.25rem;
      font-weight: 700;
      color: var(--text-heading);
      letter-spacing: -0.01em;
    }
    .header-subtext {
      font-size: 0.8rem;
      color: var(--text-muted);
      margin-top: 0.15rem;
    }
    .system-status { display: flex; align-items: center; gap: 0.75rem; flex-wrap: wrap; }
    .status-tag {
      display: inline-flex;
      align-items: center;
      gap: 0.4rem;
      background-color: var(--success-bg);
      color: var(--success-green);
      border: 1px solid #a7f3d0;
      padding: 0.3rem 0.7rem;
      border-radius: 4px;
      font-size: 0.75rem;
      font-weight: 700;
      font-family: var(--font-mono);
    }
    .status-indicator {
      width: 8px;
      height: 8px;
      background-color: var(--success-green);
      border-radius: 50%;
    }
    .status-indicator.open { background-color: var(--danger-red); }
    .status-indicator.half-open { background-color: var(--warning-amber); }
    .ticker-tag {
      font-family: var(--font-mono);
      font-size: 0.75rem;
      color: var(--text-muted);
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      padding: 0.3rem 0.6rem;
      border-radius: 4px;
    }

    /* Tabs Navigation */
    .nav-tabs {
      display: flex;
      gap: 0.25rem;
      border-bottom: 1px solid var(--border-light);
      margin-bottom: 1.25rem;
      overflow-x: auto;
      background-color: var(--bg-card);
      border: 1px solid var(--border-light);
      border-radius: 8px 8px 0 0;
      padding: 0.25rem 0.5rem 0 0.5rem;
    }
    .nav-tab {
      background: none;
      border: none;
      border-bottom: 3px solid transparent;
      color: var(--text-muted);
      font-size: 0.85rem;
      font-weight: 600;
      padding: 0.75rem 1rem;
      cursor: pointer;
      transition: all 0.15s ease;
      white-space: nowrap;
    }
    .nav-tab:hover { color: var(--text-heading); background-color: var(--bg-subtle); border-radius: 4px 4px 0 0; }
    .nav-tab.active {
      color: var(--primary-navy);
      border-bottom-color: var(--primary-navy);
      background-color: #ffffff;
    }

    /* Layout Grids */
    .grid-4 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 1rem;
      margin-bottom: 1.25rem;
    }
    .grid-3 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 1rem;
      margin-bottom: 1.25rem;
    }
    .grid-2 {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1.25rem;
      margin-bottom: 1.25rem;
    }
    @media (max-width: 960px) {
      .grid-2 { grid-template-columns: 1fr; }
    }

    /* Cards */
    .card {
      background-color: var(--bg-card);
      border: 1px solid var(--border-light);
      border-radius: 8px;
      padding: 1.25rem;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
      margin-bottom: 1.25rem;
    }
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
      padding-bottom: 0.6rem;
      border-bottom: 1px solid var(--border-light);
    }
    .card-title {
      font-size: 0.95rem;
      font-weight: 700;
      color: var(--text-heading);
      letter-spacing: -0.01em;
    }
    .account-id {
      font-family: var(--font-mono);
      font-size: 0.75rem;
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      color: var(--text-muted);
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
    }
    .metric-value {
      font-family: var(--font-mono);
      font-size: 1.65rem;
      font-weight: 800;
      color: var(--primary-navy);
      margin: 0.35rem 0;
    }
    .account-type {
      font-size: 0.75rem;
      color: var(--text-muted);
    }

    /* Circuit Breaker Visual State Box */
    .breaker-machine {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 1.25rem;
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      border-radius: 6px;
      margin-bottom: 1.25rem;
      gap: 0.5rem;
      flex-wrap: wrap;
    }
    .breaker-node {
      flex: 1;
      min-width: 140px;
      text-align: center;
      padding: 0.75rem 0.5rem;
      background-color: #ffffff;
      border: 2px solid var(--border-light);
      border-radius: 6px;
      font-family: var(--font-mono);
      font-weight: 700;
      font-size: 0.85rem;
      color: var(--text-muted);
      transition: all 0.2s;
    }
    .breaker-node.active-closed {
      border-color: var(--success-green);
      background-color: var(--success-bg);
      color: var(--success-green);
      box-shadow: 0 0 0 3px rgba(5, 150, 105, 0.15);
    }
    .breaker-node.active-open {
      border-color: var(--danger-red);
      background-color: var(--danger-bg);
      color: var(--danger-red);
      box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.15);
    }
    .breaker-node.active-halfopen {
      border-color: var(--warning-amber);
      background-color: var(--warning-bg);
      color: var(--warning-amber);
      box-shadow: 0 0 0 3px rgba(217, 119, 6, 0.15);
    }
    .breaker-arrow {
      font-family: var(--font-mono);
      color: var(--text-muted);
      font-size: 1.1rem;
      font-weight: 800;
    }

    /* Forms */
    .form-group { margin-bottom: 1rem; }
    label {
      display: block;
      font-size: 0.75rem;
      font-weight: 700;
      color: var(--text-heading);
      margin-bottom: 0.35rem;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    input, select {
      width: 100%;
      background-color: #ffffff;
      border: 1px solid var(--border-dark);
      border-radius: 5px;
      padding: 0.6rem 0.8rem;
      color: var(--text-heading);
      font-family: inherit;
      font-size: 0.85rem;
      transition: border-color 0.15s;
    }
    input:focus, select:focus {
      outline: none;
      border-color: var(--accent-blue);
      box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.12);
    }
    .chip-group {
      display: flex;
      gap: 0.4rem;
      margin-top: 0.4rem;
      flex-wrap: wrap;
    }
    .amount-chip {
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      color: var(--text-heading);
      font-family: var(--font-mono);
      font-size: 0.75rem;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.15s;
    }
    .amount-chip:hover {
      background-color: var(--border-dark);
    }

    /* Buttons */
    .action-group {
      display: flex;
      flex-wrap: wrap;
      gap: 0.6rem;
      margin-top: 1.25rem;
    }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      font-size: 0.8rem;
      font-weight: 700;
      padding: 0.55rem 1rem;
      border-radius: 5px;
      border: 1px solid transparent;
      cursor: pointer;
      transition: all 0.15s;
      font-family: var(--font-sans);
    }
    .btn-primary { background-color: var(--primary-navy); color: #ffffff; }
    .btn-primary:hover { background-color: var(--primary-hover); }
    .btn-secondary { background-color: #ffffff; border-color: var(--border-dark); color: var(--text-heading); }
    .btn-secondary:hover { background-color: var(--bg-subtle); }
    .btn-warning { background-color: #ffffff; border-color: var(--warning-amber); color: var(--warning-amber); }
    .btn-warning:hover { background-color: var(--warning-bg); }
    .btn-danger { background-color: #ffffff; border-color: var(--danger-red); color: var(--danger-red); }
    .btn-danger:hover { background-color: var(--danger-bg); }
    .btn-success { background-color: #ffffff; border-color: var(--success-green); color: var(--success-green); }
    .btn-success:hover { background-color: var(--success-bg); }

    /* Console Boxes */
    .console-box {
      background-color: #0f172a;
      color: #38bdf8;
      font-family: var(--font-mono);
      font-size: 0.78rem;
      line-height: 1.6;
      padding: 1rem;
      border-radius: 6px;
      border: 1px solid #1e293b;
      min-height: 220px;
      max-height: 360px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
    .badge-execution {
      display: inline-block;
      font-family: var(--font-mono);
      font-size: 0.7rem;
      font-weight: 700;
      padding: 0.2rem 0.5rem;
      border-radius: 4px;
    }
    .badge-success { background-color: var(--success-bg); color: var(--success-green); border: 1px solid #a7f3d0; }
    .badge-warning { background-color: var(--warning-bg); color: var(--warning-amber); border: 1px solid #fde68a; }
    .badge-danger  { background-color: var(--danger-bg);  color: var(--danger-red);    border: 1px solid #fecaca; }

    /* Tables */
    .table-container { overflow-x: auto; margin-top: 0.75rem; }
    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
    }
    th, td {
      padding: 0.75rem 0.85rem;
      border-bottom: 1px solid var(--border-light);
      font-size: 0.8rem;
    }
    th {
      background-color: var(--bg-subtle);
      color: var(--text-heading);
      font-weight: 700;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    tr:hover td { background-color: #fafafa; }
    .spec-badge {
      display: inline-block;
      font-family: var(--font-mono);
      font-size: 0.72rem;
      font-weight: 600;
      padding: 0.15rem 0.45rem;
      border-radius: 3px;
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      color: var(--text-heading);
    }
    .spec-badge.green { background-color: var(--success-bg); color: var(--success-green); border-color: #a7f3d0; }
    .spec-badge.blue  { background-color: #e0f2fe; color: #0369a1; border-color: #bae6fd; }

    /* Callouts */
    .callout {
      padding: 0.85rem 1rem;
      background-color: var(--bg-subtle);
      border-left: 4px solid var(--primary-navy);
      border-radius: 0 4px 4px 0;
      margin: 1rem 0;
      font-size: 0.8rem;
      color: var(--text-body);
    }
    .callout.warning {
      background-color: var(--warning-bg);
      border-left-color: var(--warning-amber);
    }
    .callout.danger {
      background-color: var(--danger-bg);
      border-left-color: var(--danger-red);
    }

    /* Tab Switcher */
    .tab-content { display: none; }
    .tab-content.active { display: block; }
  </style>
</head>
<body>
  <div class="container">
    <!-- Header -->
    <header>
      <div class="brand-section">
        <span class="brand-badge">A BANK</span>
        <div>
          <h1>Mobile Wallet Operations &amp; SRE Platform</h1>
          <p class="header-subtext">
            Production Reference Architecture • Multi-AZ AWS ALB &bull; ECS Fargate Spot &bull; Cloudflare WAF &bull; PCI-DSS v4.0
          </p>
        </div>
      </div>
      <div class="system-status">
        <span class="status-tag" id="clusterStatusBadge"><span class="status-indicator" id="statusDot"></span> CLUSTER ACTIVE</span>
        <span class="ticker-tag" id="latencyTicker">Latency: -- ms</span>
        <span class="ticker-tag" id="circuitTicker">Rail: CLOSED</span>
      </div>
    </header>

    <!-- Navigation Tabs -->
    <div class="nav-tabs">
      <button class="nav-tab active" onclick="showTab('tab-resilience')">SRE Resilience &amp; Chaos Engine</button>
      <button class="nav-tab" onclick="showTab('tab-ledger')">Core Banking Ledger &amp; Idempotency</button>
      <button class="nav-tab" onclick="showTab('tab-cbs')">CBS Decoupling &amp; Payment Rails</button>
      <button class="nav-tab" onclick="showTab('tab-cloud')">AWS Cloud Infrastructure (27 Resources)</button>
      <button class="nav-tab" onclick="showTab('tab-devsecops')">DevSecOps &amp; PCI-DSS v4.0</button>
      <button class="nav-tab" onclick="showTab('tab-k8s')">Kubernetes, GitOps &amp; Fleet</button>
      <button class="nav-tab" onclick="showTab('tab-slo')">SRE Golden Signals &amp; SLOs</button>
      <button class="nav-tab" onclick="showTab('tab-runbook')">Verification Runbook &amp; CLI</button>
    </div>

    <!-- TAB 1: SRE RESILIENCE & CHAOS ENGINE -->
    <div id="tab-resilience" class="tab-content active">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Downstream Payment Rail 3-State Circuit Breaker (CBM-Net 2 Central Bank Clearing)</span>
          <span class="spec-badge blue">THREAD-SAFE GOROUTINE PROTECTION</span>
        </div>

        <p style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 1rem;">
          In high-volume banking systems, external interbank payment switches (CBM-Net 2, MPU, Visa) periodically experience latency spikes or downtime. Without a circuit breaker, incoming mobile payment goroutines exhaust the thread pool waiting on slow sockets, causing cascading gateway timeouts (HTTP 504) across all retail banking services.
        </p>

        <!-- Visual Circuit Breaker State Machine -->
        <div class="breaker-machine">
          <div class="breaker-node active-closed" id="node-closed">
            <div>CLOSED</div>
            <div style="font-size: 0.7rem; font-weight: normal; margin-top: 0.25rem;">Normal Operations<br>Traffic Flowing</div>
          </div>
          <div class="breaker-arrow">&rarr; (5 Failures) &rarr;</div>
          <div class="breaker-node" id="node-open">
            <div>OPEN</div>
            <div style="font-size: 0.7rem; font-weight: normal; margin-top: 0.25rem;">Fail-Fast &lt;1ms<br>HTTP 503 Service Unavailable</div>
          </div>
          <div class="breaker-arrow">&rarr; (10s Window) &rarr;</div>
          <div class="breaker-node" id="node-halfopen">
            <div>HALF-OPEN</div>
            <div style="font-size: 0.7rem; font-weight: normal; margin-top: 0.25rem;">Canary Probing<br>2 Consecutive Passes</div>
          </div>
        </div>

        <!-- Telemetry Counters -->
        <div class="grid-4">
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Active State</span>
            <div class="metric-value" id="cb-state-display" style="font-size: 1.25rem; color: var(--success-green);">CLOSED</div>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Consecutive Failures</span>
            <div class="metric-value" id="cb-failures-display" style="font-size: 1.25rem;">0 / 5</div>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Trip Threshold</span>
            <div class="metric-value" style="font-size: 1.25rem;">5 Failures</div>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Recovery Cooldown</span>
            <div class="metric-value" style="font-size: 1.25rem;">10 Seconds</div>
          </div>
        </div>

        <!-- Interactive Chaos Engineering Actions -->
        <div class="action-group">
          <button class="btn btn-danger" onclick="triggerCircuitTrip()">Trip Circuit Breaker (Simulate CBM-Net Outage)</button>
          <button class="btn btn-success" onclick="triggerCircuitReset()">Reset Circuit Breaker (Restore Payment Rail)</button>
          <button class="btn btn-primary" onclick="probeReadinessAndLiveness()">Probe Liveness &amp; Readiness Status</button>
          <button class="btn btn-secondary" onclick="simulateTrafficBurst(10)">Simulate 10x Concurrent Burst Requests</button>
        </div>

        <div class="callout warning">
          <strong>SRE Anti-Crash-Loop Architecture Win</strong>: When the circuit breaker trips to <code>OPEN</code>, the <strong>Readiness Probe (<code>/readyz</code>)</strong> immediately returns <code>HTTP 503 Service Unavailable</code> so load balancers (AWS ALB / Kubernetes Ingress) stop sending traffic to the degraded instance. Concurrently, the <strong>Liveness Probe (<code>/healthz</code>)</strong> remains <code>HTTP 200 OK</code>, preventing Kubernetes/ECS from needlessly restarting the container in an endless crash-loop during external partner downtime.
        </div>

        <!-- Chaos Audit Log -->
        <div class="card" style="margin-top: 1rem; margin-bottom: 0;">
          <div class="card-header">
            <span class="card-title">Live Chaos Telemetry Log</span>
            <span class="spec-badge">GET /api/v1/resilience/circuit-breaker</span>
          </div>
          <div class="console-box" id="chaosConsole">// Resilience engine online. Click 'Trip Circuit Breaker' to execute a live Chaos Engineering test.</div>
        </div>
      </div>
    </div>

    <!-- TAB 2: CORE BANKING LEDGER & IDEMPOTENCY -->
    <div id="tab-ledger" class="tab-content">
      <!-- Balance Cards -->
      <div class="grid-3">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Customer Account A</span>
            <span class="account-id">ACC-1001</span>
          </div>
          <div class="metric-value" id="display-ACC-1001">-- MMK</div>
          <p class="account-type">Retail Consumer Profile</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Merchant Partner B</span>
            <span class="account-id">ACC-2002</span>
          </div>
          <div class="metric-value" id="display-ACC-2002">-- MMK</div>
          <p class="account-type">Commercial POS Clearing Account</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Central Reserve Pool</span>
            <span class="account-id">ACC-9999</span>
          </div>
          <div class="metric-value" id="display-ACC-9999">-- MMK</div>
          <p class="account-type">Central Bank Liquidity Reserve</p>
        </div>
      </div>

      <!-- Transaction Form & Audit Output -->
      <div class="grid-2">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Atomic Double-Entry Transaction</span>
            <span class="spec-badge">POST /api/v1/wallets/transfer</span>
          </div>

          <div class="form-group">
            <label>Source Account (Debit)</label>
            <select id="sourceAccount">
              <option value="ACC-1001">ACC-1001 (Customer Account A)</option>
              <option value="ACC-2002">ACC-2002 (Merchant Partner B)</option>
              <option value="ACC-9999">ACC-9999 (Central Reserve Pool)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Destination Account (Credit)</label>
            <select id="targetAccount">
              <option value="ACC-2002">ACC-2002 (Merchant Partner B)</option>
              <option value="ACC-1001">ACC-1001 (Customer Account A)</option>
              <option value="ACC-9999">ACC-9999 (Central Reserve Pool)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Transaction Amount (MMK)</label>
            <input type="number" id="transferAmount" value="25000" min="1" step="500">
            <div class="chip-group">
              <span class="amount-chip" onclick="selectAmount(5000)">5,000</span>
              <span class="amount-chip" onclick="selectAmount(25000)">25,000</span>
              <span class="amount-chip" onclick="selectAmount(50000)">50,000</span>
              <span class="amount-chip" onclick="selectAmount(100000)">100,000</span>
              <span class="amount-chip" onclick="selectAmount(500000)">500,000</span>
            </div>
          </div>

          <div class="form-group">
            <label>Idempotency Key (X-Idempotency-Key)</label>
            <div style="display: flex; gap: 0.5rem;">
              <input type="text" id="idempotencyKey" style="font-family: var(--font-mono); font-size: 0.85rem;">
              <button class="btn btn-secondary" onclick="generateNewKey()" type="button" style="white-space: nowrap;">Generate Key</button>
            </div>
            <p style="font-size: 0.75rem; color: var(--text-muted); margin-top: 0.35rem;">
              Guarantees zero duplicate debits across flaky cellular handovers and retried mobile payment submissions.
            </p>
          </div>

          <div class="action-group">
            <button class="btn btn-primary" onclick="executeTransfer(false)">Execute Transfer</button>
            <button class="btn btn-warning" onclick="executeTransfer(true)" title="Submits the identical key to verify idempotent replay">Simulate Duplicate Request</button>
            <button class="btn btn-danger" onclick="simulateInsufficientFunds()" title="Attempts to debit more than available balance">Simulate Insufficient Funds</button>
            <button class="btn btn-secondary" onclick="updateAllBalances()">Refresh Balances</button>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Transaction Receipt &amp; Audit Log</span>
            <span id="responseStatusBadge" class="badge-execution" style="display: none;"></span>
          </div>
          <div class="console-box" id="auditConsole">// Ready. Select accounts and click 'Execute Transfer' above.</div>
          <div style="margin-top: 1rem; padding: 0.75rem; background-color: var(--bg-subtle); border-radius: 4px; font-size: 0.75rem; color: var(--text-muted);">
            <strong>Mathematical Invariant</strong>: <code>Debit Sum == Credit Sum</code>. The double-entry ledger verifies balance conservation atomically under mutex locking with SHA-256 payload hashing stored in DynamoDB PITR.
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 3: CBS DECOUPLING & PAYMENT RAILS -->
    <div id="tab-cbs" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Why Wallets Must Never Touch Core Banking System (CBS) Databases Directly</span>
          <span class="spec-badge green">TRANSACTIONAL OUTBOX PATTERN</span>
        </div>

        <div class="callout danger">
          <strong>The General Ledger Lock Contention Trap</strong>: Core Banking Systems (Oracle FLEXCUBE, Finacle, or Apache Fineract) are engineered for nightly End-Of-Day (EOD) batch accruals, daily interest computations, and GL reconciliation. If 50,000 mobile wallet requests per second write directly to the CBS database during flash promotions or salary payout mornings, row-level locks on the <code>GL_ACCOUNTS</code> table freeze branch tellers, ATM switches, and SWIFT wires nationwide.
        </div>

        <p style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 1rem;">
          Our architecture implements an <strong>Edge Ledger Pattern</strong>: the Go microservice authorizes transactions locally in sub-milliseconds with DynamoDB idempotency locking, emitting settled entries into an <strong>AWS SQS FIFO Transactional Outbox</strong> for asynchronous batch reconciliation with the Core Banking System.
        </p>

        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Financial Rail / Standard</th>
                <th>Protocol &amp; Format</th>
                <th>Role in A Bank Ecosystem</th>
                <th>SRE Decoupling Implementation</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><strong>Apache Fineract (Mifos X)</strong></td>
                <td>REST API / Double-Entry GL</td>
                <td>Master General Ledger &amp; CASA accounts</td>
                <td>Asynchronous outbox workers batch-reconcile settled wallet debits into Fineract GL accounts</td>
              </tr>
              <tr>
                <td><strong>CBM-Net 2 (Central Bank RTGS)</strong></td>
                <td>ISO 20022 XML (<code>pacs.008</code>, <code>pain.001</code>)</td>
                <td>Interbank high-value clearing &amp; settlement</td>
                <td>Protected via thread-safe 3-state Circuit Breaker; fails fast if CBM-Net RTGS switch times out</td>
              </tr>
              <tr>
                <td><strong>Myanmar Payment Union (MPU)</strong></td>
                <td>ISO 8583 Bitmap Protocol</td>
                <td>Card transactions, ATM switches &amp; POS terminals</td>
                <td>Processed via ISO 8583 message broker; queued to dead-letter queue (DLQ) upon partner timeout</td>
              </tr>
              <tr>
                <td><strong>MMQR (EMVCo Merchant QR)</strong></td>
                <td>EMVCo Merchant-Presented QR</td>
                <td>Cross-wallet interoperability (A Bank, KBZPay, CB Pay, WavePay)</td>
                <td>Locally verified cryptographic signature; offloaded to clearing queue for end-of-day interbank netting</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB 4: AWS CLOUD INFRASTRUCTURE (27 RESOURCES) -->
    <div id="tab-cloud" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Physical Cloud Infrastructure Registry (27 Active Terraform Resources)</span>
          <span class="spec-badge blue">TERRAFORM v1.9.5 MANAGED</span>
        </div>

        <p style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 1rem;">
          The platform infrastructure is codified in <code>terraform/main.tf</code> with state persistence. Below is the live inventory of active cloud resources:
        </p>

        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Category</th>
                <th>Terraform Resource</th>
                <th>Physical ID / Identifier</th>
                <th>Architectural Purpose</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><strong>Edge Ingress</strong></td>
                <td><code>cloudflare_record.wallet_dns</code></td>
                <td><code>fc356b9c20d80baa6c4412c2a028d718</code></td>
                <td>Proxied CNAME pointing <code>wallet.thaw-zin-2k77.de5.net</code> to AWS ALB</td>
              </tr>
              <tr>
                <td><strong>Perimeter WAF</strong></td>
                <td><code>aws_wafv2_web_acl.wallet_waf</code></td>
                <td><code>505f7bb7-621c-4fbf-afe5-6d3cf6903f87</code></td>
                <td>CF-Connecting-IP rate-limiting (300 req/5m) + OWASP Top 10 rule sets</td>
              </tr>
              <tr>
                <td><strong>Load Balancing</strong></td>
                <td><code>aws_lb.wallet_alb</code></td>
                <td><code>a-bank-wallet-alb-301681684</code></td>
                <td>Multi-AZ internet-facing ALB spanning 6 public subnets in us-east-1</td>
              </tr>
              <tr>
                <td><strong>ALB Listeners</strong></td>
                <td><code>aws_lb_listener.http</code>, <code>http_8080</code></td>
                <td>Port 80 &amp; Port 8080 Listeners</td>
                <td>Dual-port forwarding to Target Group with connection draining</td>
              </tr>
              <tr>
                <td><strong>Target Group</strong></td>
                <td><code>aws_lb_target_group.wallet_tg</code></td>
                <td><code>a-bank-wallet-tg</code></td>
                <td>IP-target mode forwarding to ECS Fargate containers on port 8080</td>
              </tr>
              <tr>
                <td><strong>Compute Cluster</strong></td>
                <td><code>aws_ecs_cluster.wallet_cluster</code></td>
                <td><code>a-bank-wallet-cluster</code></td>
                <td>Serverless ECS Fargate cluster with Capacity Provider Strategies</td>
              </tr>
              <tr>
                <td><strong>Service Tier</strong></td>
                <td><code>aws_ecs_service.wallet_service</code></td>
                <td><code>a-bank-wallet-service-svc</code></td>
                <td>Hybrid compute: 1 Base On-Demand task + 4 Burst Spot tasks (70% FinOps savings)</td>
              </tr>
              <tr>
                <td><strong>Database Tier</strong></td>
                <td><code>aws_dynamodb_table.idempotency</code></td>
                <td><code>a-bank-wallet-idempotency</code></td>
                <td>On-Demand ACID idempotency table with continuous Point-in-Time Recovery (PITR)</td>
              </tr>
              <tr>
                <td><strong>KMS Encryption</strong></td>
                <td><code>aws_kms_key.wallet_kms</code></td>
                <td><code>alias/a-bank-wallet-cmk</code></td>
                <td>Customer Managed Key (CMK) with automated 365-day rotation (PCI-DSS 3.5)</td>
              </tr>
              <tr>
                <td><strong>Outbox Queue</strong></td>
                <td><code>aws_sqs_queue.tx_outbox</code></td>
                <td><code>a-bank-transaction-outbox</code></td>
                <td>KMS-encrypted FIFO queue decoupling mobile ledger from Core Banking System</td>
              </tr>
              <tr>
                <td><strong>Dead-Letter</strong></td>
                <td><code>aws_sqs_queue.tx_outbox_dlq</code></td>
                <td><code>a-bank-transaction-outbox-dlq</code></td>
                <td>14-day compliance audit retention for failed or malformed transaction events</td>
              </tr>
              <tr>
                <td><strong>Backup Vault</strong></td>
                <td><code>aws_backup_vault.banking_vault</code></td>
                <td><code>a-bank-financial-audit-vault</code></td>
                <td>WORM compliance vault enforcing 35-day financial cycle backup retention</td>
              </tr>
              <tr>
                <td><strong>SLO Alarms</strong></td>
                <td><code>aws_cloudwatch_metric_alarm</code> (x3)</td>
                <td><code>alb_5xx_errors</code>, <code>latency</code>, <code>unhealthy</code></td>
                <td>Real-time CloudWatch alerting on 5xx breaches, P99 latency &gt;250ms, and pod faults</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div class="callout">
          <strong>FinOps Run-Rate Analysis</strong>: Using ECS Fargate Spot for burst capacity reduces compute costs by 70%. Total steady-state cloud spend for this multi-AZ resilient banking architecture is <strong>~$38.81 / month</strong>.
        </div>
      </div>
    </div>

    <!-- TAB 5: DEVSECOPS & PCI-DSS V4.0 -->
    <div id="tab-devsecops" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Four-Stage DevSecOps CI/CD Pipeline &amp; Compliance Audit</span>
          <span class="spec-badge green">GITHUB ACTIONS VERIFIED</span>
        </div>

        <div class="grid-4">
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Stage 1: Quality Gate</span>
            <div class="metric-value" style="font-size: 1.1rem; color: var(--success-green);">PASSED</div>
            <p class="account-type">Go vet, -race tests, Helm lint, promtool &amp; kubeconform</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Stage 2: Container Security</span>
            <div class="metric-value" style="font-size: 1.1rem; color: var(--success-green);">0 CVEs</div>
            <p class="account-type">Aqua Security Trivy Scan (0 Critical / High)</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Stage 3: IaC Security</span>
            <div class="metric-value" style="font-size: 1.1rem; color: var(--success-green);">PASSED</div>
            <p class="account-type">Aqua Security tfsec Static Code Analysis</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Stage 4: Supply Chain SBOM</span>
            <div class="metric-value" style="font-size: 1.1rem; color: var(--success-green);">SPDX-JSON</div>
            <p class="account-type">Anchore Syft SBOM generated &amp; archived</p>
          </div>
        </div>

        <div class="card" style="margin-top: 1rem; margin-bottom: 0;">
          <div class="card-header">
            <span class="card-title">PCI-DSS v4.0 Compliance Verification Matrix</span>
          </div>
          <div class="table-container">
            <table>
              <thead>
                <tr>
                  <th>PCI-DSS Requirement</th>
                  <th>Compliance Specification</th>
                  <th>Codified Implementation</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td><strong>Requirement 1.2</strong></td>
                  <td>Network security controls &amp; micro-segmentation</td>
                  <td>ECS Fargate security group permits inbound traffic <em>strictly</em> from ALB security group ID</td>
                </tr>
                <tr>
                  <td><strong>Requirement 3.5</strong></td>
                  <td>Cryptographic architecture &amp; key management</td>
                  <td>AWS KMS Customer Managed Key (CMK) with automated 365-day rotation (<code>alias/a-bank-wallet-cmk</code>)</td>
                </tr>
                <tr>
                  <td><strong>Requirement 6.3.2</strong></td>
                  <td>Software supply chain inventory &amp; SBOM</td>
                  <td>Anchore Syft scans distroless image and publishes SPDX-JSON SBOM artifact with 30-day retention</td>
                </tr>
                <tr>
                  <td><strong>Requirement 6.4.3</strong></td>
                  <td>Runtime attack surface minimization</td>
                  <td>Distroless <code>scratch</code> container (&lt;15MB) with zero shell binaries and unprivileged UID 10001</td>
                </tr>
                <tr>
                  <td><strong>Requirement 10.5</strong></td>
                  <td>Audit trail integrity &amp; write-once retention</td>
                  <td>AWS Backup Vault with 35-day retention lifecycle and continuous DynamoDB PITR streaming</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 6: KUBERNETES, GITOPS & FLEET -->
    <div id="tab-k8s" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Hybrid Sovereign Cloud &amp; Datacenter Parity (Helm, ArgoCD &amp; Ansible)</span>
          <span class="spec-badge blue">CLOUD-AGNOSTIC ARCHITECTURE</span>
        </div>

        <div class="grid-3">
          <div class="card" style="margin-bottom: 0;">
            <div class="card-header">
              <span class="card-title">Helm v3 Chart</span>
              <span class="spec-badge">charts/wallet-service</span>
            </div>
            <ul style="font-size: 0.8rem; color: var(--text-body); list-style-position: inside; line-height: 1.8;">
              <li><strong>PodDisruptionBudget</strong>: <code>minAvailable: 1</code></li>
              <li><strong>Multi-AZ Topology Spread</strong>: Anti-affinity across zones</li>
              <li><strong>NetworkPolicy</strong>: Default-deny with strict ingress/egress</li>
              <li><strong>KEDA Autoscaling</strong>: Scales on SQS queue depth</li>
            </ul>
          </div>

          <div class="card" style="margin-bottom: 0;">
            <div class="card-header">
              <span class="card-title">ArgoCD GitOps</span>
              <span class="spec-badge">gitops/argocd-application.yaml</span>
            </div>
            <ul style="font-size: 0.8rem; color: var(--text-body); list-style-position: inside; line-height: 1.8;">
              <li><strong>Self-Healing</strong>: <code>selfHeal: true</code> auto-reverts drift</li>
              <li><strong>Orphan Pruning</strong>: <code>prune: true</code> deletes stale resources</li>
              <li><strong>Zero-Human-Kubectl</strong>: Git is sole production truth</li>
              <li><strong>Validation</strong>: Verified via <code>kubeconform</code></li>
            </ul>
          </div>

          <div class="card" style="margin-bottom: 0;">
            <div class="card-header">
              <span class="card-title">Ansible Fleet Hardening</span>
              <span class="spec-badge">ansible/playbooks/site.yaml</span>
            </div>
            <ul style="font-size: 0.8rem; color: var(--text-body); list-style-position: inside; line-height: 1.8;">
              <li><strong>Socket Exhaustion Defense</strong>: <code>net.ipv4.ip_local_port_range</code></li>
              <li><strong>SYN Flood Defense</strong>: <code>net.ipv4.tcp_syncookies = 1</code></li>
              <li><strong>TCP Buffer Ceiling</strong>: <code>net.core.somaxconn = 65535</code></li>
              <li><strong>K3s Bootstrapping</strong>: Idempotent edge fleet runtime</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 7: SRE GOLDEN SIGNALS & SLOS -->
    <div id="tab-slo" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Financial SRE Golden Signals &amp; 99.95% Availability SLO Burn Rate</span>
          <span class="spec-badge green">PROMETHEUS &amp; GRAFANA CODIFIED</span>
        </div>

        <div class="grid-4">
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">P99 Latency Target</span>
            <div class="metric-value">&lt; 250 ms</div>
            <p class="account-type">SLI across all payment API routes</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Availability Target</span>
            <div class="metric-value">99.95%</div>
            <p class="account-type">30-day rolling evaluation window</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Monthly Error Budget</span>
            <div class="metric-value">21.6 min</div>
            <p class="account-type">Total allowable outage downtime</p>
          </div>
          <div class="card" style="margin-bottom: 0;">
            <span class="account-type">Critical Page Threshold</span>
            <div class="metric-value" style="color: var(--danger-red);">14.4x Burn</div>
            <p class="account-type">Consumes 2% budget in 1 hour</p>
          </div>
        </div>

        <div class="card" style="margin-top: 1rem; margin-bottom: 0;">
          <div class="card-header">
            <span class="card-title">Live Prometheus Metrics Stream (/metrics)</span>
            <span class="spec-badge">GET /metrics</span>
          </div>
          <div class="console-box" id="metricsRawView">// Polling Prometheus telemetry stream...</div>
        </div>
      </div>
    </div>

    <!-- TAB 8: VERIFICATION RUNBOOK & CLI -->
    <div id="tab-runbook" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Interviewer CLI Verification Runbook</span>
          <span class="spec-badge">ZERO-TRUST VERIFICATION</span>
        </div>
        <p style="font-size: 0.8rem; color: var(--text-muted); margin-bottom: 1rem;">
          You can independently verify every capability demonstrated on this console using standard curl commands in any terminal:
        </p>

        <div style="margin-bottom: 1rem;">
          <p style="font-size: 0.75rem; font-weight: 700; color: var(--text-heading); margin-bottom: 0.35rem;">1. Zero-Trust Liveness &amp; Readiness Health Probes</p>
          <div class="console-box" style="min-height: auto; padding: 0.75rem;">curl -i "https://wallet.thaw-zin-2k77.de5.net/healthz"
curl -i "https://wallet.thaw-zin-2k77.de5.net/readyz"</div>
        </div>

        <div style="margin-bottom: 1rem;">
          <p style="font-size: 0.75rem; font-weight: 700; color: var(--text-heading); margin-bottom: 0.35rem;">2. Trip Payment Rail Circuit Breaker (Chaos Simulation)</p>
          <div class="console-box" style="min-height: auto; padding: 0.75rem;">curl -s -X POST "https://wallet.thaw-zin-2k77.de5.net/api/v1/resilience/circuit-breaker/trip" | jq .
curl -i "https://wallet.thaw-zin-2k77.de5.net/readyz" # Notice HTTP 503 Service Unavailable</div>
        </div>

        <div style="margin-bottom: 1rem;">
          <p style="font-size: 0.75rem; font-weight: 700; color: var(--text-heading); margin-bottom: 0.35rem;">3. Reset Payment Rail Circuit Breaker (Self-Healing)</p>
          <div class="console-box" style="min-height: auto; padding: 0.75rem;">curl -s -X POST "https://wallet.thaw-zin-2k77.de5.net/api/v1/resilience/circuit-breaker/reset" | jq .
curl -i "https://wallet.thaw-zin-2k77.de5.net/readyz" # Notice HTTP 200 OK restored</div>
        </div>

        <div style="margin-bottom: 1rem;">
          <p style="font-size: 0.75rem; font-weight: 700; color: var(--text-heading); margin-bottom: 0.35rem;">4. Double-Entry Transfer with Network Idempotency Replay</p>
          <div class="console-box" style="min-height: auto; padding: 0.75rem;">curl -s -X POST "https://wallet.thaw-zin-2k77.de5.net/api/v1/wallets/transfer" \
     -H "Content-Type: application/json" \
     -H "X-Idempotency-Key: GOLDEN-TEST-001" \
     -d '{"from_account":"ACC-1001","to_account":"ACC-2002","amount":10000,"currency":"MMK","reference":"Golden Test"}' | jq .</div>
        </div>

        <div>
          <p style="font-size: 0.75rem; font-weight: 700; color: var(--text-heading); margin-bottom: 0.35rem;">5. Public Git Source &amp; Release Notes</p>
          <p style="font-size: 0.8rem; font-family: var(--font-mono); color: var(--primary-navy);">
            https://github.com/pretamane/fintech-wallet-sre-showcase/releases/tag/v1.0.0-enterprise
          </p>
        </div>
      </div>
    </div>
  </div>

  <script>
    function showTab(targetId) {
      document.querySelectorAll('.nav-tab').forEach(t => t.classList.remove('active'));
      document.querySelectorAll('.tab-content').forEach(c => c.classList.remove('active'));
      event.target.classList.add('active');
      document.getElementById(targetId).classList.add('active');
      if (targetId === 'tab-slo') fetchTelemetry();
      if (targetId === 'tab-resilience') fetchCircuitStatus();
    }

    function selectAmount(val) {
      document.getElementById('transferAmount').value = val;
    }

    function generateNewKey() {
      const generated = 'TX-' + Math.random().toString(36).substring(2, 9).toUpperCase() + '-' + Date.now().toString().slice(-4);
      document.getElementById('idempotencyKey').value = generated;
    }

    async function probeLatency() {
      const startTime = performance.now();
      try {
        await fetch('/healthz');
        const elapsed = Math.round(performance.now() - startTime);
        document.getElementById('latencyTicker').innerText = 'Latency: ' + elapsed + ' ms';
      } catch (err) {
        document.getElementById('latencyTicker').innerText = 'Latency: N/A';
      }
    }

    async function queryBalance(accId) {
      try {
        const response = await fetch('/api/v1/wallets/' + accId + '/balance');
        if (response.ok) {
          const payload = await response.json();
          document.getElementById('display-' + accId).innerText = Number(payload.balance).toLocaleString() + ' MMK';
        }
      } catch (e) {}
    }

    async function updateAllBalances() {
      await queryBalance('ACC-1001');
      await queryBalance('ACC-2002');
      await queryBalance('ACC-9999');
      probeLatency();
    }

    async function fetchCircuitStatus() {
      try {
        const response = await fetch('/api/v1/resilience/circuit-breaker');
        if (response.ok) {
          const data = await response.json();
          const state = data.state || 'CLOSED';
          const nodeClosed = document.getElementById('node-closed');
          const nodeOpen = document.getElementById('node-open');
          const nodeHalf = document.getElementById('node-halfopen');
          const stateDisplay = document.getElementById('cb-state-display');
          const circuitTicker = document.getElementById('circuitTicker');
          const statusDot = document.getElementById('statusDot');
          const failDisplay = document.getElementById('cb-failures-display');

          nodeClosed.className = 'breaker-node' + (state === 'CLOSED' ? ' active-closed' : '');
          nodeOpen.className = 'breaker-node' + (state === 'OPEN' ? ' active-open' : '');
          nodeHalf.className = 'breaker-node' + (state === 'HALF-OPEN' ? ' active-halfopen' : '');

          stateDisplay.innerText = state;
          stateDisplay.style.color = state === 'CLOSED' ? 'var(--success-green)' : (state === 'OPEN' ? 'var(--danger-red)' : 'var(--warning-amber)');
          circuitTicker.innerText = 'Rail: ' + state;

          statusDot.className = 'status-indicator' + (state === 'OPEN' ? ' open' : (state === 'HALF-OPEN' ? ' half-open' : ''));
          failDisplay.innerText = (data.consecutive_failures || 0) + ' / ' + (data.failure_threshold || 5);
        }
      } catch (e) {}
    }

    async function triggerCircuitTrip() {
      const consoleEl = document.getElementById('chaosConsole');
      consoleEl.innerText = '>>> POST /api/v1/resilience/circuit-breaker/trip\n>>> Injecting synthetic partner outage into CBM-Net clearing rail...';
      try {
        const response = await fetch('/api/v1/resilience/circuit-breaker/trip', { method: 'POST' });
        const data = await response.json();
        consoleEl.innerText = '<<< HTTP ' + response.status + ' ' + response.statusText + '\n' + JSON.stringify(data, null, 2);
        await fetchCircuitStatus();
        await probeReadinessAndLiveness();
      } catch (err) {
        consoleEl.innerText = '<<< Error tripping circuit breaker: ' + err.message;
      }
    }

    async function triggerCircuitReset() {
      const consoleEl = document.getElementById('chaosConsole');
      consoleEl.innerText = '>>> POST /api/v1/resilience/circuit-breaker/reset\n>>> Resetting circuit breaker to CLOSED...';
      try {
        const response = await fetch('/api/v1/resilience/circuit-breaker/reset', { method: 'POST' });
        const data = await response.json();
        consoleEl.innerText = '<<< HTTP ' + response.status + ' ' + response.statusText + '\n' + JSON.stringify(data, null, 2);
        await fetchCircuitStatus();
        await probeReadinessAndLiveness();
      } catch (err) {
        consoleEl.innerText = '<<< Error resetting circuit breaker: ' + err.message;
      }
    }

    async function probeReadinessAndLiveness() {
      const consoleEl = document.getElementById('chaosConsole');
      let log = consoleEl.innerText + '\n\n>>> Probing Liveness & Readiness Endpoints:';
      try {
        const livenessRes = await fetch('/healthz');
        log += '\n[Liveness  /healthz] Status: ' + livenessRes.status + ' ' + livenessRes.statusText + ' (Container Healthy)';
        
        const readinessRes = await fetch('/readyz');
        const readData = await readinessRes.json();
        log += '\n[Readiness /readyz ] Status: ' + readinessRes.status + ' ' + readinessRes.statusText + ' -> ' + JSON.stringify(readData);
        consoleEl.innerText = log;
      } catch (err) {
        consoleEl.innerText = log + '\nProbe Failed: ' + err.message;
      }
    }

    async function simulateTrafficBurst(count) {
      const consoleEl = document.getElementById('chaosConsole');
      consoleEl.innerText = '>>> Firing ' + count + ' concurrent burst requests to test WAF & ALB concurrency...';
      const promises = [];
      const start = performance.now();
      for (let i = 0; i < count; i++) {
        promises.push(fetch('/healthz', { cache: 'no-store' }));
      }
      const results = await Promise.all(promises);
      const elapsed = Math.round(performance.now() - start);
      const successful = results.filter(r => r.ok).length;
      consoleEl.innerText = '<<< Completed ' + count + ' concurrent requests in ' + elapsed + ' ms.\n<<< Successful 200 OK responses: ' + successful + '/' + count + '\n<<< Average latency per request: ' + (elapsed / count).toFixed(1) + ' ms';
      probeLatency();
    }

    async function fetchTelemetry() {
      try {
        const response = await fetch('/metrics');
        if (response.ok) {
          const data = await response.json();
          document.getElementById('metricsRawView').innerText = JSON.stringify(data, null, 2);
        }
      } catch (e) {}
    }

    async function executeTransfer(isDuplicateSimulation) {
      const fromAcc = document.getElementById('sourceAccount').value;
      const toAcc = document.getElementById('targetAccount').value;
      const amount = parseFloat(document.getElementById('transferAmount').value);
      let key = document.getElementById('idempotencyKey').value.trim();

      if (!key) {
        generateNewKey();
        key = document.getElementById('idempotencyKey').value;
      }

      const consoleEl = document.getElementById('auditConsole');
      const badgeEl = document.getElementById('responseStatusBadge');
      
      consoleEl.innerText = '>>> POST /api/v1/wallets/transfer\n>>> X-Idempotency-Key: ' + key + '\n>>> Submitting transaction payload to ledger...';
      badgeEl.style.display = 'none';

      try {
        const response = await fetch('/api/v1/wallets/transfer', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Idempotency-Key': key
          },
          body: JSON.stringify({
            from_account: fromAcc,
            to_account: toAcc,
            amount: amount,
            currency: 'MMK',
            reference: 'Console Transfer'
          })
        });

        const result = await response.json();
        consoleEl.innerText = '<<< HTTP ' + response.status + ' ' + response.statusText + '\n' + JSON.stringify(result, null, 2);
        badgeEl.style.display = 'inline-block';

        if (response.ok) {
          if (result.idempotent_replay) {
            badgeEl.className = 'badge-execution badge-warning';
            badgeEl.innerText = 'STATUS: 200 OK (IDEMPOTENT_REPLAY)';
          } else {
            badgeEl.className = 'badge-execution badge-success';
            badgeEl.innerText = 'STATUS: 200 OK (NEW_TRANSACTION)';
            if (!isDuplicateSimulation) {
              generateNewKey();
            }
          }
        } else {
          badgeEl.className = 'badge-execution badge-danger';
          badgeEl.innerText = 'STATUS: ' + response.status + ' ERROR';
        }

        updateAllBalances();
      } catch (err) {
        consoleEl.innerText = '<<< Network Connection Failure: ' + err.message;
        badgeEl.style.display = 'inline-block';
        badgeEl.className = 'badge-execution badge-danger';
        badgeEl.innerText = 'CONNECTION ERROR';
      }
    }

    async function simulateInsufficientFunds() {
      document.getElementById('sourceAccount').value = 'ACC-1001';
      document.getElementById('targetAccount').value = 'ACC-2002';
      document.getElementById('transferAmount').value = '9999999999';
      generateNewKey();
      await executeTransfer(false);
    }

    window.onload = function() {
      generateNewKey();
      updateAllBalances();
      fetchCircuitStatus();
      setInterval(probeLatency, 8000);
      setInterval(fetchCircuitStatus, 5000);
    };
  </script>
</body>
</html>`
