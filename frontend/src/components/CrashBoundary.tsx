import { Component, type ErrorInfo, type ReactNode } from 'react';
import { SystemService, toApiError } from '../utils/api';

interface CrashBoundaryProps {
  children: ReactNode;
}

interface CrashBoundaryState {
  error: string | null;
  report: string;
}

function buildReport(error: string, platform: string, version: string): string {
  return [
    'Talus Echo crash report',
    `time: ${new Date().toISOString()}`,
    `platform: ${platform || 'unknown'}`,
    `userAgent: ${typeof navigator !== 'undefined' ? navigator.userAgent : 'n/a'}`,
    `version: ${version || 'dev'}`,
    '',
    error,
  ].join('\n');
}

export class CrashBoundary extends Component<CrashBoundaryProps, CrashBoundaryState> {
  state: CrashBoundaryState = { error: null, report: '' };

  static getDerivedStateFromError(error: Error): Partial<CrashBoundaryState> {
    return { error: error.message || String(error) };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    void this.capture(error.stack || error.message, info.componentStack || '');
  }

  componentDidMount() {
    window.addEventListener('error', this.onWindowError);
    window.addEventListener('unhandledrejection', this.onUnhandledRejection);
  }

  componentWillUnmount() {
    window.removeEventListener('error', this.onWindowError);
    window.removeEventListener('unhandledrejection', this.onUnhandledRejection);
  }

  onWindowError = (event: ErrorEvent) => {
    void this.capture(event.error?.stack || event.message, '');
  };

  onUnhandledRejection = (event: PromiseRejectionEvent) => {
    void this.capture(toApiError(event.reason), '');
  };

  capture = async (message: string, extra: string) => {
    let platform = '';
    let version = '';
    try {
      [platform, version] = await Promise.all([SystemService.Platform(), SystemService.Version()]);
    } catch {
      // Native bindings may be unavailable in a crashed renderer.
    }
    const error = extra ? `${message}\n${extra}` : message;
    this.setState({ error, report: buildReport(error, platform, version) });
  };

  render() {
    if (!this.state.error) {
      return this.props.children;
    }

    return (
      <div style={{ fontFamily: 'sans-serif', maxWidth: 720, margin: '48px auto', padding: 24 }}>
        <h1 style={{ fontSize: 22 }}>Something went wrong</h1>
        <p style={{ color: '#666' }}>
          Copy this report when asking for help. The desktop app has no public DevTools by default.
        </p>
        <pre
          style={{
            whiteSpace: 'pre-wrap',
            background: '#111',
            color: '#eee',
            padding: 16,
            borderRadius: 8,
          }}
        >
          {this.state.report || this.state.error}
        </pre>
        <button
          type="button"
          onClick={() => {
            void navigator.clipboard?.writeText(this.state.report || this.state.error || '');
          }}
        >
          Copy report
        </button>
      </div>
    );
  }
}
