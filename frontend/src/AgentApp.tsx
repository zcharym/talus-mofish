import { useCallback, useEffect, useMemo, useState } from 'react';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import { Box, MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { AppService } from '../bindings/github.com/songwei.ma/talus-mofish';
import { ChatInput } from './components/agent/ChatInput';
import { ChatMessageItem, ChatThread } from './components/agent/ChatThread';
import { ChatSessionItem, SessionSidebar } from './components/agent/SessionSidebar';
import { notify } from './services/notifications';
import classes from './AgentApp.module.css';

type ThemeOption = 'auto' | 'light' | 'dark';

function AgentApp() {
  const [sessions, setSessions] = useState<ChatSessionItem[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessageItem[]>([]);
  const [sending, setSending] = useState(false);
  const [colorScheme, setColorScheme] = useState<ThemeOption>('auto');

  const activeSession = useMemo(
    () => sessions.find((session) => session.id === activeSessionId) ?? null,
    [sessions, activeSessionId],
  );

  const loadSessions = useCallback(async () => {
    try {
      const items = await AppService.ListChatSessions();
      setSessions(items as ChatSessionItem[]);
    } catch (err) {
      notify.failed('Failed to load chat sessions', String(err));
    }
  }, []);

  const loadMessages = useCallback(async (sessionId: string) => {
    try {
      const items = await AppService.ListChatMessages(sessionId);
      setMessages(items as ChatMessageItem[]);
    } catch (err) {
      notify.failed('Failed to load messages', String(err));
    }
  }, []);

  useEffect(() => {
    AppService.GetConfig()
      .then((cfg) => {
        const theme = (cfg.theme as ThemeOption) || 'auto';
        setColorScheme(theme);
      })
      .catch((err: unknown) => {
        console.error(err);
      });
  }, []);

  useEffect(() => {
    void loadSessions();
  }, [loadSessions]);

  useEffect(() => {
    if (activeSessionId) {
      void loadMessages(activeSessionId);
    } else {
      setMessages([]);
    }
  }, [activeSessionId, loadMessages]);

  const handleNewChat = async () => {
    try {
      const session = (await AppService.CreateChatSession('')) as ChatSessionItem;
      await loadSessions();
      setActiveSessionId(session.id);
      setMessages([]);
    } catch (err) {
      notify.failed('Failed to create chat', String(err));
    }
  };

  const handleRenameSession = async (sessionId: string, title: string) => {
    try {
      await AppService.RenameChatSession(sessionId, title);
      await loadSessions();
    } catch (err) {
      notify.failed('Failed to rename chat', String(err));
    }
  };

  const handleDeleteSession = async (sessionId: string) => {
    try {
      await AppService.DeleteChatSession(sessionId);
      if (activeSessionId === sessionId) {
        setActiveSessionId(null);
        setMessages([]);
      }
      await loadSessions();
    } catch (err) {
      notify.failed('Failed to delete chat', String(err));
    }
  };

  const handleSend = async (content: string) => {
    if (!activeSessionId) {
      return;
    }

    setSending(true);
    try {
      const result = await AppService.SendChatMessage(activeSessionId, content);
      setMessages((current) => [
        ...current,
        result.user_message as ChatMessageItem,
        result.assistant_message as ChatMessageItem,
      ]);
      await loadSessions();
    } catch (err) {
      notify.failed('Failed to send message', String(err));
    } finally {
      setSending(false);
    }
  };

  const handleOpenEditor = () => {
    AppService.ShowEditorWindow().catch((err: unknown) => {
      notify.failed('Failed to open editor', String(err));
    });
  };

  return (
    <MantineProvider
      defaultColorScheme="auto"
      forceColorScheme={colorScheme === 'auto' ? undefined : colorScheme}
    >
      <Notifications position="top-right" limit={5} />
      <Box className={classes.app}>
        <SessionSidebar
          sessions={sessions}
          activeSessionId={activeSessionId}
          onSelectSession={setActiveSessionId}
          onNewChat={handleNewChat}
          onRenameSession={handleRenameSession}
          onDeleteSession={handleDeleteSession}
          onOpenEditor={handleOpenEditor}
        />

        <Box className={classes.main}>
          <ChatThread
            messages={messages}
            sessionTitle={activeSession?.title ?? null}
            hasActiveSession={activeSessionId !== null}
          />
          <ChatInput
            disabled={activeSessionId === null}
            sending={sending}
            onSend={handleSend}
          />
        </Box>
      </Box>
    </MantineProvider>
  );
}

export default AgentApp;
