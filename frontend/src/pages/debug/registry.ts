import type { ComponentType } from "react";
import { ChatBubbleDemo } from "./demos/ChatBubbleDemo";
import { ChatModalDemo } from "./demos/ChatModalDemo";
import { FlipCardDemo } from "./demos/FlipCardDemo";

export interface DebugDemoEntry {
  id: string;
  title: string;
  description: string;
  Component: ComponentType;
}

export const debugDemos: DebugDemoEntry[] = [
  {
    id: "chat-bubble",
    title: "ChatBubble",
    description: "Chat message bubbles with simulated streaming content.",
    Component: ChatBubbleDemo,
  },
  {
    id: "floating-chat",
    title: "FloatingChat",
    description: "Bottom-right floating chat widget wrapping a reusable ChatPanel.",
    Component: ChatModalDemo,
  },
  {
    id: "flip-card",
    title: "FlipCard",
    description: "3D flip cards for vocabulary, articles, and SRS content.",
    Component: FlipCardDemo,
  },
];
