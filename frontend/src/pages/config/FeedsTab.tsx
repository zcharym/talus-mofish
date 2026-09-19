import { useState } from "react";
import { Box, Button, Fieldset, LoadingOverlay, PasswordInput, Stack, Text, TextInput } from "@mantine/core";
import { FeedsService, toApiError } from "../../utils/api";
import { notify } from "../../services/notifications";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface FeedsTabProps {
  youtubeApiKey: string;
  bilibiliSessdata: string;
  feedsProxyUrl: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function FeedsTab({
  youtubeApiKey,
  bilibiliSessdata,
  feedsProxyUrl,
  onChange,
}: FeedsTabProps) {
  const [testingYouTube, setTestingYouTube] = useState(false);
  const [testingBilibili, setTestingBilibili] = useState(false);
  const testing = testingYouTube || testingBilibili;

  const testYouTube = async () => {
    setTestingYouTube(true);
    try {
      const result = await FeedsService.PingYouTube();
      notify.success("YouTube", result.detail || "Data API key works.");
    } catch (err) {
      notify.failed("YouTube", toApiError(err));
    } finally {
      setTestingYouTube(false);
    }
  };

  const testBilibili = async () => {
    setTestingBilibili(true);
    try {
      const result = await FeedsService.PingBilibili();
      notify.success("Bilibili", result.detail || "SESSDATA accepted.");
    } catch (err) {
      notify.failed("Bilibili", toApiError(err));
    } finally {
      setTestingBilibili(false);
    }
  };

  return (
    <Box pos="relative" p="md">
      <LoadingOverlay visible={testing} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
      <Stack gap="md">
        <Text size="sm" c="dimmed">
          Optional credentials and proxy for the Agent Feeds tab. YouTube / Google APIs often need a
          local HTTP proxy. Leave proxy empty to use the OS system proxy; set e.g. http://127.0.0.1:7890
          for Clash / V2Ray. Save configuration, then test.
        </Text>
        <Fieldset legend="Proxy">
          <TextInput
            label="HTTP proxy"
            description="Used for YouTube, RSS, and Bilibili fetches. Empty = system proxy."
            placeholder="http://127.0.0.1:7890"
            value={feedsProxyUrl}
            onChange={(event) => onChange("feedsProxyUrl", event.currentTarget.value)}
          />
        </Fieldset>
        <Fieldset legend="YouTube">
          <Stack gap="sm">
            <PasswordInput
              label="YouTube Data API key"
              description="Google Cloud Console → YouTube Data API v3. Used to resolve @handles to channel RSS."
              value={youtubeApiKey}
              onChange={(event) => onChange("youtubeApiKey", event.currentTarget.value)}
            />
            <Button variant="light" onClick={() => void testYouTube()} loading={testingYouTube}>
              Test YouTube
            </Button>
          </Stack>
        </Fieldset>
        <Fieldset legend="Bilibili">
          <Stack gap="sm">
            <PasswordInput
              label="Bilibili SESSDATA"
              description="From browser cookies on bilibili.com. Required for 稍后再看 / watch later."
              value={bilibiliSessdata}
              onChange={(event) => onChange("bilibiliSessdata", event.currentTarget.value)}
            />
            <Button variant="light" onClick={() => void testBilibili()} loading={testingBilibili}>
              Test Bilibili
            </Button>
          </Stack>
        </Fieldset>
      </Stack>
    </Box>
  );
}
