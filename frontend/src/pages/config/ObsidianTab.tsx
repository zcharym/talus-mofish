import { useState } from "react";
import { Box, Button, Fieldset, LoadingOverlay, PasswordInput, Stack, Text, TextInput } from "@mantine/core";
import { ObsidianService, toApiError } from "../../utils/api";
import { notify } from "../../services/notifications";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface ObsidianTabProps {
  obsidianBaseUrl: string;
  obsidianAPIKey: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function ObsidianTab({ obsidianBaseUrl, obsidianAPIKey, onChange }: ObsidianTabProps) {
  const [testing, setTesting] = useState(false);

  const testConnection = async () => {
    setTesting(true);
    try {
      const status = await ObsidianService.Ping();
      const version = status.versions?.self || status.ok || "ok";
      if (status.authenticated) {
        notify.success("Obsidian", `Connected (${version}).`);
      } else {
        notify.warning(
          "Obsidian",
          "Server reached, but the API key was not accepted. Save a valid key, then test again.",
        );
      }
    } catch (err) {
      notify.failed("Obsidian", toApiError(err));
    } finally {
      setTesting(false);
    }
  };

  return (
    <Box pos="relative" p="md">
      <LoadingOverlay visible={testing} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
      <Stack gap="md">
        <Text size="sm" c="dimmed">
          Connect to Obsidian Local REST API. Obsidian must be running with the plugin enabled. Save
          configuration, then test.
        </Text>
        <Fieldset legend="Connection">
          <Stack gap="sm">
            <TextInput
              label="Base URL"
              description="HTTPS default is https://127.0.0.1:27124 (self-signed). Use http://127.0.0.1:27123 if you enabled the HTTP server in the plugin."
              placeholder="https://127.0.0.1:27124"
              value={obsidianBaseUrl}
              onChange={(event) => onChange("obsidianBaseUrl", event.currentTarget.value)}
            />
            <PasswordInput
              label="API key"
              description="From Obsidian Settings → Local REST API."
              value={obsidianAPIKey}
              onChange={(event) => onChange("obsidianAPIKey", event.currentTarget.value)}
            />
            <Button variant="light" onClick={() => void testConnection()} loading={testing}>
              Test connection
            </Button>
          </Stack>
        </Fieldset>
      </Stack>
    </Box>
  );
}
