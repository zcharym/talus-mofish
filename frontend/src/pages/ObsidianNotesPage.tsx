import { useCallback, useEffect, useRef, useState } from "react";
import {
  Badge,
  Button,
  Group,
  Paper,
  ScrollArea,
  SegmentedControl,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { ObsidianService } from "../../bindings/github.com/songwei.ma/talus-mofish/backend/services";
import { Note } from "../../bindings/github.com/songwei.ma/talus-mofish/backend/obsidian/models";
import { VaultTree } from "../components/management/VaultTree";
import { useDynamicScrollHeight } from "../hooks/useDynamicScrollHeight";
import { notify } from "../services/notifications";

function errorMessage(err: unknown): string {
  if (typeof err === "string") {
    return err;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return String(err);
}

export interface ObsidianNotesPageProps {
  focusPath?: string | null;
  onFocusConsumed?: () => void;
}

export function ObsidianNotesPage({ focusPath, onFocusConsumed }: ObsidianNotesPageProps) {
  const [selectedPath, setSelectedPath] = useState<string | null>(null);
  const [note, setNote] = useState<Note | null>(null);
  const [draft, setDraft] = useState("");
  const [preview, setPreview] = useState<"edit" | "preview">("edit");
  const [notMarkdown, setNotMarkdown] = useState(false);
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const dirty = note !== null && draft !== (note.content ?? "");

  const editorAnchorRef = useRef<HTMLDivElement>(null);
  const editorFooterRef = useRef<HTMLDivElement>(null);
  const editorHeight = useDynamicScrollHeight(editorAnchorRef, editorFooterRef, [selectedPath, preview]);
  const [treeFocus, setTreeFocus] = useState<string | null>(null);
  const dirtyRef = useRef(false);
  dirtyRef.current = dirty;

  const confirmLeave = useCallback(() => {
    if (!dirtyRef.current) {
      return true;
    }
    return window.confirm("This note has unsaved changes. Discard them?");
  }, []);

  const loadNote = useCallback(async (path: string) => {
    setLoading(true);
    setNotMarkdown(false);
    try {
      const loaded = (await ObsidianService.ReadNote(path)) as Note;
      setSelectedPath(path);
      setNote(loaded);
      setDraft(loaded.content ?? "");
      setPreview("edit");
    } catch (err) {
      notify.failed("Obsidian", errorMessage(err));
    } finally {
      setLoading(false);
    }
  }, []);

  const handleSelectFile = useCallback(
    (path: string, isMarkdown: boolean) => {
      if (path === selectedPath) {
        return;
      }
      if (!confirmLeave()) {
        return;
      }
      if (!isMarkdown) {
        setSelectedPath(path);
        setNote(null);
        setDraft("");
        setNotMarkdown(true);
        return;
      }
      void loadNote(path);
    },
    [confirmLeave, loadNote, selectedPath],
  );

  useEffect(() => {
    if (!focusPath) {
      return;
    }
    if (!confirmLeave()) {
      onFocusConsumed?.();
      return;
    }
    setTreeFocus(focusPath);
    void loadNote(focusPath);
    onFocusConsumed?.();
  }, [confirmLeave, focusPath, loadNote, onFocusConsumed]);

  const save = async () => {
    if (!selectedPath || !note) {
      return;
    }
    setSaving(true);
    try {
      await ObsidianService.WriteNote(selectedPath, draft);
      setNote({ ...note, content: draft });
      notify.success("Obsidian", "Note saved.");
    } catch (err) {
      notify.failed("Obsidian", errorMessage(err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <>
      <div ref={editorAnchorRef} />
      <Group align="stretch" gap="md" wrap="nowrap" mt="sm">
        <Paper withBorder p="xs" w={280} miw={220} style={{ flexShrink: 0 }}>
          <Text size="xs" c="dimmed" mb="xs">
            Vault
          </Text>
          <ScrollArea h={editorHeight} offsetScrollbars>
            <VaultTree
              selectedPath={selectedPath}
              expandToPath={treeFocus}
              onSelectFile={handleSelectFile}
            />
          </ScrollArea>
        </Paper>

        <Stack gap="sm" style={{ flex: 1, minWidth: 0 }}>
          {selectedPath ? (
            <>
              <Group justify="space-between" align="flex-start">
                <Stack gap={4} style={{ minWidth: 0, flex: 1 }}>
                  <Text fw={600} lineClamp={1}>
                    {selectedPath}
                  </Text>
                  {note?.tags && note.tags.length > 0 ? (
                    <Group gap={4}>
                      {note.tags.map((tag) => (
                        <Badge key={tag} size="sm" variant="light">
                          #{tag}
                        </Badge>
                      ))}
                    </Group>
                  ) : null}
                </Stack>
                {!notMarkdown ? (
                  <Group gap="xs">
                    <SegmentedControl
                      size="xs"
                      value={preview}
                      onChange={(value) => setPreview(value as "edit" | "preview")}
                      data={[
                        { label: "Edit", value: "edit" },
                        { label: "Preview", value: "preview" },
                      ]}
                    />
                    <Button size="xs" onClick={() => void save()} loading={saving} disabled={!dirty}>
                      Save
                    </Button>
                  </Group>
                ) : null}
              </Group>
              {notMarkdown ? (
                <Text c="dimmed">This file is not a markdown note.</Text>
              ) : loading ? (
                <Text c="dimmed">Loading note…</Text>
              ) : preview === "preview" ? (
                <Paper withBorder p="md" h={editorHeight} style={{ overflow: "auto" }}>
                  <ReactMarkdown remarkPlugins={[remarkGfm]}>{draft || "_Empty note_"}</ReactMarkdown>
                </Paper>
              ) : (
                <Textarea
                  autosize={false}
                  minRows={12}
                  styles={{ input: { height: editorHeight, fontFamily: "var(--mantine-font-family-monospace)" } }}
                  value={draft}
                  onChange={(event) => setDraft(event.currentTarget.value)}
                />
              )}
            </>
          ) : (
            <Text c="dimmed" mt="sm">
              Select a markdown note from the vault. Empty folders are hidden by the Local REST API.
              Obsidian must be running.
            </Text>
          )}
        </Stack>
      </Group>
    </>
  );
}
