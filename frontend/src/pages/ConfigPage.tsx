import { useState } from "react";
import { Box, Button, Code, Group, HoverCard, LoadingOverlay, Stack, Tabs, Text, ThemeIcon } from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";
import {
  IconBook,
  IconCloud,
  IconGrid3x3,
  IconKey,
  IconMarkdown,
  IconRobot,
  IconRss,
  IconSettings,
} from "@tabler/icons-react";
import { useAppConfig } from "../hooks/useAppConfig";
import type { ThemeOption } from "../types/theme";
import { AITab } from "./config/AITab";
import { EnglishLearningTab } from "./config/EnglishLearningTab";
import { GeneralTab } from "./config/GeneralTab";
import { OAuthTab } from "./config/OAuthTab";
import { CloudflareTab } from "./config/CloudflareTab";
import { FeedsTab } from "./config/FeedsTab";
import { ObsidianTab } from "./config/ObsidianTab";
import { SudokuTab } from "./config/SudokuTab";
import classes from "./ConfigPage.module.css";

interface ConfigPageProps {
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange?: (enabled: boolean) => void;
}

const CONFIG_TABS = [
  { value: "general", label: "General", icon: IconSettings },
  { value: "english", label: "English Learning", icon: IconBook },
  { value: "ai", label: "AI", icon: IconRobot },
  { value: "oauth", label: "OAuth", icon: IconKey },
  { value: "sudoku", label: "Sudoku", icon: IconGrid3x3 },
  { value: "obsidian", label: "Obsidian", icon: IconMarkdown },
  { value: "cloudflare", label: "Cloudflare", icon: IconCloud },
  { value: "feeds", label: "Feeds", icon: IconRss },
] as const;

export function ConfigPage({ onThemeChange, onDebugModeChange }: ConfigPageProps) {
  const [activeTab, setActiveTab] = useState<string | null>("general");
  const isNarrow = useMediaQuery("(max-width: 48em)");
  const { form, updateForm, configPath, loading, saving, save } = useAppConfig({
    onThemeChange,
    onDebugModeChange,
  });

  if (loading) {
    return <Text c="dimmed">Loading configuration...</Text>;
  }

  return (
    <Box pos="relative" className={classes.page}>
      <LoadingOverlay visible={saving} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
      <Tabs
        value={activeTab}
        onChange={setActiveTab}
        orientation={isNarrow ? "horizontal" : "vertical"}
        className={classes.tabs}
      >
        <Tabs.List grow={isNarrow} className={classes.list}>
          {CONFIG_TABS.map((tab) => (
            <Tabs.Tab
              key={tab.value}
              value={tab.value}
              leftSection={
                <ThemeIcon size={22} variant="light" radius="sm">
                  <tab.icon size={14} stroke={1.5} />
                </ThemeIcon>
              }
            >
              {tab.label}
            </Tabs.Tab>
          ))}
        </Tabs.List>

        <Stack gap="md" className={classes.content}>

          <Tabs.Panel value="general" className={classes.panel} pt="md">
            <GeneralTab
              theme={form.theme}
              autoStart={form.autoStart}
              debugMode={form.debugMode}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="english" className={classes.panel} pt="md">
            <EnglishLearningTab
              dailyGoalMinutes={form.dailyGoalMinutes}
              wordsPerSession={form.wordsPerSession}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="ai" className={classes.panel} pt="md">
            <AITab
              aiProvider={form.aiProvider}
              aiModel={form.aiModel}
              aiAPIKey={form.aiAPIKey}
              aiBaseURL={form.aiBaseURL}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="oauth" className={classes.panel} pt="md">
            <OAuthTab
              githubClientId={form.githubClientId}
              githubClientSecret={form.githubClientSecret}
              googleClientId={form.googleClientId}
              googleClientSecret={form.googleClientSecret}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="sudoku" className={classes.panel} pt="md">
            <SudokuTab sudokuAPIKey={form.sudokuAPIKey} onChange={updateForm} />
          </Tabs.Panel>

          <Tabs.Panel value="obsidian" className={classes.panel} pt="md">
            <ObsidianTab
              obsidianBaseUrl={form.obsidianBaseUrl}
              obsidianAPIKey={form.obsidianAPIKey}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="cloudflare" className={classes.panel} pt="md">
            <CloudflareTab
              cloudflareAccountId={form.cloudflareAccountId}
              cloudflareAPIToken={form.cloudflareAPIToken}
              onChange={updateForm}
            />
          </Tabs.Panel>

          <Tabs.Panel value="feeds" className={classes.panel} pt="md">
            <FeedsTab
              youtubeApiKey={form.youtubeApiKey}
              bilibiliSessdata={form.bilibiliSessdata}
              feedsProxyUrl={form.feedsProxyUrl}
              onChange={updateForm}
            />
          </Tabs.Panel>

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
      </Tabs>
    </Box>
  );
}
