import { useState } from 'react';
import { IconArrowUp, IconPencil, IconPlayerStop } from '@tabler/icons-react';
import { ActionIcon, Box, Group, Textarea } from '@mantine/core';
import classes from './ChatInput.module.css';

interface ChatInputProps {
  disabled: boolean;
  sending: boolean;
  placeholder?: string;
  autoFocus?: boolean;
  compact?: boolean;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
}

export function ChatInput({
  disabled,
  sending,
  placeholder,
  autoFocus = false,
  compact = false,
  onSend,
  onCancel,
}: ChatInputProps) {
  const [value, setValue] = useState('');

  const handleSend = async () => {
    const content = value.trim();
    if (!content || disabled || sending) {
      return;
    }

    setValue('');
    await onSend(content);
  };

  const resolvedPlaceholder =
    placeholder ?? (disabled ? 'Select or create a chat to begin' : 'Message Talus Agent…');

  return (
    <Box className={classes.wrapper} data-compact={compact || undefined}>
      <Group align="flex-end" gap="sm" className={classes.inputRow}>
        <Textarea
          className={classes.textarea}
          placeholder={resolvedPlaceholder}
          value={value}
          onChange={(event) => setValue(event.currentTarget.value)}
          disabled={disabled || sending}
          autosize
          minRows={1}
          maxRows={8}
          data-autofocus={autoFocus || undefined}
          onKeyDown={(event) => {
            if (event.key === 'Enter' && !event.shiftKey) {
              event.preventDefault();
              void handleSend();
            }
          }}
        />
        {sending && onCancel ? (
          <ActionIcon
            size="lg"
            radius="xl"
            variant="light"
            color="red"
            className={classes.sendButton}
            aria-label="Stop generating"
            onClick={() => void onCancel()}
          >
            <IconPlayerStop size={18} />
          </ActionIcon>
        ) : (
          <ActionIcon
            size="lg"
            radius="xl"
            variant="filled"
            className={classes.sendButton}
            aria-label="Send message"
            disabled={disabled || sending || !value.trim()}
            onClick={() => void handleSend()}
          >
            <IconArrowUp size={18} />
          </ActionIcon>
        )}
      </Group>
      <Box className={classes.hint}>
        <IconPencil size={12} />
        <span>Enter to send, Shift+Enter for a new line</span>
      </Box>
    </Box>
  );
}
