export interface UserProfile {
  id: string;
  provider: 'github' | 'google';
  display_name: string;
  email: string;
  avatar_url: string;
}

export function getTimeGreeting(date = new Date()): string {
  const hour = date.getHours();
  if (hour < 12) {
    return 'Good morning';
  }
  if (hour < 18) {
    return 'Good afternoon';
  }
  return 'Good evening';
}

export function getWelcomeMessage(displayName: string, date = new Date()): string {
  const trimmed = displayName.trim();
  const name = trimmed || 'there';
  return `${getTimeGreeting(date)}, ${name}!`;
}
