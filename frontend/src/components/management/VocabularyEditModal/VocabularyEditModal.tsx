import { useCallback, useEffect, useMemo, useState } from "react";
import {
  Badge,
  Button,
  Checkbox,
  Divider,
  Group,
  Modal,
  Select,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
} from "@mantine/core";
import { AppService } from "../../../../bindings/github.com/songwei.ma/talus-mofish";
import {
  Card,
  Deck,
  UpdateCardContentParams,
  UpdateVocabularyParams,
  Vocabulary,
} from "../../../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { FlipCard } from "../FlipCard";
import { notify } from "../../../services/notifications";

const LEVEL_OPTIONS = [
  { value: "custom", label: "Custom" },
  { value: "cet4", label: "CET-4" },
  { value: "cet6", label: "CET-6" },
  { value: "ielts", label: "IELTS" },
  { value: "toefl", label: "TOEFL" },
  { value: "gre", label: "GRE" },
];

interface CardDraft {
  id: string;
  deck_id: string;
  front: string;
  back: string;
  example_sentence: string;
  hints: string;
  card_type: string;
  model_css: string;
}

interface VocabularyEditModalProps {
  vocabId: string | null;
  opened: boolean;
  onClose: () => void;
  onSaved: () => void;
  onDeleted: () => void;
}

function cardToDraft(card: Card): CardDraft {
  return {
    id: card.id,
    deck_id: card.deck_id,
    front: card.front,
    back: card.back,
    example_sentence: card.example_sentence,
    hints: card.hints,
    card_type: card.card_type || "vocabulary",
    model_css: card.model_css,
  };
}

