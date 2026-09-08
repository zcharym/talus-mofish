import { UserProfile } from '../../../utils/userProfile';

export function canShowCloudflareTab(
  user: UserProfile | null,
  cloudflareConfigured: boolean,
): boolean {
  return user?.provider === 'debug' || cloudflareConfigured;
}
