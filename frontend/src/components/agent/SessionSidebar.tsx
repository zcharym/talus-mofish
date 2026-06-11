import { useState } from 'react';
import {
  IconLayoutSidebar,
  IconMessagePlus,
  IconPencil,
  IconTrash,
} from '@tabler/icons-react';
import {
  ActionIcon,
  Box,
  Button,
  Group,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  UnstyledButton,
} from '@mantine/core';
import classes from './SessionSidebar.module.css';

export interface ChatSessionItem {
  id: string;
  title: string;
  created_at: string;
  updated_at: string;
}

interface SessionSidebarProps {
  sessions: ChatSessionItem[];
  activeSessionId: string | null;
  onSelectSession: (sessionId: string) => void;
  onNewChat: () => Promise<void>;
  onRenameSession: (sessionId: string, title: string) => Promise<void>;
  onDeleteSession: (sessionId: string) => Promise<void>;
  onOpenManagement: () => void;
}

export function SessionSidebar({
  sessions,
  activeSessionId,
  onSelectSession,
  onNewChat,
  onRenameSession,
  onDeleteSession,
  onOpenManagement,
}: SessionSidebarProps) {
  const [renameSessionId, setRenameSessionId] = useState<string | null>(null);
  const [renameTitle, setRenameTitle] = useState('');
  const [creating, setCreating] = useState(false);

  const openRename = (session: ChatSessionItem) => {
    setRenameSessionId(session.id);
    setRenameTitle(session.title);
  };

  const closeRename = () => {
    setRenameSessionId(null);
    setRenameTitle('');
  };

  const handleRename = async () => {
    if (!renameSessionId) {
      return;
    }
    const title = renameTitle.trim();
    if (!title) {
      return;
    }
    await onRenameSession(renameSessionId, title);
    closeRename();
  };

  const handleNewChat = async () => {
    setCreating(true);
    try {
      await onNewChat();
    } finally {
      setCreating(false);
    }
  };

  return (
    <Box className={classes.sidebar}>
      <Box className={classes.header}>
        <Text fw={600} size="sm">
          Chats
        </Text>
        <Button
          leftSection={<IconMessagePlus size={16} />}
          size="xs"
          variant="light"
          loading={creating}
          onClick={() => void handleNewChat()}
        >
          New chat
        </Button>
      </Box>

      <ScrollArea className={classes.list} type="auto">
        <Stack gap={4} p="xs">
          {sessions.length === 0 ? (
            <Text c="dimmed" size="sm" p="sm">
              No conversations yet.
            </Text>
          ) : (
            sessions.map((session) => {
              const isActive = session.id === activeSessionId;
              return (
                <Group key={session.id} gap={4} wrap="nowrap" className={classes.sessionRow}>
                  <UnstyledButton
                    className={classes.sessionButton}
                    data-active={isActive || undefined}
                    onClick={() => onSelectSession(session.id)}
                  >
                    <Text size="sm" lineClamp={1}>
                      {session.title}
                    </Text>
                  </UnstyledButton>
                  <ActionIcon
                    variant="subtle"
                    size="sm"
                    aria-label="Rename chat"
                    onClick={() => openRename(session)}
                  >
                    <IconPencil size={14} />
                  </ActionIcon>
                  <ActionIcon
                    variant="subtle"
                    size="sm"
                    color="red"
                    aria-label="Delete chat"
                    onClick={() => {
                      if (window.confirm(`Delete "${session.title}"?`)) {
                        void onDeleteSession(session.id);
                      }
                    }}
                  >
                    <IconTrash size={14} />
                  </ActionIcon>
                </Group>
              );
            })
          )}
        </Stack>
      </ScrollArea>

      <Box className={classes.footer}>
        <UnstyledButton className={classes.footerButton} onClick={onOpenManagement}>
          <IconLayoutSidebar size={18} />
          <span>Open Management</span>
        </UnstyledButton>
      </Box>

      <Modal opened={renameSessionId !== null} onClose={closeRename} title="Rename chat" centered>
        <Stack>
          <TextInput
            label="Title"
            value={renameTitle}
            onChange={(event) => setRenameTitle(event.currentTarget.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                event.preventDefault();
                void handleRename();
              }
            }}
          />
          <Group justify="flex-end">
            <Button variant="default" onClick={closeRename}>
              Cancel
            </Button>
            <Button onClick={() => void handleRename()} disabled={!renameTitle.trim()}>
              Save
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Box>
  );
}
