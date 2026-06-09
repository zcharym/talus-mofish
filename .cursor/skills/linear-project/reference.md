# Linear MCP reference

## Prerequisites

1. Install the **Linear** plugin in Cursor (Settings → Plugins → Marketplace).
2. Authenticate when prompted (OAuth). If auth fails: delete saved auth (`rm -rf ~/.mcp-auth` on Unix; remove `%USERPROFILE%\.mcp-auth` on Windows) and reconnect.
3. Confirm the `linear` MCP server appears in the workspace MCP list.

Endpoint: `https://mcp.linear.app/mcp` (HTTP streamable — do not use deprecated `/sse`).

## Tool discovery (required before every call)

Linear MCP tools are dynamic. **Always** read the tool schema before calling:

```
mcps/linear/tools/<tool-name>.json
```

If the `linear` folder is missing from `mcps/`, the plugin is not connected — tell the user to install/authenticate before proceeding.

Use `CallMcpTool` with `server: "linear"` and the exact tool name from the schema.

## Capability map

Linear MCP supports finding, creating, and updating workspace objects. Capabilities include (exact tool names vary — read schemas):

| Domain | Typical operations |
|--------|-------------------|
| Issues | Search, create, update, comment, sub-issues |
| Projects | List, create, update, link issues |
| Initiatives | Create, edit, post updates |
| Milestones | Create, edit within projects |
| Updates | Project and initiative status updates |
| Labels | Manage project labels |
| URLs | Load Linear resources from pasted issue/project URLs |

## Workflows

### A. Plan a feature (initiative → project → milestones → issues)

```
1. Read project-config.md for team/project IDs
2. If planning spans multiple cycles → create or update initiative
3. Create or select project under initiative
4. Add milestones (phases) to project
5. For each milestone, create parent issues with sub-issues
6. Post initial project update summarizing scope and timeline
```

### B. Create a single issue from conversation

```
1. Read project-config.md defaults
2. Draft title using naming conventions
3. Write description:
   - Context / user story
   - Acceptance criteria (checkboxes)
   - Technical notes (files, APIs)
   - Links (designs, related issues)
4. Create issue with teamId, projectId, labels, priority
5. Return issue URL and identifier (e.g. ENG-42) to user
```

### C. Break a plan into issues

When the user shares a plan or spec:

```
1. Identify epics (→ parent issues or milestones)
2. Split into implementable issues (~1–3 day scope each)
3. Set dependencies via "blocked by" / parent-child in descriptions until tool supports links
4. Batch-create issues; confirm count and links with user
5. Optionally update project-config.md if new labels/areas emerge
```

### D. Sync status from development

```
1. User mentions issue ID (ENG-123) or pastes Linear URL → load issue
2. On branch/PR creation → add comment or update description with links
3. On merge → move status to Done; add brief completion note
4. Weekly → draft project update from closed issues + in-progress list
```

### E. Triage and search

```
1. Search issues by keyword, label, assignee, or status
2. Summarize backlog health: counts by status/priority
3. Surface blockers and stale issues (>14 days in progress)
```

## Issue description template

```markdown
## Context
[Why this work exists]

## Acceptance criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Technical notes
- Affected paths: `frontend/src/...`
- Dependencies: ENG-XX

## Links
- Branch: `feature/ENG-123-slug`
- PR: https://github.com/...
```

## Project update template

```markdown
## Summary
[1–2 sentences on overall progress]

## Completed
- ENG-XX: ...

## In progress
- ENG-YY: ...

## Next
- ...

## Risks / blockers
- ...
```

## Error handling

| Symptom | Action |
|---------|--------|
| `MCP server does not exist: linear` | User must install Linear plugin and restart MCP |
| Auth / 401 errors | Re-authenticate; clear `~/.mcp-auth` if stuck |
| Unknown tool name | Re-list `mcps/linear/tools/` — names may have changed |
| Missing team/project ID | Run discovery workflow; update project-config.md |
| Retired team | Issues are read-only; pick active team from config |

## Linear URL patterns

Paste these to load resources when supported:

- Issue: `https://linear.app/<workspace>/issue/<ID>/<slug>`
- Project: `https://linear.app/<workspace>/project/<id>/<slug>`
