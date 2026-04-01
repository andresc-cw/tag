# tag

Bookmark Claude Code sessions so you can find and resume them later.

`tag` stores lightweight pins (session ID, project path, description) in a local JSON file and lets you browse them with [fzf](https://github.com/junegunn/fzf). Select a session, and `tag` prints the `cd` + `claude --resume` command to pick up where you left off.

## Getting started

### Prerequisites

- Go 1.22+
- [fzf](https://github.com/junegunn/fzf) (`brew install fzf`)
- [Claude Code](https://docs.anthropic.com/en/docs/claude-code) installed

### Install

```sh
go install github.com/andres-cw/tag@latest
```

### Tag your first session

1. Open a Claude Code session and note the session ID (shown in the status bar or via `claude --print-session-id`).

2. Pin it:

   ```sh
   tag add <session-id> /path/to/project "what I was working on"
   ```

3. Later, browse your tags:

   ```sh
   tag ls
   ```

   fzf opens with your pinned sessions. Press **Enter** to print the resume command, then paste it into your shell.

4. To wire this into a single keystroke, add a shell function:

   ```sh
   # fish
   function tg
       set cmd (tag ls)
       and eval $cmd
   end

   # bash / zsh
   tg() { eval "$(tag ls)"; }
   ```

   Now `tg` opens the picker and resumes the selected session in one step.

## CLI reference

### `tag add <session-id> <project-path> <description>`

Pin a Claude Code session. If a pin with the same session ID already exists, it is replaced (upsert).

| Argument | Description |
|---|---|
| `session-id` | The Claude Code session identifier |
| `project-path` | Absolute path to the project directory |
| `description` | Free-text note describing the work |

### `tag rm <session-id>`

Remove a pinned session by its full ID. Exits with code 1 if the session is not found.

### `tag ls`

Open an fzf picker showing all pinned sessions. Each line shows:

```
<short-id>  <project-name>          <description>  <date>
```

**Keybindings inside the picker:**

| Key | Action |
|---|---|
| Enter | Print `cd <project> && claude --resume <id>` to stdout |
| Ctrl-D | Delete the selected tag (list reloads automatically) |
| Escape | Quit without action |

**Flag:** `--raw` — Print lines to stdout without launching fzf. Used internally by the Ctrl-D reload binding.

## About the store

Pins are saved as JSON at:

```
$XDG_CONFIG_HOME/pml/pins.json    # default: ~/.config/pml/pins.json
```

The file format is versioned (`"version": 1`) and writes are atomic (write to `.tmp`, then rename). You can safely back up or edit this file by hand.

Each pin records:

| Field | Type | Description |
|---|---|---|
| `id` | string | Claude Code session ID |
| `project` | string | Absolute project path |
| `description` | string | Free-text note |
| `pinned_at` | RFC 3339 timestamp | When the pin was created or last updated |

## Claude Code skill

The [`skill/SKILL.md`](skill/SKILL.md) file is a reference copy of the Claude Code skill that wraps this CLI. When installed as a Claude Code skill, typing `/tag` in a session automatically generates a description, finds the session ID from `~/.claude/history.jsonl`, and runs `tag add` — no manual copy-pasting required.

The canonical skill lives in the [acrucetta-core](https://github.com/acrucetta/agent-stuff) plugin. This copy is here for reference only.
