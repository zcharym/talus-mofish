import { PasswordInput, Select, Stack, Text, TextInput } from "@mantine/core";
import { Provider } from "../../../bindings/github.com/songwei.ma/talus-mofish/backend/utils/aiclient/models";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface AITabProps {
  aiProvider: string;
  aiModel: string;
  aiAPIKey: string;
  aiBaseURL: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function AITab({ aiProvider, aiModel, aiAPIKey, aiBaseURL, onChange }: AITabProps) {
  return (
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        Provider settings for Talus Agent chat responses.
      </Text>

      <Select
        label="Provider"
        value={aiProvider}
        onChange={(value) => onChange("aiProvider", value ?? Provider.ProviderOpenAI)}
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
        onChange={(event) => onChange("aiModel", event.currentTarget.value)}
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
        onChange={(event) => onChange("aiAPIKey", event.currentTarget.value)}
      />

      <TextInput
        label="Base URL"
        description="Optional API root, e.g. https://api.openai.com/v1 or https://api.moonshot.cn/v1 (not /anthropic)"
        value={aiBaseURL}
        onChange={(event) => onChange("aiBaseURL", event.currentTarget.value)}
      />
    </Stack>
  );
}
