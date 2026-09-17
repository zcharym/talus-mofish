import type { ChatMessageItem } from '../ChatThread';
import type {
  AgentStreamChunkEvent,
  AgentTurnCancelledEvent,
  AgentTurnDoneEvent,
  AgentTurnErrorEvent,
} from '../../../utils/api';

export type { ChatMessageItem };

export interface ChatModalSuggestion {
  id: string;
  label: string;
  prompt: string;
}

export interface ChatModalDomain {
  id: string;
  title: string;
  placeholder: string;
  emptyHint: string;
  suggestions: ChatModalSuggestion[];
  buildContext: () => string;
}

export interface OverlayStreamHandlers {
  onChunk: (event: AgentStreamChunkEvent) => void;
  onDone: (event: AgentTurnDoneEvent) => void;
  onError: (event: AgentTurnErrorEvent) => void;
  onCancelled: (event: AgentTurnCancelledEvent) => void;
}

export interface ChatTransportSendInput {
  conversationId: string;
  content: string;
  history: ChatMessageItem[];
  domainContext: string;
}

export interface ChatTransport {
  send: (input: ChatTransportSendInput) => Promise<{ user: ChatMessageItem; assistant: ChatMessageItem }>;
  cancel: (messageId: string) => Promise<void>;
  subscribe?: (conversationId: string, handlers: OverlayStreamHandlers) => () => void;
}
