import { useState } from "react";
import { Button, PasswordInput, Stack, Text } from "@mantine/core";
import { FeedsService, toApiError } from "../../utils/api";
import { notify } from "../../services/notifications";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface FeedsTabProps {
  youtubeApiKey: string;
  bilibiliSessdata: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function FeedsTab({ youtubeApiKey, bilibiliSessdata, onChange }: FeedsTabProps) {
  const [testingYouTube, setTestingYouTube] = useState(false);
  const [testingBilibili, setTestingBilibili] = useState(false);

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
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        Optional credentials for the Agent Feeds tab. RSS works with no keys. YouTube handles resolve
        more reliably with a Data API key. Bilibili watch later needs the SESSDATA cookie from a
        logged-in browser session. Secrets are stored in the OS keyring. Save configuration, then
        test.
      </Text>
      <PasswordInput
        label="YouTube Data API key"
        description="Google Cloud Console → YouTube Data API v3. Used to resolve @handles to channel RSS."
        value={youtubeApiKey}
        onChange={(event) => onChange("youtubeApiKey", event.currentTarget.value)}
      />
      <Button variant="light" onClick={() => void testYouTube()} loading={testingYouTube}>
        Test YouTube
      </Button>
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
  );
}
