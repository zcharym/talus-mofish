import { Select, Stack, Switch } from "@mantine/core";
import type { AppConfigForm } from "../../hooks/useAppConfig";
import type { ThemeOption } from "../../types/theme";

interface GeneralTabProps {
  theme: ThemeOption;
  autoStart: boolean;
  debugMode: boolean;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function GeneralTab({ theme, autoStart, debugMode, onChange }: GeneralTabProps) {
  return (
    <Stack gap="md">
      <Select
        label="Theme"
        description="Application color scheme"
        value={theme}
        onChange={(value) => onChange("theme", (value as ThemeOption) ?? "auto")}
        data={[
          { value: "auto", label: "Auto (system)" },
          { value: "light", label: "Light" },
          { value: "dark", label: "Dark" },
        ]}
      />

      <Switch
        label="Start at login"
        description="Launch Talus Echo automatically when you sign in"
        checked={autoStart}
        onChange={(event) => onChange("autoStart", event.currentTarget.checked)}
      />

      <Switch
        label="Debug mode"
        description="Show a Debug tab in the management sidebar for component previews"
        checked={debugMode}
        onChange={(event) => onChange("debugMode", event.currentTarget.checked)}
      />
    </Stack>
  );
}
