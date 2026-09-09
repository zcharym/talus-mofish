import { useCallback, useEffect, useState } from 'react';
import { Box, Drawer } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { SystemService, toApiError } from './utils/api';
import { AgentHome, QuickActionId } from './components/agent/AgentHome';
import { ChatInput } from './components/agent/ChatInput';
import { ChatThread } from './components/agent/ChatThread';
import { CloudflareDashboard, canShowCloudflareTab } from './components/agent/CloudflareDashboard';
import { SessionSidebar } from './components/agent/SessionSidebar';
import { SudokuBoard } from './components/agent/SudokuBoard';
import { useAgentChat } from './hooks/useAgentChat';
import { useAgentSessions } from './hooks/useAgentSessions';
import { useConfigFlags } from './hooks/useConfigFlags';
import { useCurrentUser } from './hooks/useCurrentUser';
import { isMacPlatform, usePlatform } from './hooks/usePlatform';
import { notify } from './services/notifications';
import { ThemeRoot } from './theme';
import classes from './AgentApp.module.css';

function AgentApp() {
  const {
    sessions,
    activeSessionId,
    activeSession,
    view,
    loadSessions,
    goHome,
    selectCloudflare,
    selectSession,
    renameSession,
    deleteSession,
    startSudoku,
    createChatSession,
  } = useAgentSessions();
  const {
    messages,
    sending,
    setSending,
    streamingMessageId,
    send,
    cancel,
    clearMessages,
  } = useAgentChat(activeSessionId, activeSession?.kind, loadSessions);
  const { theme: colorScheme, cloudflareConfigured } = useConfigFlags();
  const platform = usePlatform();
  const [sidebarOpened, setSidebarOpened] = useState(false);
  const isOverlay = useMediaQuery('(max-width: 56.24em)', false, {
    getInitialValueInEffect: false,
  });
  const { user, loading: userLoading, signingIn, signInWithEmail, signIn, signOut } = useCurrentUser();
  const showCloudflareTab = canShowCloudflareTab(user, cloudflareConfigured);

  useEffect(() => {
    if (view === 'cloudflare' && !showCloudflareTab) {
      goHome();
    }
  }, [view, showCloudflareTab, goHome]);

  const closeSidebar = useCallback(() => {
    setSidebarOpened(false);
  }, [setSidebarOpened]);

  const handleGoHome = () => {
    goHome();
    clearMessages();
    closeSidebar();
  };

  const handleSelectCloudflare = () => {
    selectCloudflare();
    clearMessages();
    closeSidebar();
  };

  const handleSelectSession = (sessionId: string) => {
    selectSession(sessionId);
    closeSidebar();
  };

  const ensureActiveSession = useCallback(async (): Promise<string> => {
    if (activeSessionId && view === 'session') {
      return activeSessionId;
    }
    clearMessages();
    return createChatSession();
  }, [activeSessionId, view, createChatSession, clearMessages]);

  const handleSend = async (content: string) => {
    const sessionId = await ensureActiveSession();
    await send(content, sessionId);
  };

  const handleQuickAction = async (actionId: QuickActionId, prompt: string, autoSend: boolean) => {
    if (actionId === 'play_sudoku') {
      setSending(true);
      try {
        await startSudoku();
        clearMessages();
      } catch (err) {
        notify.failed('Could not start Sudoku', toApiError(err));
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
      notify.failed('Sign-in failed', toApiError(err));
    }
  };

  const handleSignIn = async (provider: 'github' | 'google') => {
    try {
      await signIn(provider);
    } catch (err) {
      notify.failed('Sign-in failed', toApiError(err));
    }
  };

  const handleSignOut = async () => {
    try {
      await signOut();
      goHome();
      clearMessages();
    } catch (err) {
      notify.failed('Sign-out failed', toApiError(err));
    }
  };

  const handleOpenManagement = () => {
    SystemService.ShowManagementWindow().catch((err: unknown) => {
      notify.failed('Failed to open management', toApiError(err));
    });
  };

  const isHomeView = view === 'home';
  const isCloudflareView = view === 'cloudflare';
  const isSudokuView = view === 'session' && activeSession?.kind === 'sudoku';

  const sidebar = (
    <SessionSidebar
      sessions={sessions}
      activeSessionId={isCloudflareView ? null : activeSessionId}
      user={user}
      showCloudflareTab={showCloudflareTab}
      cloudflareActive={isCloudflareView}
      overlay={Boolean(isOverlay)}
      onSelectSession={handleSelectSession}
      onSelectCloudflare={handleSelectCloudflare}
      onNewChat={handleGoHome}
      onRenameSession={renameSession}
      onDeleteSession={async (sessionId) => {
        await deleteSession(sessionId);
        if (activeSessionId === sessionId) {
          clearMessages();
        }
      }}
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
          data-platform={isMacPlatform(platform) ? 'darwin' : undefined}
        >
          {isHomeView ? (
            <AgentHome
              user={user}
              userLoading={userLoading}
              signingIn={signingIn}
              sending={sending}
              onOpenSidebar={openSidebar}
              onSend={handleSend}
              onCancel={streamingMessageId ? cancel : undefined}
              onSignInWithEmail={handleSignInWithEmail}
              onSignIn={handleSignIn}
              onQuickAction={(actionId, prompt, autoSend) => {
                void handleQuickAction(actionId, prompt, autoSend);
              }}
            />
          ) : isCloudflareView ? (
            <CloudflareDashboard
              configured={cloudflareConfigured}
              onOpenSidebar={openSidebar}
              onOpenManagement={handleOpenManagement}
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
                onCancel={streamingMessageId ? cancel : undefined}
              />
            </>
          )}
        </Box>
      </Box>
    </ThemeRoot>
  );
}

export default AgentApp;
