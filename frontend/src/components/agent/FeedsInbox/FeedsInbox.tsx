import { useState } from 'react';
import {
  IconBookmark,
  IconBookmarkFilled,
  IconExternalLink,
  IconPlus,
  IconRefresh,
  IconRss,
  IconTrash,
} from '@tabler/icons-react';
import {
  ActionIcon,
  Alert,
  Badge,
  Box,
  Burger,
  Button,
  Group,
  Loader,
  Modal,
  Paper,
  ScrollArea,
  SegmentedControl,
  Stack,
  Text,
  TextInput,
  Title,
} from '@mantine/core';
import { useFeedsInbox, type FeedsFilter } from '../../../hooks/useFeedsInbox';
import { toApiError, type FeedItem, type FeedSource } from '../../../utils/api';
import { openExternalURL } from '../../../utils/links';
import { notify } from '../../../services/notifications';
import classes from './FeedsInbox.module.css';

interface FeedsInboxProps {
  onOpenSidebar?: () => void;
  onOpenManagement: () => void;
}

const FILTERS: { value: FeedsFilter; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'saved', label: 'Read later' },
  { value: 'youtube', label: 'YouTube' },
  { value: 'bilibili', label: 'Bilibili' },
  { value: 'rss', label: 'RSS' },
];

export function FeedsInbox({ onOpenSidebar, onOpenManagement }: FeedsInboxProps) {
  const [filter, setFilter] = useState<FeedsFilter>('all');
  const inbox = useFeedsInbox(filter);
  const [addOpened, setAddOpened] = useState(false);
  const [addMode, setAddMode] = useState<'subscribe' | 'save'>('subscribe');
  const [addInput, setAddInput] = useState('');
  const [addTitle, setAddTitle] = useState('');
  const [adding, setAdding] = useState(false);

  const header = (
    <Box className={classes.header}>
      <Group gap="sm" wrap="nowrap" className={classes.headerTitle}>
        {onOpenSidebar ? (
          <Burger opened={false} onClick={onOpenSidebar} size="sm" aria-label="Open chats" />
        ) : null}
        <IconRss size={20} />
        <Title order={4}>Feeds</Title>
      </Group>
      <Group gap="sm" wrap="wrap">
        <Button size="xs" variant="light" leftSection={<IconPlus size={14} />} onClick={() => setAddOpened(true)}>
          Add
        </Button>
        <Button
          size="xs"
          variant="light"
          leftSection={<IconRefresh size={14} />}
          loading={inbox.refreshing}
          onClick={() => void inbox.refresh()}
        >
          Refresh
        </Button>
      </Group>
    </Box>
  );

  const submitAdd = async () => {
    const input = addInput.trim();
    if (!input) {
      return;
    }
    setAdding(true);
    try {
      if (addMode === 'save') {
        await inbox.saveLink(addTitle, input);
      } else {
        await inbox.addSource(input);
      }
      setAddOpened(false);
      setAddInput('');
      setAddTitle('');
      notify.success('Feeds', addMode === 'save' ? 'Saved for later.' : 'Source added.');
    } catch (err) {
      notify.failed('Feeds', toApiError(err));
    } finally {
      setAdding(false);
    }
  };

  return (
    <Box className={classes.page}>
      {header}
      {inbox.loading ? (
        <Box className={classes.empty}>
          <Loader size="sm" />
        </Box>
      ) : inbox.loadError && inbox.items.length === 0 && inbox.sources.length === 0 ? (
        <Box className={classes.empty}>
          <Alert color="red" title="Could not load feeds">
            {inbox.loadError}
          </Alert>
          <Button variant="light" onClick={() => void inbox.loadInbox()}>
            Try again
          </Button>
        </Box>
      ) : (
        <ScrollArea className={classes.body} type="auto">
          <Stack gap="md" className={classes.scrollInner}>
            <SegmentedControl
              value={filter}
              onChange={(value) => setFilter(value as FeedsFilter)}
              data={FILTERS}
            />

            {inbox.sources.length > 0 ? (
              <Box className={classes.sources}>
                {inbox.sources.map((source) => (
                  <SourceChip
                    key={source.id}
                    source={source}
                    onDelete={() => {
                      void inbox.deleteSource(source.id).catch((err: unknown) => {
                        notify.failed('Feeds', toApiError(err));
                      });
                    }}
                  />
                ))}
              </Box>
            ) : (
              <Alert color="gray" title="No subscriptions yet">
                Add an RSS URL, a YouTube channel or playlist, a Bilibili space / favorites URL, or type
                稍后再看 for Bilibili watch later. Optional API keys live in Management → Configuration →
                Feeds.
                <Button size="compact-xs" variant="subtle" ml="xs" onClick={onOpenManagement}>
                  Open config
                </Button>
              </Alert>
            )}

            {inbox.items.length === 0 ? (
              <Text size="sm" c="dimmed">
                Nothing in this view. Refresh after adding a source, or save a link for later.
              </Text>
            ) : (
              <Stack gap="sm">
                {inbox.items.map((item) => (
                  <FeedCard
                    key={item.id}
                    item={item}
                    onOpen={() => void openExternalURL(item.url)}
                    onToggleSaved={() => {
                      void inbox.setSaved(item.id, !item.saved).catch((err: unknown) => {
                        notify.failed('Feeds', toApiError(err));
                      });
                    }}
                    onMarkRead={() => {
                      void inbox.setRead(item.id, true).catch((err: unknown) => {
                        notify.failed('Feeds', toApiError(err));
                      });
                    }}
                  />
                ))}
              </Stack>
            )}
          </Stack>
        </ScrollArea>
      )}

      <Modal
        opened={addOpened}
        onClose={() => setAddOpened(false)}
        title="Add to Feeds"
        centered
      >
        <Stack>
          <SegmentedControl
            value={addMode}
            onChange={(value) => setAddMode(value as 'subscribe' | 'save')}
            data={[
              { value: 'subscribe', label: 'Subscribe' },
              { value: 'save', label: 'Save link' },
            ]}
          />
          {addMode === 'save' ? (
            <TextInput
              label="Title"
              placeholder="Optional"
              value={addTitle}
              onChange={(event) => setAddTitle(event.currentTarget.value)}
            />
          ) : null}
          <TextInput
            label={addMode === 'save' ? 'URL' : 'Feed, channel, playlist, or shortcut'}
            description="Examples: RSS URL, youtube.com/@handle, space.bilibili.com/uid, 稍后再看"
            placeholder="https://"
            value={addInput}
            onChange={(event) => setAddInput(event.currentTarget.value)}
            onKeyDown={(event) => {
              if (event.key === 'Enter') {
                event.preventDefault();
                void submitAdd();
              }
            }}
          />
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setAddOpened(false)}>
              Cancel
            </Button>
            <Button onClick={() => void submitAdd()} loading={adding} disabled={!addInput.trim()}>
              {addMode === 'save' ? 'Save' : 'Subscribe'}
            </Button>
          </Group>
        </Stack>
      </Modal>
    </Box>
  );
}

