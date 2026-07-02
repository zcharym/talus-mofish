import { useCallback, useEffect, useRef } from 'react';
import { Box, ScrollArea, Stack, Text, Title } from '@mantine/core';
import { ChatBubble } from '../ChatBubble';
import classes from './ChatThread.module.css';

export interface ChatMessageItem {
  id: string;
  session_id: string;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
  generating?: boolean;
}

interface ChatThreadProps {
  messages: ChatMessageItem[];
  sessionTitle: string | null;
  hasActiveSession: boolean;
}

export function ChatThread({ messages, sessionTitle, hasActiveSession }: ChatThreadProps) {
  const viewportRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = useCallback(() => {
    if (viewportRef.current) {
      viewportRef.current.scrollTop = viewportRef.current.scrollHeight;
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  if (!hasActiveSession) {
    return null;
  }

  return (
    <Box className={classes.thread}>
      <Box className={classes.header}>
        <Title order={4}>{sessionTitle ?? 'Chat'}</Title>
      </Box>

      <ScrollArea className={classes.scrollArea} viewportRef={viewportRef} type="auto">
        {messages.length === 0 ? (
          <Box className={classes.emptyThread}>
            <Text c="dimmed" size="sm">
              Send a message to begin this conversation.
            </Text>
          </Box>
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
    </Box>
  );
}
