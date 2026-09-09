import { Browser } from '@wailsio/runtime';
import { SystemService } from './api';

export async function openExternalURL(raw: string): Promise<void> {
  const url = raw.trim();
  if (!url.startsWith('http://') && !url.startsWith('https://')) {
    return;
  }
  try {
    await SystemService.OpenURL(url);
    return;
  } catch {
    // Fall through to the runtime helper.
  }
  try {
    await Browser.OpenURL(url);
  } catch {
    window.open(url, '_blank', 'noopener,noreferrer');
  }
}
