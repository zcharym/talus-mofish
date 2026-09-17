import { useCallback, useEffect, useState } from "react";
import {
  Badge,
  Box,
  Button,
  Code,
  EmptyState,
  Group,
  LoadingOverlay,
  Paper,
  Select,
  Stack,
  Stepper,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { IconUpload } from "@tabler/icons-react";
import { EnglishService, SystemService, AnkiDeckPreview, AnkiPreview, ImportDeckConfig, ImportResult, AnkiImportRecord, toApiError } from "../utils/api";
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
  const [history, setHistory] = useState<AnkiImportRecord[]>([]);
  const [loading, setLoading] = useState(false);
  const [importing, setImporting] = useState(false);
  const [lastResult, setLastResult] = useState<ImportResult | null>(null);
  const [activeStep, setActiveStep] = useState(0);

  const loadHistory = useCallback(async () => {
    try {
      const items = await EnglishService.ListAnkiImports();
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
      const path = await SystemService.PickAnkiAPKG();
      if (!path) {
        return;
      }
      setApkgPath(path);
      setLastResult(null);
      const data = await EnglishService.PreviewAnkiAPKG(path);
      setPreview(data);
      initDeckConfigs(data.decks ?? []);
      setActiveStep(1);
    } catch (err) {
      console.error(err);
      notify.failed("Import", toApiError(err));
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
      const result = await EnglishService.ImportAnkiAPKG(apkgPath, configs);
      setLastResult(result);
      setActiveStep(2);
      notify.success(
        "Import complete",
        `Vocabulary: ${result.stats?.vocabCreated ?? 0}, cards: ${result.stats?.cardsCreated ?? 0}, articles: ${result.stats?.articlesCreated ?? 0}`,
      );
      await loadHistory();
    } catch (err) {
      console.error(err);
      notify.failed("Import", toApiError(err));
    } finally {
      setImporting(false);
    }
  };

  const startOver = () => {
    setApkgPath("");
    setPreview(null);
    setDeckConfigs({});
    setLastResult(null);
    setActiveStep(0);
  };

  return (
    <Box pos="relative" mt="md">
      <LoadingOverlay visible={loading || importing} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
      <Stack gap="lg">
        <Stepper active={activeStep} onStepClick={setActiveStep} allowNextStepsSelect={false}>
          <Stepper.Step label="Choose file" description="Pick an .apkg">
            <Paper withBorder p="md" mt="md">
              <Stack gap="sm">
                <Text size="sm" c="dimmed">
                  Import Anki decks (.apkg) as vocabulary cards or reading articles. Scheduling is reset to new cards.
                </Text>
                <Group>
                  <Button onClick={() => void handlePickFile()} loading={loading} leftSection={<IconUpload size={16} />}>
                    Choose APKG file
                  </Button>
                  {apkgPath && <Code>{apkgPath}</Code>}
                </Group>
              </Stack>
            </Paper>
          </Stepper.Step>

          <Stepper.Step label="Map fields" description="Configure decks">
            {preview ? (
              <Paper withBorder p="md" mt="md">
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

                <Group mt="md">
                  <Button variant="default" onClick={startOver}>
                    Choose another file
                  </Button>
                  <Button onClick={() => void handleImport()} loading={importing}>
                    Import selected decks
                  </Button>
                </Group>
              </Paper>
            ) : (
              <EmptyState
                mt="md"
                align="left"
                withIndicatorBackground
                icon={<IconUpload size={28} />}
                title="No file selected"
                description="Choose an APKG file first."
              >
                <EmptyState.Actions>
                  <Button onClick={() => setActiveStep(0)}>Back</Button>
                </EmptyState.Actions>
              </EmptyState>
            )}
          </Stepper.Step>

          <Stepper.Step label="Done" description="Review results">
            <Paper withBorder p="md" mt="md">
              {lastResult ? (
                <>
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
                  <Button mt="md" variant="light" onClick={startOver}>
                    Import another deck
                  </Button>
                </>
              ) : (
                <EmptyState
                  align="left"
                  withIndicatorBackground
                  icon={<IconUpload size={28} />}
                  title="No import yet"
                  description="Finish mapping fields and run the import."
                >
                  <EmptyState.Actions>
                    <Button onClick={() => setActiveStep(preview ? 1 : 0)}>
                      {preview ? "Back to mapping" : "Choose a file"}
                    </Button>
                  </EmptyState.Actions>
                </EmptyState>
              )}
            </Paper>
          </Stepper.Step>
        </Stepper>

        <Paper withBorder p="md">
          <Title order={4}>Import history</Title>
          {history.length === 0 ? (
            <EmptyState
              mt="sm"
              size="sm"
              align="left"
              withIndicatorBackground
              icon={<IconUpload size={22} />}
              title="No imports yet"
              description="Imported decks will appear here."
            />
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
    </Box>
  );
}
