import { useCallback, useEffect, useRef, useState } from "react";
import {
  Badge,
  Box,
  EmptyState,
  Group,
  LoadingOverlay,
  Modal,
  Pagination,
  Paper,
  ScrollArea,
  Skeleton,
  Stack,
  Text,
} from "@mantine/core";
import { IconBook } from "@tabler/icons-react";
import { EnglishService, Article, ArticlePageResult, ArticleSummary, toApiError } from "../utils/api";
import { FlipCard } from "../components/management/FlipCard";
import { SafeHTML } from "../components/SafeHTML";
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
      const result = await EnglishService.ListArticlesPage(pageNum, PAGE_SIZE);
      setPageResult(result);
      setPage(result.page);
    } catch (err) {
      console.error(err);
      notify.failed("Reading", toApiError(err));
    } finally {
      setLoadingList(false);
    }
  }, []);

  const loadArticle = useCallback(async (id: string) => {
    setLoadingArticle(true);
    try {
      const article = await EnglishService.GetArticle(id);
      setSelectedArticle(article);
      setModalOpen(true);
    } catch (err) {
      console.error(err);
      notify.failed("Reading", toApiError(err));
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
    return (
      <Stack mt="md" gap="sm">
        <Skeleton height={56} radius="md" />
        <Skeleton height={56} radius="md" />
        <Skeleton height={56} radius="md" />
      </Stack>
    );
  }

  if (pageResult?.total === 0) {
    return (
      <EmptyState
        mt="md"
        align="left"
        withIndicatorBackground
        icon={<IconBook size={28} />}
        title="No articles yet"
        description="Import reading material from an Anki deck on the Import tab."
      />
    );
  }

  return (
    <Stack mt="md" gap="md">
      <Paper withBorder p="xs" ref={scrollAnchorRef}>
        <ScrollArea h={scrollHeight} type="auto">
          <Stack gap="xs">
            {loadingList ? (
              <>
                <Skeleton height={56} radius="md" />
                <Skeleton height={56} radius="md" />
                <Skeleton height={56} radius="md" />
              </>
            ) : (
              items.map((article) => (
                <Paper
                  key={article.id}
                  withBorder
                  p="sm"
                  style={{
                    cursor: "pointer",
                    borderColor: article.id === selectedId ? "var(--mantine-color-persimmon-5)" : undefined,
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
        <Box pos="relative" mih={120}>
          <LoadingOverlay visible={loadingArticle} zIndex={10} overlayProps={{ radius: "sm", blur: 1 }} />
          {selectedArticle ? (
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
                <SafeHTML className="card" html={selectedArticle.content} />
              }
              back={
                selectedArticle.translation ? (
                  <SafeHTML html={selectedArticle.translation} />
                ) : (
                  <Text c="dimmed" size="sm">No translation available.</Text>
                )
              }
            />
          ) : null}
        </Box>
      </Modal>
    </Stack>
  );
}
