> **🇬🇧 English** | [🇷🇺 Русский](../ru/architecture.md)

[← Configuration](configuration.md) · [🏠 README](../../README.md)

# Architecture

## Overview

Relix is a single-package Go application built on the [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework, following the **Elm architecture** (Model-Update-View). The entire application is organized as a screen-based state machine where each screen has dedicated update and view handlers. User interactions produce messages, messages update the model, and the model determines the rendered view.

```
User Input → tea.Msg → Update(model, msg) → New Model → View(model) → Terminal Output
```

## Screen Flow

```
screenLoading
    │
screenAuth ─────────→ screenHome
                         │
           ┌─────────────┼─────────────┐
           │                           │
    screenMain (MR list)        screenHistoryList
           │                           │
    screenEnvSelect             screenHistoryDetail
           │
    screenVersion
           │
    screenSourceBranch
           │
    screenEnvMerge
           │
    screenRootMerge
           │
    screenConfirm
           │
    screenRelease
           │
    screenHome (complete)
```

The flow is linear for the release workflow (left branch) and separate for history browsing (right branch). Navigation between screens is controlled by the central `Update()` function, which routes messages to screen-specific handlers.

## Project Structure

All source files reside in the root package (`package main`). The codebase is organized by responsibility:

### Core

| File | Lines | Purpose |
|------|-------|---------|
| `main.go` | ~100 | Entry point, CLI flag parsing, program initialization |
| `model.go` | ~600 | Central model, Update/View routing, modal management |
| `types.go` | ~400 | All type definitions, screen constants, message types |

### Screens

Each screen has an `update{Screen}` function (handles key events and messages) and a `view{Screen}` function (renders the UI):

| File | Screen | Purpose |
|------|--------|---------|
| `auth_screen.go` | `screenAuth` | GitLab credential input |
| `home_screen.go` | `screenHome` | Project info, actions menu |
| `mrs_screen.go` | `screenMain` | MR list with multi-selection |
| `environment_screen.go` | `screenEnvSelect` | Environment picker |
| `version_screen.go` | `screenVersion` | Semantic version input |
| `source_branch_screen.go` | `screenSourceBranch` | Source branch configuration |
| `env_merge_screen.go` | `screenEnvMerge` | Merge strategy selection |
| `root_merge_screen.go` | `screenRootMerge` | Merge-back strategy |
| `confirm_screen.go` | `screenConfirm` | Release summary review |
| `release_screen.go` | `screenRelease` | Release execution (largest file) |
| `history_list_screen.go` | `screenHistoryList` | Release history browser |
| `history_detail_screen.go` | `screenHistoryDetail` | Release detail view |
| `error_screen.go` | `screenError` | Error display |

### Infrastructure

| File | Purpose |
|------|---------|
| `gitlab.go` | GitLab API client (projects, MRs, pipelines, diffs) |
| `git_executor.go` | PTY-based git execution with virtual terminal emulation |
| `config.go` | Config file I/O (`~/.relix/config.json`) |
| `keyring.go` | OS keyring for secure credential storage |
| `release_history.go` | Release history persistence (index + detail files) |

### UI

| File | Purpose |
|------|---------|
| `styles.go` | Lipgloss style definitions |
| `theme.go` | Dynamic theming with ANSI color remapping |
| `modal.go` | Modal overlay base component |
| `command_menu.go` | Command menu (`/` key) |
| `project_selector.go` | Project search/selection modal |
| `open_options_modal.go` | Browser open options |
| `settings_screen.go` | Settings modal (release + theme tabs) |
| `utils.go` | Text wrapping, version parsing, file exclusion logic |

## Key Patterns

### Message-Based Async

All long-running operations (GitLab API calls, git commands) return `tea.Cmd` functions that perform the operation asynchronously and send typed result messages back to the `Update()` loop. Loading states are tracked via boolean flags (`loadingMRs`, `loadingProjects`, etc.) to display spinners in the UI.

```
User Action → tea.Cmd → Async Operation → tea.Msg → Update() → New State
```

This pattern keeps the UI responsive during network calls and git operations.

### Modal System

Modals overlay the base screen via boolean flags (`showCommandMenu`, `showProjectSelector`, `showSettings`). When a modal is active, key events are routed to the modal handler first, then to the underlying screen handler only if the modal does not consume the event. The `closeAllModals()` function centralizes modal cleanup to prevent stale state.

### Release State Machine

The release process (`release_screen.go`) is the most complex part of the application. It is implemented as a multi-step state machine tracked by the `ReleaseStep` enum. The flow:

1. Each step executes git commands via `GitExecutor`
2. A `releaseStepCompleteMsg` signals step completion
3. The next step starts automatically (or waits for user input on certain steps)
4. On conflict or error, the process pauses for user intervention
5. State is persisted to `~/.relix/release.json` after each successful step for crash recovery
6. On completion, state is saved to release history and the release file is deleted

### Git Executor

Git commands are executed in a PTY (pseudo-terminal) to capture colored output exactly as you would see it in a regular terminal. The `VirtualTerminal` component (backed by the vt10x library) parses ANSI escape codes and maintains a cell grid. Output is streamed to the UI at approximately 20 FPS via `releaseScreenMsg`.

```
Git Command → PTY → Raw Bytes → VirtualTerminal → ANSI Parsing → Cell Grid → UI Render
```

### Two-Tier History

Release history uses a two-tier storage strategy for performance:

- **Index file** (`index.json`) -- Lightweight entries with just enough data (date, version, environment, status) for fast list rendering
- **Detail files** (`{timestamp}.json`) -- Full release data including terminal output, MR metadata (URLs, IIDs, commit SHAs), and the complete release configuration

This avoids loading potentially large terminal output blobs when the user is just browsing the history list.

## See Also

- [Configuration](configuration.md) -- config structure and theme system
- [Usage Guide](usage.md) -- the release workflow from a user perspective
- [Getting Started](getting-started.md) -- installation and first run
