import { useCallback, useState } from "react";
import {
  Badge,
  Button,
  EmptyState,
  Group,
  Highlight,
  Paper,
  Stack,
  TextInput,
  UnstyledButton,
} from "@mantine/core";
import { IconSearch } from "@tabler/icons-react";
import { ObsidianService, SearchHit, toApiError } from "../utils/api";
import { notify } from "../services/notifications";

interface ObsidianSearchPageProps {
  onOpenNote: (path: string) => void;
}

export function ObsidianSearchPage({ onOpenNote }: ObsidianSearchPageProps) {
  const [query, setQuery] = useState("");
  const [hits, setHits] = useState<SearchHit[] | null>(null);
  const [loading, setLoading] = useState(false);
  const [submittedQuery, setSubmittedQuery] = useState("");

  const runSearch = useCallback(async () => {
    const trimmed = query.trim();
    if (!trimmed) {
      setHits(null);
      setSubmittedQuery("");
      return;
    }
    setLoading(true);
    try {
      const results = (await ObsidianService.SearchSimple(trimmed, 100)) as SearchHit[];
      setHits(results);
      setSubmittedQuery(trimmed);
    } catch (err) {
      notify.failed("Obsidian", toApiError(err));
    } finally {
      setLoading(false);
    }
  }, [query]);

  return (
    <Stack gap="md" mt="sm" maw={720}>
      <Group align="flex-end" wrap="wrap">
        <TextInput
          style={{ flex: 1 }}
          label="Search vault"
          placeholder="Full-text search"
          leftSection={<IconSearch size={16} />}
          value={query}
          onChange={(event) => setQuery(event.currentTarget.value)}
          onKeyDown={(event) => {
            if (event.key === "Enter") {
              event.preventDefault();
              void runSearch();
            }
          }}
        />
        <Button onClick={() => void runSearch()} loading={loading}>
          Search
        </Button>
      </Group>

      {hits && hits.length === 0 ? (
        <EmptyState
          align="left"
          withIndicatorBackground
          icon={<IconSearch size={28} />}
          title="No matching notes"
          description="Try different keywords or check that Obsidian is running."
        />
      ) : null}

      {hits?.map((hit) => (
        <UnstyledButton key={hit.filename} onClick={() => onOpenNote(hit.filename)} w="100%">
          <Paper withBorder p="sm">
            <Group justify="space-between" mb={4}>
              <Highlight highlight={submittedQuery} fw={600} size="sm">
                {hit.filename}
              </Highlight>
              {typeof hit.score === "number" ? (
                <Badge size="sm" variant="light">
                  {hit.score.toFixed(2)}
                </Badge>
              ) : null}
            </Group>
            {(hit.matches ?? []).slice(0, 3).map((match, index) => (
              <Highlight
                key={`${hit.filename}-${index}`}
                highlight={submittedQuery}
                size="sm"
                c="dimmed"
                lineClamp={2}
              >
                {match.context}
              </Highlight>
            ))}
          </Paper>
        </UnstyledButton>
      ))}
    </Stack>
  );
}
