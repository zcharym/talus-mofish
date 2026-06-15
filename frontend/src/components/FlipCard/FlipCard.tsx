import { useState } from "react";
import {
  ActionIcon,
  Button,
  Group,
  Menu,
  Stack,
  Text,
} from "@mantine/core";
import { IconDots, IconTrash } from "@tabler/icons-react";
import classes from "./FlipCard.module.css";

export interface FlipCardProps {
  title: string;
  front: React.ReactNode;
  back: React.ReactNode;
  headerExtra?: React.ReactNode;
  modelCss?: string;
  onDelete?: () => void;
}

export function FlipCard({ title, front, back, headerExtra, modelCss, onDelete }: FlipCardProps) {
  const [flipped, setFlipped] = useState(false);

  return (
    <Stack gap="sm">
      {modelCss && <style>{modelCss}</style>}

      <Stack gap="xs">
        <Group justify="space-between" wrap="nowrap">
          <Text fw={500} lineClamp={1}>{title}</Text>
          {onDelete && (
            <Menu withinPortal position="bottom-end" shadow="sm">
              <Menu.Target>
                <ActionIcon variant="subtle" color="gray">
                  <IconDots size={16} />
                </ActionIcon>
              </Menu.Target>
              <Menu.Dropdown>
                <Menu.Item leftSection={<IconTrash size={14} />} color="red" onClick={onDelete}>
                  Delete
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          )}
        </Group>
        {headerExtra}
      </Stack>

      <div className={classes.flipContainer}>
        <div className={classes.flipInner} data-flipped={flipped}>
          <div className={classes.flipFace}>{front}</div>
          <div className={classes.flipBack}>{back}</div>
        </div>
      </div>

      <Button fullWidth variant="light" onClick={() => setFlipped((v) => !v)}>
        {flipped ? "Show front" : "Flip"}
      </Button>
    </Stack>
  );
}
