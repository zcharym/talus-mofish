import { useMemo, useState } from "react";
import { Button, Group, SegmentedControl, Stack, Text } from "@mantine/core";
import { useDisclosure } from "@mantine/hooks";
import {
  FloatingChat,
  createCloudflareChatDomain,
  createGenericChatDomain,
  createMockChatTransport,
  useChatModal,
} from "../../../components/agent/ChatModal";
import type { CloudflareDashboardSnapshot } from "../../../utils/api";

const SAMPLE_SNAPSHOT = {
  configured: true,
  fetchedAt: "2026-09-10T06:00:00Z",
  account: { id: "acct_demo", name: "Talus Echo demo" },
  kpis: {
    workerCount: 2,
    requests24h: 18420,
    errors24h: 37,
    errorRate: 0.002,
  },
  workers: [
    {
      name: "echo-watch",
      createdOn: "2026-01-12T00:00:00Z",
      modifiedOn: "2026-09-09T18:00:00Z",
      requests24h: 12010,
      errors24h: 31,
      cpuTimeP99: 4200,
      cron: ["*/5 * * * *"],
    },
    {
      name: "talus-auth",
      createdOn: "2026-02-01T00:00:00Z",
      modifiedOn: "2026-09-08T12:00:00Z",
      requests24h: 6410,
      errors24h: 6,
      cpuTimeP99: 1800,
      cron: [],
    },
  ],
  known: [
    {
      id: "talus-auth",
      title: "Talus Auth",
      present: true,
      bindingsNote: "D1 database AUTH_DB",
    },
    {
      id: "echo-watch",
      title: "Echo Watch",
      present: true,
      bindingsNote: "KV namespace",
    },
  ],
  d1: [{ id: "d1-auth", name: "AUTH_DB" }],
  kv: [{ id: "kv-watch", title: "ECHO_WATCH" }],
  errors: {
    account: "",
    workers: "",
    analytics: "",
    d1: "",
    kv: "",
  },
} as CloudflareDashboardSnapshot;

export function ChatModalDemo() {
  const [domainId, setDomainId] = useState<"generic" | "cloudflare">("generic");
  const [opened, { open, close }] = useDisclosure(false);
  const domain = useMemo(
    () =>
      domainId === "cloudflare"
        ? createCloudflareChatDomain({ configured: true, snapshot: SAMPLE_SNAPSHOT })
        : createGenericChatDomain(),
    [domainId],
  );
  const transport = useMemo(() => createMockChatTransport(), []);
  const { messages, sending, streamingMessageId, send, cancel } = useChatModal({ domain, transport });

  return (
    <Stack gap="md" maw={640}>
      <Text size="sm" c="dimmed">
        FloatingChat is a dedicated overlay widget (bottom-right), not a Modal or Drawer. ChatPanel
        is the inner thread. The Cloudflare fixture injects a sample dashboard snapshot. Replies stream
        locally — no Wails chat API.
      </Text>

      <SegmentedControl
        value={domainId}
        onChange={(value) => setDomainId(value as "generic" | "cloudflare")}
        data={[
          { value: "generic", label: "Generic" },
          { value: "cloudflare", label: "Cloudflare" },
        ]}
      />

      <Group>
        <Button onClick={open}>Open floating chat</Button>
      </Group>

      <FloatingChat
        opened={opened}
        onClose={close}
        domain={domain}
        messages={messages}
        sending={sending}
        onSend={send}
        onCancel={streamingMessageId ? cancel : undefined}
      />
    </Stack>
  );
}
