import { useEffect, useState } from 'react';
import { IconArrowUp, IconPlayerStop } from '@tabler/icons-react';
import { ActionIcon, Box, Group, Text, Textarea } from '@mantine/core';
import classes from './AgentHomeInput.module.css';

interface AgentHomeInputProps {
  disabled?: boolean;
  sending: boolean;
  value?: string;
  onValueChange?: (value: string) => void;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
}

export function AgentHomeInput({
  disabled = false,
  sending,
  value,
  onValueChange,
  onSend,
  onCancel,
}: AgentHomeInputProps) {
  const [internalValue, setInternalValue] = useState('');
  const currentValue = value ?? internalValue;

  const setValue = (next: string) => {
    if (value === undefined) {
      setInternalValue(next);
    }
    onValueChange?.(next);
  };

  useEffect(() => {
    if (value !== undefined) {
      setInternalValue(value);
    }
  }, [value]);

  const handleSend = async () => {
    const content = currentValue.trim();
    if (!content || disabled || sending) {
      return;
    }

    setValue('');
    await onSend(content);
  };

  return (
    <Box className={classes.card}>
      <Textarea
        className={classes.textarea}
        placeholder="Ask Talus Agent anything…"
        value={currentValue}
        onChange={(event) => setValue(event.currentTarget.value)}
        disabled={disabled || sending}
        autosize
        minRows={3}
        maxRows={10}
        onKeyDown={(event) => {
          if (event.key === 'Enter' && !event.shiftKey) {
            event.preventDefault();
            void handleSend();
          }
        }}
      />

      <Group justify="space-between" align="center" className={classes.footer}>
        <Text size="xs" c="dimmed">
          Enter to send, Shift+Enter for a new line
        </Text>
        {sending && onCancel ? (
          <ActionIcon
            size="lg"
            radius="xl"
            variant="light"
            color="red"
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
            aria-label="Send message"
            disabled={disabled || sending || !currentValue.trim()}
            onClick={() => void handleSend()}
          >
            <IconArrowUp size={18} />
          </ActionIcon>
        )}
      </Group>
    </Box>
  );
}
