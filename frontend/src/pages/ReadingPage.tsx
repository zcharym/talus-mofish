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
  Text,
} from "@mantine/core";
import { AppService } from "../../bindings/github.com/songwei.ma/talus-mofish";
import {
  Article,
  ArticlePageResult,
  ArticleSummary,
} from "../../bindings/github.com/songwei.ma/talus-mofish/internal/store/models";
import { FlipCard } from "../components/management/FlipCard";
import { useDynamicScrollHeight } from "../hooks/useDynamicScrollHeight";
import { notify } from "../services/notifications";

const PAGE_SIZE = 10;

export function ReadingPage() {
  const [pageResult, setPageResult] = useState<ArticlePageResult | null>(null);
  const [page, setPage] = useState(1);
  const [loadingList, setLoadingList] = useState(true);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedArticle, setSelectedArticle] = useState<Article | null>(null);
  const [loadingArticle, setLoadingArticle] = useState(false);
  const [modalOpen, setModalOpen] = useState(false);
  const scrollAnchorRef = useRef<HTMLDivElement>(null);
  const scrollFooterRef = useRef<HTMLDivElement>(null);

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
      setModalOpen(true);
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
    setModalOpen(false);
    setSelectedId(null);
    setSelectedArticle(null);
  }, [page]);

  const handleClose = () => {
    setModalOpen(false);
    setSelectedId(null);
    setSelectedArticle(null);
  };

  const handleSelect = (item: ArticleSummary) => {
    setSelectedId(item.id);
    void loadArticle(item.id);
  };

  const totalPages = pageResult ? Math.max(1, Math.ceil(pageResult.total / pageResult.page_size)) : 1;
  const items = pageResult?.items ?? [];
  const scrollHeight = useDynamicScrollHeight(scrollAnchorRef, scrollFooterRef, [
    loadingList,
    pageResult?.total,
    items.length,
  ]);

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
      <Paper withBorder p="xs" ref={scrollAnchorRef}>
        <ScrollArea h={scrollHeight} type="auto">
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
        </ScrollArea>
      </Paper>

      <div ref={scrollFooterRef}>
        {pageResult && pageResult.total > pageResult.page_size && (
          <Group justify="space-between" align="center">
            <Text size="sm" c="dimmed">
              {pageResult.total} articles
            </Text>
            <Pagination value={page} onChange={setPage} total={totalPages} />
          </Group>
        )}
      </div>

      <Modal opened={modalOpen} onClose={handleClose} size="lg" title={null} padding="md">
        {loadingArticle ? (
          <Center py="xl">
            <Loader size="sm" />
          </Center>
        ) : selectedArticle ? (
          <FlipCard
            key={selectedArticle.id}
            title={selectedArticle.title}
            modelCss={selectedArticle.model_css}
            headerExtra={
              <Group gap="xs">
                {selectedArticle.source === "import:anki" && (
                  <Badge size="sm" variant="light">Anki</Badge>
                )}
                <Text size="xs" c="dimmed">{selectedArticle.word_count} words</Text>
              </Group>
            }
            front={
              <div className="card" dangerouslySetInnerHTML={{ __html: selectedArticle.content }} />
            }
            back={
              selectedArticle.translation ? (
                <div dangerouslySetInnerHTML={{ __html: selectedArticle.translation }} />
              ) : (
                <Text c="dimmed" size="sm">No translation available.</Text>
              )
            }
          />
        ) : null}
      </Modal>
    </Stack>
  );
}
