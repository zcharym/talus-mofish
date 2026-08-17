import { Code, PasswordInput, Stack, Text, TextInput } from "@mantine/core";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface OAuthTabProps {
  githubClientId: string;
  githubClientSecret: string;
  googleClientId: string;
  googleClientSecret: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function OAuthTab({
  githubClientId,
  githubClientSecret,
  googleClientId,
  googleClientSecret,
  onChange,
}: OAuthTabProps) {
  return (
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        GitHub uses the web application flow with a loopback redirect on 127.0.0.1. Register callback
        URL <Code>http://127.0.0.1/callback</Code> on your GitHub OAuth app; the app may use any
        local port at runtime.
      </Text>

      <TextInput
        label="GitHub client ID"
        value={githubClientId}
        onChange={(event) => onChange("githubClientId", event.currentTarget.value)}
        placeholder="Ov23li..."
      />

      <PasswordInput
        label="GitHub client secret"
        description="Required by GitHub when exchanging the authorization code for an access token"
        value={githubClientSecret}
        onChange={(event) => onChange("githubClientSecret", event.currentTarget.value)}
      />

      <TextInput
        label="Google client ID"
        value={googleClientId}
        onChange={(event) => onChange("googleClientId", event.currentTarget.value)}
        placeholder="1234567890-abc.apps.googleusercontent.com"
      />

      <PasswordInput
        label="Google client secret"
        description="Use a Desktop OAuth client in Google Cloud Console"
        value={googleClientSecret}
        onChange={(event) => onChange("googleClientSecret", event.currentTarget.value)}
      />
    </Stack>
  );
}
