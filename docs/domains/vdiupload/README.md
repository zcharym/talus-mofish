# VDI Upload domain

**ID:** `vdiupload` · **Kind:** sidecar · **Constant:** `domain.VDIUpload`

Automated host→guest file transfer for **H3C Workspace** (extensible to RDP/Citrix/Horizon via YAML profiles).

## Docs

- **[DESIGN.md](./DESIGN.md)** — full engineering design (architecture, clipboard, Win32, roadmap, risks, testing)

## Code

| Path | Role |
|------|------|
| `backend/vdiupload` | Domain package (orchestrator, config, ports) |
| `cmd/vdi-upload` | CLI adapter |
| `backend/vdiupload/vdi-upload.example.yaml` | Example profile |

## Commands

```bash
task vdiupload:test
task vdiupload:build
task vdiupload:run CONFIG=path/to/vdi-upload.yaml FILE=C:\payload\a.zip
```

Phase 0–1 scaffold is in place; Win32 clipboard/window/input backends land in later phases per the design roadmap.
