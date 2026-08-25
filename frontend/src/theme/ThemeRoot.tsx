import { MantineProvider } from '@mantine/core';
import { Notifications } from '@mantine/notifications';
import type { ReactNode } from 'react';
import type { ThemeOption } from '../types/theme';
import { talusCssVariablesResolver, talusTheme } from './theme';
import '@mantine/core/styles.css';
import '@mantine/notifications/styles.css';
import './fonts.css';
import './atmosphere.css';

interface ThemeRootProps {
  colorScheme: ThemeOption;
  children: ReactNode;
}

export function ThemeRoot({ colorScheme, children }: ThemeRootProps) {
  return (
    <MantineProvider
      theme={talusTheme}
      cssVariablesResolver={talusCssVariablesResolver}
      defaultColorScheme="auto"
      forceColorScheme={colorScheme === 'auto' ? undefined : colorScheme}
    >
      <Notifications position="top-right" limit={5} />
      {children}
    </MantineProvider>
  );
}
