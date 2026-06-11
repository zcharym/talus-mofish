import { RefObject, useLayoutEffect, useState } from 'react';

const MIN_HEIGHT = 160;
const VIEWPORT_BOTTOM_PADDING = 16;
const FOOTER_GAP = 16;

export function useDynamicScrollHeight(
  anchorRef: RefObject<HTMLElement | null>,
  footerRef: RefObject<HTMLElement | null>,
  deps: unknown[] = [],
): number {
  const [height, setHeight] = useState(MIN_HEIGHT);

  useLayoutEffect(() => {
    const update = () => {
      const anchor = anchorRef.current;
      if (!anchor) {
        return;
      }

      const { top } = anchor.getBoundingClientRect();
      const footerHeight = footerRef.current?.getBoundingClientRect().height ?? 0;
      const footerGap = footerHeight > 0 ? FOOTER_GAP : 0;
      const next =
        window.innerHeight - top - footerHeight - footerGap - VIEWPORT_BOTTOM_PADDING;

      setHeight(Math.max(MIN_HEIGHT, Math.floor(next)));
    };

    update();

    const observer = new ResizeObserver(update);
    observer.observe(document.documentElement);

    const anchor = anchorRef.current;
    if (anchor) {
      observer.observe(anchor);
    }

    const footer = footerRef.current;
    if (footer) {
      observer.observe(footer);
    }

    window.addEventListener('resize', update);

    return () => {
      observer.disconnect();
      window.removeEventListener('resize', update);
    };
  }, [anchorRef, footerRef, ...deps]);

  return height;
}
