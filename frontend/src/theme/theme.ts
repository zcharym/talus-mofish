import { createTheme, CSSVariablesResolver, MantineColorsTuple } from '@mantine/core';

/** Persimmon accent. Chrome uses radius md; composer/bubbles lg; chips are pills. */
const persimmon: MantineColorsTuple = [
  '#fdf4ee',
  '#f8e4d4',
  '#f0c9ab',
  '#e6a87a',
  '#dc8a56',
  '#d97845',
  '#c46232',
  '#a34e28',
  '#7e3d22',
  '#5c2e1b',
];

/** Moss-tinted dark scale so Agent Mantine surfaces follow the theme, not default blue-gray. */
const mossDark: MantineColorsTuple = [
  '#ece8de',
  '#cfcbbf',
  '#a8a496',
  '#7a7c70',
  '#3f433b',
  '#2a2e26',
  '#1c1f19',
  '#141612',
  '#10120e',
  '#0c0d0a',
];

const uiFont = '"Bricolage Grotesque Variable", "Bricolage Grotesque", sans-serif';
const monoFont = '"JetBrains Mono", ui-monospace, monospace';

export const talusTheme = createTheme({
  primaryColor: 'persimmon',
  defaultRadius: 'md',
  fontFamily: uiFont,
  fontFamilyMonospace: monoFont,
  white: '#fcfbf7',
  black: '#1c2118',
  colors: {
    persimmon,
    dark: mossDark,
  },
  headings: {
    fontFamily: uiFont,
    fontWeight: '600',
  },
  cursorType: 'pointer',
});

export const talusCssVariablesResolver: CSSVariablesResolver = () => ({
  variables: {
    '--talus-accent': '#d97845',
    '--talus-moss': '#5c7a54',
  },
  light: {
    '--mantine-color-body': '#eef0ea',
    '--mantine-color-default': '#f6f5f0',
    '--mantine-color-default-hover': '#eceae3',
    '--mantine-color-default-border': 'rgba(28, 33, 24, 0.12)',
    '--talus-bg': '#eef0ea',
    '--talus-surface': '#f6f5f0',
    '--talus-surface-raised': '#fcfbf7',
    '--talus-hairline': 'rgba(28, 33, 24, 0.12)',
    '--talus-ink': '#1c2118',
  },
  dark: {
    '--mantine-color-body': '#141612',
    '--mantine-color-default': '#1c1f19',
    '--mantine-color-default-hover': '#242820',
    '--mantine-color-default-border': 'rgba(236, 232, 222, 0.1)',
    '--talus-bg': '#141612',
    '--talus-surface': '#1c1f19',
    '--talus-surface-raised': '#242820',
    '--talus-hairline': 'rgba(236, 232, 222, 0.1)',
    '--talus-ink': '#ece8de',
  },
});
