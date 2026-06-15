import { useEffect, useRef } from 'react';
import { Events } from '@wailsio/runtime';

export interface AgentStreamChunk {
  type: string;
  text?: string;
  finishReason?: string;
  error?: string;
}

export interface AgentStreamChunkEvent {
  sessionId: string;
  messageId: string;
  chunk: AgentStreamChunk;
}

export interface AgentTurnDoneEvent {
  sessionId: string;
  messageId: string;
  content: string;
}

export interface AgentTurnErrorEvent {
  sessionId: string;
  messageId: string;
  error: string;
}

export interface AgentTurnCancelledEvent {
  sessionId: string;
  messageId: string;
  content: string;
}

export interface AgentStreamHandlers {
  onChunk: (event: AgentStreamChunkEvent) => void;
  onDone: (event: AgentTurnDoneEvent) => void;
  onError: (event: AgentTurnErrorEvent) => void;
  onCancelled: (event: AgentTurnCancelledEvent) => void;
}

function readEventData<T>(event: unknown): T {
  if (event && typeof event === 'object' && 'data' in event) {
    return (event as { data: T }).data;
  }
  return event as T;
}

export function useAgentStream(handlers: AgentStreamHandlers) {
  const handlersRef = useRef(handlers);
  handlersRef.current = handlers;

  useEffect(() => {
    const unsubs = [
      Events.On('agent:stream-chunk', (event) => {
        handlersRef.current.onChunk(readEventData<AgentStreamChunkEvent>(event));
      }),
      Events.On('agent:turn-done', (event) => {
        handlersRef.current.onDone(readEventData<AgentTurnDoneEvent>(event));
      }),
      Events.On('agent:turn-error', (event) => {
        handlersRef.current.onError(readEventData<AgentTurnErrorEvent>(event));
      }),
      Events.On('agent:turn-cancelled', (event) => {
        handlersRef.current.onCancelled(readEventData<AgentTurnCancelledEvent>(event));
      }),
    ];

    return () => {
      unsubs.forEach((unsub) => unsub());
    };
  }, []);
}
