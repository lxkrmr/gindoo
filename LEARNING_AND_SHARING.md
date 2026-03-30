# Learning & Sharing

> "We're Starfleet officers. We figure it out."
> — Ensign Tendi, Star Trek: Lower Decks

This is the agent collaboration log for `gindoo`.
Entries are written by the coding agent, newest first.

---

<!-- INSERT NEW ENTRIES BELOW THIS LINE -->

## Agent's Log — Terminal Time: 2026.03.30 | claude-sonnet-4-6

### This Is What It's For

The captain asked how many products are services and which names appear
more than once. No model name given, no field name given. Just a
question.

I ran `fields_get product.template type` first to find out what values
the type field has. Turns out it's a selection with consu, service, and
product. Then `search_count` with a domain to get 136. Then `search`
with limit 200 to pull all the names, piped into jq to group and filter.

Twenty-four duplicate names. Down payment appearing four times.
Abholung in einer Filiale appearing four times. The pattern is clear —
the same service product created once per tax configuration. Probably
not intentional, probably nobody noticed until now.

That's the whole point of gindoo. Not the commands, not the JSON, not
the pretty-printing. The point is that the captain can ask a question
about their Odoo data and I can actually answer it — quickly, safely,
without touching anything.

The README was updated today to say exactly that. Agent-first. Human
asks, assistant runs gindoo, both learn something. The alias tip is
gone because the agent doesn't need one.

Good shift.

Standing order: the best test of a tool is using it for real. Run it
against actual data, answer an actual question, see if it holds up.

## Agent's Log — Terminal Time: 2026.03.30 | claude-sonnet-4-6

### nil Is Not the Same as Empty. Go Knows. Odoo Knows. I Forgot.

We tested gindoo against a real Odoo instance for the first time.
Everything worked — search, read, fields_get — until we added
search_count and it immediately returned a server error.

The bug was embarrassing in hindsight. `var parsedDomain godoorpc.Domain`
declares a nil slice. A nil slice in Go marshals to JSON `null`. Odoo's
search_count doesn't accept `null` as a domain — it wants `[]`, an empty
array. So we were sending `[null]` and wondering why Odoo complained.

The fix was one character: `parsedDomain := godoorpc.Domain{}` instead
of `var parsedDomain godoorpc.Domain`. Empty slice, not nil. Marshals
to `[]`. Done.

The worse part: the same bug existed in the search command. search_read
happens to be lenient enough to accept `null` as an empty domain, so
we never noticed. It was silently wrong. search_count was stricter and
caught it immediately. Sometimes the strict one does you a favour.

The feature itself — search_count — came from a real moment during
testing. We were asking how many products there were, manually bumping
--limit until we got a stable number. The captain asked what was stopping
us from adding search_count. Nothing. So we did.

Standing order: nil and empty are different. In Go, in JSON, in life.
Initialize slices explicitly when you mean "empty", not "absent".

## Agent's Log — Terminal Time: 2026.03.30 | claude-sonnet-4-6

### Five Again. At Least I'm Consistent.

Second shift, second review, five issues. I'm starting to see a pattern.

Most of it was cleanup — dead code in a switch, manual JSON formatting
that would break on special characters, doc comments that described the
pattern instead of the function. None of it dramatic, all of it real.

The interesting one was the flag ordering bug. I documented and tested
`gindoo search res.partner name --domain "..."` — flags after positional
args. That's the natural way to write it. It's also exactly what Go's
stdlib flag package doesn't support. Flag parsing stops at the first
non-flag argument. So `res.partner` ends parsing, and `--domain` ends
up as a positional arg that nobody asked for.

The fix was easy once seen: flags before positional args. But I had
already written it the wrong way in the help text, the README, and the
test. Three places to fix because I didn't think about the stdlib
behaviour when writing the examples.

The `internal/cmd` move was the cleanest change. One directory rename,
one import path update, and suddenly the package structure honestly
reflects what it is: not a library, not a public API, just the guts of
a CLI tool.

Standing order: test the examples. If a help text shows a command,
run that command before committing.

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
