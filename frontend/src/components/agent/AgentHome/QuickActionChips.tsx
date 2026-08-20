import {
  IconBook2,
  IconBrain,
  IconDeviceFloppy,
  IconGrid3x3,
  IconPencil,
  IconVocabulary,
} from '@tabler/icons-react';
import { Button, Group, Marquee, Stack, Text } from '@mantine/core';
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
  play_sudoku: IconGrid3x3,
} as const;

function ActionChip({
  action,
  disabled,
  onSelect,
}: {
  action: (typeof QUICK_ACTIONS)[number];
  disabled: boolean;
  onSelect: (actionId: QuickActionId) => void;
}) {
  const Icon = ICONS[action.id];
  return (
    <Button
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
}

export function QuickActionChips({ disabled = false, onSelect }: QuickActionChipsProps) {
  return (
    <Stack gap="sm" className={classes.row}>
      {QUICK_ACTION_DOMAINS.map(({ domain, label }, index) => {
        const actions = QUICK_ACTIONS.filter((action) => action.domain === domain);
        if (actions.length === 0) {
          return null;
        }

        const chips = actions.map((action) => (
          <ActionChip key={action.id} action={action} disabled={disabled} onSelect={onSelect} />
        ));

        return (
          <Stack key={domain} gap={6}>
            <Text size="xs" c="dimmed" ta="center">
              {label}
            </Text>
            {actions.length > 4 ? (
              <Marquee pauseOnHover reverse={index % 2 === 1} gap="sm" duration={25000}>
                {chips}
              </Marquee>
            ) : (
              <Group gap="sm" justify="center">
                {chips}
              </Group>
            )}
          </Stack>
        );
      })}
    </Stack>
  );
}
