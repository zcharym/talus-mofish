import { useCallback, useEffect, useRef, useState } from "react";
import {
  Badge,
  Center,
  Group,
  Loader,
  Modal,
  Pagination,
  Paper,
  ScrollArea,
  Stack,
  Table,
  Text,
  Title,
} from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { Vocabulary, VocabularyPageResult } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { FlipCard } from "../components/FlipCard";
import { useDynamicScrollHeight } from "../hooks/useDynamicScrollHeight";
import { notify } from "../services/notifications";

const PAGE_SIZE = 10;

export function VocabularyPage() {
  const [pageResult, setPageResult] = useState<VocabularyPageResult | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [selectedVocab, setSelectedVocab] = useState<Vocabulary | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const scrollAnchorRef = useRef<HTMLDivElement>(null);
  const scrollFooterRef = useRef<HTMLDivElement>(null);

  const load = useCallback(async (pageNum: number) => {
    setLoading(true);
    try {
      const result = await AppService.ListVocabularyPage(pageNum, PAGE_SIZE);
      setPageResult(result);
      setPage(result.page);
    } catch (err) {
      console.error(err);
      notify.failed("Vocabulary", "Failed to load vocabulary.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load(page);
  }, [load, page]);

  useEffect(() => {
    setModalOpen(false);
    setSelectedVocab(null);
  }, [page]);

  const handleClose = () => {
    setModalOpen(false);
    setSelectedVocab(null);
  };

  const handleSelect = (row: Vocabulary) => {
    setSelectedVocab(row);
    setModalOpen(true);
  };

  const items: Vocabulary[] = pageResult?.items ?? [];
  const totalPages = pageResult ? Math.max(1, Math.ceil(pageResult.total / pageResult.page_size)) : 1;
  const scrollHeight = useDynamicScrollHeight(scrollAnchorRef, scrollFooterRef, [
    loading,
    pageResult?.total,
    items.length,
  ]);

  if (loading && !pageResult) {
    return <Text c="dimmed" mt="md">Loading vocabulary…</Text>;
  }

  if (pageResult?.total === 0) {
    return (
      <Text c="dimmed" mt="md">
        Your vocabulary bank is empty. Import an Anki deck or add words from Reading.
      </Text>
    );
  }

  return (
    <Stack mt="md" gap="md">
      <Paper withBorder ref={scrollAnchorRef}>
        <ScrollArea h={scrollHeight} type="auto">
          <Table striped highlightOnHover>
            <Table.Thead>
              <Table.Tr>
                <Table.Th>Word</Table.Th>
                <Table.Th>Phonetic</Table.Th>
                <Table.Th>Definition</Table.Th>
                <Table.Th>Source</Table.Th>
              </Table.Tr>
            </Table.Thead>
            <Table.Tbody>
              {loading ? (
                <Table.Tr>
                  <Table.Td colSpan={4}>
                    <Center py="md">
                      <Loader size="sm" />
                    </Center>
                  </Table.Td>
                </Table.Tr>
              ) : (
                items.map((row) => (
                  <Table.Tr
                    key={row.id}
                    onClick={() => handleSelect(row)}
                    style={{
                      cursor: "pointer",
                      backgroundColor:
                        selectedVocab?.id === row.id && modalOpen
                          ? "var(--mantine-color-blue-light)"
                          : undefined,
                    }}
                  >
                    <Table.Td fw={500}>{row.word}</Table.Td>
                    <Table.Td>{row.phonetic}</Table.Td>
                    <Table.Td>
                      <Text size="sm" lineClamp={2}>{row.definition}</Text>
                    </Table.Td>
                    <Table.Td>
                      {row.source === "import:anki" ? (
                        <Badge size="sm" variant="light">Anki</Badge>
                      ) : (
                        row.source
                      )}
                    </Table.Td>
                  </Table.Tr>
                ))
              )}
            </Table.Tbody>
          </Table>
        </ScrollArea>
      </Paper>

      <div ref={scrollFooterRef}>
        {pageResult && pageResult.total > pageResult.page_size && (
          <Group justify="space-between" align="center">
            <Text size="sm" c="dimmed">
              {pageResult.total} words
            </Text>
            <Pagination value={page} onChange={setPage} total={totalPages} />
          </Group>
        )}
      </div>

      <Modal opened={modalOpen} onClose={handleClose} size="lg" title={null} padding="md">
        {selectedVocab && (
          <FlipCard
            key={selectedVocab.id}
            title={selectedVocab.word}
            headerExtra={
              <Group gap="xs">
                {selectedVocab.phonetic && (
                  <Text size="sm" c="dimmed">{selectedVocab.phonetic}</Text>
                )}
                {selectedVocab.pos && (
                  <Badge size="sm" variant="outline">{selectedVocab.pos}</Badge>
                )}
                {selectedVocab.source === "import:anki" ? (
                  <Badge size="sm" variant="light">Anki</Badge>
                ) : (
                  <Text size="xs" c="dimmed">{selectedVocab.source}</Text>
                )}
              </Group>
            }
            front={
              <Stack gap="xs" align="center" justify="center" h="100%">
                <Title order={2}>{selectedVocab.word}</Title>
                {selectedVocab.phonetic && (
                  <Text c="dimmed">{selectedVocab.phonetic}</Text>
                )}
                {selectedVocab.pos && (
                  <Badge variant="outline">{selectedVocab.pos}</Badge>
                )}
              </Stack>
            }
            back={
              <Stack gap="sm">
                <Text>{selectedVocab.definition}</Text>
                {selectedVocab.definition_en && (
                  <Text size="sm" c="dimmed">{selectedVocab.definition_en}</Text>
                )}
                {selectedVocab.examples && (
                  <Text size="sm" c="dimmed">{selectedVocab.examples}</Text>
                )}
              </Stack>
            }
          />
        )}
      </Modal>
    </Stack>
  );
}
