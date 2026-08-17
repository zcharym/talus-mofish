function base64UrlEncode(data: ArrayBuffer | Uint8Array): string {
  const bytes = data instanceof Uint8Array ? data : new Uint8Array(data);
  return btoa(String.fromCharCode(...bytes))
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=+$/, '');
}

function base64UrlDecode(value: string): Uint8Array {
  const padded = value.replace(/-/g, '+').replace(/_/g, '/');
  const padLen = (4 - (padded.length % 4)) % 4;
  const binary = atob(padded + '='.repeat(padLen));
  return Uint8Array.from(binary, (c) => c.charCodeAt(0));
}

async function importHmacKey(secret: string): Promise<CryptoKey> {
  return crypto.subtle.importKey(
    'raw',
    new TextEncoder().encode(secret),
    { name: 'HMAC', hash: 'SHA-256' },
    false,
    ['sign', 'verify'],
  );
}

export async function signJWT(
  secret: string,
  payload: Record<string, unknown>,
): Promise<string> {
  const header = { alg: 'HS256', typ: 'JWT' };
  const encodedHeader = base64UrlEncode(new TextEncoder().encode(JSON.stringify(header)));
  const encodedPayload = base64UrlEncode(new TextEncoder().encode(JSON.stringify(payload)));
  const signingInput = `${encodedHeader}.${encodedPayload}`;
  const key = await importHmacKey(secret);
  const signature = await crypto.subtle.sign('HMAC', key, new TextEncoder().encode(signingInput));
  return `${signingInput}.${base64UrlEncode(signature)}`;
}

export async function verifyJWT(
  secret: string,
  token: string,
): Promise<Record<string, unknown> | null> {
  const parts = token.split('.');
  if (parts.length !== 3) {
    return null;
  }

  const [encodedHeader, encodedPayload, encodedSignature] = parts;
  const signingInput = `${encodedHeader}.${encodedPayload}`;
  const key = await importHmacKey(secret);
  const signature = base64UrlDecode(encodedSignature);
  const valid = await crypto.subtle.verify(
    'HMAC',
    key,
    signature,
    new TextEncoder().encode(signingInput),
  );
  if (!valid) {
    return null;
  }

  const payload = JSON.parse(new TextDecoder().decode(base64UrlDecode(encodedPayload))) as Record<
    string,
    unknown
  >;
  const exp = payload.exp;
  if (typeof exp === 'number' && exp * 1000 < Date.now()) {
    return null;
  }
  return payload;
}

export async function createAccessToken(
  secret: string,
  userId: string,
  email: string,
): Promise<string> {
  const now = Math.floor(Date.now() / 1000);
  return signJWT(secret, {
    sub: userId,
    email,
    iat: now,
    exp: now + 30 * 24 * 60 * 60,
  });
}
