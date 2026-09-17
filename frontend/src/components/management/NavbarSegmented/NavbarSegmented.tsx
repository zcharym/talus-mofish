import { useState } from 'react';
import {
  IconAdjustments,
  IconBook,
  IconBug,
  IconComponents,
  IconFileText,
  IconInfoCircle,
  IconLanguage,
  IconMarkdown,
  IconMessageChatbot,
  IconSearch,
  IconUpload,
  IconVocabulary,
} from '@tabler/icons-react';
import { NavLink, Stack, Text } from '@mantine/core';
import { SystemService } from '../../../utils/api';
import {
  isEnglishLearningRoute,
  isObsidianRoute,
  ManagementRoute,
  type ManagementRouteId,
} from '../../../navigation/routes';
import classes from './NavbarSegmented.module.css';

interface NavItem {
  id: ManagementRouteId;
  label: string;
  icon: typeof IconUpload;
}

const englishLearningItems: NavItem[] = [
  { id: ManagementRoute.EnglishImport, label: 'Import', icon: IconUpload },
  { id: ManagementRoute.EnglishReading, label: 'Reading', icon: IconBook },
  { id: ManagementRoute.EnglishVocabulary, label: 'Vocabulary', icon: IconVocabulary },
];

const obsidianItems: NavItem[] = [
  { id: ManagementRoute.ObsidianNotes, label: 'Notes', icon: IconFileText },
  { id: ManagementRoute.ObsidianSearch, label: 'Search', icon: IconSearch },
];

const debugItems: NavItem[] = [
  { id: ManagementRoute.DebugPlayground, label: 'Component playground', icon: IconComponents },
];

interface NavbarSegmentedProps {
  activeItem: string;
  debugMode: boolean;
  onActiveItemChange: (itemId: ManagementRouteId) => void;
}

export function NavbarSegmented({ activeItem, debugMode, onActiveItemChange }: NavbarSegmentedProps) {
  const [englishExpanded, setEnglishExpanded] = useState(
    () => isEnglishLearningRoute(activeItem) || activeItem === '',
  );
  const [obsidianExpanded, setObsidianExpanded] = useState(() => isObsidianRoute(activeItem));
  const [debugExpanded, setDebugExpanded] = useState(true);

  const renderLink = (item: NavItem) => (
    <NavLink
      key={item.id}
      component="button"
      type="button"
      className={classes.link}
      label={item.label}
      leftSection={<item.icon size={18} stroke={1.5} />}
      active={item.id === activeItem}
      onClick={() => onActiveItemChange(item.id)}
    />
  );

  return (
    <nav className={classes.navbar}>
      <div className={classes.navbarMain}>
        <Text fw={600} size="sm" className={classes.title}>
          Talus Echo
        </Text>
        <Text size="xs" c="dimmed" mt={4}>
          Manage
        </Text>

        <Stack gap={4} className={classes.links}>
          <NavLink
            component="button"
            type="button"
            className={classes.sectionHeader}
            label="English Learning"
            leftSection={<IconLanguage size={18} stroke={1.5} />}
            opened={englishExpanded}
            onChange={setEnglishExpanded}
            childrenOffset={16}
          >
            {englishLearningItems.map((item) => renderLink(item))}
          </NavLink>

          <NavLink
            component="button"
            type="button"
            className={classes.sectionHeader}
            label="Obsidian"
            leftSection={<IconMarkdown size={18} stroke={1.5} />}
            opened={obsidianExpanded}
            onChange={setObsidianExpanded}
            childrenOffset={16}
          >
            {obsidianItems.map((item) => renderLink(item))}
          </NavLink>

          {debugMode && (
            <NavLink
              component="button"
              type="button"
              className={classes.sectionHeader}
              label="Debug"
              leftSection={<IconBug size={18} stroke={1.5} />}
              opened={debugExpanded}
              onChange={setDebugExpanded}
              childrenOffset={16}
            >
              {debugItems.map((item) => renderLink(item))}
            </NavLink>
          )}
        </Stack>
      </div>

      <Stack gap={4} className={classes.footer}>
        <NavLink
          component="button"
          type="button"
          className={classes.link}
          label="Agent Chat"
          leftSection={<IconMessageChatbot size={18} stroke={1.5} />}
          onClick={() => {
            SystemService.ShowAgentWindow().catch((err: unknown) => {
              console.error(err);
            });
          }}
        />

        <NavLink
          component="button"
          type="button"
          className={classes.link}
          label="Config"
          leftSection={<IconAdjustments size={18} stroke={1.5} />}
          active={activeItem === ManagementRoute.Config}
          onClick={() => onActiveItemChange(ManagementRoute.Config)}
        />

        <NavLink
          component="button"
          type="button"
          className={classes.link}
          label="About"
          leftSection={<IconInfoCircle size={18} stroke={1.5} />}
          active={activeItem === ManagementRoute.About}
          onClick={() => onActiveItemChange(ManagementRoute.About)}
        />
      </Stack>
    </nav>
  );
}
