import { useCallback, useEffect, useState } from 'react';
import { Events } from '@wailsio/runtime';
import { ConfigService } from '../../bindings/github.com/songwei.ma/talus-mofish/backend/services';
import type { ThemeOption } from '../types/theme';

function isThemeOption(value: unknown): value is ThemeOption {
  return value === 'auto' || value === 'light' || value === 'dark';
}

function readEventTheme(event: unknown): ThemeOption | null {
  const payload =
    event && typeof event === 'object' && 'data' in event
      ? (event as { data: unknown }).data
      : event;

  if (isThemeOption(payload)) {
    return payload;
  }
  if (payload && typeof payload === 'object' && 'theme' in payload) {
    const theme = (payload as { theme: unknown }).theme;
    if (isThemeOption(theme)) {
      return theme;
    }
  }
  return null;
}

export function useConfigColorScheme() {
  const [colorScheme, setColorScheme] = useState<ThemeOption>('auto');

  const applyTheme = useCallback((theme: ThemeOption) => {
    setColorScheme(theme);
  }, []);

  useEffect(() => {
    ConfigService.GetConfig()
      .then((cfg) => {
        applyTheme((cfg.theme as ThemeOption) || 'auto');
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, [applyTheme]);

  useEffect(() => {
    const unsubscribe = Events.On('config:changed', (event) => {
      const theme = readEventTheme(event);
      if (theme) {
        applyTheme(theme);
      }
    });
    return () => unsubscribe();
  }, [applyTheme]);

  return { colorScheme, applyTheme };
}
