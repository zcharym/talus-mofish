import { useCallback, useEffect, useState } from 'react';
import { CloudflareService, toApiError, type CloudflareDashboardSnapshot } from '../utils/api';
import { notify } from '../services/notifications';

export function useCloudflareDashboard(configured: boolean) {
  const [loading, setLoading] = useState(configured);
  const [refreshing, setRefreshing] = useState(false);
  const [snapshot, setSnapshot] = useState<CloudflareDashboardSnapshot | null>(null);
  const [loadError, setLoadError] = useState('');

  const loadDashboard = useCallback(
    async (isRefresh = false) => {
      if (!configured) {
        setSnapshot(null);
        setLoadError('');
        setLoading(false);
        return;
      }
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      try {
        const next = await CloudflareService.GetDashboard();
        setSnapshot(next);
        setLoadError('');
      } catch (err) {
        const message = toApiError(err);
        setLoadError(message);
        notify.failed('Cloudflare', message);
      } finally {
        setLoading(false);
        setRefreshing(false);
      }
    },
    [configured],
  );

  useEffect(() => {
    void loadDashboard();
  }, [loadDashboard]);

  return { loading, refreshing, snapshot, loadError, loadDashboard };
}
