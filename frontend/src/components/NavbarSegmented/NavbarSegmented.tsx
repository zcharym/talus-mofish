import { useEffect, useState } from 'react';
import {
  IconAdjustments,
  IconBook,
  IconBrain,
  IconClipboard,
  IconDatabase,
  IconHeadphones,
  IconInfoCircle,
  IconLanguage,
  IconSettings,
  IconTool,
  IconUpload,
  IconVocabulary,
} from '@tabler/icons-react';
import { Divider, SegmentedControl, Text } from '@mantine/core';
import classes from './NavbarSegmented.module.css';

type Section = 'tools' | 'english';

interface NavItem {
  id: string;
  label: string;
  icon: typeof IconTool;
}

const toolsItems: NavItem[] = [
  { id: 'import', label: 'Import', icon: IconUpload },
  { id: 'database', label: 'Database', icon: IconDatabase },
  { id: 'clipboard', label: 'Clipboard', icon: IconClipboard },
  { id: 'notes', label: 'Notes', icon: IconBook },
  { id: 'settings', label: 'Settings', icon: IconSettings },
];

const englishInteractiveItems: NavItem[] = [
  { id: 'recite', label: 'Recite Words', icon: IconBrain },
  { id: 'listening', label: 'Listening', icon: IconHeadphones },
  { id: 'grammar', label: 'Grammar', icon: IconLanguage },
];

const englishManagementItems: NavItem[] = [
  { id: 'reading', label: 'Reading', icon: IconBook },
  { id: 'vocabulary', label: 'Vocabulary', icon: IconVocabulary },
];

const englishItemIds = new Set([
  ...englishInteractiveItems.map((item) => item.id),
  ...englishManagementItems.map((item) => item.id),
]);

const toolsItemIds = new Set(toolsItems.map((item) => item.id));

function sectionForItem(itemId: string): Section | null {
  if (englishItemIds.has(itemId)) {
    return 'english';
  }
  if (toolsItemIds.has(itemId)) {
    return 'tools';
  }
  return null;
}

function defaultItemForSection(section: Section): string {
  return section === 'tools' ? toolsItems[0].id : englishInteractiveItems[0].id;
}

interface NavbarSegmentedProps {
  activeItem: string;
  onActiveItemChange: (itemId: string) => void;
}

export function NavbarSegmented({ activeItem, onActiveItemChange }: NavbarSegmentedProps) {
  const [section, setSection] = useState<Section>(() => sectionForItem(activeItem) ?? 'tools');

  useEffect(() => {
    const nextSection = sectionForItem(activeItem);
    if (nextSection) {
      setSection(nextSection);
    }
  }, [activeItem]);

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

  const sectionLinks =
    section === 'tools' ? (
      toolsItems.map(renderLink)
    ) : (
      <>
        {englishInteractiveItems.map(renderLink)}
        <Divider className={classes.divider} />
        {englishManagementItems.map(renderLink)}
      </>
    );

  return (
    <nav className={classes.navbar}>
      <div className={classes.navbarMain}>
        <Text fw={500} size="sm" className={classes.title} c="dimmed">
          Talus Echo
        </Text>

        <SegmentedControl
          value={section}
          onChange={(value) => {
            const nextSection = value as Section;
            setSection(nextSection);
            onActiveItemChange(defaultItemForSection(nextSection));
          }}
          transitionTimingFunction="ease"
          fullWidth
          mt="md"
          data={[
            { label: 'Tools', value: 'tools' },
            { label: 'English Study', value: 'english' },
          ]}
        />

        <div className={classes.links}>{sectionLinks}</div>
      </div>

      <div className={classes.footer}>
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
