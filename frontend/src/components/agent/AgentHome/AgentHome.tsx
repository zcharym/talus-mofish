import { useMemo, useState } from 'react';
import { IconMessageChatbot } from '@tabler/icons-react';
import { Anchor, Avatar, Box, Burger, Group, Loader, Stack, Title } from '@mantine/core';
import { AuthSignIn } from '../AuthSignIn';
import { UserProfile, getWelcomeMessage } from '../../../utils/userProfile';
import { AgentHomeInput } from './AgentHomeInput';
import { QuickActionChips } from './QuickActionChips';
import { QUICK_ACTIONS, QuickActionId } from './quickActions';
import classes from './AgentHome.module.css';

interface AgentHomeProps {
  user: UserProfile | null;
  userLoading: boolean;
  signingIn: 'email' | 'github' | 'google' | null;
  sending: boolean;
  onSend: (content: string) => Promise<void>;
  onCancel?: () => Promise<void>;
  onSignInWithEmail: (email: string) => Promise<void>;
  onSignIn: (provider: 'github' | 'google') => Promise<void>;
  onQuickAction: (actionId: QuickActionId, prompt: string, autoSend: boolean) => void;
  onOpenSidebar?: () => void;
}

export function AgentHome({
  user,
  userLoading,
  signingIn,
  sending,
  onSend,
  onCancel,
  onSignInWithEmail,
  onSignIn,
  onQuickAction,
  onOpenSidebar,
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
    if (!action.autoSend && action.id !== 'play_sudoku') {
      setDraft(action.prompt);
    }
  };

  const menuButton = onOpenSidebar ? (
    <Burger
      opened={false}
      onClick={onOpenSidebar}
      size="sm"
      aria-label="Open chats"
    />
  ) : null;

  if (userLoading) {
    return (
      <Box className={classes.page} data-with-menu={onOpenSidebar ? true : undefined}>
        {menuButton ? <Box className={classes.menuButton}>{menuButton}</Box> : null}
        <Loader size="sm" />
      </Box>
    );
  }

  if (!user) {
    return (
      <Box className={classes.page} data-with-menu={onOpenSidebar ? true : undefined}>
        {menuButton ? <Box className={classes.menuButton}>{menuButton}</Box> : null}
        <AuthSignIn
          signingIn={signingIn}
          onSignInWithEmail={onSignInWithEmail}
          onSignIn={async (provider) => {
            await onSignIn(provider);
          }}
        />
      </Box>
    );
  }

  return (
    <Box className={classes.page} data-with-menu={onOpenSidebar ? true : undefined}>
      {menuButton ? <Box className={classes.menuButton}>{menuButton}</Box> : null}
      <Stack gap="xl" className={classes.content}>
        <Group
          className={`${classes.greetingRow} ${classes.reveal}`}
          justify="space-between"
          align="flex-start"
          wrap="wrap"
        >
          <Group gap="sm" align="center" wrap="nowrap">
            {user.avatar_url ? (
              <Avatar src={user.avatar_url} alt={user.display_name} radius="xl" size={28} />
            ) : (
              <IconMessageChatbot size={28} stroke={1.25} className={classes.icon} />
            )}
            <Title order={2} className={classes.greeting} lineClamp={2}>
              {welcomeMessage}
            </Title>
          </Group>
          <Anchor size="sm" c="dimmed" className={classes.feedbackLink}>
            Give feedback
          </Anchor>
        </Group>

        <Box className={`${classes.reveal} ${classes.revealDelay1}`}>
          <AgentHomeInput
            sending={sending}
            value={draft}
            onValueChange={setDraft}
            onSend={onSend}
            onCancel={onCancel}
          />
        </Box>

        <Box className={`${classes.reveal} ${classes.revealDelay2}`}>
          <QuickActionChips disabled={sending} onSelect={handleQuickAction} />
        </Box>
      </Stack>
    </Box>
  );
}
