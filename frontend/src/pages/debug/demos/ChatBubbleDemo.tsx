import { useCallback, useEffect, useRef, useState } from "react";
import { Button, Group, Stack, Text } from "@mantine/core";
import { ChatBubble } from "../../../components/agent/ChatBubble";

const SAMPLE_REPLY = `Streaming lets the assistant reply **word by word**, so you see progress instead of waiting for the full message.

This demo simulates that behavior locally — no API call needed.

Example list:
- Markdown **bold** and *italic*
- \`inline code\` and fenced blocks
- [Links](https://example.com)`;

/** Split into words while preserving spaces for natural chunk simulation. */
function splitIntoChunks(text: string): string[] {
  const chunks: string[] = [];
  const parts = text.split(/(\s+)/);
  for (const part of parts) {
    if (!part) continue;
    if (/^\s+$/.test(part)) {
      chunks.push(part);
    } else {
      chunks.push(part);
    }
  }
  return chunks;
}

const SAMPLE_CHUNKS = splitIntoChunks(SAMPLE_REPLY);

export function ChatBubbleDemo() {
  const [streamedContent, setStreamedContent] = useState("");
  const [streaming, setStreaming] = useState(false);
  const timerRef = useRef<number | null>(null);
  const chunkIndexRef = useRef(0);

  const clearTimer = useCallback(() => {
    if (timerRef.current !== null) {
      window.clearInterval(timerRef.current);
      timerRef.current = null;
    }
  }, []);

  const startStream = useCallback(() => {
    clearTimer();
    setStreamedContent("");
    setStreaming(true);
    chunkIndexRef.current = 0;

    timerRef.current = window.setInterval(() => {
      const index = chunkIndexRef.current;
      if (index >= SAMPLE_CHUNKS.length) {
        clearTimer();
        setStreaming(false);
        return;
      }

      const chunk = SAMPLE_CHUNKS[index];
      chunkIndexRef.current = index + 1;
      setStreamedContent((current) => current + chunk);
    }, 110);
  }, [clearTimer]);

  useEffect(() => () => clearTimer(), [clearTimer]);

  return (
    <Stack gap="md" maw={640}>
      <Text size="sm" c="dimmed">
        Words appear at a steady, readable pace. Chunks arrive from the simulated API every ~110ms;
        the bubble reveals them smoothly instead of jumping in all at once.
      </Text>

      <Stack gap="xs">
        <ChatBubble role="user" content="Can you explain streaming responses?" />
        <ChatBubble role="assistant" content={streamedContent} streaming={streaming} />
      </Stack>

      <Group>
        <Button onClick={startStream} loading={streaming} disabled={streaming}>
          Play stream
        </Button>
        <Button
          variant="default"
          onClick={() => {
            clearTimer();
            setStreamedContent("");
            setStreaming(false);
            chunkIndexRef.current = 0;
          }}
          disabled={streaming && streamedContent.length === 0}
        >
          Reset
        </Button>
      </Group>
    </Stack>
  );
}
