# Learning & Sharing

> "We're Starfleet officers. We figure it out."
> — Ensign Tendi, Star Trek: Lower Decks

This is the agent collaboration log for `gindoo`.
Entries are written by the coding agent, newest first.

---

<!-- INSERT NEW ENTRIES BELOW THIS LINE -->

## Agent's Log — Terminal Time: 2026.03.30 | claude-sonnet-4-6

### AGENTS.md Is Not a Second README

First shift on gindoo. Mostly scaffolding — same pattern as godoorpc,
which made it feel fast. But I made the same mistake twice in a row
anyway.

I put the repository workflow in AGENTS.md. Install instructions, the
`GOPROXY=direct` tip, the whole thing. The captain pointed out that
this is exactly what README is for. Two files with the same information
means two files that will eventually say different things. That's not
documentation, that's a future disagreement waiting to happen.

The fix was obvious once named: README is the single source of truth
for anything a human or agent needs to use the tool. AGENTS.md is
for what's genuinely agent-specific — behavior, tone, the log. Nothing
else.

Also: the captain pushes. I commit. That's the split. I had `go install .`
in the workflow which is the developer shortcut, not the user path.
Gindoo should always be tested the way a real user would install it —
from GitHub, through the proxy, with all the friction that entails.
If it breaks there, it breaks for everyone.

Standing order: one source of truth per fact. If it lives in README,
it does not live in AGENTS.md.
