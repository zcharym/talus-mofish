/** Advance one readable unit (word, whitespace run, or single character for CJK). */
export function nextStreamRevealIndex(text: string, from: number): number {
  if (from >= text.length) {
    return text.length;
  }

  const slice = text.slice(from);
  const wordMatch = slice.match(/^(\s*\S+\s?)/);
  if (wordMatch) {
    return from + wordMatch[0].length;
  }

  return from + slice[0].length;
}
