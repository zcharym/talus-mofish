import { useCallback, useEffect, useRef, useState } from "react";
import {
  Badge,
  Box,
  Button,
  EmptyState,
  Group,
  LoadingOverlay,
  Paper,
  ScrollArea,
  SegmentedControl,
  Splitter,
  Stack,
  Text,
  Textarea,
} from "@mantine/core";
import { IconFileText } from "@tabler/icons-react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { ObsidianService, Note, toApiError } from "../utils/api";
import { VaultTree } from "../components/management/VaultTree";
import { useDynamicScrollHeight } from "../hooks/useDynamicScrollHeight";
import { notify } from "../services/notifications";
import classes from "./ObsidianNotesPage.module.css";

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
      notify.failed("Obsidian", toApiError(err));
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
      notify.failed("Obsidian", toApiError(err));
    } finally {
      setSaving(false);
    }
  };

  return (
    <>
      <div ref={editorAnchorRef} />
      <Splitter className={classes.layout} h={editorHeight + 48}>
        <Splitter.Pane defaultSize="280px" min="200px" max="50%" className={classes.vaultPane}>
          <Paper withBorder p="xs" className={classes.vault}>
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
        </Splitter.Pane>

        <Splitter.Pane defaultSize={100} min="40%" className={classes.editorPane}>
          <Box pos="relative" className={classes.editor}>
            <LoadingOverlay visible={loading || saving} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
            {selectedPath ? (
              <Stack gap="sm">
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
              </Stack>
            ) : (
              <EmptyState
                mt="sm"
                align="left"
                withIndicatorBackground
                icon={<IconFileText size={28} />}
                title="Select a note"
                description="Pick a markdown note from the vault. Empty folders are hidden by the Local REST API. Obsidian must be running."
              />
            )}
          </Box>
        </Splitter.Pane>
      </Splitter>
      <div ref={editorFooterRef} />
    </>
  );
}
