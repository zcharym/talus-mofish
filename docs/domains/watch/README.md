# Echo Watch domain

**ID:** `watch` · **Kind:** sidecar · **Constant:** `domain.Watch`

## Responsibility

Run on the **physical PC** (never inside the VDI guest). Locate the local VDI client window, capture a configured region, OCR + rule-match, push minimal alerts through Cloudflare to an iOS PWA.

## Code ownership

| Layer | Path |
|-------|------|
| Domain logic | `backend/watch` |
| CLI adapter | `cmd/echo-watch` |
| Cloud + PWA | `cloud/echo-watch` |
| Ops docs | `cloud/echo-watch/README.md` |

## Boundaries

- **Owns:** target window resolution, capture/OCR, rule engine, alert notifier contract.
- **Uses:** Win32 via `golang.org/x/sys/windows` (today in-package); planned extract to `backend/platform/win32` when shared with `vdiupload`.
- **Does not own:** guest-side automation or file transfer (`vdiupload`).

## Related

- Domain map: [../README.md](../README.md)
- Sibling domain (same VDI host surface, different job): [vdiupload](../vdiupload/DESIGN.md)
