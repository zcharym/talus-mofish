import { useCallback, useEffect, useMemo, useState } from 'react';
import { ChatService, SudokuService, toApiError } from '../utils/api';
import type { ChatSessionItem } from '../components/agent/SessionSidebar';
import { notify } from '../services/notifications';

export type AgentView = 'home' | 'cloudflare' | 'session';

export function useAgentSessions() {
  const [sessions, setSessions] = useState<ChatSessionItem[]>([]);
  const [activeSessionId, setActiveSessionId] = useState<string | null>(null);
  const [view, setView] = useState<AgentView>('home');

  const activeSession = useMemo(
    () => sessions.find((session) => session.id === activeSessionId) ?? null,
    [sessions, activeSessionId],
  );

  const loadSessions = useCallback(async () => {
    try {
      const items = await ChatService.ListChatSessions();
      const next = (items ?? []) as ChatSessionItem[];
      setSessions(next);
    } catch (err) {
      notify.failed('Failed to load chat sessions', toApiError(err));
    }
  }, []);

  useEffect(() => {
    void loadSessions();
  }, [loadSessions]);

  const goHome = useCallback(() => {
    setView('home');
    setActiveSessionId(null);
  }, []);

  const selectCloudflare = useCallback(() => {
    setView('cloudflare');
    setActiveSessionId(null);
  }, []);

  const selectSession = useCallback((sessionId: string) => {
    setView('session');
    setActiveSessionId(sessionId);
  }, []);

  const renameSession = useCallback(
    async (sessionId: string, title: string) => {
      try {
        await ChatService.RenameChatSession(sessionId, title);
        await loadSessions();
      } catch (err) {
        notify.failed('Failed to rename chat', toApiError(err));
      }
    },
    [loadSessions],
  );

  const deleteSession = useCallback(
    async (sessionId: string) => {
      try {
        await ChatService.DeleteChatSession(sessionId);
        setActiveSessionId((current) => {
          if (current === sessionId) {
            setView('home');
            return null;
          }
          return current;
        });
        await loadSessions();
      } catch (err) {
        notify.failed('Failed to delete chat', toApiError(err));
      }
    },
    [loadSessions],
  );

  const startSudoku = useCallback(async () => {
    const result = await SudokuService.NewSudokuGame('easy');
    await loadSessions();
    setView('session');
    setActiveSessionId(result.session.id);
    return result.session.id;
  }, [loadSessions]);

  const createChatSession = useCallback(async () => {
    const session = (await ChatService.CreateChatSession('')) as ChatSessionItem;
    setView('session');
    setActiveSessionId(session.id);
    await loadSessions();
    return session.id;
  }, [loadSessions]);

  return {
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
  };
}
