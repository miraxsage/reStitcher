# Design: New Step [5] — Env Merge

## Overview

Add a new release configuration step that lets the user choose how the release root branch content is applied to the environment release branch: squash merge (current behavior) or regular git merge.

## Flow Position

```
[1] MRs → [2] Env → [3] Version → [4] Source Branch → [5] Env Merge (NEW) → [6] Root Merge → Confirm → Release
```

Current [5] Root Merge becomes [6]. All `[5]`/`[6]` labels shift across screens and sidebar sections.

## New Screen: `screenEnvMerge`

A screen with two vertically listed options (arrow/j/k to switch, enter to confirm):

**Option 1 — Squash merge** (selected by default):
> All current release changes will be accumulated in one common commit from the current {ENV} branch and just copied as independent changes via "git checkout {release branch} -- .". That makes the merge safe from conflicts; containing commits will be mentioned in the commit message.

**Option 2 — Regular merge:**
> All selected commits will be merged straight to {ENV} branch. There is a risk of conflicts; you should resolve them manually.

When "Regular merge" is selected, after the description:
1. Empty line
2. Spinner with "Merging commits calculation" while running `git rev-list --count {ENV}..{release root branch}` + summing `CommitsCount` from selected MRs (GitLab API data already fetched)
3. Once calculated, replace spinner with: "To {ENV} branch per current release will be merged {n} new commits"

Navigation: `j`/`k`/`↑`/`↓` to switch options, `Enter` to confirm, `Ctrl+Q` to go back to Source Branch.

## Model Changes

### New screen constant

`screenEnvMerge` — inserted between `screenSourceBranch` and `screenRootMerge` in the iota.

### New model fields

- `envMergeOptionIndex int` — currently focused option (0 = squash, 1 = regular)
- `envMergeSelection int` — confirmed choice (0 = squash, 1 = regular)
- `envMergeCommitCount int` — calculated commit count for regular merge
- `envMergeCountLoading bool` — true while commit count is being calculated

### New ReleaseState field

`EnvMergeMode string` (`"squash"` or `"regular"`) — persisted in release.json for crash recovery.

### New tea message

`envMergeCommitCountMsg { count int; err error }` — returned by async commit count calculation.

### New ReleaseHistoryEntry field

`EnvMergeMode string` — for record-keeping in release history.

## Screen Transitions

- `screenSourceBranch` → Enter → `screenEnvMerge` (was `screenRootMerge`)
- `screenEnvMerge` → Enter → `screenRootMerge`
- `screenEnvMerge` → Ctrl+Q → `screenSourceBranch`
- `screenRootMerge` → Ctrl+Q → `screenEnvMerge` (was `screenSourceBranch`)

## Sidebar Progression

Each screen's sidebar shows all previously completed steps:

| Screen | Sidebar sections |
|--------|-----------------|
| [5] Env Merge | quad sidebar: [1] MRs, [2] Env, [3] Version, [4] Source Branch |
| [6] Root Merge | quint sidebar: [1] MRs, [2] Env, [3] Version, [4] Source Branch, [5] Env Merge |
| Confirm | six sidebar: [1] MRs, [2] Env, [3] Version, [4] Source Branch, [5] Env Merge, [6] Root Merge |

### New sidebar section

`renderEnvMergeSidebarSection(width, contentHeight)` — shows "Squash" or "Regular" with appropriate styling.

### Sidebar layout functions

- Root Merge screen: `renderQuadSidebar` → new `renderQuintSidebar` (5 sections)
- Confirm screen: `renderFiveSidebar` → new `renderSixSidebar` (6 sections)

## Release Execution Changes

### When `EnvMergeMode == "squash"` (default)

No changes. Current behavior preserved:
1. `git rm -rf .`
2. `git checkout {sourceBranch} -- .`
3. Apply exclusion patterns
4. `git commit` with release message

### When `EnvMergeMode == "regular"`

Step `ReleaseStepCopyContent` changes behavior:
1. `git checkout {envReleaseBranch}`
2. `git merge {sourceBranch}` (actual git merge, no squash)
3. If merge conflict detected → suspend release with retry button (same pattern as `ReleaseStepMergeBranches`)
4. On success, the merge commit is created automatically by git

Step `ReleaseStepCommit` is skipped — git merge already created the commit.

Exclusion patterns are NOT applied in regular merge mode.

## Confirm Screen Markdown Changes

Step 4 text changes based on env merge selection:
- **Squash**: current text (copy via `git checkout -- .` as independent commit)
- **Regular**: "Merge {sourceBranch} to release/rpb-{ver}-{env} (regular merge, may require conflict resolution)"

Step 5 (exclusions) shown only for squash mode.

## Commit Count Calculation

Triggered when user selects "Regular merge" option on the Env Merge screen:
1. Run `git rev-list --count origin/{ENV}..{sourceBranch}` (async via tea.Cmd)
2. Sum `CommitsCount` from all selected MRs (already available from GitLab API)
3. Display total as: "To {ENV} branch per current release will be merged {n} new commits"

## Files Affected

| File | Changes |
|------|---------|
| `types.go` | Add `screenEnvMerge`, `envMergeCommitCountMsg`, `EnvMergeMode` to `ReleaseState` and `ReleaseHistoryEntry` |
| `model.go` | Add env merge model fields, handle new message type, update defaults |
| `env_merge_screen.go` | New file: `updateEnvMerge`, `viewEnvMerge`, `renderEnvMergeContent`, `renderEnvMergeSidebarSection`, commit count command |
| `source_branch_screen.go` | Change Enter transition target to `screenEnvMerge` |
| `root_merge_screen.go` | Change back navigation to `screenEnvMerge`, update `[5]` → `[6]` labels, `renderQuadSidebar` → `renderQuintSidebar` |
| `confirm_screen.go` | Update back navigation, `renderFiveSidebar` → `renderSixSidebar`, update markdown for step 4/5, add env merge sidebar section |
| `release_screen.go` | Branch on `EnvMergeMode` in `ReleaseStepCopyContent` and `ReleaseStepCommit`, add `EnvMergeMode` to state creation in `startRelease`, update `calculateReleaseTotalSteps` |
| `release_history.go` | Persist `EnvMergeMode` in history entry |
