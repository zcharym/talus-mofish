export interface CloudflareAccount {
  id: string;
  name: string;
}

export interface CloudflareWorker {
  name: string;
  createdOn: string;
  modifiedOn: string;
  requests24h: number;
  errors24h: number;
  cpuTimeP99: number;
  cron: string[];
}

export interface CloudflareD1Database {
  id: string;
  name: string;
}

export interface CloudflareKVNamespace {
  id: string;
  title: string;
}

export interface CloudflareKnownService {
  id: string;
  title: string;
  present: boolean;
  bindingsNote: string;
  worker?: CloudflareWorker | null;
}

export interface CloudflareKPIs {
  workerCount: number;
  requests24h: number;
  errors24h: number;
  errorRate: number;
}

export interface CloudflareDashboardErrors {
  account: string;
  workers: string;
  analytics: string;
  d1: string;
  kv: string;
}

export interface CloudflareDashboardSnapshot {
  configured: boolean;
  fetchedAt: string;
  account: CloudflareAccount;
  kpis: CloudflareKPIs;
  workers: CloudflareWorker[];
  known: CloudflareKnownService[];
  d1: CloudflareD1Database[];
  kv: CloudflareKVNamespace[];
  errors: CloudflareDashboardErrors;
}
