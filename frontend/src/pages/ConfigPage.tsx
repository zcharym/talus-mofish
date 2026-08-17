import { useState } from "react";
import { Button, Code, Group, HoverCard, Stack, Tabs, Text } from "@mantine/core";
import {
  IconBook,
  IconKey,
  IconRobot,
  IconSettings,
} from "@tabler/icons-react";
import { useAppConfig } from "../hooks/useAppConfig";
import type { ThemeOption } from "../types/theme";
import { AITab } from "./config/AITab";
import { EnglishLearningTab } from "./config/EnglishLearningTab";
import { GeneralTab } from "./config/GeneralTab";
import { OAuthTab } from "./config/OAuthTab";

interface ConfigPageProps {
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange?: (enabled: boolean) => void;
}

const CONFIG_TABS = [
  { value: "general", label: "General", icon: IconSettings },
  { value: "english", label: "English Learning", icon: IconBook },
  { value: "ai", label: "AI", icon: IconRobot },
  { value: "oauth", label: "OAuth", icon: IconKey },
] as const;

export function ConfigPage({ onThemeChange, onDebugModeChange }: ConfigPageProps) {
  const [activeTab, setActiveTab] = useState<string | null>("general");
  const { form, updateForm, configPath, loading, saving, save } = useAppConfig({
    onThemeChange,
    onDebugModeChange,
  });

  if (loading) {
    return <Text c="dimmed">Loading configuration...</Text>;
  }

  return (
    <Stack maw={560} gap="md">
      <Tabs value={activeTab} onChange={setActiveTab}>
        <Tabs.List>
          {CONFIG_TABS.map((tab) => (
            <Tabs.Tab
              key={tab.value}
              value={tab.value}
              leftSection={<tab.icon size={14} stroke={1.5} />}
            >
              {tab.label}
            </Tabs.Tab>
          ))}
        </Tabs.List>

        <Tabs.Panel value="general" pt="md">
          <GeneralTab
            theme={form.theme}
            autoStart={form.autoStart}
            debugMode={form.debugMode}
            onChange={updateForm}
          />
        </Tabs.Panel>

        <Tabs.Panel value="english" pt="md">
          <EnglishLearningTab
            dailyGoalMinutes={form.dailyGoalMinutes}
            wordsPerSession={form.wordsPerSession}
            onChange={updateForm}
          />
        </Tabs.Panel>

        <Tabs.Panel value="ai" pt="md">
          <AITab
            aiProvider={form.aiProvider}
            aiModel={form.aiModel}
            aiAPIKey={form.aiAPIKey}
            aiBaseURL={form.aiBaseURL}
            onChange={updateForm}
          />
        </Tabs.Panel>

        <Tabs.Panel value="oauth" pt="md">
          <OAuthTab
            githubClientId={form.githubClientId}
            githubClientSecret={form.githubClientSecret}
            googleClientId={form.googleClientId}
            googleClientSecret={form.googleClientSecret}
            onChange={updateForm}
          />
        </Tabs.Panel>
      </Tabs>

      <Group>
        {configPath ? (
          <HoverCard width={320} shadow="md" withArrow openDelay={200}>
            <HoverCard.Target>
              <Button onClick={() => void save()} loading={saving}>
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
          <Button onClick={() => void save()} loading={saving}>
            Save configuration
          </Button>
        )}
      </Group>
    </Stack>
  );
}
