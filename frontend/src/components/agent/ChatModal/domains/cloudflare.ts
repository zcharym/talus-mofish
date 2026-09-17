import type { CloudflareDashboardSnapshot } from '../../../../utils/api';
import type { ChatModalDomain } from '../types';

const MAX_WORKERS = 12;
const MAX_NAMED = 8;

export function createCloudflareChatDomain({
  configured,
  snapshot,
  loadError,
}: {
  configured: boolean;
  snapshot: CloudflareDashboardSnapshot | null;
  loadError?: string;
}): ChatModalDomain {
  return {
    id: 'cloudflare',
    title: 'Cloudflare assistant',
    placeholder: 'Ask about this dashboard…',
    emptyHint: configured
      ? 'Ask about Workers, 24h traffic, D1, or KV. This chat uses the latest dashboard snapshot and is not saved in the sidebar.'
      : 'Cloudflare is not connected yet. Ask how to add an account ID and API token, or open Management → Configuration → Cloudflare.',
    suggestions: configured
      ? [
          {
            id: 'health',
            label: 'Summarize account health',
            prompt: 'Summarize this Cloudflare account’s health from the current dashboard snapshot.',
          },
          {
            id: 'errors',
            label: 'Which worker has the most errors?',
            prompt: 'Which Worker has the most errors in the last 24 hours, and what stands out?',
          },
          {
            id: 'known',
            label: 'What’s missing among known services?',
            prompt: 'Which known Talus services are missing on this account, and what does that mean?',
          },
        ]
      : [
          {
            id: 'setup',
            label: 'How do I connect Cloudflare?',
            prompt: 'How do I connect a Cloudflare account in Talus Echo?',
          },
        ],
    buildContext: () => formatCloudflareContext({ configured, snapshot, loadError }),
  };
}

function formatCloudflareContext({
  configured,
  snapshot,
  loadError,
}: {
  configured: boolean;
  snapshot: CloudflareDashboardSnapshot | null;
  loadError?: string;
}): string {
  const lines = [
    'You are helping the user with a read-only Cloudflare operations dashboard in Talus Echo.',
    'Do not suggest deploying, tailing logs, rolling back, or mutating Cloudflare resources.',
    'Answer from the snapshot below. If a figure is missing, say so instead of inventing it.',
  ];

  if (!configured) {
    lines.push(
      'Status: credentials are not configured.',
      'Tell the user to add cloudflare.accountId and cloudflare.apiToken in Management → Configuration → Cloudflare, then save and reopen the Agent Cloudflare tab.',
    );
    return lines.join('\n');
  }

  if (loadError && !snapshot) {
    lines.push(`Status: dashboard failed to load. Error: ${loadError}`);
    return lines.join('\n');
  }

  if (!snapshot) {
    lines.push('Status: dashboard snapshot is not available yet.');
    return lines.join('\n');
  }

  const accountName = snapshot.account?.name || snapshot.account?.id || 'unknown';
  const kpis = snapshot.kpis;
  lines.push(
    `Account: ${accountName}${snapshot.account?.id ? ` (${snapshot.account.id})` : ''}`,
    `Fetched at: ${snapshot.fetchedAt || 'unknown'}`,
    `KPIs (24h): workers=${kpis?.workerCount ?? 0} requests=${Math.round(kpis?.requests24h || 0)} errors=${Math.round(kpis?.errors24h || 0)} errorRate=${((kpis?.errorRate || 0) * 100).toFixed(2)}%`,
  );

  const errors = snapshot.errors;
  const scopedErrors = [
    errors?.account ? `account: ${errors.account}` : '',
    errors?.workers ? `workers: ${errors.workers}` : '',
    errors?.analytics ? `analytics: ${errors.analytics}` : '',
    errors?.d1 ? `d1: ${errors.d1}` : '',
    errors?.kv ? `kv: ${errors.kv}` : '',
  ].filter(Boolean);
  if (scopedErrors.length > 0) {
    lines.push(`Dashboard errors: ${scopedErrors.join('; ')}`);
  }

  const known = snapshot.known ?? [];
  if (known.length > 0) {
    lines.push(
      'Known services: ' +
        known
          .map((service) => `${service.title} (${service.id})=${service.present ? 'found' : 'missing'}`)
          .join('; '),
    );
  }

  const workers = [...(snapshot.workers ?? [])]
    .sort((a, b) => (b.errors24h || 0) - (a.errors24h || 0))
    .slice(0, MAX_WORKERS);
  if (workers.length === 0) {
    lines.push('Workers: none listed.');
  } else {
    lines.push('Workers (top by 24h errors):');
    for (const worker of workers) {
      const cron = (worker.cron ?? []).filter(Boolean).join(', ') || 'none';
      lines.push(
        `- ${worker.name}: req=${Math.round(worker.requests24h || 0)} err=${Math.round(worker.errors24h || 0)} cpuP99us=${Math.round(worker.cpuTimeP99 || 0)} cron=${cron}`,
      );
    }
  }

  const d1 = (snapshot.d1 ?? []).slice(0, MAX_NAMED);
  lines.push(
    d1.length > 0
      ? `D1: ${d1.map((item) => item.name || item.id).join(', ')}`
      : 'D1: none visible.',
  );

  const kv = (snapshot.kv ?? []).slice(0, MAX_NAMED);
  lines.push(
    kv.length > 0
      ? `KV: ${kv.map((item) => item.title || item.id).join(', ')}`
      : 'KV: none visible.',
  );

  return lines.join('\n');
}
