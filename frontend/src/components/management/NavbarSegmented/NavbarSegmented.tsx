import { useState } from 'react';
import {
  IconAdjustments,
  IconBook,
  IconBug,
  IconChevronDown,
  IconChevronRight,
  IconFileText,
  IconInfoCircle,
  IconLanguage,
  IconMarkdown,
  IconMessageChatbot,
  IconSearch,
  IconUpload,
  IconVocabulary,
} from '@tabler/icons-react';
import { Collapse, Text, UnstyledButton } from '@mantine/core';
import { SystemService } from '../../../../bindings/github.com/songwei.ma/talus-mofish/backend/services';
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

  const renderLink = (item: NavItem, nested = false) => (
    <a
      className={nested ? classes.nestedLink : classes.link}
      data-active={item.id === activeItem || undefined}
      href="#"
      key={item.id}
      onClick={(event) => {
        event.preventDefault();
        onActiveItemChange(item.id);
      }}
    >
      <item.icon className={classes.linkIcon} stroke={1.5} />
      <span>{item.label}</span>
    </a>
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

        <div className={classes.links}>
          <UnstyledButton
            className={classes.sectionHeader}
            onClick={() => setEnglishExpanded((expanded) => !expanded)}
          >
            <IconLanguage className={classes.linkIcon} stroke={1.5} />
            <span className={classes.sectionLabel}>English Learning</span>
            {englishExpanded ? (
              <IconChevronDown className={classes.chevron} size={16} />
            ) : (
              <IconChevronRight className={classes.chevron} size={16} />
            )}
          </UnstyledButton>

          <Collapse in={englishExpanded}>
            <div className={classes.nestedLinks}>
              {englishLearningItems.map((item) => renderLink(item, true))}
            </div>
          </Collapse>

          <UnstyledButton
            className={classes.sectionHeader}
            onClick={() => setObsidianExpanded((expanded) => !expanded)}
          >
            <IconMarkdown className={classes.linkIcon} stroke={1.5} />
            <span className={classes.sectionLabel}>Obsidian</span>
            {obsidianExpanded ? (
              <IconChevronDown className={classes.chevron} size={16} />
            ) : (
              <IconChevronRight className={classes.chevron} size={16} />
            )}
          </UnstyledButton>

          <Collapse in={obsidianExpanded}>
            <div className={classes.nestedLinks}>
              {obsidianItems.map((item) => renderLink(item, true))}
            </div>
          </Collapse>

          {debugMode && renderLink({ id: ManagementRoute.Debug, label: 'Debug', icon: IconBug })}
        </div>
      </div>

      <div className={classes.footer}>
        <a
          href="#"
          className={classes.link}
          onClick={(event) => {
            event.preventDefault();
            SystemService.ShowAgentWindow().catch((err: unknown) => {
              console.error(err);
            });
          }}
        >
          <IconMessageChatbot className={classes.linkIcon} stroke={1.5} />
          <span>Agent Chat</span>
        </a>

        <a
          href="#"
          className={classes.link}
          data-active={activeItem === ManagementRoute.Config || undefined}
          onClick={(event) => {
            event.preventDefault();
            onActiveItemChange(ManagementRoute.Config);
          }}
        >
          <IconAdjustments className={classes.linkIcon} stroke={1.5} />
          <span>Config</span>
        </a>

        <a
          href="#"
          className={classes.link}
          data-active={activeItem === ManagementRoute.About || undefined}
          onClick={(event) => {
            event.preventDefault();
            onActiveItemChange(ManagementRoute.About);
          }}
        >
          <IconInfoCircle className={classes.linkIcon} stroke={1.5} />
          <span>About</span>
        </a>
      </div>
    </nav>
  );
}
