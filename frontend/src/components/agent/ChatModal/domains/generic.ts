import type { ChatModalDomain } from '../types';

export function createGenericChatDomain(): ChatModalDomain {
  return {
    id: 'generic',
    title: 'Chat',
    placeholder: 'Message Talus Agent…',
    emptyHint: 'Ask anything. This overlay chat is not saved as a sidebar session.',
    suggestions: [
      {
        id: 'help',
        label: 'What can you help with?',
        prompt: 'What can you help me with in Talus Echo?',
      },
      {
        id: 'concise',
        label: 'Give a concise overview',
        prompt: 'Give a concise overview of how this agent chat works.',
      },
    ],
    buildContext: () =>
      'You are Talus Agent in an overlay chat. You have no extra domain snapshot. Be concise.',
  };
}
