import { useCallback, useEffect, useState } from 'react';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import { AppShell, MantineProvider, Text, Title } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { AppService } from '../bindings/github.com/songwei.ma/talus-mofish';
import { NavbarSegmented } from './components/management/NavbarSegmented';
import {
  DEFAULT_MANAGEMENT_ROUTE,
  ManagementRoute,
  PAGE_TITLES,
  type ManagementRouteId,
} from './navigation/routes';
import { ConfigPage } from './pages/ConfigPage';
import { DebugPage } from './pages/DebugPage';
import { ImportPage } from './pages/ImportPage';
import { ReadingPage } from './pages/ReadingPage';
import { VocabularyPage } from './pages/VocabularyPage';

type ThemeOption = 'auto' | 'light' | 'dark';

function MainContent({
  activeItem,
  onThemeChange,
  onDebugModeChange,
}: {
  activeItem: ManagementRouteId;
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange: (enabled: boolean) => void;
}) {
  const title = PAGE_TITLES[activeItem] ?? 'Talus Echo';

  if (activeItem === ManagementRoute.Config) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ConfigPage onThemeChange={onThemeChange} onDebugModeChange={onDebugModeChange} />
      </>
    );
  }

  if (activeItem === ManagementRoute.Debug) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <DebugPage />
      </>
    );
  }

  if (activeItem === ManagementRoute.EnglishImport) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ImportPage />
      </>
    );
  }

  if (activeItem === ManagementRoute.EnglishReading) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ReadingPage />
      </>
    );
  }

  if (activeItem === ManagementRoute.EnglishVocabulary) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <VocabularyPage />
      </>
    );
  }

  if (activeItem === ManagementRoute.About) {
    return (
      <>
        <Title order={2}>{title}</Title>
        <Text c="dimmed" mt="sm">
          Talus Echo — a chat-oriented desktop agent for multiple domains. English Learning is the
          first domain: manage vocabulary, reading, and Anki imports here; use Agent Chat for
          interactive sessions.
        </Text>
      </>
    );
  }

  return null;
}

function ManagementApp() {
  const [activeItem, setActiveItem] = useState<ManagementRouteId>(DEFAULT_MANAGEMENT_ROUTE);
  const [colorScheme, setColorScheme] = useState<ThemeOption>('auto');
  const [debugMode, setDebugMode] = useState(false);

  const applyTheme = useCallback((theme: ThemeOption) => {
    setColorScheme(theme);
  }, []);

  const applyDebugMode = useCallback((enabled: boolean) => {
    setDebugMode(enabled);
    setActiveItem((current) =>
      current === ManagementRoute.Debug && !enabled ? DEFAULT_MANAGEMENT_ROUTE : current,
    );
  }, []);

  useEffect(() => {
    AppService.GetConfig()
      .then((cfg) => {
        const theme = (cfg.theme as ThemeOption) || 'auto';
        applyTheme(theme);
        applyDebugMode(cfg.debugMode ?? false);
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, [applyTheme, applyDebugMode]);

  return (
    <MantineProvider
      defaultColorScheme="auto"
      forceColorScheme={colorScheme === 'auto' ? undefined : colorScheme}
    >
      <Notifications position="top-right" limit={5} />
      <AppShell navbar={{ width: 300, breakpoint: 'sm' }} padding="md">
        <AppShell.Navbar p={0}>
          <NavbarSegmented
            activeItem={activeItem}
            debugMode={debugMode}
            onActiveItemChange={setActiveItem}
          />
        </AppShell.Navbar>

        <AppShell.Main>
          <MainContent
            activeItem={activeItem}
            onThemeChange={applyTheme}
            onDebugModeChange={applyDebugMode}
          />
        </AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
}

export default ManagementApp;
