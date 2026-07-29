import {
  buildPushPayload,
  type PushMessage,
  type PushSubscription,
  type VapidKeys,
} from "@block65/webcrypto-web-push";

export interface Env {
  KV: KVNamespace;
  ASSETS: Fetcher;
  VAPID_PUBLIC_KEY: string;
  VAPID_PRIVATE_KEY: string;
  VAPID_SUBJECT: string;
}

interface StoredSubscription {
  endpoint: string;
  keys: { p256dh: string; auth: string };
  created_at: string;
  user_agent?: string;
}

interface SubscribeBody {
  token: string;
  subscription: PushSubscription;
}

interface AlertBody {
  rule_id: string;
  title: string;
  body: string;
  target?: string;
}

const CORS_HEADERS: Record<string, string> = {
  "Access-Control-Allow-Origin": "*",
  "Access-Control-Allow-Methods": "GET, POST, DELETE, OPTIONS",
  "Access-Control-Allow-Headers": "Content-Type, Authorization",
};

async function tokenHash(token: string): Promise<string> {
  const data = new TextEncoder().encode(token);
  const digest = await crypto.subtle.digest("SHA-256", data);
  return [...new Uint8Array(digest)]
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}

function json(data: unknown, status = 200): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json", ...CORS_HEADERS },
  });
}

function unauthorized(): Response {
  return json({ error: "unauthorized" }, 401);
}

function vapidKeys(env: Env): VapidKeys {
  return {
    subject: env.VAPID_SUBJECT,
    publicKey: env.VAPID_PUBLIC_KEY,
    privateKey: env.VAPID_PRIVATE_KEY,
  };
}

async function getSubscription(env: Env, token: string): Promise<StoredSubscription | null> {
  const key = `sub:${await tokenHash(token)}`;
  const raw = await env.KV.get(key);
  if (!raw) return null;
  return JSON.parse(raw) as StoredSubscription;
}

async function saveSubscription(
  env: Env,
  token: string,
  subscription: PushSubscription,
  userAgent: string | null,
): Promise<void> {
  const key = `sub:${await tokenHash(token)}`;
  const stored: StoredSubscription = {
    endpoint: subscription.endpoint,
    keys: {
      p256dh: subscription.keys.p256dh,
      auth: subscription.keys.auth,
    },
    created_at: new Date().toISOString(),
    user_agent: userAgent ?? undefined,
  };
  await env.KV.put(key, JSON.stringify(stored));
}

async function deleteSubscription(env: Env, token: string): Promise<void> {
  const key = `sub:${await tokenHash(token)}`;
  await env.KV.delete(key);
}

async function sendPush(
  env: Env,
  stored: StoredSubscription,
  message: PushMessage,
): Promise<{ ok: boolean; status: number; stale: boolean }> {
  const subscription: PushSubscription = {
    endpoint: stored.endpoint,
    expirationTime: null,
    keys: stored.keys,
  };

  const payload = await buildPushPayload(message, subscription, vapidKeys(env));
  const res = await fetch(subscription.endpoint, payload);
  return { ok: res.ok, status: res.status, stale: res.status === 410 };
}

async function handleSubscribe(request: Request, env: Env): Promise<Response> {
  const body = (await request.json()) as SubscribeBody;
  if (!body.token || !body.subscription?.endpoint || !body.subscription.keys) {
    return json({ error: "invalid subscription payload" }, 400);
  }
  await saveSubscription(env, body.token, body.subscription, request.headers.get("user-agent"));
  return json({ ok: true });
}

async function handleUnsubscribe(request: Request, env: Env): Promise<Response> {
  const body = (await request.json()) as { token?: string };
  if (!body.token) {
    return json({ error: "token required" }, 400);
  }
  await deleteSubscription(env, body.token);
  return json({ ok: true });
}

async function handleAlert(request: Request, env: Env): Promise<Response> {
  const auth = request.headers.get("authorization") ?? "";
  const token = auth.startsWith("Bearer ") ? auth.slice(7).trim() : "";
  if (!token) {
    return unauthorized();
  }

  const stored = await getSubscription(env, token);
  if (!stored) {
    return json({ error: "no subscription for token" }, 404);
  }

  const body = (await request.json()) as AlertBody;
  if (!body.title || !body.body) {
    return json({ error: "title and body required" }, 400);
  }

  const ruleID = body.rule_id || "alert";
  const stateKey = `state:${ruleID}`;
  const now = Date.now();
  const lastRaw = await env.KV.get(stateKey);
  if (lastRaw) {
    const last = Number.parseInt(lastRaw, 10);
    if (!Number.isNaN(last) && now-last < 60_000) {
      return json({ ok: true, skipped: "debounced" });
    }
  }

  const message: PushMessage = {
    data: JSON.stringify({
      title: body.title,
      body: body.body,
      rule_id: ruleID,
      target: body.target,
      ts: new Date().toISOString(),
    }),
  };

  const result = await sendPush(env, stored, message);
  if (result.stale) {
    await deleteSubscription(env, token);
    return json({ error: "subscription expired", stale: true }, 410);
  }
  if (!result.ok) {
    return json({ error: "push failed", status: result.status }, 502);
  }

  await env.KV.put(stateKey, String(now), { expirationTtl: 86400 });
  return json({ ok: true });
}

async function handleTestPush(request: Request, env: Env): Promise<Response> {
  const auth = request.headers.get("authorization") ?? "";
  const token = auth.startsWith("Bearer ") ? auth.slice(7).trim() : "";
  if (!token) {
    return unauthorized();
  }

  const stored = await getSubscription(env, token);
  if (!stored) {
    return json({ error: "no subscription for token" }, 404);
  }

  const message: PushMessage = {
    data: JSON.stringify({
      title: "Echo Watch",
      body: "Test notification — pairing works.",
      rule_id: "test",
      ts: new Date().toISOString(),
    }),
  };

  const result = await sendPush(env, stored, message);
  if (result.stale) {
    await deleteSubscription(env, token);
    return json({ error: "subscription expired", stale: true }, 410);
  }
  if (!result.ok) {
    return json({ error: "push failed", status: result.status }, 502);
  }
  return json({ ok: true });
}

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);

    if (request.method === "OPTIONS") {
      return new Response(null, { status: 204, headers: CORS_HEADERS });
    }

    if (url.pathname === "/health") {
      return json({ ok: true, service: "echo-watch" });
    }

    if (url.pathname === "/api/vapid-public-key" && request.method === "GET") {
      if (!env.VAPID_PUBLIC_KEY) {
        return json({ error: "VAPID not configured" }, 503);
      }
      return json({ publicKey: env.VAPID_PUBLIC_KEY });
    }

    if (url.pathname === "/api/subscribe" && request.method === "POST") {
      return handleSubscribe(request, env);
    }

    if (url.pathname === "/api/subscribe" && request.method === "DELETE") {
      return handleUnsubscribe(request, env);
    }

    if (url.pathname === "/api/alert" && request.method === "POST") {
      return handleAlert(request, env);
    }

    if (url.pathname === "/api/test-push" && request.method === "POST") {
      return handleTestPush(request, env);
    }

    return env.ASSETS.fetch(request);
  },
};
