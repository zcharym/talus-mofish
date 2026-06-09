---
name: linear-project
description: Manages the talus_echo_loop project in Linear via MCP — creates plans, initiatives, projects, milestones, issues, and status updates. Use when the user mentions Linear, issues, backlog, roadmap, milestones, project planning, triage, or syncing work with this repo.
---

# Linear Project Management

Manage **Talus Mofish** (`talus_echo_loop`) in Linear using the Cursor Linear plugin MCP server.

## Quick start

```
Task Progress:
- [ ] Confirm Linear MCP connected (`linear` server available)
- [ ] Read [project-config.md](project-config.md) for team/project IDs
- [ ] Read tool schemas in `mcps/linear/tools/` before any CallMcpTool
- [ ] Execute requested workflow
- [ ] Return issue URLs and identifiers to the user
```

## Before every session

1. **Check MCP**: If `linear` is not in available MCP servers, stop and ask the user to install the Linear plugin (Cursor Settings → Plugins) and authenticate.
2. **Read config**: Load [project-config.md](project-config.md). If Team ID or Project ID are `_TODO_`, run discovery first (list teams → find Talus Mofish project → update config).
3. **Read schemas**: Before each `CallMcpTool`, read the matching JSON descriptor under `mcps/linear/tools/`.

## Core workflows

### Create a plan

When the user asks to plan work (feature, sprint, release):

1. Clarify scope and timeline if ambiguous.
2. Choose structure:
   - **Large/multi-quarter** → initiative + project + milestones
   - **Single feature/release** → project milestones or epic issues
   - **Small task** → single issue
3. Create hierarchy in Linear (initiative → project → milestones → issues).
4. Post a project or initiative update with the plan summary.
5. Present a table of created items with Linear links.

### Create issues

1. Apply naming from project-config: `[area] imperative title`.
2. Use the issue description template in [reference.md](reference.md).
3. Set `teamId`, `projectId`, labels, and priority from config defaults unless user overrides.
4. For multi-step work, create a parent issue with sub-issues.

### Break down existing plan/spec

1. Parse the plan into epics and tasks (~1–3 day issues).
2. Batch-create in Linear; group under milestones or parent issues.
3. Confirm the breakdown with the user before creating if >5 issues.

### Search, triage, update

- Search/filter issues by status, label, assignee, or keyword.
- Update status, add comments, link branches/PRs in descriptions.
- Draft project updates from recent activity.

## Output conventions

Always include in responses:

- Linear identifier (e.g. `ENG-42`) and clickable URL
- What was created/updated and why
- Next suggested step (e.g. "assign ENG-42 and start branch `feature/ENG-42-...`")

When creating multiple issues, summarize in a list:

```markdown
| ID | Title | Priority |
|----|-------|----------|
| ENG-42 | [frontend] Add review screen | High |
```

## Repo integration

When work ties to code in this repo:

- Reference paths under `frontend/`, `internal/`, `db/`, etc.
- Mention Wails/sqlc constraints from README when relevant
- Link GitHub PRs in issue descriptions; never commit secrets

## Additional resources

- Tool details, templates, error handling: [reference.md](reference.md)
- Team/project IDs and naming rules: [project-config.md](project-config.md)
