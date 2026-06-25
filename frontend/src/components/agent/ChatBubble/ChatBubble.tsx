import { useMemo, useRef } from "react";
import { Box, Loader, Paper } from "@mantine/core";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useSmoothStreamText } from "../../../hooks/useSmoothStreamText";
import classes from "./ChatBubble.module.css";

export interface ChatBubbleProps {
  role: "user" | "assistant";
  content: string;
  /** When true, shows a loading state before content arrives or a cursor while content streams in. */
  streaming?: boolean;
}

export function ChatBubble({ role, content, streaming = false }: ChatBubbleProps) {
  const isUser = role === "user";
  const wasStreamingRef = useRef(false);
  if (streaming) {
    wasStreamingRef.current = true;
  }

  const smoothEnabled = !isUser && wasStreamingRef.current;

  const { text: displayContent, revealing } = useSmoothStreamText(content, {
    enabled: smoothEnabled,
    streaming,
  });

  const markdownSource = useMemo(() => {
    if (isUser) {
      return content || "…";
    }
    if (smoothEnabled) {
      return displayContent;
    }
    return content || "…";
  }, [isUser, content, smoothEnabled, displayContent]);

  const showLoader = !isUser && streaming && !content;
  const showCursor = !isUser && revealing && Boolean(markdownSource);

  return (
    <Box className={`${classes.row} ${isUser ? classes.rowUser : classes.rowAssistant}`}>
      <Paper
        className={`${classes.bubble} ${isUser ? classes.bubbleUser : classes.bubbleAssistant}`}
        radius="md"
        p="sm"
        shadow="xs"
      >
        {showLoader ? (
          <Loader size="xs" type="dots" />
        ) : (
          <div className={classes.content}>
            <ReactMarkdown remarkPlugins={[remarkGfm]}>{markdownSource}</ReactMarkdown>
            {showCursor ? <span className={classes.cursor}>▍</span> : null}
          </div>
        )}
      </Paper>
    </Box>
  );
}
