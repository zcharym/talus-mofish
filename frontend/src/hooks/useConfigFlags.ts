import { useCallback, useEffect, useState } from 'react';
import { ConfigService, onConfigChanged } from '../utils/api';
import type { ThemeOption } from '../types/theme';

function isThemeOption(value: unknown): value is ThemeOption {
  return value === 'auto' || value === 'light' || value === 'dark';
}

export function useConfigFlags() {
  const [theme, setTheme] = useState<ThemeOption>('auto');
  const [debugMode, setDebugMode] = useState(false);
  const [cloudflareConfigured, setCloudflareConfigured] = useState(false);

  const applyTheme = useCallback((next: ThemeOption) => {
    setTheme(next);
  }, []);

  useEffect(() => {
    ConfigService.GetConfig()
      .then((cfg) => {
        applyTheme(isThemeOption(cfg.theme) ? cfg.theme : 'auto');
        setDebugMode(Boolean(cfg.debugMode));
        setCloudflareConfigured(Boolean(cfg.cloudflare?.accountId && cfg.cloudflare?.apiToken));
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, [applyTheme]);

  useEffect(() => {
    return onConfigChanged((payload) => {
      if (isThemeOption(payload.theme)) {
        applyTheme(payload.theme);
      }
      setDebugMode(Boolean(payload.debugMode));
      setCloudflareConfigured(Boolean(payload.cloudflareConfigured));
    });
  }, [applyTheme]);

  return { theme, debugMode, cloudflareConfigured, applyTheme };
}
