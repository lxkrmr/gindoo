# Contributing

## Commits

Use Conventional Commits.

Format:

```text
type(scope): short description
```

Examples:

```text
feat(search): add domain filter flag
fix(output): handle false values from Odoo correctly
docs(adr): add decision for json-only output
refactor(cmd): extract connection flags to shared helper
test(search): cover empty result set
```

Rules:
- keep commits small and meaningful
- write commit messages in English
- prefer one focused change per commit
- use a scope that matches the main area you changed

Common types:
- `feat`
- `fix`
- `docs`
- `refactor`
- `test`
- `chore`
