import {
  IconBook2,
  IconBrain,
  IconDeviceFloppy,
  IconPencil,
  IconVocabulary,
} from '@tabler/icons-react';
import { Button, Group } from '@mantine/core';
import { QUICK_ACTIONS, QuickActionId } from './quickActions';
import classes from './QuickActionChips.module.css';

interface QuickActionChipsProps {
  disabled?: boolean;
  onSelect: (actionId: QuickActionId) => void;
}

const ICONS = {
  recite_words: IconVocabulary,
  read_article: IconBook2,
  ielts_writing: IconPencil,
  vocabulary_quiz: IconBrain,
  save_as_flow: IconDeviceFloppy,
} as const;

export function QuickActionChips({ disabled = false, onSelect }: QuickActionChipsProps) {
  return (
    <Group gap="sm" justify="center" className={classes.row}>
      {QUICK_ACTIONS.map((action) => {
        const Icon = ICONS[action.id];
        return (
          <Button
            key={action.id}
            className={classes.chip}
            variant="default"
            radius="xl"
            size="sm"
            leftSection={<Icon size={16} />}
            disabled={disabled}
            onClick={() => onSelect(action.id)}
          >
            {action.label}
          </Button>
        );
      })}
    </Group>
  );
}
