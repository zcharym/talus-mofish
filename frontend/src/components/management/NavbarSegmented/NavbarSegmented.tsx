import {
  IconAdjustments,
  IconBook,
  IconBug,
  IconInfoCircle,
  IconMessageChatbot,
  IconUpload,
  IconVocabulary,
} from '@tabler/icons-react';
import { Text } from '@mantine/core';
import { AppService } from '../../../../bindings/github.com/songwei.ma/talus-mofish';
import classes from './NavbarSegmented.module.css';

interface NavItem {
  id: string;
  label: string;
  icon: typeof IconUpload;
}

const libraryItems: NavItem[] = [
  { id: 'import', label: 'Import', icon: IconUpload },
  { id: 'reading', label: 'Reading', icon: IconBook },
  { id: 'vocabulary', label: 'Vocabulary', icon: IconVocabulary },
  { id: 'debug', label: 'Debug', icon: IconBug },
];

interface NavbarSegmentedProps {
  activeItem: string;
  debugMode: boolean;
  onActiveItemChange: (itemId: string) => void;
}

export function NavbarSegmented({ activeItem, debugMode, onActiveItemChange }: NavbarSegmentedProps) {
  const visibleItems = debugMode
    ? libraryItems
    : libraryItems.filter((item) => item.id !== 'debug');

  const renderLink = (item: NavItem) => (
    <a
      className={classes.link}
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

        <div className={classes.links}>{visibleItems.map(renderLink)}</div>
      </div>

      <div className={classes.footer}>
        <a
          href="#"
          className={classes.link}
          onClick={(event) => {
            event.preventDefault();
            AppService.ShowAgentWindow().catch((err: unknown) => {
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
          data-active={activeItem === 'config' || undefined}
          onClick={(event) => {
            event.preventDefault();
            onActiveItemChange('config');
          }}
        >
          <IconAdjustments className={classes.linkIcon} stroke={1.5} />
          <span>Config</span>
        </a>

        <a
          href="#"
          className={classes.link}
          data-active={activeItem === 'about' || undefined}
          onClick={(event) => {
            event.preventDefault();
            onActiveItemChange('about');
          }}
        >
          <IconInfoCircle className={classes.linkIcon} stroke={1.5} />
          <span>About</span>
        </a>
      </div>
    </nav>
  );
}
