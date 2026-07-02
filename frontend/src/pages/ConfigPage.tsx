import { useCallback, useEffect, useState } from "react";
import {
  Button,
  Code,
  Group,
  HoverCard,
  NumberInput,
  PasswordInput,
  Select,
  Stack,
  Switch,
  Text,
  TextInput,
} from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { App as AppConfig, OAuth } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/config/models";
import { Config as AIConfig, Provider } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/aiclient/models";
import { notify } from "../services/notifications";

type ThemeOption = "auto" | "light" | "dark";

interface ConfigPageProps {
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange?: (enabled: boolean) => void;
}

export function ConfigPage({ onThemeChange, onDebugModeChange }: ConfigPageProps) {
  const [theme, setTheme] = useState<ThemeOption>("auto");
  const [dailyGoalMinutes, setDailyGoalMinutes] = useState(30);
  const [wordsPerSession, setWordsPerSession] = useState(20);
  const [autoStart, setAutoStart] = useState(false);
  const [debugMode, setDebugMode] = useState(false);
  const [aiProvider, setAIProvider] = useState<string>(Provider.ProviderOpenAI);
  const [aiModel, setAIModel] = useState("gpt-4o-mini");
  const [aiAPIKey, setAIAPIKey] = useState("");
  const [aiBaseURL, setAIBaseURL] = useState("");
  const [githubClientId, setGithubClientId] = useState("");
  const [githubClientSecret, setGithubClientSecret] = useState("");
  const [googleClientId, setGoogleClientId] = useState("");
  const [googleClientSecret, setGoogleClientSecret] = useState("");
  const [oauthHttpProxy, setOauthHttpProxy] = useState("");
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
      setDebugMode(cfg.debugMode ?? false);
      onDebugModeChange?.(cfg.debugMode ?? false);
      setAIProvider(cfg.ai?.provider || Provider.ProviderOpenAI);
      setAIModel(cfg.ai?.model || "gpt-4o-mini");
      setAIAPIKey(cfg.ai?.apiKey || "");
      setAIBaseURL(cfg.ai?.baseURL || "");
      setGithubClientId(cfg.oauth?.githubClientId || "");
      setGithubClientSecret(cfg.oauth?.githubClientSecret || "");
      setGoogleClientId(cfg.oauth?.googleClientId || "");
      setGoogleClientSecret(cfg.oauth?.googleClientSecret || "");
      setOauthHttpProxy(cfg.oauth?.httpProxy || "");
      setConfigPath(path);
      onThemeChange(nextTheme);
    } catch (err) {
      console.error(err);
      notify.failed("Error", "Failed to load configuration.");
    } finally {
      setLoading(false);
    }
  }, [onThemeChange, onDebugModeChange]);

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
      debugMode,
      ai: new AIConfig({
        provider: aiProvider as Provider,
        model: aiModel,
        apiKey: aiAPIKey,
        baseURL: aiBaseURL,
      }),
      oauth: new OAuth({
        githubClientId,
        githubClientSecret,
        googleClientId,
        googleClientSecret,
        httpProxy: oauthHttpProxy,
      }),
    });

    try {
      await AppService.SaveConfig(payload);
      onThemeChange(theme);
      onDebugModeChange?.(debugMode);
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

      <Text fw={600} mt="md">
        English Learning
      </Text>
      <Text size="sm" c="dimmed">
        Settings for the English Learning domain in Management and Agent quick actions.
      </Text>

      <NumberInput
        label="Daily study goal"
        description="Target minutes per day for English Learning sessions"
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
        description="Launch Talus Echo automatically when you sign in"
        checked={autoStart}
        onChange={(event) => setAutoStart(event.currentTarget.checked)}
      />

      <Switch
        label="Debug mode"
        description="Show a Debug tab in the management sidebar for component previews"
        checked={debugMode}
        onChange={(event) => setDebugMode(event.currentTarget.checked)}
      />

      <Text fw={600} mt="md">
        AI agent
      </Text>
      <Text size="sm" c="dimmed">
        Provider settings for Talus Agent chat responses.
      </Text>

      <Select
        label="Provider"
        value={aiProvider}
        onChange={(value) => setAIProvider(value ?? Provider.ProviderOpenAI)}
        data={[
          { value: Provider.ProviderOpenAI, label: "OpenAI" },
          { value: Provider.ProviderDeepSeek, label: "DeepSeek" },
          { value: Provider.ProviderMoonshot, label: "Moonshot (Kimi)" },
          { value: Provider.ProviderOllama, label: "Ollama (local)" },
        ]}
      />

      <TextInput
        label="Model"
        description="Model name for the selected provider"
        value={aiModel}
        onChange={(event) => setAIModel(event.currentTarget.value)}
        placeholder="gpt-4o-mini"
      />

      <PasswordInput
        label="API key"
        description={
          aiProvider === Provider.ProviderOllama
            ? "Optional for local Ollama"
            : "Required for cloud providers"
        }
        value={aiAPIKey}
        onChange={(event) => setAIAPIKey(event.currentTarget.value)}
      />

      <TextInput
        label="Base URL"
        description="Optional API root, e.g. https://api.openai.com/v1 or https://api.moonshot.cn/v1 (not /anthropic)"
        value={aiBaseURL}
        onChange={(event) => setAIBaseURL(event.currentTarget.value)}
      />

      <Text fw={600} mt="md">
        Account / OAuth
      </Text>
      <Text size="sm" c="dimmed">
        GitHub uses the web application flow with a loopback redirect on 127.0.0.1. Register
        callback URL <Code>http://127.0.0.1/callback</Code> on your GitHub OAuth app; the app may
        use any local port at runtime.
      </Text>

      <TextInput
        label="GitHub client ID"
        value={githubClientId}
        onChange={(event) => setGithubClientId(event.currentTarget.value)}
        placeholder="Ov23li..."
      />

      <PasswordInput
        label="GitHub client secret"
        description="Required by GitHub when exchanging the authorization code for an access token"
        value={githubClientSecret}
        onChange={(event) => setGithubClientSecret(event.currentTarget.value)}
      />

      <TextInput
        label="Google client ID"
        value={googleClientId}
        onChange={(event) => setGoogleClientId(event.currentTarget.value)}
        placeholder="1234567890-abc.apps.googleusercontent.com"
      />

      <PasswordInput
        label="Google client secret"
        description="Use a Desktop OAuth client in Google Cloud Console"
        value={googleClientSecret}
        onChange={(event) => setGoogleClientSecret(event.currentTarget.value)}
      />

      <TextInput
        label="OAuth HTTP proxy"
        description="Optional override, e.g. http://127.0.0.1:7890. Leave empty to use Windows system proxy settings."
        value={oauthHttpProxy}
        onChange={(event) => setOauthHttpProxy(event.currentTarget.value)}
        placeholder="http://127.0.0.1:7890"
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
