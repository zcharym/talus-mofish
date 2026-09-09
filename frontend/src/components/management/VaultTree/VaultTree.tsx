import { Loader, Text, UnstyledButton } from "@mantine/core";
import { IconChevronDown, IconChevronRight, IconFile, IconFileText, IconFolder } from "@tabler/icons-react";
import type { FileEntry } from "../../../utils/api";
import { useVaultTree, VAULT_ROOT_KEY } from "../../../hooks/useVaultTree";
import classes from "./VaultTree.module.css";

export function isMarkdownPath(path: string): boolean {
  return path.toLowerCase().endsWith(".md");
}

interface VaultTreeProps {
  selectedPath: string | null;
  expandToPath?: string | null;
  onSelectFile: (path: string, isMarkdown: boolean) => void;
}

export function VaultTree({ selectedPath, expandToPath, onSelectFile }: VaultTreeProps) {
  const { dirs, toggleDir } = useVaultTree(expandToPath);

  const renderEntries = (dirPath: string, depth: number) => {
    const state = dirs[dirPath];
    if (!state) {
      return null;
    }
    if (state.loading && !state.entries) {
      return (
        <div className={classes.status} style={{ paddingLeft: 12 + depth * 16 }}>
          <Loader size="xs" />
        </div>
      );
    }
    if (state.error && !state.entries) {
      return (
        <Text size="xs" c="red" px="xs" py={4} style={{ paddingLeft: 12 + depth * 16 }}>
          {state.error}
        </Text>
      );
    }
    if (!state.expanded) {
      return null;
    }

    return (state.entries ?? []).map((entry: FileEntry) => {
      const childPath = entry.path;
      if (entry.isDir) {
        const child = dirs[childPath];
        const expanded = child?.expanded ?? false;
        return (
          <div key={childPath}>
            <UnstyledButton
              className={classes.row}
              style={{ paddingLeft: 8 + depth * 16 }}
              onClick={() => toggleDir(childPath)}
            >
              {expanded ? (
                <IconChevronDown className={classes.chevron} size={14} />
              ) : (
                <IconChevronRight className={classes.chevron} size={14} />
              )}
              <IconFolder className={classes.icon} size={16} stroke={1.5} />
              <span className={classes.label}>{entry.name}</span>
            </UnstyledButton>
            {renderEntries(childPath, depth + 1)}
          </div>
        );
      }

      const markdown = isMarkdownPath(childPath);
      const FileIcon = markdown ? IconFileText : IconFile;
      return (
        <UnstyledButton
          key={childPath}
          className={classes.row}
          data-active={childPath === selectedPath || undefined}
          style={{ paddingLeft: 8 + depth * 16 }}
          onClick={() => onSelectFile(childPath, markdown)}
        >
          <span className={classes.chevronSpacer} />
          <FileIcon className={classes.icon} size={16} stroke={1.5} />
          <span className={classes.label}>{entry.name}</span>
        </UnstyledButton>
      );
    });
  };

  return <div className={classes.tree}>{renderEntries(VAULT_ROOT_KEY, 0)}</div>;
}
