import { useCallback, useEffect, useState } from 'react';
import { AppShell, Burger, Group, Text, Title } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { ConfigService } from '../bindings/github.com/songwei.ma/talus-mofish/backend/services';
import { NavbarSegmented } from './components/management/NavbarSegmented';
import { useConfigColorScheme } from './hooks/useConfigColorScheme';
import {
  DEFAULT_MANAGEMENT_ROUTE,
  ManagementRoute,
  PAGE_TITLES,
  type ManagementRouteId,
} from './navigation/routes';
import { ConfigPage } from './pages/ConfigPage';
import { DebugPage } from './pages/DebugPage';
import { ImportPage } from './pages/ImportPage';
import { ObsidianNotesPage } from './pages/ObsidianNotesPage';
import { ObsidianSearchPage } from './pages/ObsidianSearchPage';
import { ReadingPage } from './pages/ReadingPage';
import { VocabularyPage } from './pages/VocabularyPage';
import { ThemeRoot } from './theme';
import type { ThemeOption } from './types/theme';
import { isMacOS } from './utils/platform';
import classes from './ManagementApp.module.css';

function MainContent({
  activeItem,
  onThemeChange,
  onDebugModeChange,
  obsidianFocusPath,
  onObsidianFocusConsumed,
  onOpenObsidianNote,
}: {
  activeItem: ManagementRouteId;
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange: (enabled: boolean) => void;
  obsidianFocusPath: string | null;
  onObsidianFocusConsumed: () => void;
  onOpenObsidianNote: (path: string) => void;
}) {
  if (activeItem === ManagementRoute.Config) {
    return <ConfigPage onThemeChange={onThemeChange} onDebugModeChange={onDebugModeChange} />;
  }

  if (activeItem === ManagementRoute.Debug) {
    return <DebugPage />;
  }

  if (activeItem === ManagementRoute.EnglishImport) {
    return <ImportPage />;
  }

  if (activeItem === ManagementRoute.EnglishReading) {
    return <ReadingPage />;
  }

  if (activeItem === ManagementRoute.EnglishVocabulary) {
    return <VocabularyPage />;
  }

  if (activeItem === ManagementRoute.ObsidianNotes) {
    return (
      <ObsidianNotesPage
        focusPath={obsidianFocusPath}
        onFocusConsumed={onObsidianFocusConsumed}
      />
    );
  }

  if (activeItem === ManagementRoute.ObsidianSearch) {
    return <ObsidianSearchPage onOpenNote={onOpenObsidianNote} />;
  }

  if (activeItem === ManagementRoute.About) {
    return (
      <Text c="dimmed" mt="sm">
        Talus Echo — a chat-oriented desktop agent for multiple domains. English Learning is the
        first domain: manage vocabulary, reading, and Anki imports here; use Agent Chat for
        interactive sessions.
      </Text>
    );
  }

  return null;
}

function ManagementApp() {
  const [activeItem, setActiveItem] = useState<ManagementRouteId>(DEFAULT_MANAGEMENT_ROUTE);
  const { colorScheme, applyTheme } = useConfigColorScheme();
  const [debugMode, setDebugMode] = useState(false);
  const [obsidianFocusPath, setObsidianFocusPath] = useState<string | null>(null);
  const [navOpened, setNavOpened] = useState(false);
  const isMobile = useMediaQuery('(max-width: 48em)', false, {
    getInitialValueInEffect: false,
  });

  const applyDebugMode = useCallback((enabled: boolean) => {
    setDebugMode(enabled);
    setActiveItem((current) =>
      current === ManagementRoute.Debug && !enabled ? DEFAULT_MANAGEMENT_ROUTE : current,
    );
  }, []);

  const openObsidianNote = useCallback((path: string) => {
    setObsidianFocusPath(path);
    setActiveItem(ManagementRoute.ObsidianNotes);
    setNavOpened(false);
  }, []);

  const consumeObsidianFocus = useCallback(() => {
    setObsidianFocusPath(null);
  }, []);

  const handleActiveItemChange = useCallback((itemId: ManagementRouteId) => {
    setActiveItem(itemId);
    setNavOpened(false);
  }, []);

  useEffect(() => {
    ConfigService.GetConfig()
      .then((cfg) => {
        applyDebugMode(cfg.debugMode ?? false);
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, [applyDebugMode]);

  const title = PAGE_TITLES[activeItem] ?? 'Talus Echo';

  return (
    <ThemeRoot colorScheme={colorScheme}>
      <AppShell
        header={{ height: 52 }}
        navbar={{
          width: 280,
          breakpoint: 'sm',
          collapsed: { mobile: !navOpened },
        }}
        padding={{ base: 'sm', sm: 'md' }}
        className={classes.shell}
      >
        <AppShell.Header className={classes.header} data-platform={isMacOS ? 'darwin' : undefined}>
          <Group h="100%" px="md" gap="sm" wrap="nowrap">
            {isMobile ? (
              <Burger
                opened={navOpened}
                onClick={() => setNavOpened((opened) => !opened)}
                size="sm"
                aria-label="Toggle navigation"
              />
            ) : null}
            <Title order={4} className={classes.headerTitle}>
              {title}
            </Title>
          </Group>
        </AppShell.Header>

        <AppShell.Navbar p={0}>
          <NavbarSegmented
            activeItem={activeItem}
            debugMode={debugMode}
            onActiveItemChange={handleActiveItemChange}
          />
        </AppShell.Navbar>

        <AppShell.Main>
          <MainContent
            activeItem={activeItem}
            onThemeChange={applyTheme}
            onDebugModeChange={applyDebugMode}
            obsidianFocusPath={obsidianFocusPath}
            onObsidianFocusConsumed={consumeObsidianFocus}
            onOpenObsidianNote={openObsidianNote}
          />
        </AppShell.Main>
      </AppShell>
    </ThemeRoot>
  );
}

export default ManagementApp;
