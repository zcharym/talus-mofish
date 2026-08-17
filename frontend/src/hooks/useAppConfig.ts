import { useCallback, useEffect, useState } from "react";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { App as AppConfig, OAuth } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/config/models";
import { Config as AIConfig, Provider } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/aiclient/models";
import { notify } from "../services/notifications";
import type { ThemeOption } from "../types/theme";

export interface AppConfigForm {
  theme: ThemeOption;
  dailyGoalMinutes: number;
  wordsPerSession: number;
  autoStart: boolean;
  debugMode: boolean;
  aiProvider: string;
  aiModel: string;
  aiAPIKey: string;
  aiBaseURL: string;
  githubClientId: string;
  githubClientSecret: string;
  googleClientId: string;
  googleClientSecret: string;
}

const defaultForm: AppConfigForm = {
  theme: "auto",
  dailyGoalMinutes: 30,
  wordsPerSession: 20,
  autoStart: false,
  debugMode: false,
  aiProvider: Provider.ProviderOpenAI,
  aiModel: "gpt-4o-mini",
  aiAPIKey: "",
  aiBaseURL: "",
  githubClientId: "",
  githubClientSecret: "",
  googleClientId: "",
  googleClientSecret: "",
};

export interface UseAppConfigOptions {
  onThemeChange: (theme: ThemeOption) => void;
  onDebugModeChange?: (enabled: boolean) => void;
}

export function useAppConfig({ onThemeChange, onDebugModeChange }: UseAppConfigOptions) {
  const [form, setForm] = useState<AppConfigForm>(defaultForm);
  const [configPath, setConfigPath] = useState("");
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const updateForm = useCallback(<K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  }, []);

  const loadConfig = useCallback(async () => {
    setLoading(true);

    try {
      const [cfg, path] = await Promise.all([AppService.GetConfig(), AppService.ConfigPath()]);

      const nextTheme = (cfg.theme as ThemeOption) || "auto";
      const nextDebugMode = cfg.debugMode ?? false;

      setForm({
        theme: nextTheme,
        dailyGoalMinutes: cfg.dailyGoalMinutes,
        wordsPerSession: cfg.wordsPerSession,
        autoStart: cfg.autoStart,
        debugMode: nextDebugMode,
        aiProvider: cfg.ai?.provider || Provider.ProviderOpenAI,
        aiModel: cfg.ai?.model || "gpt-4o-mini",
        aiAPIKey: cfg.ai?.apiKey || "",
        aiBaseURL: cfg.ai?.baseURL || "",
        githubClientId: cfg.oauth?.githubClientId || "",
        githubClientSecret: cfg.oauth?.githubClientSecret || "",
        googleClientId: cfg.oauth?.googleClientId || "",
        googleClientSecret: cfg.oauth?.googleClientSecret || "",
      });
      setConfigPath(path);
      onThemeChange(nextTheme);
      onDebugModeChange?.(nextDebugMode);
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

  const save = useCallback(async () => {
    setSaving(true);

    const payload = new AppConfig({
      theme: form.theme,
      dailyGoalMinutes: form.dailyGoalMinutes,
      wordsPerSession: form.wordsPerSession,
      autoStart: form.autoStart,
      debugMode: form.debugMode,
      ai: new AIConfig({
        provider: form.aiProvider as Provider,
        model: form.aiModel,
        apiKey: form.aiAPIKey,
        baseURL: form.aiBaseURL,
      }),
      oauth: new OAuth({
        githubClientId: form.githubClientId,
        githubClientSecret: form.githubClientSecret,
        googleClientId: form.googleClientId,
        googleClientSecret: form.googleClientSecret,
      }),
    });

    try {
      await AppService.SaveConfig(payload);
      onThemeChange(form.theme);
      onDebugModeChange?.(form.debugMode);
      notify.success("Saved", "Configuration saved to config.json.");
    } catch (err) {
      console.error(err);
      notify.failed("Error", "Failed to save configuration.");
    } finally {
      setSaving(false);
    }
  }, [form, onThemeChange, onDebugModeChange]);

  return {
    form,
    updateForm,
    configPath,
    loading,
    saving,
    save,
  };
}
