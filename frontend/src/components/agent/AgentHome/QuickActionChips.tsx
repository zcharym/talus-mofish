import {
  IconBook2,
  IconBrain,
  IconDeviceFloppy,
  IconPencil,
  IconVocabulary,
} from '@tabler/icons-react';
import { Button, Group, Stack, Text } from '@mantine/core';
import { QUICK_ACTION_DOMAINS, QUICK_ACTIONS, QuickActionId } from './quickActions';
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
    <Stack gap="sm" className={classes.row}>
      {QUICK_ACTION_DOMAINS.map(({ domain, label }) => {
        const actions = QUICK_ACTIONS.filter((action) => action.domain === domain);
        if (actions.length === 0) {
          return null;
        }

        return (
          <Stack key={domain} gap={6}>
            <Text size="xs" c="dimmed" ta="center">
              {label}
            </Text>
            <Group gap="sm" justify="center">
              {actions.map((action) => {
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
          </Stack>
        );
      })}
    </Stack>
  );
}
