import { useCallback, useState } from "react";
import {
  Badge,
  Button,
  Group,
  Paper,
  Stack,
  Text,
  TextInput,
  UnstyledButton,
} from "@mantine/core";
import { IconSearch } from "@tabler/icons-react";
import { ObsidianService } from "../../bindings/github.com/songwei.ma/talus-mofish/backend/services";
import { SearchHit } from "../../bindings/github.com/songwei.ma/talus-mofish/backend/obsidian/models";
import { notify } from "../services/notifications";

function errorMessage(err: unknown): string {
  if (typeof err === "string") {
    return err;
  }
  if (err instanceof Error) {
    return err.message;
  }
  return String(err);
}

interface ObsidianSearchPageProps {
  onOpenNote: (path: string) => void;
}

export function ObsidianSearchPage({ onOpenNote }: ObsidianSearchPageProps) {
  const [query, setQuery] = useState("");
  const [hits, setHits] = useState<SearchHit[] | null>(null);
  const [loading, setLoading] = useState(false);

  const runSearch = useCallback(async () => {
    const trimmed = query.trim();
    if (!trimmed) {
      setHits(null);
      return;
    }
    setLoading(true);
    try {
      const results = (await ObsidianService.SearchSimple(trimmed, 100)) as SearchHit[];
      setHits(results);
    } catch (err) {
      notify.failed("Obsidian", errorMessage(err));
    } finally {
      setLoading(false);
    }
  }, [query]);

  return (
    <Stack gap="md" mt="sm" maw={720}>
      <Group align="flex-end">
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
        <Text c="dimmed">No matching notes.</Text>
      ) : null}

      {hits?.map((hit) => (
        <UnstyledButton key={hit.filename} onClick={() => onOpenNote(hit.filename)} w="100%">
          <Paper withBorder p="sm">
            <Group justify="space-between" mb={4}>
              <Text fw={600} size="sm">
                {hit.filename}
              </Text>
              {typeof hit.score === "number" ? (
                <Badge size="sm" variant="light">
                  {hit.score.toFixed(2)}
                </Badge>
              ) : null}
            </Group>
            {(hit.matches ?? []).slice(0, 3).map((match, index) => (
              <Text key={`${hit.filename}-${index}`} size="sm" c="dimmed" lineClamp={2}>
                {match.context}
              </Text>
            ))}
          </Paper>
        </UnstyledButton>
      ))}
    </Stack>
  );
}
