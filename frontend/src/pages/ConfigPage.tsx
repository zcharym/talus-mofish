import { useCallback, useEffect, useState } from 'react';
import {
  Alert,
  Button,
  Code,
  Group,
  NumberInput,
  Select,
  Stack,
  Text,
} from '@mantine/core';
import { AppService } from '../../bindings/github.com/songwei.ma/talus-mofish';
import { App as AppConfig } from '../../bindings/github.com/songwei.ma/talus-mofish/internal/config/models';

type ThemeOption = 'auto' | 'light' | 'dark';

interface ConfigPageProps {
  onThemeChange: (theme: ThemeOption) => void;
}

export function ConfigPage({ onThemeChange }: ConfigPageProps) {
  const [theme, setTheme] = useState<ThemeOption>('auto');
  const [dailyGoalMinutes, setDailyGoalMinutes] = useState(30);
  const [wordsPerSession, setWordsPerSession] = useState(20);
  const [configPath, setConfigPath] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [saved, setSaved] = useState(false);

  const loadConfig = useCallback(async () => {
    setLoading(true);
    setError(null);
    setSaved(false);

    try {
      const [cfg, path] = await Promise.all([
        AppService.GetConfig(),
        AppService.ConfigPath(),
      ]);

      const nextTheme = (cfg.theme as ThemeOption) || 'auto';
      setTheme(nextTheme);
      setDailyGoalMinutes(cfg.dailyGoalMinutes);
      setWordsPerSession(cfg.wordsPerSession);
      setConfigPath(path);
      onThemeChange(nextTheme);
    } catch (err) {
      console.error(err);
      setError('Failed to load configuration.');
    } finally {
      setLoading(false);
    }
  }, [onThemeChange]);

  useEffect(() => {
    void loadConfig();
  }, [loadConfig]);

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    setSaved(false);

    const payload = new AppConfig({
      theme,
      dailyGoalMinutes,
      wordsPerSession,
    });

    try {
      await AppService.SaveConfig(payload);
      onThemeChange(theme);
      setSaved(true);
    } catch (err) {
      console.error(err);
      setError('Failed to save configuration.');
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <Text c="dimmed">Loading configuration...</Text>;
  }

  return (
    <Stack maw={480} gap="md">
      {error ? (
        <Alert color="red" title="Error">
          {error}
        </Alert>
      ) : null}

      {saved ? (
        <Alert color="green" title="Saved">
          Configuration saved to config.json.
        </Alert>
      ) : null}

      <Select
        label="Theme"
        description="Application color scheme"
        value={theme}
        onChange={(value) => setTheme((value as ThemeOption) ?? 'auto')}
        data={[
          { value: 'auto', label: 'Auto (system)' },
          { value: 'light', label: 'Light' },
          { value: 'dark', label: 'Dark' },
        ]}
      />

      <NumberInput
        label="Daily study goal"
        description="Target minutes per day for English study"
        value={dailyGoalMinutes}
        onChange={(value) => setDailyGoalMinutes(Number(value) || 30)}
        min={1}
        max={480}
        suffix=" min"
      />

      <NumberInput
        label="Words per session"
        description="Number of words to practice in each recite session"
        value={wordsPerSession}
        onChange={(value) => setWordsPerSession(Number(value) || 20)}
        min={1}
        max={200}
      />

      <Group>
        <Button onClick={() => void handleSave()} loading={saving}>
          Save configuration
        </Button>
      </Group>

      {configPath ? (
        <Text size="sm" c="dimmed">
          Config file: <Code>{configPath}</Code>
        </Text>
      ) : null}
    </Stack>
  );
}