export function VocabularyEditModal({
  vocabId,
  opened,
  onClose,
  onSaved,
  onDeleted,
}: VocabularyEditModalProps) {
  const [mode, setMode] = useState<"view" | "edit">("view");
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [vocab, setVocab] = useState<Vocabulary | null>(null);
  const [cards, setCards] = useState<CardDraft[]>([]);
  const [decks, setDecks] = useState<Deck[]>([]);
  const [deleteCardsChecked, setDeleteCardsChecked] = useState(false);
  const [deleteConfirmOpen, setDeleteConfirmOpen] = useState(false);

  const deckNames = useMemo(() => {
    const map = new Map<string, string>();
    for (const deck of decks) {
      map.set(deck.id, deck.name);
    }
    return map;
  }, [decks]);

  const load = useCallback(async () => {
    if (!vocabId) {
      return;
    }
    setLoading(true);
    try {
      const [vocabItem, linkedCards, deckList] = await Promise.all([
        AppService.GetVocabulary(vocabId),
        AppService.ListCardsForVocab(vocabId),
        AppService.ListDecks(),
      ]);
      setVocab(vocabItem);
      setCards(linkedCards.map((card) => cardToDraft(card as Card)));
      setDecks(deckList as Deck[]);
      setMode("view");
      setDeleteCardsChecked(false);
    } catch (err) {
      notify.failed("Vocabulary", "Failed to load entry.");
      console.error(err);
      onClose();
    } finally {
      setLoading(false);
    }
  }, [vocabId, onClose]);

  useEffect(() => {
    if (opened && vocabId) {
      void load();
    }
  }, [opened, vocabId, load]);

  const updateVocabField = <K extends keyof Vocabulary>(key: K, value: Vocabulary[K]) => {
    setVocab((current) => (current ? new Vocabulary({ ...current, [key]: value }) : current));
  };

  const updateCardField = (index: number, key: keyof CardDraft, value: string) => {
    setCards((current) =>
      current.map((card, i) => (i === index ? { ...card, [key]: value } : card)),
    );
  };

  const handleSave = async () => {
    if (!vocab) {
      return;
    }
    setSaving(true);
    try {
      await AppService.UpdateVocabulary(
        new UpdateVocabularyParams({
          id: vocab.id,
          word: vocab.word,
          phonetic: vocab.phonetic,
          pos: vocab.pos,
          definition: vocab.definition,
          definition_en: vocab.definition_en,
          examples: vocab.examples,
          level: vocab.level,
          source: vocab.source,
        }),
      );

      for (const card of cards) {
        await AppService.UpdateCardContent(
          new UpdateCardContentParams({
            id: card.id,
            front: card.front,
            back: card.back,
            example_sentence: card.example_sentence,
            hints: card.hints,
            card_type: card.card_type,
            model_css: card.model_css,
          }),
        );
      }

      notify.success("Saved", "Vocabulary and linked cards updated.");
      setMode("view");
      onSaved();
      await load();
    } catch (err) {
      notify.failed("Save failed", String(err));
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!vocab) {
      return;
    }
    setSaving(true);
    try {
      if (deleteCardsChecked) {
        for (const card of cards) {
          await AppService.DeleteCard(card.id);
        }
      }
      await AppService.DeleteVocabulary(vocab.id);
      notify.success("Deleted", "Vocabulary entry removed.");
      setDeleteConfirmOpen(false);
      onDeleted();
      onClose();
    } catch (err) {
      notify.failed("Delete failed", String(err));
      console.error(err);
    } finally {
      setSaving(false);
    }
  };

  if (!vocab && loading) {
    return (
      <Modal opened={opened} onClose={onClose} title="Loading…" size="lg">
        <Text c="dimmed">Loading vocabulary…</Text>
      </Modal>
    );
  }

  if (!vocab) {
    return null;
  }

  return (
    <>
      <Modal
        opened={opened}
        onClose={onClose}
        size="xl"
        title={mode === "view" ? vocab.word : `Edit: ${vocab.word}`}
      >
        {mode === "view" ? (
          <Stack gap="md">
            <FlipCard
              title={vocab.word}
              headerExtra={
                <Group gap="xs">
                  {vocab.phonetic && <Text size="sm" c="dimmed">{vocab.phonetic}</Text>}
                  {vocab.pos && <Badge size="sm" variant="outline">{vocab.pos}</Badge>}
                  {vocab.source === "import:anki" ? (
                    <Badge size="sm" variant="light">Anki</Badge>
                  ) : (
                    <Text size="xs" c="dimmed">{vocab.source}</Text>
                  )}
                </Group>
              }
              front={
                <Stack gap="xs" align="center" justify="center" h="100%">
                  <Title order={2}>{vocab.word}</Title>
                  {vocab.phonetic && <Text c="dimmed">{vocab.phonetic}</Text>}
                  {vocab.pos && <Badge variant="outline">{vocab.pos}</Badge>}
                </Stack>
              }
              back={
                <Stack gap="sm">
                  <Text>{vocab.definition}</Text>
                  {vocab.definition_en && <Text size="sm" c="dimmed">{vocab.definition_en}</Text>}
                  {vocab.examples && <Text size="sm" c="dimmed">{vocab.examples}</Text>}
                </Stack>
              }
            />

            {cards.length > 0 && (
              <>
                <Divider label="Linked SRS cards" labelPosition="left" />
                {cards.map((card) => (
                  <FlipCard
                    key={card.id}
                    title={deckNames.get(card.deck_id) ?? "SRS Card"}
                    modelCss={card.model_css}
                    front={<div dangerouslySetInnerHTML={{ __html: card.front }} />}
                    back={<div dangerouslySetInnerHTML={{ __html: card.back }} />}
                  />
                ))}
              </>
            )}

            <Group justify="flex-end">
              <Button variant="default" onClick={onClose}>Close</Button>
              <Button onClick={() => setMode("edit")}>Edit</Button>
            </Group>
          </Stack>
        ) : (
          <Stack gap="md">
            <Text fw={500} size="sm">Vocabulary</Text>
            <Group grow align="flex-start">
              <TextInput
                label="Word"
                value={vocab.word}
                onChange={(e) => updateVocabField("word", e.currentTarget.value)}
                required
              />
              <TextInput
                label="Phonetic"
                value={vocab.phonetic}
                onChange={(e) => updateVocabField("phonetic", e.currentTarget.value)}
              />
            </Group>
            <Group grow align="flex-start">
              <TextInput
                label="Part of speech"
                value={vocab.pos}
                onChange={(e) => updateVocabField("pos", e.currentTarget.value)}
              />
              <Select
                label="Level"
                data={LEVEL_OPTIONS}
                value={vocab.level}
                onChange={(value) => updateVocabField("level", value ?? "custom")}
              />
            </Group>
            <Textarea
              label="Definition"
              value={vocab.definition}
              onChange={(e) => updateVocabField("definition", e.currentTarget.value)}
              minRows={2}
              autosize
            />
            <Textarea
              label="Definition (English)"
              value={vocab.definition_en}
              onChange={(e) => updateVocabField("definition_en", e.currentTarget.value)}
              minRows={2}
              autosize
            />
            <Textarea
              label="Examples"
              value={vocab.examples}
              onChange={(e) => updateVocabField("examples", e.currentTarget.value)}
              minRows={2}
              autosize
            />

            {cards.length > 0 && (
              <>
                <Divider label="Linked SRS cards" labelPosition="left" />
                {cards.map((card, index) => (
                  <Stack key={card.id} gap="xs" p="sm" style={{ border: "1px solid var(--mantine-color-gray-3)", borderRadius: 8 }}>
                    <Text size="sm" fw={500}>
                      {deckNames.get(card.deck_id) ?? "SRS Card"}
                    </Text>
                    <Textarea
                      label="Front (HTML ok)"
                      value={card.front}
                      onChange={(e) => updateCardField(index, "front", e.currentTarget.value)}
                      minRows={3}
                      autosize
                    />
                    <Textarea
                      label="Back (HTML ok)"
                      value={card.back}
                      onChange={(e) => updateCardField(index, "back", e.currentTarget.value)}
                      minRows={3}
                      autosize
                    />
                    <Textarea
                      label="Example sentence"
                      value={card.example_sentence}
                      onChange={(e) => updateCardField(index, "example_sentence", e.currentTarget.value)}
                      minRows={2}
                      autosize
                    />
                    <TextInput
                      label="Hints"
                      value={card.hints}
                      onChange={(e) => updateCardField(index, "hints", e.currentTarget.value)}
                    />
                    <FlipCard
                      title="Preview"
                      modelCss={card.model_css}
                      front={<div dangerouslySetInnerHTML={{ __html: card.front }} />}
                      back={<div dangerouslySetInnerHTML={{ __html: card.back }} />}
                    />
                  </Stack>
                ))}
              </>
            )}

            <Group justify="space-between">
              <Button color="red" variant="light" onClick={() => setDeleteConfirmOpen(true)}>
                Delete vocabulary
              </Button>
              <Group>
                <Button variant="default" onClick={() => setMode("view")}>Cancel</Button>
                <Button loading={saving} onClick={() => void handleSave()}>Save</Button>
              </Group>
            </Group>
          </Stack>
        )}
      </Modal>

      <Modal
        opened={deleteConfirmOpen}
        onClose={() => setDeleteConfirmOpen(false)}
        title="Delete vocabulary?"
        centered
      >
        <Stack gap="md">
          <Text size="sm">
            Delete &quot;{vocab.word}&quot;? This removes the vocabulary entry.
          </Text>
          {cards.length > 0 && (
            <Checkbox
              label={`Also delete linked cards (${cards.length})`}
              checked={deleteCardsChecked}
              onChange={(e) => setDeleteCardsChecked(e.currentTarget.checked)}
            />
          )}
          <Group justify="flex-end">
            <Button variant="default" onClick={() => setDeleteConfirmOpen(false)}>Cancel</Button>
            <Button color="red" loading={saving} onClick={() => void handleDelete()}>Delete</Button>
          </Group>
        </Stack>
      </Modal>
    </>
  );
}
