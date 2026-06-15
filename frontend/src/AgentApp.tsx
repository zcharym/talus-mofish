import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import { Box, MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import { AppService } from '../bindings/github.com/songwei.ma/talus-mofish';
import { ChatInput } from './components/agent/ChatInput';
import { ChatMessageItem, ChatThread } from './components/agent/ChatThread';
import { ChatSessionItem, SessionSidebar } from './components/agent/SessionSidebar';
import { useAgentStream } from './hooks/useAgentStream';
import { notify } from './services/notifications';
import classes from './AgentApp.module.css';

type ThemeOption = 'auto' | 'light' | 'dark';

function AgentApp() {
  const [sessions, setSessions] = useState<ChatSessionItem[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessageItem[]>([]);
  const [sending, setSending] = useState(false);
  const [streamingMessageId, setStreamingMessageId] = useState<string | null>(null);
  const [colorScheme, setColorScheme] = useState<ThemeOption>('auto');
  const activeSessionIdRef = useRef<string | null>(null);

  useEffect(() => {
    activeSessionIdRef.current = activeSessionId;
  }, [activeSessionId]);

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

  const updateAssistantMessage = useCallback((messageId: string, patch: Partial<ChatMessageItem>) => {
    setMessages((current) =>
      current.map((message) => (message.id === messageId ? { ...message, ...patch } : message)),
    );
  }, []);

  useAgentStream({
    onChunk: ({ sessionId, messageId, chunk }) => {
      if (sessionId !== activeSessionIdRef.current || chunk.type !== 'text-delta' || !chunk.text) {
        return;
      }
      setMessages((current) =>
        current.map((message) =>
          message.id === messageId
            ? {
                ...message,
                content: message.content + chunk.text,
                generating: true,
              }
            : message,
        ),
      );
    },
    onDone: ({ sessionId, messageId, content }) => {
      if (sessionId !== activeSessionIdRef.current) {
        return;
      }
      updateAssistantMessage(messageId, { content, generating: false });
      setSending(false);
      setStreamingMessageId(null);
      void loadSessions();
    },
    onError: ({ sessionId, messageId, error }) => {
      if (sessionId !== activeSessionIdRef.current) {
        return;
      }
      updateAssistantMessage(messageId, {
        content: error,
        generating: false,
      });
      setSending(false);
      setStreamingMessageId(null);
      notify.failed('Agent error', error);
    },
    onCancelled: ({ sessionId, messageId, content }) => {
      if (sessionId !== activeSessionIdRef.current) {
        return;
      }
      updateAssistantMessage(messageId, { content, generating: false });
      setSending(false);
      setStreamingMessageId(null);
    },
  });

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
      const result = await AppService.StartChatTurn(activeSessionId, content);
      const userMessage = result.user_message as ChatMessageItem;
      const assistantMessage = {
        ...(result.assistant_message as ChatMessageItem),
        generating: true,
      };
      setMessages((current) => [...current, userMessage, assistantMessage]);
      setStreamingMessageId(assistantMessage.id);
      await loadSessions();
    } catch (err) {
      notify.failed('Failed to send message', String(err));
      setSending(false);
      setStreamingMessageId(null);
    }
  };

  const handleCancel = async () => {
    if (!activeSessionId || !streamingMessageId) {
      return;
    }
    try {
      await AppService.CancelChatTurn(activeSessionId, streamingMessageId);
    } catch (err) {
      notify.failed('Failed to cancel response', String(err));
    }
  };

  const handleOpenManagement = () => {
    AppService.ShowManagementWindow().catch((err: unknown) => {
      notify.failed('Failed to open management', String(err));
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
          onOpenManagement={handleOpenManagement}
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
            onCancel={streamingMessageId ? handleCancel : undefined}
          />
        </Box>
      </Box>
    </MantineProvider>
  );
}

export default AgentApp;
