import { createAccessToken, verifyJWT } from './jwt';
import {
  displayNameFromEmail,
  isValidEmail,
  normalizeEmail,
  randomId,
  randomToken,
  sha256Hex,
} from './crypto';
import { sendMagicLinkEmail } from './resend';

export interface Env {
  AUTH_DB: D1Database;
  AUTH_JWT_SECRET: string;
  RESEND_API_KEY: string;
  AUTH_FROM_EMAIL: string;
  AUTH_PUBLIC_URL: string;
  AUTH_APP_NAME: string;
}

const MAGIC_LINK_TTL_MS = 10 * 60 * 1000;

type UserRow = {
  id: string;
  email: string;
  display_name: string;
};

type MagicLinkRow = {
  id: string;
  email: string;
  token_hash: string;
  status: string;
  user_id: string | null;
  access_token: string | null;
  expires_at: string;
  verified_at: string | null;
};

export default {
  async fetch(request: Request, env: Env): Promise<Response> {
    const url = new URL(request.url);
    const { pathname } = url;

    try {
      if (request.method === 'POST' && pathname === '/v1/auth/magic-link') {
        return handleMagicLinkRequest(request, env);
      }
      if (request.method === 'GET' && pathname.startsWith('/v1/auth/magic-link/status/')) {
        const requestId = pathname.slice('/v1/auth/magic-link/status/'.length);
        return handleMagicLinkStatus(requestId, env);
      }
      if (request.method === 'GET' && pathname === '/v1/auth/verify') {
        return handleVerify(url, env);
      }
      if (request.method === 'GET' && pathname === '/v1/auth/me') {
        return handleMe(request, env);
      }
      if (request.method === 'POST' && pathname === '/v1/auth/signout') {
        return jsonResponse({ ok: true });
      }
      return jsonResponse({ error: 'not found' }, 404);
    } catch (err) {
      console.error(err);
      return jsonResponse({ error: err instanceof Error ? err.message : 'internal error' }, 500);
    }
  },
};

async function handleMagicLinkRequest(request: Request, env: Env): Promise<Response> {
  const body = (await request.json()) as { email?: string };
  const email = normalizeEmail(body.email ?? '');
  if (!isValidEmail(email)) {
    return jsonResponse({ error: 'invalid email' }, 400);
  }

  const requestId = randomId();
  const rawToken = randomToken();
  const tokenHash = await sha256Hex(rawToken);
  const expiresAt = new Date(Date.now() + MAGIC_LINK_TTL_MS).toISOString();

  await env.AUTH_DB.prepare(
    `INSERT INTO magic_link_requests (id, email, token_hash, status, expires_at)
     VALUES (?, ?, ?, 'pending', ?)`,
  )
    .bind(requestId, email, tokenHash, expiresAt)
    .run();

  const verifyURL = `${env.AUTH_PUBLIC_URL.replace(/\/$/, '')}/v1/auth/verify?token=${encodeURIComponent(rawToken)}&request=${encodeURIComponent(requestId)}`;
  await sendMagicLinkEmail(
    env.RESEND_API_KEY,
    env.AUTH_FROM_EMAIL,
    email,
    env.AUTH_APP_NAME,
    verifyURL,
  );

  return jsonResponse({ requestId, expiresAt });
}

async function handleMagicLinkStatus(requestId: string, env: Env): Promise<Response> {
  if (!requestId) {
    return jsonResponse({ error: 'missing request id' }, 400);
  }

  const row = await env.AUTH_DB.prepare(
    `SELECT id, email, token_hash, status, user_id, access_token, expires_at, verified_at
     FROM magic_link_requests WHERE id = ?`,
  )
    .bind(requestId)
    .first<MagicLinkRow>();

  if (!row) {
    return jsonResponse({ error: 'request not found' }, 404);
  }

  if (row.status === 'pending' && new Date(row.expires_at).getTime() < Date.now()) {
    await env.AUTH_DB.prepare(`UPDATE magic_link_requests SET status = 'expired' WHERE id = ?`)
      .bind(requestId)
      .run();
    return jsonResponse({ status: 'expired' });
  }

  if (row.status !== 'verified' || !row.user_id || !row.access_token) {
    return jsonResponse({ status: row.status });
  }

  const user = await env.AUTH_DB.prepare(
    `SELECT id, email, display_name FROM users WHERE id = ?`,
  )
    .bind(row.user_id)
    .first<UserRow>();

  if (!user) {
    return jsonResponse({ error: 'user not found' }, 404);
  }

  return jsonResponse({
    status: 'verified',
    accessToken: row.access_token,
    user: {
      id: user.id,
      email: user.email,
      displayName: user.display_name,
    },
  });
}

