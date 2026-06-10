import { useCallback, useEffect, useState } from "react";
import { Badge, Group, Paper, Stack, Text, Title } from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import { Article } from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { notify } from "../services/notifications";

export function ReadingPage() {
  const [items, setItems] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedId, setSelectedId] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    try {
      const rows = await AppService.ListArticles();
      setItems(rows ?? []);
      if (rows?.length && !selectedId) {
        setSelectedId(rows[0].id);
      }
    } catch (err) {
      console.error(err);
      notify.failed("Reading", "Failed to load articles.");
    } finally {
      setLoading(false);
    }
  }, [selectedId]);

  useEffect(() => {
    void load();
  }, [load]);

  const selected = items.find((a) => a.id === selectedId);

  if (loading) {
    return <Text c="dimmed" mt="md">Loading articles…</Text>;
  }

  if (items.length === 0) {
    return (
      <Text c="dimmed" mt="md">
        No articles yet. Import reading material from an Anki deck on the Import tab.
      </Text>
    );
  }

  return (
    <Stack mt="md" gap="md">
      <Stack gap="xs">
        {items.map((article) => (
          <Paper
            key={article.id}
            withBorder
            p="sm"
            style={{
              cursor: "pointer",
              borderColor: article.id === selectedId ? "var(--mantine-color-blue-5)" : undefined,
            }}
            onClick={() => setSelectedId(article.id)}
          >
            <Group gap="xs">
              <Text fw={600}>{article.title}</Text>
              {article.source === "import:anki" && (
                <Badge size="sm" variant="light">Anki</Badge>
              )}
              <Text size="xs" c="dimmed">{article.word_count} words</Text>
            </Group>
          </Paper>
        ))}
      </Stack>

      {selected && (
        <Paper withBorder p="md">
          <Title order={4}>{selected.title}</Title>
          {selected.model_css && (
            <style>{selected.model_css}</style>
          )}
          <div className="card" dangerouslySetInnerHTML={{ __html: selected.content }} />
          {selected.translation && (
            <>
              <Title order={5} mt="lg">Translation</Title>
              <div dangerouslySetInnerHTML={{ __html: selected.translation }} />
            </>
          )}
        </Paper>
      )}
    </Stack>
  );
}
