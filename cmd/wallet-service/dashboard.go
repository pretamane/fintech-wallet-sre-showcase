package main

const rootDashboardHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>A Bank Mobile Wallet • FinTech SRE Operations Console</title>
  <style>
    :root {
      --bg-base: #0a0f1d;
      --bg-surface: #111827;
      --bg-surface-elevated: #1f2937;
      --border-color: #374151;
      --primary: #10b981;
      --primary-hover: #059669;
      --accent: #38bdf8;
      --warning: #f59e0b;
      --danger: #ef4444;
      --text-main: #f9fafb;
      --text-muted: #9ca3af;
      --font-mono: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace;
      --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    }
    * { box-sizing: border-box; margin: 0; padding: 0; }
    body {
      background-color: var(--bg-base);
      color: var(--text-main);
      font-family: var(--font-sans);
      line-height: 1.5;
      padding: 1.5rem;
    }
    .container { max-width: 1200px; margin: 0 auto; }
    header {
      display: flex;
      flex-wrap: wrap;
      justify-content: space-between;
      align-items: center;
      padding-bottom: 1.5rem;
      border-bottom: 1px solid var(--border-color);
      margin-bottom: 2rem;
      gap: 1rem;
    }
    .brand-title { display: flex; align-items: center; gap: 0.75rem; }
    .brand-badge {
      background: linear-gradient(135deg, #10b981, #0284c7);
      color: #fff;
      font-weight: 800;
      font-size: 0.85rem;
      padding: 0.35rem 0.65rem;
      border-radius: 6px;
      letter-spacing: 0.05em;
    }
    h1 { font-size: 1.5rem; font-weight: 700; color: #fff; }
    .status-pill {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: #064e3b;
      color: #6ee7b7;
      padding: 0.35rem 0.85rem;
      border-radius: 9999px;
      font-size: 0.8rem;
      font-weight: 600;
      border: 1px solid #059669;
    }
    .status-dot {
      width: 8px;
      height: 8px;
      background-color: #34d399;
      border-radius: 50%;
      box-shadow: 0 0 8px #34d399;
      animation: pulse 2s infinite;
    }
    @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.4; } 100% { opacity: 1; } }
    
    /* Navigation Tabs */
    .tabs {
      display: flex;
      gap: 0.5rem;
      border-bottom: 1px solid var(--border-color);
      margin-bottom: 2rem;
      overflow-x: auto;
    }
    .tab-btn {
      background: none;
      border: none;
      color: var(--text-muted);
      font-size: 0.95rem;
      font-weight: 600;
      padding: 0.75rem 1.25rem;
      cursor: pointer;
      border-bottom: 2px solid transparent;
      transition: all 0.2s ease;
      white-space: nowrap;
    }
    .tab-btn:hover { color: var(--text-main); }
    .tab-btn.active {
      color: var(--primary);
      border-bottom-color: var(--primary);
    }
    
    /* Grid Layouts */
    .grid-3 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 1.5rem;
      margin-bottom: 2rem;
    }
    .grid-2 {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
      gap: 1.5rem;
      margin-bottom: 2rem;
    }
    @media (max-width: 768px) {
      .grid-2 { grid-template-columns: 1fr; }
    }
    
    /* Cards */
    .card {
      background-color: var(--bg-surface);
      border: 1px solid var(--border-color);
      border-radius: 10px;
      padding: 1.5rem;
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
    }
    .card-title {
      font-size: 1.1rem;
      font-weight: 600;
      margin-bottom: 1rem;
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .balance-val {
      font-size: 1.8rem;
      font-weight: 700;
      color: var(--primary);
      font-family: var(--font-mono);
      margin: 0.5rem 0;
    }
    .acc-label {
      font-size: 0.8rem;
      color: var(--text-muted);
      text-transform: uppercase;
      letter-spacing: 0.05em;
    }
    
    /* Forms */
    .form-group { margin-bottom: 1rem; }
    label { display: block; font-size: 0.85rem; font-weight: 500; margin-bottom: 0.35rem; color: var(--text-muted); }
    input, select {
      width: 100%;
      background-color: var(--bg-base);
      border: 1px solid var(--border-color);
      border-radius: 6px;
      padding: 0.65rem 0.85rem;
      color: var(--text-main);
      font-family: inherit;
      font-size: 0.95rem;
    }
    input:focus, select:focus {
      outline: none;
      border-color: var(--accent);
      box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
    }
    .btn-group { display: flex; gap: 0.75rem; flex-wrap: wrap; margin-top: 1.25rem; }
    .btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.65rem 1.25rem;
      font-size: 0.9rem;
      font-weight: 600;
      border-radius: 6px;
      border: none;
      cursor: pointer;
      transition: all 0.2s;
    }
    .btn-primary { background-color: var(--primary); color: #fff; }
    .btn-primary:hover { background-color: var(--primary-hover); }
    .btn-warning { background-color: #d97706; color: #fff; }
    .btn-warning:hover { background-color: #b45309; }
    .btn-outline { background: transparent; border: 1px solid var(--border-color); color: var(--text-main); }
    .btn-outline:hover { background-color: var(--bg-surface-elevated); }
    
    /* Quick Chips */
    .chip-container { display: flex; gap: 0.5rem; margin-top: 0.5rem; }
    .chip {
      background-color: var(--bg-surface-elevated);
      border: 1px solid var(--border-color);
      color: var(--accent);
      padding: 0.25rem 0.6rem;
      border-radius: 4px;
      font-size: 0.8rem;
      font-family: var(--font-mono);
      cursor: pointer;
    }
    .chip:hover { background-color: var(--border-color); }
    
    /* Terminal Console Output */
    .terminal-box {
      background-color: #030712;
      border: 1px solid var(--border-color);
      border-radius: 6px;
      padding: 1rem;
      font-family: var(--font-mono);
      font-size: 0.85rem;
      color: #38bdf8;
      max-height: 280px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }
    
    /* Architecture Comparison Table */
    table { width: 100%; border-collapse: collapse; margin-top: 1rem; }
    th, td {
      padding: 0.85rem 1rem;
      text-align: left;
      border-bottom: 1px solid var(--border-color);
      font-size: 0.9rem;
    }
    th { background-color: var(--bg-surface-elevated); color: var(--text-muted); font-weight: 600; }
    tr:hover td { background-color: rgba(255, 255, 255, 0.02); }
    .badge-aws { background: #ff9900; color: #111; padding: 2px 6px; border-radius: 4px; font-weight: 700; font-size: 0.75rem; }
    .badge-huawei { background: #ed1c24; color: #fff; padding: 2px 6px; border-radius: 4px; font-weight: 700; font-size: 0.75rem; }
    .badge-wechat { background: #07c160; color: #fff; padding: 2px 6px; border-radius: 4px; font-weight: 700; font-size: 0.75rem; }
    
    /* Hidden Tab Content */
    .tab-pane { display: none; }
    .tab-pane.active { display: block; }
  </style>
</head>
<body>
  <div class="container">
    <header>
      <div class="brand-title">
        <span class="brand-badge">A-BANK SRE</span>
        <div>
          <h1>Digital Wallet & Super-App Operations Console</h1>
          <p style="font-size: 0.85rem; color: var(--text-muted);">High-Throughput Double-Entry Ledger • Cloudflare Edge WAF • AWS ECS Fargate</p>
        </div>
      </div>
      <div style="display: flex; gap: 0.75rem; align-items: center;">
        <span class="status-pill"><span class="status-dot"></span> LIVE ON FARGATE</span>
        <span id="latencyTag" style="font-size: 0.8rem; font-family: var(--font-mono); color: var(--accent);">Edge RTT: -- ms</span>
      </div>
    </header>

    <!-- Navigation Tabs -->
    <div class="tabs">
      <button class="tab-btn active" onclick="switchTab('tab-ops')">💳 Live Wallet Operations</button>
      <button class="tab-btn" onclick="switchTab('tab-metrics')">📊 SRE Observability & Metrics</button>
      <button class="tab-btn" onclick="switchTab('tab-arch')">☁️ Super-App Cloud Architecture (KBZPay / WeChat vs AWS)</button>
      <button class="tab-btn" onclick="switchTab('tab-runbook')">📖 Golden Test Handover Runbook</button>
    </div>

    <!-- TAB 1: LIVE WALLET OPERATIONS -->
    <div id="tab-ops" class="tab-pane active">
      <!-- Balance Cards -->
      <div class="grid-3">
        <div class="card">
          <div class="card-title">
            <span>Customer Account A</span>
            <span class="acc-label">ACC-1001</span>
          </div>
          <div class="balance-val" id="bal-ACC-1001">-- MMK</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Primary Consumer Wallet Profile</div>
        </div>

        <div class="card">
          <div class="card-title">
            <span>Merchant Partner B</span>
            <span class="acc-label">ACC-2002</span>
          </div>
          <div class="balance-val" id="bal-ACC-2002">-- MMK</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Merchant POS Settlement Account</div>
        </div>

        <div class="card">
          <div class="card-title">
            <span>Central Bank Reserve</span>
            <span class="acc-label">ACC-9999</span>
          </div>
          <div class="balance-val" id="bal-ACC-9999">-- MMK</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Settlement Liquidity Pool</div>
        </div>
      </div>

      <!-- Execution Form & Console -->
      <div class="grid-2">
        <div class="card">
          <div class="card-title">⚡ Execute Double-Entry Atomic Transfer</div>
          
          <div class="form-group">
            <label>Source Account (Debit)</label>
            <select id="fromAcc">
              <option value="ACC-1001">ACC-1001 (Customer A)</option>
              <option value="ACC-2002">ACC-2002 (Merchant B)</option>
              <option value="ACC-9999">ACC-9999 (Central Reserve)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Destination Account (Credit)</label>
            <select id="toAcc">
              <option value="ACC-2002">ACC-2002 (Merchant B)</option>
              <option value="ACC-1001">ACC-1001 (Customer A)</option>
              <option value="ACC-9999">ACC-9999 (Central Reserve)</option>
            </select>
          </div>

          <div class="form-group">
            <label>Amount (MMK)</label>
            <input type="number" id="txAmount" value="25000" min="1" step="500">
            <div class="chip-container">
              <span class="chip" onclick="setAmount(5000)">5,000</span>
              <span class="chip" onclick="setAmount(25000)">25,000</span>
              <span class="chip" onclick="setAmount(50000)">50,000</span>
              <span class="chip" onclick="setAmount(100000)">100,000</span>
            </div>
          </div>

          <div class="form-group">
            <label>Idempotency Key (X-Idempotency-Key)</label>
            <div style="display: flex; gap: 0.5rem;">
              <input type="text" id="idempKey" value="" style="font-family: var(--font-mono);">
              <button class="btn btn-outline" onclick="generateUUID()">🎲 New</button>
            </div>
            <p style="font-size: 0.75rem; color: var(--text-muted); margin-top: 0.25rem;">
              PCI-DSS Mandate: Prevents double-debiting on mobile network timeouts.
            </p>
          </div>

          <div class="btn-group">
            <button class="btn btn-primary" onclick="submitTransfer(false)">🚀 Execute Transfer</button>
            <button class="btn btn-warning" onclick="submitTransfer(true)" title="Resends the EXACT same key to prove idempotency">🔁 Replay Same Key (Attack Test)</button>
            <button class="btn btn-outline" onclick="refreshBalances()">🔄 Refresh</button>
          </div>
        </div>

        <div class="card">
          <div class="card-title">
            <span>Audit Trail & Terminal Output</span>
            <span id="txStatusBadge" style="font-size: 0.8rem; font-family: var(--font-mono); color: var(--text-muted);">Idle</span>
          </div>
          <div class="terminal-box" id="terminalOutput">// Ready for execution. Click 'Execute Transfer' above.</div>
          <div style="margin-top: 1rem; font-size: 0.8rem; color: var(--text-muted);">
            💡 <strong>SRE Pro-Tip</strong>: Click <em>"Replay Same Key"</em> to simulate a customer double-tapping the transfer button or a 4G packet replay attack. You will see <code>"idempotent_replay": true</code> and the balance will NOT debit twice.
          </div>
        </div>
      </div>
    </div>

    <!-- TAB 2: SRE OBSERVABILITY & METRICS -->
    <div id="tab-metrics" class="tab-pane">
      <div class="grid-3">
        <div class="card">
          <div class="card-title">Total Transactions</div>
          <div class="balance-val" id="metric-txCount">0</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Executed in-memory & logged</div>
        </div>
        <div class="card">
          <div class="card-title">Total Volume Settled</div>
          <div class="balance-val" id="metric-volume">0 MMK</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Cumulative ledger throughput</div>
        </div>
        <div class="card">
          <div class="card-title">Active Idempotency Keys</div>
          <div class="balance-val" id="metric-keys">0</div>
          <div style="font-size: 0.8rem; color: var(--text-muted);">Deduplication cache records</div>
        </div>
      </div>

      <div class="card">
        <div class="card-title">Raw Prometheus SRE Telemetry Stream (/metrics)</div>
        <div class="terminal-box" id="metricsOutput">// Loading telemetry data...</div>
      </div>
    </div>

    <!-- TAB 3: SUPER-APP CLOUD ARCHITECTURE -->
    <div id="tab-arch" class="tab-pane">
      <div class="card">
        <div class="card-title">Super-App Financial Cloud Mapping: KBZPay / WeChat Pay vs A Bank Showcase</div>
        <p style="font-size: 0.9rem; color: var(--text-muted); margin-bottom: 1rem;">
          How this showcase replicates the cloud-native patterns of Myanmar's KBZPay (Huawei Cloud) and China's WeChat Pay (Tencent/Alipay SOFAStack) using AWS Serverless & Cloudflare:
        </p>

        <table>
          <thead>
            <tr>
              <th>Architecture Layer</th>
              <th>KBZPay (Huawei Cloud)</th>
              <th>WeChat Pay / Alipay</th>
              <th>A Bank SRE Portfolio (AWS Live)</th>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td><strong>Edge Security & DDoS</strong></td>
              <td><span class="badge-huawei">Huawei WAF</span> Anti-DDoS</td>
              <td><span class="badge-wechat">Tencent Dayu</span> Shield WAF</td>
              <td><span class="badge-aws">Cloudflare Anycast</span> Edge Proxy + WAF</td>
            </tr>
            <tr>
              <td><strong>Microservice Runtime</strong></td>
              <td><span class="badge-huawei">CCE (K8s)</span> + <span class="badge-huawei">CCI Serverless</span></td>
              <td><span class="badge-wechat">Tencent TKE</span> Matrix Containers</td>
              <td><span class="badge-aws">AWS ECS Fargate</span> (Serverless Container Tier)</td>
            </tr>
            <tr>
              <td><strong>Event Bursts / Webhooks</strong></td>
              <td><span class="badge-huawei">FunctionGraph</span></td>
              <td><span class="badge-wechat">Serverless SCF</span></td>
              <td><span class="badge-aws">AWS Lambda</span> + EventBridge</td>
            </tr>
            <tr>
              <td><strong>Financial Ledger DB</strong></td>
              <td><span class="badge-huawei">GaussDB</span> (Distributed Multi-AZ)</td>
              <td><span class="badge-wechat">TDSQL / OceanBase</span></td>
              <td><span class="badge-aws">AWS DynamoDB</span> (On-Demand PITR + ACID)</td>
            </tr>
            <tr>
              <td><strong>Infrastructure Code</strong></td>
              <td>Huawei Cloud Terraform Provider</td>
              <td>Terraform / Internal IaC</td>
              <td><span class="badge-aws">Terraform</span> (AWS + Cloudflare Providers)</td>
            </tr>
            <tr>
              <td><strong>On-Prem / Edge Fleet</strong></td>
              <td>Bare-Metal Bank Branch Nodes</td>
              <td>Bare-Metal Hybrid Private Cloud</td>
              <td><span class="badge-aws">Ansible</span> (Kernel Hardening + K3s Fleet)</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- TAB 4: GOLDEN TEST HANDOVER RUNBOOK -->
    <div id="tab-runbook" class="tab-pane">
      <div class="card">
        <div class="card-title">SRE Golden Test CLI Verification Runbook</div>
        <p style="font-size: 0.9rem; color: var(--text-muted); margin-bottom: 1rem;">
          Interviewers and engineers can test this exact running cluster from any terminal using standard CLI tools:
        </p>

        <div style="margin-bottom: 1.5rem;">
          <h3 style="font-size: 1rem; color: var(--accent); margin-bottom: 0.5rem;">1. Zero-Trust Health & Readiness Probes</h3>
          <div class="terminal-box">curl -i "https://wallet.thaw-zin-2k77.de5.net/healthz"
curl -i "https://wallet.thaw-zin-2k77.de5.net/readyz"</div>
        </div>

        <div style="margin-bottom: 1.5rem;">
          <h3 style="font-size: 1rem; color: var(--accent); margin-bottom: 0.5rem;">2. Live Account Balance Query</h3>
          <div class="terminal-box">curl -s "https://wallet.thaw-zin-2k77.de5.net/api/v1/wallets/ACC-1001/balance" | jq .</div>
        </div>

        <div style="margin-bottom: 1.5rem;">
          <h3 style="font-size: 1rem; color: var(--accent); margin-bottom: 0.5rem;">3. Atomic Transfer with Strict Idempotency Key</h3>
          <div class="terminal-box">curl -s -X POST "https://wallet.thaw-zin-2k77.de5.net/api/v1/wallets/transfer" \
     -H "Content-Type: application/json" \
     -H "X-Idempotency-Key: GOLDEN-TEST-001" \
     -d '{"from_account":"ACC-1001","to_account":"ACC-2002","amount":10000,"currency":"MMK","reference":"Interview Golden Test"}' | jq .</div>
        </div>

        <div>
          <h3 style="font-size: 1rem; color: var(--accent); margin-bottom: 0.5rem;">4. GitHub Source Repository & CI/CD Pipeline</h3>
          <p style="font-size: 0.85rem; color: var(--text-muted);">
            Source Code: <a href="https://github.com/pretamane/fintech-wallet-sre-showcase" target="_blank" style="color: var(--accent);">github.com/pretamane/fintech-wallet-sre-showcase</a>
          </p>
        </div>
      </div>
    </div>
  </div>

  <script>
    function switchTab(tabId) {
      document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
      document.querySelectorAll('.tab-pane').forEach(pane => pane.classList.remove('active'));
      event.target.classList.add('active');
      document.getElementById(tabId).classList.add('active');
      if (tabId === 'tab-metrics') refreshMetrics();
    }

    function setAmount(val) {
      document.getElementById('txAmount').value = val;
    }

    function generateUUID() {
      const uuid = 'TX-' + Math.random().toString(36).substring(2, 9).toUpperCase() + '-' + Date.now().toString().slice(-4);
      document.getElementById('idempKey').value = uuid;
    }

    async function measureLatency() {
      const start = performance.now();
      try {
        await fetch('/healthz');
        const duration = Math.round(performance.now() - start);
        document.getElementById('latencyTag').innerText = 'Edge RTT: ' + duration + ' ms';
      } catch (e) {
        document.getElementById('latencyTag').innerText = 'Edge RTT: err';
      }
    }

    async function fetchAccountBalance(accId) {
      try {
        const res = await fetch('/api/v1/wallets/' + accId + '/balance');
        if (res.ok) {
          const data = await res.json();
          document.getElementById('bal-' + accId).innerText = Number(data.balance).toLocaleString() + ' MMK';
        }
      } catch (e) {}
    }

    async function refreshBalances() {
      await fetchAccountBalance('ACC-1001');
      await fetchAccountBalance('ACC-2002');
      await fetchAccountBalance('ACC-9999');
      measureLatency();
    }

    async function refreshMetrics() {
      try {
        const res = await fetch('/metrics');
        if (res.ok) {
          const data = await res.json();
          document.getElementById('metric-txCount').innerText = data.total_transactions || 0;
          document.getElementById('metric-volume').innerText = (data.total_volume || 0).toLocaleString() + ' MMK';
          document.getElementById('metric-keys').innerText = data.idempotency_keys || 0;
          document.getElementById('metricsOutput').innerText = JSON.stringify(data, null, 2);
        }
      } catch (e) {}
    }

    async function submitTransfer(isReplay) {
      const fromAcc = document.getElementById('fromAcc').value;
      const toAcc = document.getElementById('toAcc').value;
      const amount = parseFloat(document.getElementById('txAmount').value);
      let key = document.getElementById('idempKey').value.trim();

      if (!key) {
        generateUUID();
        key = document.getElementById('idempKey').value;
      }

      const terminal = document.getElementById('terminalOutput');
      const badge = document.getElementById('txStatusBadge');
      terminal.innerText = '>>> POST /api/v1/wallets/transfer\n>>> X-Idempotency-Key: ' + key + '\n>>> Submitting transaction payload...';

      try {
        const res = await fetch('/api/v1/wallets/transfer', {
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

        const data = await res.json();
        terminal.innerText = '<<< HTTP ' + res.status + ' ' + res.statusText + '\n' + JSON.stringify(data, null, 2);

        if (res.ok) {
          if (data.idempotent_replay) {
            badge.innerText = '🟡 IDEMPOTENT REPLAY (200 OK)';
            badge.style.color = 'var(--warning)';
          } else {
            badge.innerText = '🟢 NEW TRANSACTION (200 OK)';
            badge.style.color = 'var(--primary)';
            if (!isReplay) generateUUID(); // Prepare next key
          }
        } else {
          badge.innerText = '🔴 ERROR ' + res.status;
          badge.style.color = 'var(--danger)';
        }

        refreshBalances();
      } catch (err) {
        terminal.innerText = '<<< Network Error: ' + err.message;
        badge.innerText = '🔴 FAILED';
        badge.style.color = 'var(--danger)';
      }
    }

    // Initialize on load
    window.onload = function() {
      generateUUID();
      refreshBalances();
      setInterval(measureLatency, 8000);
    };
  </script>
</body>
</html>`
