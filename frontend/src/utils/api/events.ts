import { Events } from '@wailsio/runtime';
import type {
  StreamChunkEvent,
  TurnCancelledEvent,
  TurnDoneEvent,
  TurnErrorEvent,
} from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/agent/models';
import type { ConfigChanged } from '../../../bindings/github.com/songwei.ma/talus-mofish/backend/types/models';

export const ConfigChangedEvent = 'config:changed';
export const AgentStreamChunkEventName = 'agent:stream-chunk';
export const AgentTurnDoneEventName = 'agent:turn-done';
export const AgentTurnErrorEventName = 'agent:turn-error';
export const AgentTurnCancelledEventName = 'agent:turn-cancelled';

export type ConfigChangedPayload = ConfigChanged;
export type AgentStreamChunkEvent = StreamChunkEvent;
export type AgentTurnDoneEvent = TurnDoneEvent;
export type AgentTurnErrorEvent = TurnErrorEvent;
export type AgentTurnCancelledEvent = TurnCancelledEvent;

export function readEventData<T>(event: unknown): T {
  if (event && typeof event === 'object' && 'data' in event) {
    return (event as { data: T }).data;
  }
  return event as T;
}

export function onConfigChanged(handler: (payload: ConfigChangedPayload) => void): () => void {
  return Events.On(ConfigChangedEvent, (event) => {
    handler(readEventData<ConfigChangedPayload>(event));
  });
}

export function onAgentStreamChunk(handler: (payload: AgentStreamChunkEvent) => void): () => void {
  return Events.On(AgentStreamChunkEventName, (event) => {
    handler(readEventData<AgentStreamChunkEvent>(event));
  });
}

export function onAgentTurnDone(handler: (payload: AgentTurnDoneEvent) => void): () => void {
  return Events.On(AgentTurnDoneEventName, (event) => {
    handler(readEventData<AgentTurnDoneEvent>(event));
  });
}

export function onAgentTurnError(handler: (payload: AgentTurnErrorEvent) => void): () => void {
  return Events.On(AgentTurnErrorEventName, (event) => {
    handler(readEventData<AgentTurnErrorEvent>(event));
  });
}

export function onAgentTurnCancelled(handler: (payload: AgentTurnCancelledEvent) => void): () => void {
  return Events.On(AgentTurnCancelledEventName, (event) => {
    handler(readEventData<AgentTurnCancelledEvent>(event));
  });
}
