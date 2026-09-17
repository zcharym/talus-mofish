import { useCallback, useEffect, useRef } from 'react';
import { Button, Group, ScrollArea, Stack, Text } from '@mantine/core';
import { ChatBubble } from '../ChatBubble';
import { ChatInput } from '../ChatInput';
import type { ChatMessageItem, ChatModalDomain } from './types';
import classes from './ChatPanel.module.css';

export interface ChatPanelProps {
  domain: ChatModalDomain;
  messages: ChatMessageItem[];
  sending: boolean;
  autoFocus?: boolean;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
}

export function ChatPanel({
  domain,
  messages,
  sending,
  autoFocus = false,
  onSend,
  onCancel,
}: ChatPanelProps) {
  const viewportRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = useCallback(() => {
    if (viewportRef.current) {
      viewportRef.current.scrollTop = viewportRef.current.scrollHeight;
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  return (
    <div className={classes.panel}>
      <ScrollArea className={classes.thread} viewportRef={viewportRef} type="auto">
        {messages.length === 0 ? (
          <Text className={classes.empty} c="dimmed" size="sm">
            {domain.emptyHint}
          </Text>
        ) : (
          <Stack gap="xs" className={classes.messageList}>
            {messages.map((message) => (
              <ChatBubble
                key={message.id}
                role={message.role}
                content={message.content}
                streaming={message.generating}
              />
            ))}
          </Stack>
        )}
      </ScrollArea>

      {messages.length === 0 && domain.suggestions.length > 0 ? (
        <Group className={classes.suggestions} gap="xs" wrap="wrap">
          {domain.suggestions.map((suggestion) => (
            <Button
              key={suggestion.id}
              variant="light"
              size="compact-sm"
              radius="xl"
              disabled={sending}
              onClick={() => void onSend(suggestion.prompt)}
            >
              {suggestion.label}
            </Button>
          ))}
        </Group>
      ) : null}

      <div className={classes.composer}>
        <ChatInput
          disabled={false}
          sending={sending}
          placeholder={domain.placeholder}
          autoFocus={autoFocus}
          compact
          onSend={onSend}
          onCancel={onCancel}
        />
      </div>
    </div>
  );
}
