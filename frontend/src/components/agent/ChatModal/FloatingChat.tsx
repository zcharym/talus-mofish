import { useEffect } from 'react';
import { ActionIcon, Box, FocusTrap, Paper, Portal, Text, Transition } from '@mantine/core';
import { IconX } from '@tabler/icons-react';
import { ChatPanel } from './ChatPanel';
import type { ChatMessageItem, ChatModalDomain } from './types';
import classes from './FloatingChat.module.css';

export interface FloatingChatProps {
  opened: boolean;
  onClose: () => void;
  domain: ChatModalDomain;
  messages: ChatMessageItem[];
  sending: boolean;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
}

export function FloatingChat({
  opened,
  onClose,
  domain,
  messages,
  sending,
  onSend,
  onCancel,
}: FloatingChatProps) {
  useEffect(() => {
    if (!opened) {
      return undefined;
    }
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === 'Escape') {
        onClose();
      }
    };
    window.addEventListener('keydown', onKeyDown);
    return () => window.removeEventListener('keydown', onKeyDown);
  }, [opened, onClose]);

  return (
    <Portal>
      <Transition mounted={opened} transition="fade" duration={160}>
        {(fadeStyle) => (
          <Box className={classes.backdrop} style={fadeStyle} onClick={onClose} />
        )}
      </Transition>
      <Transition mounted={opened} transition="slide-up" duration={180}>
        {(panelStyle) => (
          <Paper
            className={classes.shell}
            style={panelStyle}
            shadow="xl"
            radius="md"
            withBorder
            role="dialog"
            aria-modal="true"
            aria-label={domain.title}
            onClick={(event) => event.stopPropagation()}
          >
            <FocusTrap active={opened}>
              <div className={classes.bodyWrap}>
                <div className={classes.header}>
                  <Text className={classes.title} fw={600} size="sm">
                    {domain.title}
                  </Text>
                  <ActionIcon
                    variant="subtle"
                    color="gray"
                    aria-label="Close chat"
                    onClick={onClose}
                  >
                    <IconX size={16} />
                  </ActionIcon>
                </div>
                <div className={classes.body}>
                  <ChatPanel
                    domain={domain}
                    messages={messages}
                    sending={sending}
                    autoFocus
                    onSend={onSend}
                    onCancel={onCancel}
                  />
                </div>
              </div>
            </FocusTrap>
          </Paper>
        )}
      </Transition>
    </Portal>
  );
}
