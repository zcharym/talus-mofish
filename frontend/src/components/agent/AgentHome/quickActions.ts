export type QuickActionDomain = 'english' | 'general';

export type QuickActionId =
  | 'recite_words'
  | 'read_article'
  | 'ielts_writing'
  | 'vocabulary_quiz'
  | 'save_as_flow';

export interface QuickAction {
  id: QuickActionId;
  label: string;
  prompt: string;
  autoSend: boolean;
  domain: QuickActionDomain;
}

export const QUICK_ACTIONS: QuickAction[] = [
  {
    id: 'recite_words',
    label: 'Recite words',
    prompt: 'Start a vocabulary recitation session with 20 words from my decks.',
    autoSend: true,
    domain: 'english',
  },
  {
    id: 'read_article',
    label: 'Read article',
    prompt: 'Help me read and analyze an article from my reading library.',
    autoSend: true,
    domain: 'english',
  },
  {
    id: 'ielts_writing',
    label: 'IELTS writing',
    prompt: 'Start an IELTS Writing Task 2 practice session with feedback.',
    autoSend: true,
    domain: 'english',
  },
  {
    id: 'vocabulary_quiz',
    label: 'Vocabulary quiz',
    prompt: 'Quiz me on vocabulary using typed recall.',
    autoSend: true,
    domain: 'english',
  },
  {
    id: 'save_as_flow',
    label: 'Save as flow',
    prompt: 'Help me design a reusable agent flow template from this goal:',
    autoSend: false,
    domain: 'general',
  },
];

export const QUICK_ACTION_DOMAINS: { domain: QuickActionDomain; label: string }[] = [
  { domain: 'english', label: 'English Learning' },
  { domain: 'general', label: 'Agent' },
];
