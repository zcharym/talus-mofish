import { useCallback, useEffect, useRef, useState } from 'react';
import { ChatService, toApiError } from '../utils/api';
import type { ChatMessageItem } from '../components/agent/ChatThread';
import { notify } from '../services/notifications';
import { useAgentStream } from './useAgentStream';

export function useAgentChat(activeSessionId: string | null, sessionKind: string | undefined, onSessionsChanged: () => Promise<void>) {
  const [messages, setMessages] = useState<ChatMessageItem[]>([]);
  const [sending, setSending] = useState(false);
  const [streamingMessageId, setStreamingMessageId] = useState<string | null>(null);
  const activeSessionIdRef = useRef<string | null>(activeSessionId);

  useEffect(() => {
    activeSessionIdRef.current = activeSessionId;
  }, [activeSessionId]);

  const loadMessages = useCallback(async (sessionId: string) => {
    try {
      const items = await ChatService.ListChatMessages(sessionId);
      setMessages((items ?? []) as ChatMessageItem[]);
    } catch (err) {
      notify.failed('Failed to load messages', toApiError(err));
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
      void onSessionsChanged();
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
    if (!activeSessionId || sessionKind === 'sudoku') {
      setMessages([]);
      return;
    }
    void loadMessages(activeSessionId);
  }, [activeSessionId, sessionKind, loadMessages]);

  const clearMessages = useCallback(() => {
    setMessages([]);
  }, []);

  const send = useCallback(
    async (content: string, sessionId: string) => {
      setSending(true);
      try {
        const result = await ChatService.StartChatTurn(sessionId, content);
        const userMessage = result.user_message as ChatMessageItem;
        const assistantMessage = {
          ...(result.assistant_message as ChatMessageItem),
          generating: true,
        };
        setMessages((current) => [...current, userMessage, assistantMessage]);
        setStreamingMessageId(assistantMessage.id);
        await onSessionsChanged();
      } catch (err) {
        notify.failed('Failed to send message', toApiError(err));
        setSending(false);
        setStreamingMessageId(null);
      }
    },
    [onSessionsChanged],
  );

  const cancel = useCallback(async () => {
    if (!activeSessionId || !streamingMessageId) {
      return;
    }
    try {
      await ChatService.CancelChatTurn(activeSessionId, streamingMessageId);
    } catch (err) {
      notify.failed('Failed to cancel response', toApiError(err));
    }
  }, [activeSessionId, streamingMessageId]);

  return {
    messages,
    sending,
    setSending,
    streamingMessageId,
    send,
    cancel,
    clearMessages,
  };
}