function SourceChip({ source, onDelete }: { source: FeedSource; onDelete: () => void }) {
  return (
    <Badge
      variant="light"
      rightSection={
        <ActionIcon size="xs" variant="transparent" aria-label={`Remove ${source.title}`} onClick={onDelete}>
          <IconTrash size={10} />
        </ActionIcon>
      }
    >
      {sourceKindLabel(source.kind)} · {source.title}
      {source.lastError ? ' · error' : ''}
    </Badge>
  );
}

function FeedCard({
  item,
  onOpen,
  onToggleSaved,
  onMarkRead,
}: {
  item: FeedItem;
  onOpen: () => void;
  onToggleSaved: () => void;
  onMarkRead: () => void;
}) {
  return (
    <Paper withBorder p="sm" className={classes.itemCard} opacity={item.read ? 0.72 : 1}>
      {item.thumbnailUrl ? (
        <img className={classes.thumb} src={item.thumbnailUrl} alt="" />
      ) : (
        <Box className={classes.thumbFallback}>
          <IconRss size={18} />
        </Box>
      )}
      <Stack gap={4}>
        <Group justify="space-between" wrap="nowrap" gap="xs">
          <Text fw={600} size="sm" lineClamp={2}>
            {item.title}
          </Text>
          <Group gap={4} wrap="nowrap">
            <ActionIcon variant="subtle" aria-label={item.saved ? 'Remove from read later' : 'Save for later'} onClick={onToggleSaved}>
              {item.saved ? <IconBookmarkFilled size={16} /> : <IconBookmark size={16} />}
            </ActionIcon>
            <ActionIcon variant="subtle" aria-label="Open" onClick={onOpen}>
              <IconExternalLink size={16} />
            </ActionIcon>
          </Group>
        </Group>
        <Group gap={6}>
          <Badge size="xs" variant="light">
            {sourceKindLabel(item.sourceKind)}
          </Badge>
          {item.sourceTitle ? (
            <Text size="xs" c="dimmed" lineClamp={1}>
              {item.sourceTitle}
            </Text>
          ) : null}
          {item.publishedAt ? (
            <Text size="xs" c="dimmed">
              {formatTimestamp(item.publishedAt)}
            </Text>
          ) : null}
        </Group>
        {item.summary ? (
          <Text size="xs" c="dimmed" lineClamp={2}>
            {item.summary}
          </Text>
        ) : null}
        {!item.read ? (
          <Button size="compact-xs" variant="subtle" onClick={onMarkRead}>
            Mark read
          </Button>
        ) : null}
      </Stack>
    </Paper>
  );
}

function sourceKindLabel(kind: string | undefined): string {
  switch (kind) {
    case 'youtube':
      return 'YouTube';
    case 'bilibili':
      return 'Bilibili';
    case 'readlater':
      return 'Saved';
    default:
      return 'RSS';
  }
}

function formatTimestamp(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleString();
}
