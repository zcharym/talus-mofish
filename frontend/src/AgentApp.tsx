import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { Box, Drawer } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { ChatService, SudokuService, SystemService } from '../bindings/github.com/songwei.ma/talus-mofish/backend/services';
import { AgentHome, QuickActionId } from './components/agent/AgentHome';
import { ChatInput } from './components/agent/ChatInput';
import { ChatMessageItem, ChatThread } from './components/agent/ChatThread';
import { ChatSessionItem, SessionSidebar } from './components/agent/SessionSidebar';
import { SudokuBoard } from './components/agent/SudokuBoard';
import { useAgentStream } from './hooks/useAgentStream';
import { useConfigColorScheme } from './hooks/useConfigColorScheme';
import { useCurrentUser } from './hooks/useCurrentUser';
import { notify } from './services/notifications';
import { ThemeRoot } from './theme';
import { isMacOS } from './utils/platform';
import classes from './AgentApp.module.css';

function AgentApp() {
  const [sessions, setSessions] = useState<ChatSessionItem[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [messages, setMessages] = useState<ChatMessageItem[]>([]);
  const [sending, setSending] = useState(false);
  const [streamingMessageId, setStreamingMessageId] = useState<string | null>(null);
  const { colorScheme } = useConfigColorScheme();
  const [sidebarOpened, setSidebarOpened] = useState(false);
  const isOverlay = useMediaQuery('(max-width: 56.24em)', false, {
    getInitialValueInEffect: false,
  });
  const activeSessionIdRef = useRef<string | null>(null);
  const { user, loading: userLoading, signingIn, signInWithEmail, signIn, signOut } = useCurrentUser();

  useEffect(() => {
    activeSessionIdRef.current = activeSessionId;
  }, [activeSessionId]);

  const activeSession = useMemo(
    () => sessions.find((session) => session.id === activeSessionId) ?? null,
    [sessions, activeSessionId],
  );

  const loadSessions = useCallback(async () => {
    try {
      const items = await ChatService.ListChatSessions();
      setSessions(items as ChatSessionItem[]);
    } catch (err) {
      notify.failed('Failed to load chat sessions', String(err));
    }
  }, []);

  const loadMessages = useCallback(async (sessionId: string) => {
    try {
      const items = await ChatService.ListChatMessages(sessionId);
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
    void loadSessions();
  }, [loadSessions]);

  useEffect(() => {
    if (!activeSessionId) {
      setMessages([]);
      return;
    }
    if (activeSession?.kind === 'sudoku') {
      setMessages([]);
      return;
    }
    void loadMessages(activeSessionId);
  }, [activeSessionId, activeSession?.kind, loadMessages]);

  const ensureActiveSession = useCallback(async (): Promise<string> => {
    if (activeSessionId) {
      return activeSessionId;
    }

    const session = (await ChatService.CreateChatSession('')) as ChatSessionItem;
    setActiveSessionId(session.id);
    setMessages([]);
    await loadSessions();
    return session.id;
  }, [activeSessionId, loadSessions]);

  const closeSidebar = useCallback(() => {
    setSidebarOpened(false);
  }, []);

  const handleGoHome = () => {
    setActiveSessionId(null);
    setMessages([]);
    closeSidebar();
  };

  const handleSelectSession = (sessionId: string) => {
    setActiveSessionId(sessionId);
    closeSidebar();
  };

  const handleRenameSession = async (sessionId: string, title: string) => {
    try {
      await ChatService.RenameChatSession(sessionId, title);
      await loadSessions();
    } catch (err) {
      notify.failed('Failed to rename chat', String(err));
    }
  };

  const handleDeleteSession = async (sessionId: string) => {
    try {
      await ChatService.DeleteChatSession(sessionId);
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
    setSending(true);
    try {
      const sessionId = await ensureActiveSession();
      const result = await ChatService.StartChatTurn(sessionId, content);
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

  const handleQuickAction = async (actionId: QuickActionId, prompt: string, autoSend: boolean) => {
    if (actionId === 'play_sudoku') {
      setSending(true);
      try {
        const result = await SudokuService.NewSudokuGame('easy');
        await loadSessions();
        setActiveSessionId(result.session.id);
        setMessages([]);
      } catch (err) {
        notify.failed('Could not start Sudoku', String(err));
      } finally {
        setSending(false);
      }
      return;
    }
    if (!autoSend) {
      return;
    }
    await handleSend(prompt);
  };

  const handleSignInWithEmail = async (email: string) => {
    try {
      await signInWithEmail(email);
    } catch (err) {
      notify.failed('Sign-in failed', String(err));
    }
  };

  const handleSignIn = async (provider: 'github' | 'google') => {
    try {
      await signIn(provider);
    } catch (err) {
      notify.failed('Sign-in failed', String(err));
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
      setActiveSessionId(null);
      setMessages([]);
    } catch (err) {
      notify.failed('Sign-out failed', String(err));
    }
  };

  const handleCancel = async () => {
    if (!activeSessionId || !streamingMessageId) {
      return;
    }
    try {
      await ChatService.CancelChatTurn(activeSessionId, streamingMessageId);
    } catch (err) {
      notify.failed('Failed to cancel response', String(err));
    }
  };

  const handleOpenManagement = () => {
    SystemService.ShowManagementWindow().catch((err: unknown) => {
      notify.failed('Failed to open management', String(err));
    });
  };

  const isHomeView = activeSessionId === null;
  const isSudokuView = activeSession?.kind === 'sudoku';

  const sidebar = (
    <SessionSidebar
      sessions={sessions}
      activeSessionId={activeSessionId}
      user={user}
      overlay={Boolean(isOverlay)}
      onSelectSession={handleSelectSession}
      onNewChat={handleGoHome}
      onRenameSession={handleRenameSession}
      onDeleteSession={handleDeleteSession}
      onOpenManagement={handleOpenManagement}
      onSignOut={handleSignOut}
    />
  );

  const openSidebar = isOverlay ? () => setSidebarOpened(true) : undefined;

  return (
    <ThemeRoot colorScheme={colorScheme}>
      <Box className={classes.app}>
        {isOverlay ? (
          <Drawer
            opened={sidebarOpened}
            onClose={closeSidebar}
            padding={0}
            size={260}
            withCloseButton={false}
            styles={{
              content: { height: '100%' },
              body: { height: '100%', padding: 0, display: 'flex' },
            }}
          >
            {sidebar}
          </Drawer>
        ) : (
          sidebar
        )}

        <Box
          className={classes.main}
          data-overlay={isOverlay || undefined}
          data-platform={isMacOS ? 'darwin' : undefined}
        >
          {isHomeView ? (
            <AgentHome
              user={user}
              userLoading={userLoading}
              signingIn={signingIn}
              sending={sending}
              onOpenSidebar={openSidebar}
              onSend={handleSend}
              onCancel={streamingMessageId ? handleCancel : undefined}
              onSignInWithEmail={handleSignInWithEmail}
              onSignIn={handleSignIn}
              onQuickAction={(actionId, prompt, autoSend) => {
                void handleQuickAction(actionId, prompt, autoSend);
              }}
            />
          ) : isSudokuView && activeSessionId ? (
            <SudokuBoard
              sessionId={activeSessionId}
              sessionTitle={activeSession?.title ?? null}
              onSessionUpdated={loadSessions}
              onOpenSidebar={openSidebar}
            />
          ) : (
            <>
              <ChatThread
                messages={messages}
                sessionTitle={activeSession?.title ?? null}
                hasActiveSession
                onOpenSidebar={openSidebar}
              />
              <ChatInput
                disabled={false}
                sending={sending}
                onSend={handleSend}
                onCancel={streamingMessageId ? handleCancel : undefined}
              />
            </>
          )}
        </Box>
      </Box>
    </ThemeRoot>
  );
}

export default AgentApp;
