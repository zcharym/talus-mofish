import { useMemo } from 'react';
import {
  FloatingChat,
  createAgentOverlayTransport,
  createCloudflareChatDomain,
  useChatModal,
} from '../ChatModal';
import type { CloudflareDashboardSnapshot } from '../../../utils/api';

interface CloudflareFloatingChatProps {
  opened: boolean;
  onClose: () => void;
  configured: boolean;
  snapshot: CloudflareDashboardSnapshot | null;
  loadError?: string;
}

export function CloudflareFloatingChat({
  opened,
  onClose,
  configured,
  snapshot,
  loadError,
}: CloudflareFloatingChatProps) {
  const domain = useMemo(
    () => createCloudflareChatDomain({ configured, snapshot, loadError }),
    [configured, snapshot, loadError],
  );
  const transport = useMemo(() => createAgentOverlayTransport(), []);
  const { messages, sending, streamingMessageId, send, cancel } = useChatModal({ domain, transport });

  return (
    <FloatingChat
      opened={opened}
      onClose={onClose}
      domain={domain}
      messages={messages}
      sending={sending}
      onSend={send}
      onCancel={streamingMessageId ? cancel : undefined}
    />
  );
}
