import { useCallback, useEffect, useState } from 'react';
import type { FileEntry } from '../utils/api';
import { ObsidianService, toApiError } from '../utils/api';

export const VAULT_ROOT_KEY = '';

export interface VaultDirState {
  expanded: boolean;
  loading: boolean;
  entries: FileEntry[] | null;
  error: string | null;
}

export function useVaultTree(expandToPath?: string | null) {
  const [dirs, setDirs] = useState<Record<string, VaultDirState>>({});

  const loadDir = useCallback(async (dirPath: string) => {
    setDirs((prev) => ({
      ...prev,
      [dirPath]: {
        expanded: true,
        loading: true,
        entries: prev[dirPath]?.entries ?? null,
        error: null,
      },
    }));
    try {
      const entries = await ObsidianService.ListDirectory(dirPath);
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { expanded: true, loading: false, entries, error: null },
      }));
    } catch (err) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: {
          expanded: true,
          loading: false,
          entries: prev[dirPath]?.entries ?? null,
          error: toApiError(err),
        },
      }));
    }
  }, []);

  useEffect(() => {
    void loadDir(VAULT_ROOT_KEY);
  }, [loadDir]);

  useEffect(() => {
    if (!expandToPath) {
      return;
    }
    const parts = expandToPath.split('/').filter(Boolean);
    let cancelled = false;

    const expandParents = async () => {
      await loadDir(VAULT_ROOT_KEY);
      let acc = '';
      for (let i = 0; i < parts.length - 1; i += 1) {
        if (cancelled) {
          return;
        }
        acc = acc ? `${acc}/${parts[i]}` : parts[i];
        await loadDir(acc);
      }
    };

    void expandParents();
    return () => {
      cancelled = true;
    };
  }, [expandToPath, loadDir]);

  const toggleDir = (dirPath: string) => {
    const current = dirs[dirPath];
    if (current?.expanded) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { ...current, expanded: false },
      }));
      return;
    }
    if (current?.entries) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { ...current, expanded: true },
      }));
      return;
    }
    void loadDir(dirPath);
  };

  return { dirs, toggleDir };
}
