---
name: backlog
description: >-
  Manage project tasks, status, and documentation using the Backlog.md CLI.
  Use when tasks, backlog, todo, planning, or project management are mentioned.
  Handles task creation, editing, acceptance criteria, and kanban workflows.
  Interacts with the Backlog.md CLI tool only, never edits task markdown files directly.
---

# Backlog.md Skill

Manage project tasks via the Backlog.md CLI (`backlog-bun`). All task operations MUST go through the CLI. Direct file editing of task markdown files is strictly forbidden.

## When to Use Me

- User mentions tasks, backlog, todo items, or project planning
- User wants to create, edit, check, or manage project tasks
- User references task IDs (e.g., "TASK-001", "task 42")
- User asks about project status, kanban board, or sprint progress
- User mentions acceptance criteria, definition of done, or task completion
- User asks to list, search, view, or filter tasks
- User wants to add notes, plans, or final summaries to tasks

## When NOT to Use Me

- User wants to edit task files directly (forbidden — use CLI only)
- User wants general project advice without task management
- User asks about Backlog.md installation or setup (handled by AGENTS.md)

## Core Principle

**NEVER edit task files directly.** Always use the Backlog.md CLI.

## Workflow

### Reading Tasks
```bash
# List all tasks
backlog-bun task list --plain

# View specific task
backlog-bun task <id> --plain

# Search tasks
backlog-bun search "<query>" --plain

# Filter by status
backlog-bun task list -s "In Progress" --plain

# Filter by assignee
backlog-bun task list -a @sara --plain
```

### Creating Tasks
```bash
backlog-bun task create "Title" -d "Description" --ac "Criterion 1" --ac "Criterion 2"
```

### Updating Tasks
```bash
# Change status and assign
backlog-bun task edit <id> -s "In Progress" -a @myself

# Mark AC complete (supports multiple)
backlog-bun task edit <id> --check-ac 1 --check-ac 2

# Add notes
backlog-bun task edit <id> --notes "Progress update"

# Add implementation plan
backlog-bun task edit <id> --plan "1. Research\n2. Implement\n3. Test"

# Append notes progressively
backlog-bun task edit <id> --append-notes "Another update"

# Add final summary (PR description)
backlog-bun task edit <id> --final-summary "Implemented X, updated Y"
```

### Board & Browser
```bash
# Terminal kanban board
backlog-bun board

# Web UI
backlog-bun browser
```

## Critical Rules

- **NEVER** edit markdown files in `backlog/tasks/` directly
- **ALWAYS** use `backlog-bun task edit` for any task changes
- **USE** `--plain` flag for AI-readable output when viewing/listing
- **CHECK** task ID with `backlog-bun task list --plain` before referencing
- **CREATE** tasks via CLI only: `backlog-bun task create "Title" ...`
- **MULTIPLE** operations in single command: `--check-ac 1 --uncheck-ac 2 --remove-ac 3`

## Multi-line Input

Prefer repeating flags for each line (works in all shells):
```bash
backlog-bun task edit 42 --notes "First line"
backlog-bun task edit 42 --append-notes "Second line"
backlog-bun task edit 42 --append-notes "Final paragraph"
```

Or real newlines inside double quotes:
```bash
backlog-bun task edit 42 --notes "First line
Second line

Final paragraph"
```

## Phase Discipline

| Phase | What to add |
|-------|-------------|
| Creation | Title, Description, Acceptance Criteria, labels/priority |
| Implementation | Plan (after moving to In Progress), Notes (progress log) |
| Wrap-up | Final Summary, verify AC and DoD, set status Done |

## Definition of Done

A task is Done only when ALL of the following are complete:
1. All acceptance criteria checked via CLI
2. All DoD items checked via CLI
3. Final Summary added via CLI
4. Status set to Done via CLI
