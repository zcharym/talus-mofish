import { useEffect, useRef } from 'react';
import {
  onAgentStreamChunk,
  onAgentTurnCancelled,
  onAgentTurnDone,
  onAgentTurnError,
  type AgentStreamChunkEvent,
  type AgentTurnCancelledEvent,
  type AgentTurnDoneEvent,
  type AgentTurnErrorEvent,
} from '../utils/api';

export type {
  AgentStreamChunkEvent,
  AgentTurnCancelledEvent,
  AgentTurnDoneEvent,
  AgentTurnErrorEvent,
} from '../utils/api';

export interface AgentStreamHandlers {
  onChunk: (event: AgentStreamChunkEvent) => void;
  onDone: (event: AgentTurnDoneEvent) => void;
  onError: (event: AgentTurnErrorEvent) => void;
  onCancelled: (event: AgentTurnCancelledEvent) => void;
}

export function useAgentStream(handlers: AgentStreamHandlers) {
  const handlersRef = useRef(handlers);
  handlersRef.current = handlers;

  useEffect(() => {
    const unsubs = [
      onAgentStreamChunk((event) => handlersRef.current.onChunk(event)),
      onAgentTurnDone((event) => handlersRef.current.onDone(event)),
      onAgentTurnError((event) => handlersRef.current.onError(event)),
      onAgentTurnCancelled((event) => handlersRef.current.onCancelled(event)),
    ];
    return () => {
      unsubs.forEach((unsub) => unsub());
    };
  }, []);
}
