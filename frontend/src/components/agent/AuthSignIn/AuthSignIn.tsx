import { IconBrandGithub, IconBrandGoogle } from '@tabler/icons-react';
import { Box, Button, Loader, Overlay, Stack, Text } from '@mantine/core';
import classes from './AuthSignIn.module.css';

interface AuthSignInProps {
  signingIn: 'github' | 'google' | null;
  onSignIn: (provider: 'github' | 'google') => Promise<void>;
}

export function AuthSignIn({ signingIn, onSignIn }: AuthSignInProps) {
  const waiting = signingIn !== null;

  return (
    <Box className={classes.wrapper}>
      <Stack gap="md" className={classes.card}>
        <Stack gap={4} align="center">
          <Text fw={600} size="lg">
            Sign in to Talus Agent
          </Text>
          <Text c="dimmed" size="sm" ta="center">
            Connect with GitHub or Google to personalize your workspace and sync your data.
          </Text>
        </Stack>

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
              Finish signing in in your browser
            </Text>
            <Text c="dimmed" size="xs" ta="center" maw={260}>
              Return here after approving access. This window will update automatically.
            </Text>
          </Stack>
        </Overlay>
      ) : null}
    </Box>
  );
}
