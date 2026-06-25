import { useCallback, useEffect, useRef, useState } from "react";
import {
  Badge,
  Center,
  Group,
  Loader,
  Pagination,
  Paper,
  ScrollArea,
  Stack,
  Table,
  Text,
  TextInput,
} from "@mantine/core";
import { IconSearch } from "@tabler/icons-react";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { Vocabulary, VocabularyPageResult } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { VocabularyEditModal } from "../components/management/VocabularyEditModal";
import { useDynamicScrollHeight } from "../hooks/useDynamicScrollHeight";
import { notify } from "../services/notifications";

const PAGE_SIZE = 10;
const SEARCH_LIMIT = 50;

export function VocabularyPage() {
  const [pageResult, setPageResult] = useState<VocabularyPageResult | null>(null);
  const [searchResults, setSearchResults] = useState<Vocabulary[] | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);
  const [selectedVocabId, setSelectedVocabId] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const scrollAnchorRef = useRef<HTMLDivElement>(null);
  const scrollFooterRef = useRef<HTMLDivElement>(null);
  const searchAnchorRef = useRef<HTMLDivElement>(null);

  const isSearching = searchQuery.trim().length > 0;

  const loadPage = useCallback(async (pageNum: number) => {
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

  const runSearch = useCallback(async (query: string) => {
    const trimmed = query.trim();
    if (!trimmed) {
      setSearchResults(null);
      return;
    }
    setLoading(true);
    try {
      const results = await AppService.SearchVocabulary(trimmed, SEARCH_LIMIT);
      setSearchResults(results as Vocabulary[]);
    } catch (err) {
      console.error(err);
      notify.failed("Vocabulary", "Search failed.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (isSearching) {
      return;
    }
    void loadPage(page);
  }, [loadPage, page, isSearching]);

  useEffect(() => {
    if (!isSearching) {
      setSearchResults(null);
      return;
    }
    const timer = window.setTimeout(() => {
      void runSearch(searchQuery);
    }, 300);
    return () => window.clearTimeout(timer);
  }, [searchQuery, isSearching, runSearch]);

  useEffect(() => {
    setModalOpen(false);
    setSelectedVocabId(null);
  }, [page, searchQuery]);

  const handleClose = () => {
    setModalOpen(false);
    setSelectedVocabId(null);
  };

  const handleSelect = (row: Vocabulary) => {
    setSelectedVocabId(row.id);
    setModalOpen(true);
  };

  const handleRefresh = () => {
    if (isSearching) {
      void runSearch(searchQuery);
    } else {
      void loadPage(page);
    }
  };

  const items: Vocabulary[] = isSearching ? (searchResults ?? []) : (pageResult?.items ?? []);
  const totalPages = pageResult ? Math.max(1, Math.ceil(pageResult.total / pageResult.page_size)) : 1;
  const scrollHeight = useDynamicScrollHeight(scrollAnchorRef, scrollFooterRef, [
    loading,
    pageResult?.total,
    searchResults?.length,
    items.length,
    searchQuery,
  ]);

  const showEmpty = !loading && items.length === 0;

  return (
    <Stack mt="md" gap="md">
      <div ref={searchAnchorRef}>
        <TextInput
          placeholder="Search words or definitions…"
          leftSection={<IconSearch size={16} />}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.currentTarget.value)}
        />
      </div>

      {showEmpty && !isSearching && pageResult?.total === 0 ? (
        <Text c="dimmed">
          Your vocabulary bank is empty. Import an Anki deck or add words from Reading.
        </Text>
      ) : showEmpty ? (
        <Text c="dimmed">No matches found.</Text>
      ) : (
        <>
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
                            selectedVocabId === row.id && modalOpen
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
            {!isSearching && pageResult && pageResult.total > pageResult.page_size && (
              <Group justify="space-between" align="center">
                <Text size="sm" c="dimmed">
                  {pageResult.total} words
                </Text>
                <Pagination value={page} onChange={setPage} total={totalPages} />
              </Group>
            )}
            {isSearching && searchResults && (
              <Text size="sm" c="dimmed">
                {searchResults.length} result{searchResults.length === 1 ? "" : "s"}
              </Text>
            )}
          </div>
        </>
      )}

      <VocabularyEditModal
        vocabId={selectedVocabId}
        opened={modalOpen}
        onClose={handleClose}
        onSaved={handleRefresh}
        onDeleted={handleRefresh}
      />
    </Stack>
  );
}
