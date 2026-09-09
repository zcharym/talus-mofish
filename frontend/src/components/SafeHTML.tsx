import type { HTMLAttributes } from 'react';
import { sanitizeHTML } from '../utils/html';

export function SafeHTML({ html, className, ...rest }: { html: string } & HTMLAttributes<HTMLDivElement>) {
  return <div className={className} {...rest} dangerouslySetInnerHTML={{ __html: sanitizeHTML(html) }} />;
}
