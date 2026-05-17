# Issue tracker: Backlog.md

Issues and PRDs for this repo live in [Backlog.md](https://github.com/benwisler/backlog), managed via the `backlog` CLI.

## Conventions

- Tasks are created with `backlog-bun task create "Title" -d "Description" --ac "Criterion 1"`
- Tasks live under `backlog/tasks/`
- Triage state is recorded as a `Label:` line in each task file (see `triage-labels.md` for the role strings)
- Comments and conversation history append via `backlog-bun task edit <id> --notes "..."` or `--append-notes`

## When a skill says "publish to the issue tracker"

Create a new Backlog task via the CLI: `backlog-bun task create "Title" -d "..."`.

## When a skill says "fetch the relevant ticket"

Read the task with `backlog-bun task <id> --plain`. The user will normally pass the task ID directly.
