import { useState } from 'react';
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
import { SegmentedControl, Text } from '@mantine/core';
import classes from './NavbarSegmented.module.css';

type Section = 'tools' | 'english';

interface NavItem {
  id: string;
  label: string;
  icon: typeof IconTool;
}

const tabs: Record<Section, NavItem[]> = {
  tools: [
    { id: 'import', label: 'Import', icon: IconUpload },
    { id: 'database', label: 'Database', icon: IconDatabase },
    { id: 'clipboard', label: 'Clipboard', icon: IconClipboard },
    { id: 'notes', label: 'Notes', icon: IconBook },
    { id: 'settings', label: 'Settings', icon: IconSettings },
  ],
  english: [
    { id: 'reading', label: 'Reading', icon: IconBook },
    { id: 'recite', label: 'Recite Words', icon: IconBrain },
    { id: 'vocabulary', label: 'Vocabulary', icon: IconVocabulary },
    { id: 'listening', label: 'Listening', icon: IconHeadphones },
    { id: 'grammar', label: 'Grammar', icon: IconLanguage },
  ],
};

interface NavbarSegmentedProps {
  activeItem: string;
  onActiveItemChange: (itemId: string) => void;
}

export function NavbarSegmented({ activeItem, onActiveItemChange }: NavbarSegmentedProps) {
  const [section, setSection] = useState<Section>('tools');

  const links = tabs[section].map((item) => (
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
  ));

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
            onActiveItemChange(tabs[nextSection][0].id);
          }}
          transitionTimingFunction="ease"
          fullWidth
          mt="md"
          data={[
            { label: 'Tools', value: 'tools' },
            { label: 'English Study', value: 'english' },
          ]}
        />

        <div style={{ marginTop: 'var(--mantine-spacing-md)' }}>{links}</div>
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
