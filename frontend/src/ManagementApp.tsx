import { useCallback, useEffect, useState } from 'react';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import { AppShell, MantineProvider, Text, Title } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { AppService } from '../bindings/github.com/songwei.ma/talus-mofish';
import { NavbarSegmented } from './components/NavbarSegmented';
import { ConfigPage } from './pages/ConfigPage';
import { ImportPage } from './pages/ImportPage';
import { ReadingPage } from './pages/ReadingPage';
import { VocabularyPage } from './pages/VocabularyPage';

type ThemeOption = 'auto' | 'light' | 'dark';

const pageTitles: Record<string, string> = {
  import: 'Import',
  reading: 'Reading',
  vocabulary: 'Vocabulary',
  config: 'Configuration',
  about: 'About',
};

function MainContent({ activeItem, onThemeChange }: { activeItem: string; onThemeChange: (theme: ThemeOption) => void }) {
  const title = pageTitles[activeItem] ?? 'Talus Echo';

  if (activeItem === 'config') {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ConfigPage onThemeChange={onThemeChange} />
      </>
    );
  }

  if (activeItem === 'import') {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ImportPage />
      </>
    );
  }

  if (activeItem === 'reading') {
    return (
      <>
        <Title order={2}>{title}</Title>
        <ReadingPage />
      </>
    );
  }

  if (activeItem === 'vocabulary') {
    return (
      <>
        <Title order={2}>{title}</Title>
        <VocabularyPage />
      </>
    );
  }

  if (activeItem === 'about') {
    return (
      <>
        <Title order={2}>{title}</Title>
        <Text c="dimmed" mt="sm">
          Talus Echo — English learning library and management.
        </Text>
      </>
    );
  }

  return null;
}

function ManagementApp() {
  const [activeItem, setActiveItem] = useState('vocabulary');
  const [colorScheme, setColorScheme] = useState<ThemeOption>('auto');

  const applyTheme = useCallback((theme: ThemeOption) => {
    setColorScheme(theme);
  }, []);

  useEffect(() => {
    AppService.GetConfig()
      .then((cfg) => {
        const theme = (cfg.theme as ThemeOption) || 'auto';
        applyTheme(theme);
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, [applyTheme]);

  return (
    <MantineProvider
      defaultColorScheme="auto"
      forceColorScheme={colorScheme === 'auto' ? undefined : colorScheme}
    >
      <Notifications position="top-right" limit={5} />
      <AppShell navbar={{ width: 300, breakpoint: 'sm' }} padding="md">
        <AppShell.Navbar p={0}>
          <NavbarSegmented activeItem={activeItem} onActiveItemChange={setActiveItem} />
        </AppShell.Navbar>

        <AppShell.Main>
          <MainContent activeItem={activeItem} onThemeChange={applyTheme} />
        </AppShell.Main>
      </AppShell>
    </MantineProvider>
  );
}

export default ManagementApp;
