---
name: tag
description: >
  Bookmark the current Claude Code session so you can find and resume it later.
  Generates a short description of the work done, locates the session ID in
  history.jsonl, and saves it via `tag add`. Use when the user types /tag,
  says "pin this session", "bookmark this", or "save this session".
allowed-tools: Bash(grep *), Bash(python3 *), Bash(tag add *)
---

# /tag — Pin Current Claude Code Session

## Step 1: Generate Description

Write a concise one-line description (max 10 words) of what this session is
about. Focus on the work done, not meta-commentary.

Good examples:
- "Fix subscription renewal bug in billing service"
- "Add fzf session picker to pml CLI"
- "Debug connection pool exhaustion in production"

Bad examples (too vague or meta):
- "Various code changes" — no substance
- "Worked on the project" — says nothing
- "Claude Code session" — that's every session

## Step 2: Determine Project Path

Use `$PWD` or infer from the files worked on in this session.

## Step 3: Find the Session ID

Tell the user: *"Finding session ID from history..."*

Read `~/.claude/history.jsonl` and find the most recent entry for this project.
Use `dangerouslyDisableSandbox: true` — sandbox blocks pipeline reads from
`~/.claude/`.

```bash
grep "<path-substring>" ~/.claude/history.jsonl | tail -1 | python3 -c "import sys,json; line=sys.stdin.readline().strip(); d=json.loads(line) if line else {}; print(d.get('sessionId','NOT FOUND'))"
```

Use a short, unique fragment of the project path as the grep pattern
(e.g. `pin-ml` for `/Users/foo/Documents/Git/pin-ml`). Do NOT try to match
the full JSON key — history.jsonl uses `"project": "..."` with a space after
the colon, so an exact-key grep will silently return nothing.

If the result is `NOT FOUND`, tell the user and stop — don't guess.

## Step 4: Pin It

Tell the user: *"Pinning session `<id[:8]>`..."*

```bash
tag add <sessionId> <projectPath> "<description>"
```

## Step 5: Confirm

Reply:

> Pinned `<id[:8]>` — *"<description>"*. Resume with: `tag ls`
