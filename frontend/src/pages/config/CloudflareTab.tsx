import { useState } from "react";
import { Button, PasswordInput, Stack, Text, TextInput } from "@mantine/core";
import { CloudflareService, toApiError } from "../../utils/api";
import { notify } from "../../services/notifications";
import type { AppConfigForm } from "../../hooks/useAppConfig";

interface CloudflareTabProps {
  cloudflareAccountId: string;
  cloudflareAPIToken: string;
  onChange: <K extends keyof AppConfigForm>(key: K, value: AppConfigForm[K]) => void;
}

export function CloudflareTab({
  cloudflareAccountId,
  cloudflareAPIToken,
  onChange,
}: CloudflareTabProps) {
  const [testing, setTesting] = useState(false);

  const testConnection = async () => {
    setTesting(true);
    try {
      const result = await CloudflareService.Ping();
      const name = result.accountName || result.accountId || "account";
      notify.success("Cloudflare", `Connected to ${name}.`);
    } catch (err) {
      notify.failed("Cloudflare", toApiError(err));
    } finally {
      setTesting(false);
    }
  };

  return (
    <Stack gap="md">
      <Text size="sm" c="dimmed">
        API token and account ID unlock a pinned Cloudflare dashboard in the Agent window. Save
        configuration, then test. Use a read-only token with Account Settings Read, Workers Scripts
        Read, Account Analytics Read, D1 Read, and Workers KV Storage Read.
      </Text>
      <TextInput
        label="Account ID"
        description="From Cloudflare dashboard → Overview, or any Worker’s API credentials panel."
        placeholder="0123456789abcdef0123456789abcdef"
        value={cloudflareAccountId}
        onChange={(event) => onChange("cloudflareAccountId", event.currentTarget.value)}
      />
      <PasswordInput
        label="API token"
        description="Create at Cloudflare dashboard → My Profile → API Tokens. The token is stored in the OS keyring."
        value={cloudflareAPIToken}
        onChange={(event) => onChange("cloudflareAPIToken", event.currentTarget.value)}
      />
      <Button variant="light" onClick={() => void testConnection()} loading={testing}>
        Test connection
      </Button>
    </Stack>
  );
}
