import { useCallback, useEffect, useRef, useState } from 'react';
import { toApiError } from '../../../utils/api';
import { notify } from '../../../services/notifications';
import type { ChatMessageItem, ChatModalDomain, ChatTransport } from './types';

export function useChatModal({
  domain,
  transport,
}: {
  domain: ChatModalDomain;
  transport: ChatTransport;
}) {
  const [conversationId, setConversationId] = useState(() => crypto.randomUUID());
  const [messages, setMessages] = useState<ChatMessageItem[]>([]);
  const [sending, setSending] = useState(false);
  const [streamingMessageId, setStreamingMessageId] = useState<string | null>(null);
  const conversationIdRef = useRef(conversationId);
  const domainRef = useRef(domain);
  conversationIdRef.current = conversationId;
  domainRef.current = domain;

  const updateAssistantMessage = useCallback((messageId: string, patch: Partial<ChatMessageItem>) => {
    setMessages((current) =>
      current.map((message) => (message.id === messageId ? { ...message, ...patch } : message)),
    );
  }, []);

  const prevDomainIdRef = useRef(domain.id);

  useEffect(() => {
    if (prevDomainIdRef.current === domain.id) {
      return;
    }
    prevDomainIdRef.current = domain.id;
    setConversationId(crypto.randomUUID());
    setMessages([]);
    setSending(false);
    setStreamingMessageId(null);
  }, [domain.id]);

  useEffect(() => {
    if (!transport.subscribe) {
      return undefined;
    }

    return transport.subscribe(conversationId, {
      onChunk: ({ sessionId, messageId, chunk }) => {
        if (sessionId !== conversationIdRef.current || chunk.type !== 'text-delta' || !chunk.text) {
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
        if (sessionId !== conversationIdRef.current) {
          return;
        }
        updateAssistantMessage(messageId, { content, generating: false });
        setSending(false);
        setStreamingMessageId(null);
      },
      onError: ({ sessionId, messageId, error }) => {
        if (sessionId !== conversationIdRef.current) {
          return;
        }
        updateAssistantMessage(messageId, { content: error, generating: false });
        setSending(false);
        setStreamingMessageId(null);
        notify.failed('Agent error', error);
      },
      onCancelled: ({ sessionId, messageId, content }) => {
        if (sessionId !== conversationIdRef.current) {
          return;
        }
        updateAssistantMessage(messageId, { content, generating: false });
        setSending(false);
        setStreamingMessageId(null);
      },
    });
  }, [conversationId, transport, updateAssistantMessage]);

  const send = useCallback(
    async (content: string) => {
      const trimmed = content.trim();
      if (!trimmed || sending) {
        return;
      }
      setSending(true);
      try {
        const result = await transport.send({
          conversationId: conversationIdRef.current,
          content: trimmed,
          history: messages.filter((message) => message.content.trim() !== ''),
          domainContext: domainRef.current.buildContext(),
        });
        setMessages((current) => [
          ...current,
          result.user,
          { ...result.assistant, generating: true },
        ]);
        setStreamingMessageId(result.assistant.id);
      } catch (err) {
        notify.failed('Failed to send message', toApiError(err));
        setSending(false);
        setStreamingMessageId(null);
      }
    },
    [messages, sending, transport],
  );

  const cancel = useCallback(async () => {
    if (!streamingMessageId) {
      return;
    }
    try {
      await transport.cancel(streamingMessageId);
    } catch (err) {
      notify.failed('Failed to cancel response', toApiError(err));
    }
  }, [streamingMessageId, transport]);

  return {
    messages,
    sending,
    streamingMessageId,
    send,
    cancel,
  };
}
