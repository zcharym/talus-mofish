import { useCallback, useEffect, useState } from "react";
import { Badge, Center, Group, Loader, Pagination, Paper, Stack, Text, Title } from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import {
  Article,
  ArticlePageResult,
  ArticleSummary,
} from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { notify } from "../services/notifications";

const PAGE_SIZE = 20;

export function ReadingPage() {
  const [pageResult, setPageResult] = useState<ArticlePageResult | null>(null);
  const [page, setPage] = useState(1);
  const [loadingList, setLoadingList] = useState(true);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [loadingArticle, setLoadingArticle] = useState(false);

  const loadPage = useCallback(async (pageNum: number) => {
    setLoadingList(true);
    try {
      const result = await AppService.ListArticlesPage(pageNum, PAGE_SIZE);
      setPageResult(result);
      setPage(result.page);
    } catch (err) {
      console.error(err);
      notify.failed("Reading", "Failed to load articles.");
    } finally {
      setLoadingList(false);
    }
  }, []);

  const loadArticle = useCallback(async (id: string) => {
    setLoadingArticle(true);
    try {
      const article = await AppService.GetArticle(id);
      setSelectedArticle(article);
    } catch (err) {
      console.error(err);
      notify.failed("Reading", "Failed to load article.");
      setSelectedArticle(null);
    } finally {
      setLoadingArticle(false);
    }
  }, []);

  useEffect(() => {
    void loadPage(page);
  }, [loadPage, page]);

  useEffect(() => {
    const items = pageResult?.items ?? [];
    if (items.length === 0) {
      setSelectedId(null);
      setSelectedArticle(null);
      return;
    }
    const firstId = items[0].id;
    setSelectedId(firstId);
    void loadArticle(firstId);
  }, [pageResult, loadArticle]);

  const handleSelect = (item: ArticleSummary) => {
    if (item.id === selectedId) {
      return;
    }
    setSelectedId(item.id);
    void loadArticle(item.id);
  };

  const totalPages = pageResult ? Math.max(1, Math.ceil(pageResult.total / pageResult.page_size)) : 1;
  const items = pageResult?.items ?? [];

  if (loadingList && !pageResult) {
    return <Text c="dimmed" mt="md">Loading articles…</Text>;
  }

  if (pageResult?.total === 0) {
    return (
      <Text c="dimmed" mt="md">
        No articles yet. Import reading material from an Anki deck on the Import tab.
      </Text>
    );
  }

  return (
    <Stack mt="md" gap="md">
      <Stack gap="xs">
        {loadingList ? (
          <Center py="md">
            <Loader size="sm" />
          </Center>
        ) : (
          items.map((article) => (
            <Paper
              key={article.id}
              withBorder
              p="sm"
              style={{
                cursor: "pointer",
                borderColor: article.id === selectedId ? "var(--mantine-color-blue-5)" : undefined,
              }}
              onClick={() => handleSelect(article)}
            >
              <Group gap="xs">
                <Text fw={600}>{article.title}</Text>
                {article.source === "import:anki" && (
                  <Badge size="sm" variant="light">Anki</Badge>
                )}
                <Text size="xs" c="dimmed">{article.word_count} words</Text>
              </Group>
            </Paper>
          ))
        )}
      </Stack>

      {pageResult && pageResult.total > pageResult.page_size && (
        <Group justify="space-between" align="center">
          <Text size="sm" c="dimmed">
            {pageResult.total} articles
          </Text>
          <Pagination value={page} onChange={setPage} total={totalPages} />
        </Group>
      )}

      {loadingArticle ? (
        <Center py="xl">
          <Loader size="sm" />
        </Center>
      ) : selectedArticle ? (
        <Paper withBorder p="md">
          <Title order={4}>{selectedArticle.title}</Title>
          {selectedArticle.model_css && (
            <style>{selectedArticle.model_css}</style>
          )}
          <div className="card" dangerouslySetInnerHTML={{ __html: selectedArticle.content }} />
          {selectedArticle.translation && (
            <>
              <Title order={5} mt="lg">Translation</Title>
              <div dangerouslySetInnerHTML={{ __html: selectedArticle.translation }} />
            </>
          )}
        </Paper>
      ) : null}
    </Stack>
  );
}
