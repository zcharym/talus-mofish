import { useMemo, useState } from 'react';
import { IconMessageChatbot } from '@tabler/icons-react';
import { Anchor, Avatar, Box, Group, Loader, Stack, Title } from '@mantine/core';
import { AuthSignIn } from '../AuthSignIn';
import { UserProfile, getWelcomeMessage } from '../../../utils/userProfile';
import { AgentHomeInput } from './AgentHomeInput';
import { QuickActionChips } from './QuickActionChips';
import { QUICK_ACTIONS, QuickActionId } from './quickActions';
import classes from './AgentHome.module.css';

interface AgentHomeProps {
  user: UserProfile | null;
  userLoading: boolean;
  signingIn: 'github' | 'google' | null;
  sending: boolean;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
  onSignIn: (provider: 'github' | 'google') => Promise<void>;
  onQuickAction: (actionId: QuickActionId, prompt: string, autoSend: boolean) => void;
}

export function AgentHome({
  user,
  userLoading,
  signingIn,
  sending,
  onSend,
  onCancel,
  onSignIn,
  onQuickAction,
}: AgentHomeProps) {
  const [draft, setDraft] = useState('');

  const welcomeMessage = useMemo(() => {
    if (!user) {
      return null;
    }
    return getWelcomeMessage(user.display_name);
  }, [user]);

  const handleQuickAction = (actionId: QuickActionId) => {
    const action = QUICK_ACTIONS.find((item) => item.id === actionId);
    if (!action) {
      return;
    }
    onQuickAction(action.id, action.prompt, action.autoSend);
    if (!action.autoSend) {
      setDraft(action.prompt);
    }
  };

  if (userLoading) {
    return (
      <Box className={classes.page}>
        <Loader size="sm" />
      </Box>
    );
  }

  if (!user) {
    return (
      <Box className={classes.page}>
        <AuthSignIn
          signingIn={signingIn}
          onSignIn={async (provider) => {
            await onSignIn(provider);
          }}
        />
      </Box>
    );
  }

  return (
    <Box className={classes.page}>
      <Stack gap="xl" className={classes.content}>
        <Group justify="space-between" align="flex-start" wrap="nowrap">
          <Group gap="sm" align="center">
            {user.avatar_url ? (
              <Avatar src={user.avatar_url} alt={user.display_name} radius="xl" size={28} />
            ) : (
              <IconMessageChatbot size={28} stroke={1.25} className={classes.icon} />
            )}
            <Title order={2} className={classes.greeting}>
              {welcomeMessage}
            </Title>
          </Group>
          <Anchor size="sm" c="dimmed" className={classes.feedbackLink}>
            Give feedback
          </Anchor>
        </Group>

        <AgentHomeInput
          sending={sending}
          value={draft}
          onValueChange={setDraft}
          onSend={onSend}
          onCancel={onCancel}
        />

        <QuickActionChips disabled={sending} onSelect={handleQuickAction} />
      </Stack>
    </Box>
  );
}
