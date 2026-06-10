import { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Button,
  Code,
  Group,
  Paper,
  Select,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import {
  AnkiDeckPreview,
  AnkiPreview,
  ImportDeckConfig,
  ImportResult,
} from "../../bindings/github.com/songwei.ma/talus-mofish/internal/content/models";
import { AnkiImport } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { notify } from "../services/notifications";

type TargetType = "vocabulary" | "reading" | "skip";

interface DeckConfigState {
  targetType: TargetType;
  fieldMapping: Record<string, number>;
}

const VOCAB_KEYS = ["word", "definition", "phonetic", "pos", "definition_en", "examples"] as const;
const READING_KEYS = ["title", "content", "translation"] as const;

function defaultVocabMapping(fields: string[]): Record<string, number> {
  const m: Record<string, number> = {};
  if (fields.length > 0) m.word = 0;
  if (fields.length > 1) m.definition = 1;
  return m;
}

function defaultReadingMapping(fields: string[]): Record<string, number> {
  const m: Record<string, number> = {};
  if (fields.length > 0) m.title = 0;
  if (fields.length > 1) m.content = 1;
  if (fields.length > 2) m.translation = 2;
  return m;
}

function fieldOptions(fields: string[]) {
  return fields.map((name, index) => ({
    value: String(index),
    label: `${index}: ${name}`,
  }));
}

export function ImportPage() {
  const [apkgPath, setApkgPath] = useState("");
  const [preview, setPreview] = useState<AnkiPreview | null>(null);
  const [deckConfigs, setDeckConfigs] = useState<Record<number, DeckConfigState>>({});
  const [history, setHistory] = useState<AnkiImport[]>([]);
  const [loading, setLoading] = useState(false);
  const [importing, setImporting] = useState(false);
  const [lastResult, setLastResult] = useState<ImportResult | null>(null);

  const loadHistory = useCallback(async () => {
    try {
      const items = await AppService.ListAnkiImports();
      setHistory(items ?? []);
    } catch (err) {
      console.error(err);
    }
  }, []);

  useEffect(() => {
    void loadHistory();
  }, [loadHistory]);

  const initDeckConfigs = (decks: AnkiDeckPreview[]) => {
    const next: Record<number, DeckConfigState> = {};
    for (const deck of decks) {
      next[deck.ankiDeckId] = {
        targetType: "vocabulary",
        fieldMapping: defaultVocabMapping(deck.fields ?? []),
      };
    }
    setDeckConfigs(next);
  };

  const handlePickFile = async () => {
    setLoading(true);
    try {
      const path = await AppService.PickAnkiAPKG();
      if (!path) {
        return;
      }
      setApkgPath(path);
      setLastResult(null);
      const data = await AppService.PreviewAnkiAPKG(path);
      setPreview(data);
      initDeckConfigs(data.decks ?? []);
    } catch (err) {
      console.error(err);
      notify.failed("Import", "Failed to open or preview the APKG file.");
    } finally {
      setLoading(false);
    }
  };

  const updateTarget = (deckId: number, deck: AnkiDeckPreview, targetType: TargetType) => {
    setDeckConfigs((prev) => {
      const mapping =
        targetType === "reading"
          ? defaultReadingMapping(deck.fields ?? [])
          : targetType === "vocabulary"
            ? defaultVocabMapping(deck.fields ?? [])
            : prev[deckId]?.fieldMapping ?? {};
      return {
        ...prev,
        [deckId]: { targetType, fieldMapping: mapping },
      };
    });
  };

  const updateMapping = (deckId: number, key: string, fieldIndex: string | null) => {
    if (!fieldIndex) return;
    setDeckConfigs((prev) => ({
      ...prev,
      [deckId]: {
        ...prev[deckId],
        fieldMapping: {
          ...prev[deckId].fieldMapping,
          [key]: Number(fieldIndex),
        },
      },
    }));
  };

  const handleImport = async () => {
    if (!apkgPath || !preview) return;
    setImporting(true);
    try {
      const configs: ImportDeckConfig[] = preview.decks.map((deck) => {
        const state = deckConfigs[deck.ankiDeckId];
        return new ImportDeckConfig({
          ankiDeckId: deck.ankiDeckId,
          ankiDeckName: deck.name,
          targetType: state?.targetType ?? "skip",
          ankiModelId: deck.ankiModelId,
          fieldMapping: state?.fieldMapping ?? {},
        });
      });
      const result = await AppService.ImportAnkiAPKG(apkgPath, configs);
      setLastResult(result);
      notify.success(
        "Import complete",
        `Vocabulary: ${result.stats?.vocabCreated ?? 0}, cards: ${result.stats?.cardsCreated ?? 0}, articles: ${result.stats?.articlesCreated ?? 0}`,
      );
      await loadHistory();
    } catch (err) {
      console.error(err);
      notify.failed("Import", "Failed to import the APKG file.");
    } finally {
      setImporting(false);
    }
  };

  return (
    <Stack mt="md" gap="lg">
      <Paper withBorder p="md">
        <Stack gap="sm">
          <Text size="sm" c="dimmed">
            Import Anki decks (.apkg) as vocabulary cards or reading articles. Scheduling is reset to new cards.
          </Text>
          <Group>
            <Button onClick={() => void handlePickFile()} loading={loading}>
              Choose APKG file
            </Button>
            {apkgPath && <Code>{apkgPath}</Code>}
          </Group>
        </Stack>
      </Paper>

      {preview && (
        <Paper withBorder p="md">
          <Title order={4}>Configure decks — {preview.filename}</Title>
          <Table mt="md" striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Deck</Table.Th>
                <Table.Th>Model</Table.Th>
                <Table.Th>Notes</Table.Th>
                <Table.Th>Cards</Table.Th>
                <Table.Th>Import as</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {preview.decks.map((deck) => {
                const state = deckConfigs[deck.ankiDeckId];
                return (
                  <Table.Tr key={deck.ankiDeckId}>
                    <Table.Td>{deck.name}</Table.Td>
                    <Table.Td>{deck.modelName}</Table.Td>
                    <Table.Td>{deck.noteCount}</Table.Td>
                    <Table.Td>{deck.cardCount}</Table.Td>
                    <Table.Td>
                      <Select
                        data={[
                          { value: "vocabulary", label: "Vocabulary + SRS cards" },
                          { value: "reading", label: "Reading article" },
                          { value: "skip", label: "Skip" },
                        ]}
                        value={state?.targetType ?? "vocabulary"}
                        onChange={(v) =>
                          updateTarget(deck.ankiDeckId, deck, (v as TargetType) ?? "skip")
                        }
                      />
                    </Table.Td>
                  </Table.Tr>
                );
              })}
            </Table.Tbody>
          </Table>

          {preview.decks
            .filter((d) => deckConfigs[d.ankiDeckId]?.targetType !== "skip")
            .map((deck) => {
              const state = deckConfigs[deck.ankiDeckId];
              const keys =
                state?.targetType === "reading" ? READING_KEYS : VOCAB_KEYS;
              const options = fieldOptions(deck.fields ?? []);
              return (
                <Paper key={deck.ankiDeckId} withBorder p="sm" mt="md">
                  <Text fw={500} size="sm">{deck.name} — field mapping</Text>
                  <Group mt="sm" grow>
                    {keys.map((key) => (
                      <Select
                        key={key}
                        label={key}
                        placeholder="Field"
                        data={options}
                        value={
                          state?.fieldMapping[key] !== undefined
                            ? String(state.fieldMapping[key])
                            : null
                        }
                        onChange={(v) => updateMapping(deck.ankiDeckId, key, v)}
                        clearable
                      />
                    ))}
                  </Group>
                </Paper>
              );
            })}

          <Button mt="md" onClick={() => void handleImport()} loading={importing}>
            Import selected decks
          </Button>
        </Paper>
      )}

      {lastResult && (
        <Paper withBorder p="md">
          <Title order={4}>Last import</Title>
          <Group mt="sm" gap="xs">
            <Badge color="green">{lastResult.status}</Badge>
            <Text size="sm">
              Vocabulary: {lastResult.stats?.vocabCreated ?? 0} · Cards:{" "}
              {lastResult.stats?.cardsCreated ?? 0} · Articles:{" "}
              {lastResult.stats?.articlesCreated ?? 0} · Skipped notes:{" "}
              {lastResult.stats?.skippedNotes ?? 0}
            </Text>
          </Group>
        </Paper>
      )}

      <Paper withBorder p="md">
        <Title order={4}>Import history</Title>
        {history.length === 0 ? (
          <Text c="dimmed" mt="sm" size="sm">No imports yet.</Text>
        ) : (
          <Table mt="sm" striped>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>File</Table.Th>
                <Table.Th>When</Table.Th>
                <Table.Th>Status</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {history.map((row) => (
                <Table.Tr key={row.id}>
                  <Table.Td>{row.filename}</Table.Td>
                  <Table.Td>{row.imported_at}</Table.Td>
                  <Table.Td>{row.status}</Table.Td>
                </Table.Tr>
              ))}
            </Table.Tbody>
          </Table>
        )}
      </Paper>
    </Stack>
  );
}
