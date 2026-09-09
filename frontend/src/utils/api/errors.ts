import { Call } from '@wailsio/runtime';

type RuntimeErrorConstructor = new (...args: unknown[]) => Error;

function runtimeErrorType(): RuntimeErrorConstructor | undefined {
  const ctor = (Call as { RuntimeError?: RuntimeErrorConstructor }).RuntimeError;
  return ctor;
}

export function isRuntimeError(err: unknown): boolean {
  const RuntimeError = runtimeErrorType();
  return Boolean(RuntimeError && err instanceof RuntimeError);
}

function causeMessage(cause: unknown): string {
  if (!cause || typeof cause !== 'object') {
    return '';
  }
  if ('message' in cause && typeof cause.message === 'string' && cause.message) {
    return cause.message;
  }
  if ('reason' in cause && typeof cause.reason === 'string' && cause.reason) {
    return cause.reason;
  }
  return '';
}

/** Unwrap a Wails v3 binding rejection. Prefer Error.message; use cause when present. */
export function toApiError(err: unknown): string {
  if (typeof err === 'string' && err.trim()) {
    return err;
  }
  if (err instanceof Error) {
    const fromCause = causeMessage((err as Error & { cause?: unknown }).cause);
    if (fromCause && fromCause !== err.message) {
      return err.message ? `${err.message}: ${fromCause}` : fromCause;
    }
    if (err.message) {
      return err.message;
    }
  }
  if (err && typeof err === 'object' && 'message' in err && typeof err.message === 'string' && err.message) {
    return err.message;
  }
  return String(err);
}
