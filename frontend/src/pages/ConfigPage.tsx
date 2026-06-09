import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Code,
  Group,
  HoverCard,
  NumberInput,
  Select,
  Stack,
  Switch,
  Text,
} from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { App as AppConfig } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/config/models";
import { notify } from "../services/notifications";

type ThemeOption = "auto" | "light" | "dark";

interface ConfigPageProps {
  onThemeChange: (theme: ThemeOption) => void;
}

export function ConfigPage({ onThemeChange }: ConfigPageProps) {
  const [theme, setTheme] = useState<ThemeOption>("auto");
  const [dailyGoalMinutes, setDailyGoalMinutes] = useState(30);
  const [wordsPerSession, setWordsPerSession] = useState(20);
  const [autoStart, setAutoStart] = useState(false);
  const [configPath, setConfigPath] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const loadConfig = useCallback(async () => {
    setLoading(true);

    try {
      const [cfg, path] = await Promise.all([
        AppService.GetConfig(),
        AppService.ConfigPath(),
      ]);

      const nextTheme = (cfg.theme as ThemeOption) || "auto";
      setTheme(nextTheme);
      setDailyGoalMinutes(cfg.dailyGoalMinutes);
      setWordsPerSession(cfg.wordsPerSession);
      setAutoStart(cfg.autoStart);
      setConfigPath(path);
      onThemeChange(nextTheme);
    } catch (err) {
      console.error(err);
      notify.failed("Error", "Failed to load configuration.");
    } finally {
      setLoading(false);
    }
  }, [onThemeChange]);

  useEffect(() => {
    void loadConfig();
  }, [loadConfig]);

  const handleSave = async () => {
    setSaving(true);

    const payload = new AppConfig({
      theme,
      dailyGoalMinutes,
      wordsPerSession,
      autoStart,
    });

    try {
      await AppService.SaveConfig(payload);
      onThemeChange(theme);
      notify.success("Saved", "Configuration saved to config.json.");
    } catch (err) {
      console.error(err);
      notify.failed("Error", "Failed to save configuration.");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <Text c="dimmed">Loading configuration...</Text>;
  }

  return (
    <Stack maw={480} gap="md">
      <Select
        label="Theme"
        description="Application color scheme"
        value={theme}
        onChange={(value) => setTheme((value as ThemeOption) ?? "auto")}
        data={[
          { value: "auto", label: "Auto (system)" },
          { value: "light", label: "Light" },
          { value: "dark", label: "Dark" },
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

      <Switch
        label="Start at login"
        description="Launch Talus MoFish automatically when you sign in"
        checked={autoStart}
        onChange={(event) => setAutoStart(event.currentTarget.checked)}
      />

      <Group>
        {configPath ? (
          <HoverCard width={320} shadow="md" withArrow openDelay={200}>
            <HoverCard.Target>
              <Button onClick={() => void handleSave()} loading={saving}>
                Save configuration
              </Button>
            </HoverCard.Target>
            <HoverCard.Dropdown>
              <Text size="sm">
                Config file: <Code>{configPath}</Code>
              </Text>
            </HoverCard.Dropdown>
          </HoverCard>
        ) : (
          <Button onClick={() => void handleSave()} loading={saving}>
            Save configuration
          </Button>
        )}
      </Group>
    </Stack>
  );
}
