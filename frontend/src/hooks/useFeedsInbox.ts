import { useCallback, useEffect, useState } from 'react';
import { FeedsService, toApiError } from '../utils/api';
import type { FeedInbox, FeedItem, FeedSource } from '../utils/api';
import { notify } from '../services/notifications';

export type FeedsFilter = 'all' | 'saved' | 'youtube' | 'bilibili' | 'rss';

export function useFeedsInbox(filter: FeedsFilter) {
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [sources, setSources] = useState<FeedSource[]>([]);
  const [items, setItems] = useState<FeedItem[]>([]);
  const [loadError, setLoadError] = useState('');

  const loadInbox = useCallback(
    async (isRefresh = false) => {
      if (isRefresh) {
        setRefreshing(true);
      } else {
        setLoading(true);
      }
      try {
        const inbox = (await FeedsService.GetInbox(filter)) as FeedInbox;
        setSources(inbox.sources ?? []);
        setItems(inbox.items ?? []);
        setLoadError('');
      } catch (err) {
        const message = toApiError(err);
        setLoadError(message);
        notify.failed('Feeds', message);
      } finally {
        setLoading(false);
        setRefreshing(false);
      }
    },
    [filter],
  );

  useEffect(() => {
    void loadInbox();
  }, [loadInbox]);

  const refresh = useCallback(async () => {
    setRefreshing(true);
    try {
      await FeedsService.Refresh();
      await loadInbox(true);
    } catch (err) {
      notify.failed('Feeds', toApiError(err));
      setRefreshing(false);
    }
  }, [loadInbox]);

  const addSource = useCallback(
    async (input: string) => {
      await FeedsService.AddSource(input);
      await loadInbox(true);
    },
    [loadInbox],
  );

  const saveLink = useCallback(
    async (title: string, url: string) => {
      await FeedsService.SaveLink(title, url);
      await loadInbox(true);
    },
    [loadInbox],
  );

  const deleteSource = useCallback(
    async (id: string) => {
      await FeedsService.DeleteSource(id);
      await loadInbox(true);
    },
    [loadInbox],
  );

  const setSaved = useCallback(
    async (id: string, saved: boolean) => {
      await FeedsService.SetItemSaved(id, saved);
      setItems((current) => current.map((item) => (item.id === id ? { ...item, saved } : item)));
    },
    [],
  );

  const setRead = useCallback(
    async (id: string, read: boolean) => {
      await FeedsService.SetItemRead(id, read);
      setItems((current) => current.map((item) => (item.id === id ? { ...item, read } : item)));
    },
    [],
  );

  return {
    loading,
    refreshing,
    sources,
    items,
    loadError,
    loadInbox,
    refresh,
    addSource,
    saveLink,
    deleteSource,
    setSaved,
    setRead,
  };
}
