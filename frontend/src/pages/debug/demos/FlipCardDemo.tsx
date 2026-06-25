import { useState } from "react";
import { Badge, Group, Stack, Text, Title } from "@mantine/core";
import { FlipCard } from "../../../components/management/FlipCard";
import { notify } from "../../../services/notifications";

const sampleModelCss = `
.card { font-family: serif; line-height: 1.6; }
.card b { color: var(--mantine-color-blue-6); }
`;

export function FlipCardDemo() {
  const [showDeleteDemo, setShowDeleteDemo] = useState(true);

  return (
    <Stack gap="xl" maw={640}>
      <Stack gap="md">
        <Title order={5}>Vocabulary style</Title>
        <FlipCard
          title="serendipity"
          headerExtra={
            <Group gap="xs">
              <Text size="sm" c="dimmed">/ˌserənˈdɪpəti/</Text>
              <Badge size="sm" variant="outline">n.</Badge>
            </Group>
          }
          front={
            <Stack gap="xs" align="center" justify="center" h="100%">
              <Title order={2}>serendipity</Title>
              <Text c="dimmed">/ˌserənˈdɪpəti/</Text>
              <Badge variant="outline">n.</Badge>
            </Stack>
          }
          back={
            <Stack gap="sm">
              <Text>The occurrence of events by chance in a happy or beneficial way.</Text>
              <Text size="sm" c="dimmed">
                Finding the book was pure serendipity.
              </Text>
            </Stack>
          }
        />
      </Stack>

      <Stack gap="md">
        <Title order={5}>HTML / Anki model</Title>
        <FlipCard
          title="Sample article"
          modelCss={sampleModelCss}
          headerExtra={
            <Group gap="xs">
              <Badge size="sm" variant="light">Anki</Badge>
              <Text size="xs" c="dimmed">128 words</Text>
            </Group>
          }
          front={
            <div
              className="card"
              dangerouslySetInnerHTML={{
                __html: "<p>The quick brown fox <b>jumps</b> over the lazy dog.</p>",
              }}
            />
          }
          back={
            <div
              dangerouslySetInnerHTML={{
                __html: "<p>那只敏捷的棕色狐狸跳过了懒狗。</p>",
              }}
            />
          }
        />
      </Stack>

      <Stack gap="md">
        <Title order={5}>With delete menu</Title>
        {showDeleteDemo ? (
          <FlipCard
            title="Deletable card"
            onDelete={() => {
              setShowDeleteDemo(false);
              notify.success("Deleted", "Demo card removed. Switch demos to restore.");
            }}
            front={<Text ta="center">Use the ⋮ menu to delete this card.</Text>}
            back={<Text ta="center" c="dimmed">Back side</Text>}
          />
        ) : (
          <Text c="dimmed" size="sm">
            Demo card deleted. Switch to another demo and back to restore it.
          </Text>
        )}
      </Stack>
    </Stack>
  );
}
