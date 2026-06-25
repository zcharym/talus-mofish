import { useEffect, useRef, useState } from "react";
import { nextStreamRevealIndex } from "../utils/streamReveal";

/** ~9 words/sec at ~5 chars/word — steady, readable cadence. */
const BASE_INTERVAL_MS = 55;
/** Gentle catch-up when the buffer runs far ahead of what's shown. */
const CATCHUP_INTERVAL_MS = 24;
const CATCHUP_LAG_CHARS = 60;

export interface SmoothStreamOptions {
  /** When false, text is shown immediately (e.g. user messages or history). */
  enabled: boolean;
  /** True while the upstream source is still sending content. */
  streaming: boolean;
}

export function useSmoothStreamText(
  target: string,
  { enabled, streaming }: SmoothStreamOptions,
): { text: string; revealing: boolean } {
  const streamingRef = useRef(streaming);
  streamingRef.current = streaming;

  const targetRef = useRef(target);
  targetRef.current = target;

  const [displayedLength, setDisplayedLength] = useState(() =>
    enabled ? 0 : target.length,
  );

  const displayedRef = useRef(displayedLength);
  displayedRef.current = displayedLength;

  useEffect(() => {
    if (enabled && target === "") {
      setDisplayedLength(0);
      displayedRef.current = 0;
    }
  }, [enabled, target]);

  useEffect(() => {
    if (!enabled) {
      setDisplayedLength(target.length);
      displayedRef.current = target.length;
      return;
    }

    setDisplayedLength(0);
    displayedRef.current = 0;

    let cancelled = false;
    let timeoutId = 0;

    const schedule = () => {
      if (cancelled) {
        return;
      }

      const fullText = targetRef.current;
      const shown = displayedRef.current;

      if (shown >= fullText.length && !streamingRef.current) {
        return;
      }

      if (shown >= fullText.length) {
        timeoutId = window.setTimeout(schedule, BASE_INTERVAL_MS);
        return;
      }

      const lag = fullText.length - shown;
      const delay = lag > CATCHUP_LAG_CHARS ? CATCHUP_INTERVAL_MS : BASE_INTERVAL_MS;

      timeoutId = window.setTimeout(() => {
        if (cancelled) {
          return;
        }

        const next = nextStreamRevealIndex(targetRef.current, displayedRef.current);
        displayedRef.current = next;
        setDisplayedLength(next);
        schedule();
      }, delay);
    };

    schedule();

    return () => {
      cancelled = true;
      window.clearTimeout(timeoutId);
    };
  }, [enabled]);

  const revealing = enabled && (streaming || displayedLength < target.length);
  const text = enabled ? target.slice(0, displayedLength) : target;

  return { text, revealing };
}
