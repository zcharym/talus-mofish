import { useConfigFlags } from './useConfigFlags';

export function useCloudflareConfigured() {
  return useConfigFlags().cloudflareConfigured;
}
