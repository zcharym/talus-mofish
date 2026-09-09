import { useConfigFlags } from './useConfigFlags';

export function useConfigColorScheme() {
  const { theme, applyTheme } = useConfigFlags();
  return { colorScheme: theme, applyTheme };
}
