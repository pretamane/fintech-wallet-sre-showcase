package main

const rootDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>A Bank Payment Services • Internal Operations Console</title>
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
      padding: 2rem 1.5rem;
    }
    .container { max-width: 1200px; margin: 0 auto; }
    
    /* Header */
    header {
      background-color: var(--bg-card);
      border: 1px solid var(--border-light);
      border-radius: 8px;
      padding: 1.5rem 2rem;
      margin-bottom: 1.5rem;
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
      font-weight: 700;
      font-size: 0.8rem;
      padding: 0.4rem 0.75rem;
      border-radius: 4px;
      letter-spacing: 0.08em;
    }
    h1 {
      font-size: 1.35rem;
      font-weight: 700;
      color: var(--text-heading);
      letter-spacing: -0.01em;
    }
    .system-status { display: flex; align-items: center; gap: 1rem; }
    .status-tag {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background-color: var(--success-bg);
      color: var(--success-green);
      border: 1px solid #a7f3d0;
      padding: 0.35rem 0.85rem;
      border-radius: 4px;
      font-size: 0.8rem;
      font-weight: 600;
    }
    .status-indicator {
      width: 7px;
      height: 7px;
      background-color: var(--success-green);
      border-radius: 50%;
    }
    .latency-tag {
      font-family: var(--font-mono);
      font-size: 0.8rem;
      color: var(--text-muted);
      border-left: 1px solid var(--border-light);
      padding-left: 1rem;
    }

    /* Tabs Navigation */
    .nav-tabs {
      display: flex;
      gap: 0.25rem;
      border-bottom: 1px solid var(--border-light);
      margin-bottom: 1.5rem;
      overflow-x: auto;
    }
    .nav-tab {
      background: none;
      border: none;
      border-bottom: 2px solid transparent;
      color: var(--text-muted);
      font-size: 0.9rem;
      font-weight: 600;
      padding: 0.75rem 1.25rem;
      cursor: pointer;
      transition: all 0.15s ease;
      white-space: nowrap;
    }
    .nav-tab:hover { color: var(--text-heading); }
    .nav-tab.active {
      color: var(--primary-navy);
      border-bottom-color: var(--primary-navy);
    }

    /* Layout Grids */
    .grid-3 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
      gap: 1.25rem;
      margin-bottom: 1.5rem;
    }
    .grid-2 {
      display: grid;
      grid-template-columns: 1fr 1fr;
      gap: 1.25rem;
      margin-bottom: 1.5rem;
    }
    @media (max-width: 900px) {
      .grid-2 { grid-template-columns: 1fr; }
    }

    /* Cards */
    .card {
      background-color: var(--bg-card);
      border: 1px solid var(--border-light);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.05);
    }
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1rem;
      padding-bottom: 0.75rem;
      border-bottom: 1px solid var(--border-light);
    }
    .card-title {
      font-size: 1rem;
      font-weight: 600;
      color: var(--text-heading);
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
    .balance-display {
      font-family: var(--font-mono);
      font-size: 1.75rem;
      font-weight: 700;
      color: var(--primary-navy);
      margin: 0.5rem 0;
    }
    .account-type {
      font-size: 0.8rem;
      color: var(--text-muted);
    }

    /* Forms */
    .form-group { margin-bottom: 1.1rem; }
    label {
      display: block;
      font-size: 0.8rem;
      font-weight: 600;
      color: var(--text-heading);
      margin-bottom: 0.4rem;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    input, select {
      width: 100%;
      background-color: #ffffff;
      border: 1px solid var(--border-dark);
      border-radius: 5px;
      padding: 0.65rem 0.85rem;
      color: var(--text-heading);
      font-family: inherit;
      font-size: 0.9rem;
      transition: border-color 0.15s;
    }
    input:focus, select:focus {
      outline: none;
      border-color: var(--accent-blue);
      box-shadow: 0 0 0 3px rgba(2, 132, 199, 0.12);
    }
    .chip-group {
      display: flex;
      gap: 0.5rem;
      margin-top: 0.4rem;
    }
    .amount-chip {
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      color: var(--text-heading);
      font-family: var(--font-mono);
      font-size: 0.75rem;
      padding: 0.25rem 0.6rem;
      border-radius: 4px;
      cursor: pointer;
      transition: background-color 0.15s;
    }
    .amount-chip:hover {
      background-color: var(--border-light);
    }

    /* Buttons */
    .action-group {
      display: flex;
      flex-wrap: wrap;
      gap: 0.75rem;
      margin-top: 1.5rem;
    }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      font-size: 0.85rem;
      font-weight: 600;
      padding: 0.65rem 1.25rem;
      border-radius: 5px;
      border: 1px solid transparent;
      cursor: pointer;
      transition: all 0.15s;
    }
    .btn-primary {
      background-color: var(--primary-navy);
      color: #ffffff;
    }
    .btn-primary:hover {
      background-color: var(--primary-hover);
    }
    .btn-secondary {
      background-color: #ffffff;
      border-color: var(--border-dark);
      color: var(--text-heading);
    }
    .btn-secondary:hover {
      background-color: var(--bg-subtle);
    }
    .btn-warning {
      background-color: #ffffff;
      border-color: var(--warning-amber);
      color: var(--warning-amber);
    }
    .btn-warning:hover {
      background-color: var(--warning-bg);
    }

    /* Terminal Console */
    .console-box {
      background-color: #0f172a;
      color: #38bdf8;
      font-family: var(--font-mono);
      font-size: 0.8rem;
      line-height: 1.6;
      padding: 1.25rem;
      border-radius: 6px;
      border: 1px solid #1e293b;
      min-height: 280px;
      max-height: 380px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
    .badge-execution {
      display: inline-block;
      font-family: var(--font-mono);
      font-size: 0.75rem;
      font-weight: 700;
      padding: 0.25rem 0.6rem;
      border-radius: 4px;
    }
    .badge-success { background-color: var(--success-bg); color: var(--success-green); border: 1px solid #a7f3d0; }
    .badge-warning { background-color: var(--warning-bg); color: var(--warning-amber); border: 1px solid #fde68a; }
    .badge-danger  { background-color: var(--danger-bg);  color: var(--danger-red);    border: 1px solid #fecaca; }

    /* Tables */
    .table-container { overflow-x: auto; margin-top: 1rem; }
    table {
      width: 100%;
      border-collapse: collapse;
      text-align: left;
    }
    th, td {
      padding: 0.85rem 1rem;
      border-bottom: 1px solid var(--border-light);
      font-size: 0.85rem;
    }
    th {
      background-color: var(--bg-subtle);
      color: var(--text-heading);
      font-weight: 600;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    tr:hover td { background-color: #fafafa; }
    .spec-badge {
      display: inline-block;
      font-family: var(--font-mono);
      font-size: 0.75rem;
      font-weight: 600;
      padding: 0.15rem 0.5rem;
      border-radius: 3px;
      background-color: var(--bg-subtle);
      border: 1px solid var(--border-light);
      color: var(--text-heading);
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
          <h1>Mobile Wallet Operations Console</h1>
          <p style="font-size: 0.85rem; color: var(--text-muted); margin-top: 0.15rem;">
            Core Ledger Engine &bull; AWS ECS Fargate &bull; Cloudflare WAF Ingress &bull; PCI-DSS v4.0 Baseline
          </p>
        </div>
      </div>
      <div class="system-status">
        <span class="status-tag"><span class="status-indicator"></span> CLUSTER ACTIVE</span>
        <span class="latency-tag" id="latencyTicker">Latency: -- ms</span>
      </div>
    </header>

    <!-- Navigation Tabs -->
    <div class="nav-tabs">
      <button class="nav-tab active" onclick="showTab('tab-ledger')">Ledger Operations</button>
      <button class="nav-tab" onclick="showTab('tab-telemetry')">System Telemetry &amp; Metrics</button>
      <button class="nav-tab" onclick="showTab('tab-architecture')">Cloud Architecture Specifications</button>
      <button class="nav-tab" onclick="showTab('tab-runbook')">Verification Runbook</button>
    </div>

    <!-- TAB 1: LEDGER OPERATIONS -->
    <div id="tab-ledger" class="tab-content active">
      <!-- Balance Cards -->
      <div class="grid-3">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Customer Account A</span>
            <span class="account-id">ACC-1001</span>
          </div>
          <div class="balance-display" id="display-ACC-1001">-- MMK</div>
          <p class="account-type">Retail Tier 1 User Profile</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Merchant Partner B</span>
            <span class="account-id">ACC-2002</span>
          </div>
          <div class="balance-display" id="display-ACC-2002">-- MMK</div>
          <p class="account-type">Commercial POS Settlement Account</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Central Reserve Pool</span>
            <span class="account-id">ACC-9999</span>
          </div>
          <div class="balance-display" id="display-ACC-9999">-- MMK</div>
          <p class="account-type">Master Settlement Liquidity Pool</p>
        </div>
      </div>

      <!-- Transaction Form & Audit Output -->
      <div class="grid-2">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Execute Double-Entry Transaction</span>
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
            </div>
          </div>

          <div class="form-group">
            <label>Idempotency Key (X-Idempotency-Key)</label>
            <div style="display: flex; gap: 0.5rem;">
              <input type="text" id="idempotencyKey" style="font-family: var(--font-mono); font-size: 0.85rem;">
              <button class="btn btn-secondary" onclick="generateNewKey()" type="button" style="white-space: nowrap;">Generate Key</button>
            </div>
            <p style="font-size: 0.75rem; color: var(--text-muted); margin-top: 0.35rem;">
              Enforces transaction deduplication across mobile network retries.
            </p>
          </div>

          <div class="action-group">
            <button class="btn btn-primary" onclick="executeTransfer(false)">Execute Transaction</button>
            <button class="btn btn-warning" onclick="executeTransfer(true)" title="Submits with the same key to verify idempotency replay">Simulate Duplicate Request</button>
            <button class="btn btn-secondary" onclick="updateAllBalances()">Refresh Balances</button>
          </div>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Transaction Receipt &amp; Audit Log</span>
            <span id="responseStatusBadge" class="badge-execution" style="display: none;"></span>
          </div>
          <div class="console-box" id="auditConsole">// Ready. Select accounts and click 'Execute Transaction' above.</div>
          <div style="margin-top: 1rem; padding: 0.75rem; background-color: var(--bg-subtle); border-radius: 4px; font-size: 0.75rem; color: var(--text-muted);">
            <strong>PCI-DSS Audit Trail</strong>: Every request records atomic state change, timestamp, and idempotency status. The 'Simulate Duplicate Request' action proves that duplicate network submissions receive cached receipts with zero balance alteration.
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: TELEMETRY & METRICS -->
    <div id="tab-telemetry" class="tab-content">
      <div class="grid-3">
        <div class="card">
          <div class="card-header">
            <span class="card-title">Total Transactions</span>
          </div>
          <div class="balance-display" id="metric-txCount">0</div>
          <p class="account-type">Successfully executed transactions</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Total Settled Volume</span>
          </div>
          <div class="balance-display" id="metric-totalVolume">0 MMK</div>
          <p class="account-type">Cumulative processed turnover</p>
        </div>

        <div class="card">
          <div class="card-header">
            <span class="card-title">Active Idempotency Records</span>
          </div>
          <div class="balance-display" id="metric-idempKeys">0</div>
          <p class="account-type">Cached transaction signatures</p>
        </div>
      </div>

      <div class="card">
        <div class="card-header">
          <span class="card-title">Prometheus Telemetry Feed (/metrics)</span>
          <span class="spec-badge">GET /metrics</span>
        </div>
        <div class="console-box" id="metricsRawView">// Polling telemetry stream...</div>
      </div>
    </div>

    <!-- TAB 3: ARCHITECTURE SPECIFICATIONS -->
    <div id="tab-architecture" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">Tier-1 Mobile Wallet Cloud Architecture Mapping</span>
        </div>
        <p style="font-size: 0.85rem; color: var(--text-muted); margin-bottom: 1rem;">
          Cross-platform mapping between Tier-1 payment platforms (KBZPay on Huawei Cloud, WeChat Pay on Tencent Cloud) and this AWS Managed Services implementation:
        </p>

        <div class="table-container">
          <table>
            <thead>
              <tr>
                <th>Platform Layer</th>
                <th>KBZPay (Huawei Cloud)</th>
                <th>WeChat Pay / Alipay</th>
                <th>A Bank Architecture (AWS Showcase)</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td><strong>Edge Security &amp; WAF</strong></td>
                <td>Huawei Cloud WAF + Anti-DDoS</td>
                <td>Tencent Dayu Shield WAF</td>
                <td><span class="spec-badge">Cloudflare Anycast</span> WAF + Origin Port Rewrite</td>
              </tr>
              <tr>
                <td><strong>Container Compute</strong></td>
                <td>Cloud Container Engine (CCE) + CCI</td>
                <td>Tencent TKE Elastic Matrix</td>
                <td><span class="spec-badge">AWS ECS Fargate</span> (Serverless 0.25 vCPU, 512 MB)</td>
              </tr>
              <tr>
                <td><strong>Event Callbacks &amp; Webhooks</strong></td>
                <td>Huawei FunctionGraph</td>
                <td>Serverless Cloud Functions (SCF)</td>
                <td><span class="spec-badge">AWS Lambda</span> + Amazon EventBridge</td>
              </tr>
              <tr>
                <td><strong>Financial Ledger Storage</strong></td>
                <td>Huawei GaussDB (Distributed Relational)</td>
                <td>Tencent TDSQL / Ant OceanBase</td>
                <td><span class="spec-badge">Amazon DynamoDB</span> (On-Demand ACID + PITR)</td>
              </tr>
              <tr>
                <td><strong>Infrastructure as Code</strong></td>
                <td>Huawei Cloud Terraform Provider</td>
                <td>Terraform Enterprise</td>
                <td><span class="spec-badge">HashiCorp Terraform</span> v1.5+ (AWS + Cloudflare)</td>
              </tr>
              <tr>
                <td><strong>On-Premises Edge Fleet</strong></td>
                <td>Bare-Metal Branch Infrastructure</td>
                <td>Hybrid Private Cloud Nodes</td>
                <td><span class="spec-badge">Red Hat Ansible</span> (CIS Hardening + K3s Fleet)</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- TAB 4: RUNBOOK -->
    <div id="tab-runbook" class="tab-content">
      <div class="card">
        <div class="card-header">
          <span class="card-title">SRE Golden Test CLI Verification Runbook</span>
        </div>
        <p style="font-size: 0.85rem; color: var(--text-muted); margin-bottom: 1.5rem;">
          The running cluster can be validated independently via standard terminal commands:
        </p>

        <div style="margin-bottom: 1.25rem;">
          <p style="font-size: 0.8rem; font-weight: 600; color: var(--text-heading); margin-bottom: 0.35rem;">1. Zero-Trust Health Check</p>
          <div class="console-box" style="min-height: auto; padding: 0.85rem;">curl -i "https://wallet.thaw-zin-2k77.de5.net/healthz"</div>
        </div>

        <div style="margin-bottom: 1.25rem;">
          <p style="font-size: 0.8rem; font-weight: 600; color: var(--text-heading); margin-bottom: 0.35rem;">2. Account Balance Verification</p>
          <div class="console-box" style="min-height: auto; padding: 0.85rem;">curl -s "https://wallet.thaw-zin-2k77.de5.net/api/v1/wallets/ACC-1001/balance" | jq .</div>
        </div>

        <div style="margin-bottom: 1.25rem;">
          <p style="font-size: 0.8rem; font-weight: 600; color: var(--text-heading); margin-bottom: 0.35rem;">3. Double-Entry Transfer with Idempotency Key</p>
          <div class="console-box" style="min-height: auto; padding: 0.85rem;">curl -s -X POST "https://wallet.thaw-zin-2k77.de5.net/api/v1/wallets/transfer" \
     -H "Content-Type: application/json" \
     -H "X-Idempotency-Key: GOLDEN-TEST-001" \
     -d '{"from_account":"ACC-1001","to_account":"ACC-2002","amount":10000,"currency":"MMK","reference":"Golden Test"}' | jq .</div>
        </div>

        <div>
          <p style="font-size: 0.8rem; font-weight: 600; color: var(--text-heading); margin-bottom: 0.35rem;">4. Public Git Repository &amp; Pipeline Source</p>
          <p style="font-size: 0.85rem; font-family: var(--font-mono); color: var(--primary-navy);">
            https://github.com/pretamane/fintech-wallet-sre-showcase
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
      if (targetId === 'tab-telemetry') fetchTelemetry();
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

    async function fetchTelemetry() {
      try {
        const response = await fetch('/metrics');
        if (response.ok) {
          const data = await response.json();
          document.getElementById('metric-txCount').innerText = data.total_transactions || 0;
          document.getElementById('metric-totalVolume').innerText = (data.total_volume || 0).toLocaleString() + ' MMK';
          document.getElementById('metric-idempKeys').innerText = data.idempotency_keys || 0;
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

    window.onload = function() {
      generateNewKey();
      updateAllBalances();
      setInterval(probeLatency, 10000);
    };
  </script>
</body>
</html>`
