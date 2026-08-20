import { useCallback, useEffect, useState } from "react";
import { Loader, Text, UnstyledButton } from "@mantine/core";
import { IconChevronDown, IconChevronRight, IconFile, IconFileText, IconFolder } from "@tabler/icons-react";
import { ObsidianService } from "../../../../bindings/github.com/songwei.ma/talus-mofish/backend/services";
import { FileEntry } from "../../../../bindings/github.com/songwei.ma/talus-mofish/backend/obsidian/models";
import classes from "./VaultTree.module.css";

export function isMarkdownPath(path: string): boolean {
  return path.toLowerCase().endsWith(".md");
}

function errorMessage(err: unknown): string {
  if (typeof err === "string") {
    return err;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return String(err);
}

interface DirState {
  expanded: boolean;
  loading: boolean;
  entries: FileEntry[] | null;
  error: string | null;
}

interface VaultTreeProps {
  selectedPath: string | null;
  expandToPath?: string | null;
  onSelectFile: (path: string, isMarkdown: boolean) => void;
}

const ROOT_KEY = "";

export function VaultTree({ selectedPath, expandToPath, onSelectFile }: VaultTreeProps) {
  const [dirs, setDirs] = useState<Record<string, DirState>>({});

  const loadDir = useCallback(async (dirPath: string) => {
    setDirs((prev) => ({
      ...prev,
      [dirPath]: {
        expanded: true,
        loading: true,
        entries: prev[dirPath]?.entries ?? null,
        error: null,
      },
    }));
    try {
      const entries = (await ObsidianService.ListDirectory(dirPath)) as FileEntry[];
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { expanded: true, loading: false, entries, error: null },
      }));
    } catch (err) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: {
          expanded: true,
          loading: false,
          entries: prev[dirPath]?.entries ?? null,
          error: errorMessage(err),
        },
      }));
    }
  }, []);

  useEffect(() => {
    void loadDir(ROOT_KEY);
  }, [loadDir]);

  useEffect(() => {
    if (!expandToPath) {
      return;
    }
    const parts = expandToPath.split("/").filter(Boolean);
    let cancelled = false;

    const expandParents = async () => {
      await loadDir(ROOT_KEY);
      let acc = "";
      for (let i = 0; i < parts.length - 1; i += 1) {
        if (cancelled) {
          return;
        }
        acc = acc ? `${acc}/${parts[i]}` : parts[i];
        await loadDir(acc);
      }
    };

    void expandParents();
    return () => {
      cancelled = true;
    };
  }, [expandToPath, loadDir]);

  const toggleDir = (dirPath: string) => {
    const current = dirs[dirPath];
    if (current?.expanded) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { ...current, expanded: false },
      }));
      return;
    }
    if (current?.entries) {
      setDirs((prev) => ({
        ...prev,
        [dirPath]: { ...current, expanded: true },
      }));
      return;
    }
    void loadDir(dirPath);
  };

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

    return (state.entries ?? []).map((entry) => {
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

  return <div className={classes.tree}>{renderEntries(ROOT_KEY, 0)}</div>;
}
