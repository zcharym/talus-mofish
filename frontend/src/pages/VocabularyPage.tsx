import { useCallback, useEffect, useState } from "react";
import { Badge, Paper, ScrollArea, Table, Text } from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { Vocabulary } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { notify } from "../services/notifications";

export function VocabularyPage() {
  const [items, setItems] = useState<Vocabulary[]>([]);
  const [loading, setLoading] = useState(true);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const rows = await AppService.ListVocabulary();
      setItems(rows ?? []);
    } catch (err) {
      console.error(err);
      notify.failed("Vocabulary", "Failed to load vocabulary.");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  if (loading) {
    return <Text c="dimmed" mt="md">Loading vocabulary…</Text>;
  }

  if (items.length === 0) {
    return (
      <Text c="dimmed" mt="md">
        Your vocabulary bank is empty. Import an Anki deck or add words from Reading.
      </Text>
    );
  }

  return (
    <Paper withBorder mt="md">
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
            {items.map((row) => (
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
            ))}
          </Table.Tbody>
        </Table>
      </ScrollArea>
    </Paper>
  );
}
