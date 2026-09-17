import {
  ChatMessage,
  ChatService,
  OverlayChatTurnRequest,
  onAgentStreamChunk,
  onAgentTurnCancelled,
  onAgentTurnDone,
  onAgentTurnError,
} from '../../../../utils/api';
import type { ChatMessageItem, ChatTransport } from '../types';

function asChatItem(message: ChatMessage, generating = false): ChatMessageItem {
  return {
    id: message.id,
    session_id: message.session_id,
    role: message.role as ChatMessageItem['role'],
    content: message.content,
    created_at: message.created_at,
    generating,
  };
}

export function createAgentOverlayTransport(): ChatTransport {
  let conversationId = '';

  return {
    subscribe(id, handlers) {
      const unsubs = [
        onAgentStreamChunk((event) => {
          if (event.sessionId !== id) {
            return;
          }
          handlers.onChunk(event);
        }),
        onAgentTurnDone((event) => {
          if (event.sessionId !== id) {
            return;
          }
          handlers.onDone(event);
        }),
        onAgentTurnError((event) => {
          if (event.sessionId !== id) {
            return;
          }
          handlers.onError(event);
        }),
        onAgentTurnCancelled((event) => {
          if (event.sessionId !== id) {
            return;
          }
          handlers.onCancelled(event);
        }),
      ];
      return () => {
        unsubs.forEach((unsub) => unsub());
      };
    },
    async send(input) {
      conversationId = input.conversationId;
      const result = await ChatService.StartOverlayChatTurn(
        OverlayChatTurnRequest.createFrom({
          conversation_id: input.conversationId,
          content: input.content,
          domain_context: input.domainContext,
          history: input.history.map((message) =>
            ChatMessage.createFrom({
              id: message.id,
              session_id: message.session_id,
              role: message.role,
              content: message.content,
              created_at: message.created_at,
            }),
          ),
        }),
      );
      return {
        user: asChatItem(result.user_message),
        assistant: asChatItem(result.assistant_message, true),
      };
    },
    async cancel(messageId) {
      await ChatService.CancelChatTurn(conversationId || 'overlay', messageId);
    },
  };
}
