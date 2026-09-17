import type { ChatMessageItem, ChatTransport, OverlayStreamHandlers } from '../types';

function splitIntoChunks(text: string): string[] {
  return text.split(/(\s+)/).filter(Boolean);
}

function buildMockReply(content: string, domainContext: string): string {
  const domainLabel = domainContext.includes('Cloudflare') ? 'Cloudflare' : 'Talus Agent';
  return `**${domainLabel} demo reply** — this stream is local, with no API call.

You asked: ${content}

I can see overlay domain context (${domainContext.split('\n')[0] ?? 'none'}). Switch domains on the demo page to change title, suggestions, and this reply.`;
}

export function createMockChatTransport(): ChatTransport {
  const handlersByConversation = new Map<string, OverlayStreamHandlers>();
  let timer: number | null = null;
  let activeMessageId = '';
  let activeConversationId = '';
  let partial = '';

  const clearTimer = () => {
    if (timer !== null) {
      window.clearInterval(timer);
      timer = null;
    }
  };

  return {
    subscribe(conversationId, handlers) {
      handlersByConversation.set(conversationId, handlers);
      return () => {
        handlersByConversation.delete(conversationId);
      };
    },
    async send(input) {
      clearTimer();
      partial = '';
      activeConversationId = input.conversationId;
      const now = new Date().toISOString();
      const user: ChatMessageItem = {
        id: crypto.randomUUID(),
        session_id: input.conversationId,
        role: 'user',
        content: input.content,
        created_at: now,
      };
      const assistant: ChatMessageItem = {
        id: crypto.randomUUID(),
        session_id: input.conversationId,
        role: 'assistant',
        content: '',
        created_at: now,
        generating: true,
      };
      activeMessageId = assistant.id;
      const chunks = splitIntoChunks(buildMockReply(input.content, input.domainContext));
      let index = 0;

      timer = window.setInterval(() => {
        const handlers = handlersByConversation.get(input.conversationId);
        if (index >= chunks.length) {
          clearTimer();
          handlers?.onDone({
            sessionId: input.conversationId,
            messageId: assistant.id,
            content: partial,
          });
          return;
        }
        const chunk = chunks[index];
        index += 1;
        partial += chunk;
        handlers?.onChunk({
          sessionId: input.conversationId,
          messageId: assistant.id,
          chunk: { type: 'text-delta', text: chunk },
        });
      }, 80);

      return { user, assistant };
    },
    async cancel(messageId) {
      if (messageId !== activeMessageId) {
        return;
      }
      clearTimer();
      handlersByConversation.get(activeConversationId)?.onCancelled({
        sessionId: activeConversationId,
        messageId,
        content: partial,
      });
    },
  };
}
