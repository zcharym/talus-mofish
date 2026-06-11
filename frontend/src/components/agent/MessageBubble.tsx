import { Box, Paper, Text } from '@mantine/core';
import classes from './MessageBubble.module.css';

interface MessageBubbleProps {
  role: 'user' | 'assistant';
  content: string;
}

export function MessageBubble({ role, content }: MessageBubbleProps) {
  const isUser = role === 'user';

  return (
    <Box className={`${classes.row} ${isUser ? classes.rowUser : classes.rowAssistant}`}>
      <Paper
        className={`${classes.bubble} ${isUser ? classes.bubbleUser : classes.bubbleAssistant}`}
        radius="md"
        p="sm"
        shadow="xs"
      >
        <Text size="sm" className={classes.content}>
          {content}
        </Text>
      </Paper>
    </Box>
  );
}
