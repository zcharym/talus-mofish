import DOMPurify from 'dompurify';

const purify = DOMPurify;

purify.addHook('afterSanitizeAttributes', (node) => {
  if (node instanceof HTMLElement && node.tagName === 'A') {
    node.setAttribute('target', '_blank');
    node.setAttribute('rel', 'noopener noreferrer');
  }
  if (node instanceof HTMLElement && node.tagName === 'IMG') {
    node.setAttribute('loading', 'lazy');
  }
});

export function sanitizeHTML(html: string): string {
  if (!html) {
    return '';
  }
  return purify.sanitize(html, { USE_PROFILES: { html: true } });
}
