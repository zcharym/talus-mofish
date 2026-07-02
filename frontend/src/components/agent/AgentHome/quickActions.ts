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
}

export const QUICK_ACTIONS: QuickAction[] = [
  {
    id: 'recite_words',
    label: 'Recite words',
    prompt: 'Start a vocabulary recitation session with 20 words from my decks.',
    autoSend: true,
  },
  {
    id: 'read_article',
    label: 'Read article',
    prompt: 'Help me read and analyze an article from my reading library.',
    autoSend: true,
  },
  {
    id: 'ielts_writing',
    label: 'IELTS writing',
    prompt: 'Start an IELTS Writing Task 2 practice session with feedback.',
    autoSend: true,
  },
  {
    id: 'vocabulary_quiz',
    label: 'Vocabulary quiz',
    prompt: 'Quiz me on vocabulary using typed recall.',
    autoSend: true,
  },
  {
    id: 'save_as_flow',
    label: 'Save as flow',
    prompt: 'Help me design a reusable learning flow template from this goal:',
    autoSend: false,
  },
];