async function handleVerify(url: URL, env: Env): Promise<Response> {
  const rawToken = url.searchParams.get('token') ?? '';
  const requestId = url.searchParams.get('request') ?? '';
  if (!rawToken || !requestId) {
    return htmlResponse(verifyPage(env.AUTH_APP_NAME, false, 'Invalid sign-in link.'));
  }

  const tokenHash = await sha256Hex(rawToken);
  const row = await env.AUTH_DB.prepare(
    `SELECT id, email, token_hash, status, user_id, access_token, expires_at, verified_at
     FROM magic_link_requests WHERE id = ?`,
  )
    .bind(requestId)
    .first<MagicLinkRow>();

  if (!row || row.token_hash !== tokenHash) {
    return htmlResponse(verifyPage(env.AUTH_APP_NAME, false, 'This sign-in link is invalid.'));
  }

  if (row.status === 'verified') {
    return htmlResponse(
      verifyPage(env.AUTH_APP_NAME, true, 'You are already signed in. Return to Talus Agent.'),
    );
  }

  if (new Date(row.expires_at).getTime() < Date.now()) {
    await env.AUTH_DB.prepare(`UPDATE magic_link_requests SET status = 'expired' WHERE id = ?`)
      .bind(requestId)
      .run();
    return htmlResponse(
      verifyPage(env.AUTH_APP_NAME, false, 'This sign-in link has expired. Request a new one.'),
    );
  }

  const email = normalizeEmail(row.email);
  let user = await env.AUTH_DB.prepare(`SELECT id, email, display_name FROM users WHERE email = ?`)
    .bind(email)
    .first<UserRow>();

  if (!user) {
    const userId = randomId();
    const displayName = displayNameFromEmail(email);
    await env.AUTH_DB.prepare(
      `INSERT INTO users (id, email, display_name) VALUES (?, ?, ?)`,
    )
      .bind(userId, email, displayName)
      .run();
    user = { id: userId, email, display_name: displayName };
  }

  const accessToken = await createAccessToken(env.AUTH_JWT_SECRET, user.id, user.email);
  await env.AUTH_DB.prepare(
    `UPDATE magic_link_requests
     SET status = 'verified', user_id = ?, access_token = ?, verified_at = datetime('now')
     WHERE id = ?`,
  )
    .bind(user.id, accessToken, requestId)
    .run();

  return htmlResponse(
    verifyPage(env.AUTH_APP_NAME, true, 'Sign-in successful. You can close this tab and return to Talus Agent.'),
  );
}

async function handleMe(request: Request, env: Env): Promise<Response> {
  const authHeader = request.headers.get('Authorization') ?? '';
  const token = authHeader.startsWith('Bearer ') ? authHeader.slice(7) : '';
  if (!token) {
    return jsonResponse({ error: 'missing token' }, 401);
  }

  const payload = await verifyJWT(env.AUTH_JWT_SECRET, token);
  if (!payload || typeof payload.sub !== 'string') {
    return jsonResponse({ error: 'invalid token' }, 401);
  }

  const user = await env.AUTH_DB.prepare(
    `SELECT id, email, display_name FROM users WHERE id = ?`,
  )
    .bind(payload.sub)
    .first<UserRow>();

  if (!user) {
    return jsonResponse({ error: 'user not found' }, 404);
  }

  return jsonResponse({
    id: user.id,
    email: user.email,
    displayName: user.display_name,
  });
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { 'Content-Type': 'application/json' },
  });
}

function htmlResponse(body: string): Response {
  return new Response(body, {
    status: 200,
    headers: { 'Content-Type': 'text/html; charset=utf-8' },
  });
}

function verifyPage(appName: string, success: boolean, message: string): string {
  const title = success ? 'Signed in' : 'Sign-in failed';
  return `<!DOCTYPE html>
<html lang="en">
<head><meta charset="utf-8"><title>${title}</title></head>
<body style="font-family: system-ui, sans-serif; padding: 2rem; text-align: center;">
  <h1>${title}</h1>
  <p>${message}</p>
  <p style="color: #666;">${appName}</p>
</body>
</html>`;
}
