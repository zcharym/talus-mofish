import { useCallback, useMemo, useRef, type ReactNode } from "react";
import { IconCopy, IconDots } from "@tabler/icons-react";
import { ActionIcon, Box, Group, Loader, Menu, Paper } from "@mantine/core";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import { useSmoothStreamText } from "../../../hooks/useSmoothStreamText";
import { notify } from "../../../services/notifications";
import classes from "./ChatBubble.module.css";

export interface ChatBubbleMenuItem {
  key: string;
  label: string;
  icon?: ReactNode;
  onClick: () => void;
  color?: string;
  disabled?: boolean;
}

export interface ChatBubbleProps {
  role: "user" | "assistant";
  content: string;
  /** When true, shows a loading state before content arrives or a cursor while content streams in. */
  streaming?: boolean;
  /** Extra menu entries for the assistant message options menu. */
  menuItems?: ChatBubbleMenuItem[];
}

export function ChatBubble({ role, content, streaming = false, menuItems = [] }: ChatBubbleProps) {
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
  const hasCopyableContent = Boolean(content.trim());
  const showAssistantActions = !isUser && !showLoader && hasCopyableContent;
  const hasMenuItems = menuItems.length > 0;

  const handleCopy = useCallback(async () => {
    const textToCopy = content.trim();
    if (!textToCopy) {
      return;
    }

    try {
      await navigator.clipboard.writeText(textToCopy);
      notify.success("Copied", "Message copied to clipboard");
    } catch {
      notify.failed("Copy failed", "Could not copy to clipboard");
    }
  }, [content]);

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

        {showAssistantActions ? (
          <Group className={classes.bubbleActions} gap={2} justify="flex-end" wrap="nowrap">
            {hasMenuItems ? (
              <Menu withinPortal position="bottom-end" shadow="sm">
                <Menu.Target>
                  <ActionIcon
                    variant="subtle"
                    size="sm"
                    className={classes.actionButton}
                    color="gray"
                    aria-label="Message options"
                  >
                    <IconDots size={14} />
                  </ActionIcon>
                </Menu.Target>
                <Menu.Dropdown>
                  {menuItems.map((item) => (
                    <Menu.Item
                      key={item.key}
                      leftSection={item.icon}
                      color={item.color}
                      disabled={item.disabled}
                      onClick={item.onClick}
                    >
                      {item.label}
                    </Menu.Item>
                  ))}
                </Menu.Dropdown>
              </Menu>
            ) : null}
            <ActionIcon
              variant="subtle"
              size="sm"
              className={classes.actionButton}
              color="gray"
              aria-label="Copy message"
              onClick={handleCopy}
            >
              <IconCopy size={14} />
            </ActionIcon>
          </Group>
        ) : null}
      </Paper>
    </Box>
  );
}
