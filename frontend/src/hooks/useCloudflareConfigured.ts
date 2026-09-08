import { useCallback, useEffect, useState } from 'react';
import { Events } from '@wailsio/runtime';
import { CloudflareService } from '../../bindings/github.com/songwei.ma/talus-mofish/backend/services';

function readConfigured(event: unknown): boolean | null {
  const payload =
    event && typeof event === 'object' && 'data' in event
      ? (event as { data: unknown }).data
      : event;
  if (payload && typeof payload === 'object' && 'cloudflareConfigured' in payload) {
    return Boolean((payload as { cloudflareConfigured: unknown }).cloudflareConfigured);
  }
  return null;
}

export function useCloudflareConfigured() {
  const [configured, setConfigured] = useState(false);

  const refresh = useCallback(async () => {
    try {
      setConfigured(Boolean(await CloudflareService.IsConfigured()));
    } catch (err) {
      console.error(err);
      setConfigured(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  useEffect(() => {
    const unsubscribe = Events.On('config:changed', (event) => {
      const next = readConfigured(event);
      if (next !== null) {
        setConfigured(next);
        return;
      }
      void refresh();
    });
    return () => unsubscribe();
  }, [refresh]);

  return configured;
}
