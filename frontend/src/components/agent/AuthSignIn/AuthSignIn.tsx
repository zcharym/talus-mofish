import { useState } from 'react';
import { IconBrandGithub, IconBrandGoogle, IconMail } from '@tabler/icons-react';
import {
  Box,
  Button,
  Divider,
  Loader,
  Overlay,
  Stack,
  Text,
  TextInput,
} from '@mantine/core';
import classes from './AuthSignIn.module.css';

interface AuthSignInProps {
  signingIn: 'email' | 'github' | 'google' | null;
  onSignInWithEmail: (email: string) => Promise<void>;
  onSignIn: (provider: 'github' | 'google') => Promise<void>;
}

export function AuthSignIn({ signingIn, onSignInWithEmail, onSignIn }: AuthSignInProps) {
  const [email, setEmail] = useState('');
  const waiting = signingIn !== null;
  const waitingForEmail = signingIn === 'email';

  const handleEmailSubmit = () => {
    const trimmed = email.trim();
    if (!trimmed || waiting) {
      return;
    }
    void onSignInWithEmail(trimmed);
  };

  return (
    <Box className={classes.wrapper}>
      <Stack gap="md" className={classes.card}>
        <Stack gap={4} align="center">
          <Text fw={600} size="lg">
            Sign in to Talus Agent
          </Text>
          <Text c="dimmed" size="sm" ta="center">
            Enter your email to receive a sign-in link, or continue with a connected account.
          </Text>
        </Stack>

        <Stack gap="sm">
          <TextInput
            label="Email"
            placeholder="you@example.com"
            value={email}
            disabled={waiting}
            onChange={(event) => setEmail(event.currentTarget.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                event.preventDefault();
                handleEmailSubmit();
              }
            }}
          />
          <Button
            className={classes.button}
            leftSection={<IconMail size={18} />}
            size="md"
            disabled={waiting || !email.trim()}
            loading={waitingForEmail}
            onClick={handleEmailSubmit}
          >
            Continue with email
          </Button>
        </Stack>

        <Divider label="or continue with" labelPosition="center" />

        <Stack gap="sm">
          <Button
            className={classes.button}
            leftSection={<IconBrandGithub size={18} />}
            variant="default"
            size="md"
            disabled={waiting}
            onClick={() => void onSignIn('github')}
          >
            Continue with GitHub
          </Button>
          <Button
            className={classes.button}
            leftSection={<IconBrandGoogle size={18} />}
            variant="default"
            size="md"
            disabled={waiting}
            onClick={() => void onSignIn('google')}
          >
            Continue with Google
          </Button>
        </Stack>
      </Stack>

      {waiting ? (
        <Overlay
          className={classes.overlay}
          color="#000"
          backgroundOpacity={0.35}
          blur={2}
          center
        >
          <Stack align="center" gap="sm" className={classes.overlayContent}>
            <Loader size="sm" />
            <Text fw={500} size="sm" ta="center">
              {waitingForEmail ? 'Check your email' : 'Finish signing in in your browser'}
            </Text>
            <Text c="dimmed" size="xs" ta="center" maw={260}>
              {waitingForEmail
                ? 'Click the link in your inbox, then return here. This window will update automatically.'
                : 'Return here after approving access. This window will update automatically.'}
            </Text>
          </Stack>
        </Overlay>
      ) : null}
    </Box>
  );
}
