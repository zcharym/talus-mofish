import { useCallback, useEffect, useState } from "react";
import { Badge, Center, Group, Loader, Pagination, Paper, ScrollArea, Stack, Table, Text } from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { Vocabulary, VocabularyPageResult } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { notify } from "../services/notifications";

const PAGE_SIZE = 20;

export function VocabularyPage() {
  const [pageResult, setPageResult] = useState<VocabularyPageResult | null>(null);
  const [page, setPage] = useState(1);
  const [loading, setLoading] = useState(true);

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

  const items: Vocabulary[] = pageResult?.items ?? [];
  const totalPages = pageResult ? Math.max(1, Math.ceil(pageResult.total / pageResult.page_size)) : 1;

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
      <Paper withBorder>
        <ScrollArea>
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
                  <Table.Tr key={row.id}>
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

      {pageResult && pageResult.total > pageResult.page_size && (
        <Group justify="space-between" align="center">
          <Text size="sm" c="dimmed">
            {pageResult.total} words
          </Text>
          <Pagination value={page} onChange={setPage} total={totalPages} />
        </Group>
      )}
    </Stack>
  );
}
