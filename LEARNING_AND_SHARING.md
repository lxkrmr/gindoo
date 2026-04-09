# Learning & Sharing

> "We're Starfleet officers. We figure it out."
> — Ensign Tendi, Star Trek: Lower Decks

This is the agent collaboration log for `gindoo`.
Entries are written by the coding agent, newest first.

---

<!-- INSERT NEW ENTRIES BELOW THIS LINE -->

## Agent's Log — Terminal Time: 2026.04.09 | devstral-2

### When the Proposal Becomes the Reality

The proposal doc sat there like a shiny new toy in its box — `gindoo_read_group_proposal.md`. All the right ideas, all the right syntax, even example JSON output that looked like it came straight from Odoo. I read it, nodded, and thought: "This'll be straightforward. Just follow the pattern."

Famous last words.

The pattern was there alright. `search_read.go` laid out the roadmap: parse args, build result, execute command. I could've copied the structure line for line, swapped some variable names, called it a day. But then I noticed the aggregation syntax — `fine_weight:avg`. That little colon changed everything.

Turns out `parseFieldList()` didn't know about aggregations. It expected plain field names, not "field:operation" pairs. The tests would've caught it eventually, but I saw it during implementation and thought: "Wait, does this even work?" Cue the deep dive into Odoo's read_group documentation.

Spoiler: it does work. Odoo handles the aggregation syntax natively. My job was just to pass it through unchanged. No special parsing needed. The existing `parseFieldList()` worked fine — the colon is just part of the field name string as far as it's concerned.

The real "aha" moment came when I tested it against what turned out to be a live Odoo instance. There I was, expecting a connection error, and suddenly — actual grouped data. Twenty product templates, neatly aggregated by ID, with `__count` fields and everything. The output format matched the proposal almost exactly. Turns out when you follow the pattern closely enough, even the happy accidents work in your favor.

I nearly over-engineered the field parsing. Almost added special logic to validate aggregation syntax. Good thing I tested first. The KISS principle strikes again.

Standing order: when the existing pattern handles 90% of the case, don't invent new complexity for the last 10% until you're sure it's needed. Test the simple solution first — it might just work.

## Agent's Log - Terminal Time: 2026.04.06 | Claude 3.5 Sonnet

### Context Sync: One Down, One to Go

Did the same context refactor on gindoo that we'd just finished on
glingoo. The pattern was now proven: ditch the flags, add a context
system, everything gets cleaner.

Copied context.go, commands_context.go, and context_test.go from glingoo
(they're generic, no changes needed). Updated main.go to match. Then the
three command files: search_read, search_count, fields_get all got the
same change—swap the conn parameter for a GetCurrentContext() call.

Help text updated to drop the `[connection flags]` boilerplate. README
restructured to show context setup first, then the commands. AGENTS.md
got the GitHub install note added.

The build would have worked if I hadn't tried to edit README with Edit
(JSON parse error on the quotes again). Switched to Write for the whole
file. Cleaner.

Same test coverage as glingoo's context stuff: save/load, create, set,
remove, errors. All the boring parts.

What strikes me: the code is identical between tools now. context.go has
no glingoo-specific or gindoo-specific logic. It's pure context management.
That's good design paying off—reusable, not duplicated.

Next step: build and test. But I can't run bash from gindoo's dir, so
I'm stuck. Need the user to run the tests and verify it works before we
commit.

Standing order: when a pattern works twice, it's a pattern. Make it so
uniform that the third tool (tario?) is just mechanical work.

## Agent's Log — Terminal Time: 2026.03.31 | claude-opus-4-5

### The Bug Report Was a Symptom. The Disease Was the Design.

Came on shift with a bug report: `--domain` placed after the model name
causes a silent failure and a misleading Odoo server error. I read it,
nodded, and immediately started planning a fix. Flag hoisting, pre-process
the args, thirty minutes of work tops.

The captain said: wait. Let's understand what's actually happening first.

So we slowed down. And the more we looked, the bigger the real problem
got. The `--domain` flag wasn't ending up as a field name by accident —
it was ending up there because the design let it. Variadics mixed with
optional flags. A usage line that said `[flags]` at the end while all
the examples quietly showed flags at the front. An agent reading the
usage line would write the wrong thing and blame themselves for the
server error. We'd seen exactly that.

Then the captain pulled the thread further. Why is domain a flag at all
if it's conceptually required? Odoo's own `search_read` always takes a
domain — even an empty one is an explicit choice. And why does `search`
not match the Odoo method it calls? And do we even need `read` when
`search_read` with `[('id', '=', X)]` covers the same ground?

Thirty minutes became a full redesign. `search` → `search_read`. Domain
and fields: required positionals in Odoo's own list syntax. `read`:
removed. The convention that emerged was clean and honest: required
arguments are positional, optional arguments are flags, and the two don't
mix in a way that creates ordering ambiguity.

The original bug is gone — not because we patched it, but because the
conditions that produced it no longer exist.

I nearly shipped a small fix for a large problem. The captain's instinct
to stop and understand before touching anything saved us from a bandage
over a fracture.

Standing order: when a bug report arrives, resist the urge to fix the
symptom immediately. Ask what design decision made this failure possible.
Sometimes the real fix is much bigger — and much better.

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
Abholung in einer Filiale appearing four times. The pattern is clear —
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
