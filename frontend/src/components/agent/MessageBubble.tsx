import { Box, Loader, Paper, Text } from '@mantine/core';
import classes from './MessageBubble.module.css';

interface MessageBubbleProps {
  role: 'user' | 'assistant';
  content: string;
  generating?: boolean;
}

export function MessageBubble({ role, content, generating }: MessageBubbleProps) {
  const isUser = role === 'user';
  const displayContent = content || (generating ? '' : '…');

  return (
    <Box className={`${classes.row} ${isUser ? classes.rowUser : classes.rowAssistant}`}>
      <Paper
        className={`${classes.bubble} ${isUser ? classes.bubbleUser : classes.bubbleAssistant}`}
        radius="md"
        p="sm"
        shadow="xs"
      >
        {generating && !content ? (
          <Loader size="xs" type="dots" />
        ) : (
          <Text size="sm" className={classes.content}>
            {displayContent}
            {generating && content ? ' ▍' : null}
          </Text>
        )}
      </Paper>
    </Box>
  );
}
