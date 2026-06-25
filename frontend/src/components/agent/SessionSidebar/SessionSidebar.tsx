import { useEffect, useMemo, useState } from 'react';
import {
  IconDots,
  IconLayoutSidebar,
  IconMessagePlus,
  IconPin,
  IconPinned,
  IconPencil,
  IconTrash,
} from '@tabler/icons-react';
import {
  ActionIcon,
  Box,
  Button,
  Divider,
  Group,
  Menu,
  Modal,
  ScrollArea,
  Stack,
  Text,
  TextInput,
  UnstyledButton,
} from '@mantine/core';
import { loadPinnedSessionIds, savePinnedSessionIds } from './pinnedSessions';
import classes from './SessionSidebar.module.css';

const isMacOS =
  typeof navigator !== 'undefined' &&
  (navigator.platform?.includes('Mac') || navigator.userAgent.includes('Mac'));

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

interface SessionRowProps {
  session: ChatSessionItem;
  isActive: boolean;
  isPinned: boolean;
  onSelectSession: (sessionId: string) => void;
  onTogglePin: (sessionId: string) => void;
  onRename: (session: ChatSessionItem) => void;
  onDelete: (session: ChatSessionItem) => void;
}

function SessionRow({
  session,
  isActive,
  isPinned,
  onSelectSession,
  onTogglePin,
  onRename,
  onDelete,
}: SessionRowProps) {
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
      <Group gap={2} wrap="nowrap" className={classes.sessionActions}>
        <ActionIcon
          variant="subtle"
          size="sm"
          aria-label={isPinned ? 'Unpin chat' : 'Pin chat'}
          onClick={() => onTogglePin(session.id)}
        >
          {isPinned ? <IconPinned size={14} /> : <IconPin size={14} />}
        </ActionIcon>
        <Menu withinPortal position="bottom-end" shadow="sm">
          <Menu.Target>
            <ActionIcon variant="subtle" size="sm" aria-label="Chat options">
              <IconDots size={14} />
            </ActionIcon>
          </Menu.Target>
          <Menu.Dropdown>
            <Menu.Item leftSection={<IconPencil size={14} />} onClick={() => onRename(session)}>
              Rename
            </Menu.Item>
            <Menu.Item
              leftSection={<IconTrash size={14} />}
              color="red"
              onClick={() => onDelete(session)}
            >
              Delete
            </Menu.Item>
          </Menu.Dropdown>
        </Menu>
      </Group>
    </Group>
  );
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
  const [deleteSession, setDeleteSession] = useState<ChatSessionItem | null>(null);
  const [deleting, setDeleting] = useState(false);
  const [creating, setCreating] = useState(false);
  const [pinnedIds, setPinnedIds] = useState<string[]>(() => loadPinnedSessionIds());

  useEffect(() => {
    const sessionIdSet = new Set(sessions.map((session) => session.id));
    setPinnedIds((current) => {
      const next = current.filter((id) => sessionIdSet.has(id));
      if (next.length !== current.length) {
        savePinnedSessionIds(next);
      }
      return next;
    });
  }, [sessions]);

  const { pinnedSessions, recentSessions } = useMemo(() => {
    const sessionById = new Map(sessions.map((session) => [session.id, session]));
    const pinned = pinnedIds
      .map((id) => sessionById.get(id))
      .filter((session): session is ChatSessionItem => session !== undefined);
    const pinnedIdSet = new Set(pinnedIds);
    const recent = sessions.filter((session) => !pinnedIdSet.has(session.id));
    return { pinnedSessions: pinned, recentSessions: recent };
  }, [sessions, pinnedIds]);

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

  const handleTogglePin = (sessionId: string) => {
    setPinnedIds((current) => {
      const isPinned = current.includes(sessionId);
      const next = isPinned
        ? current.filter((id) => id !== sessionId)
        : [sessionId, ...current];
      savePinnedSessionIds(next);
      return next;
    });
  };

  const openDelete = (session: ChatSessionItem) => {
    setDeleteSession(session);
  };

  const closeDelete = () => {
    if (deleting) {
      return;
    }
    setDeleteSession(null);
  };

  const handleConfirmDelete = async () => {
    if (!deleteSession) {
      return;
    }
    setDeleting(true);
    try {
      setPinnedIds((current) => {
        const next = current.filter((id) => id !== deleteSession.id);
        savePinnedSessionIds(next);
        return next;
      });
      await onDeleteSession(deleteSession.id);
      setDeleteSession(null);
    } finally {
      setDeleting(false);
    }
  };

  const renderSession = (session: ChatSessionItem) => (
    <SessionRow
      key={session.id}
      session={session}
      isActive={session.id === activeSessionId}
      isPinned={pinnedIds.includes(session.id)}
      onSelectSession={onSelectSession}
      onTogglePin={handleTogglePin}
      onRename={openRename}
      onDelete={openDelete}
    />
  );

  const hasPinned = pinnedSessions.length > 0;
  const hasRecents = recentSessions.length > 0;
  const isEmpty = sessions.length === 0;

  return (
    <Box className={classes.sidebar} data-platform={isMacOS ? 'darwin' : undefined}>
      <Box className={classes.header}>
        <Box className={classes.titlebarSpacer} aria-hidden="true" />
        <Group className={classes.headerRow} justify="space-between" wrap="nowrap">
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
        </Group>
      </Box>

      <ScrollArea className={classes.list} type="auto">
        <Stack gap={4} p="xs">
          {isEmpty ? (
            <Text c="dimmed" size="sm" p="sm">
              No conversations yet.
            </Text>
          ) : (
            <>
              {hasPinned && (
                <Stack gap={4}>
                  <Text className={classes.sectionLabel} size="xs">
                    Pinned
                  </Text>
                  {pinnedSessions.map(renderSession)}
                </Stack>
              )}

              {hasPinned && hasRecents && <Divider my="xs" />}

              <Stack gap={4}>
                <Text className={classes.sectionLabel} size="xs">
                  Recents
                </Text>
                {hasRecents ? (
                  recentSessions.map(renderSession)
                ) : (
                  <Text c="dimmed" size="sm" px="sm">
                    No recent conversations.
                  </Text>
                )}
              </Stack>
            </>
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

      <Modal
        opened={deleteSession !== null}
        onClose={closeDelete}
        title="Delete chat?"
        centered
      >
        <Stack gap="md">
          <Text size="sm">
            Delete &quot;{deleteSession?.title}&quot;? This removes the conversation and all
            messages.
          </Text>
          <Group justify="flex-end">
            <Button variant="default" onClick={closeDelete} disabled={deleting}>
              Cancel
            </Button>
            <Button color="red" loading={deleting} onClick={() => void handleConfirmDelete()}>
              Delete
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Box>
  );
}
