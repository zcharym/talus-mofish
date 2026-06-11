import { useCallback, useEffect, useRef } from 'react';
import { IconMessageChatbot } from '@tabler/icons-react';
import { Box, ScrollArea, Stack, Text, Title } from '@mantine/core';
import { MessageBubble } from './MessageBubble';
import classes from './ChatThread.module.css';

export interface ChatMessageItem {
  id: string;
  session_id: string;
  role: 'user' | 'assistant';
  content: string;
  created_at: string;
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
    return (
      <Box className={classes.emptyState}>
        <IconMessageChatbot size={48} stroke={1.25} className={classes.emptyIcon} />
        <Title order={3}>Talus Agent</Title>
        <Text c="dimmed" size="sm" mt="xs">
          Select a chat or start a new conversation.
        </Text>
      </Box>
    );
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
              <MessageBubble key={message.id} role={message.role} content={message.content} />
            ))}
          </Stack>
        )}
      </ScrollArea>
    </Box>
  );
}
