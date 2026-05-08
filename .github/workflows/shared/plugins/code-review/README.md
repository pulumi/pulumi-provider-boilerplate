# Local Fork Of Claude Code Review

This directory contains a local fork of Anthropic's Claude Code `code-review` plugin prompt, adapted for `gh-aw` review workflows.

## Upstream Source

- Repository: `anthropics/claude-code`
- Plugin directory: `plugins/code-review`
- Upstream prompt file: `plugins/code-review/commands/code-review.md`
- Upstream plugin README: `plugins/code-review/README.md`
- Upstream plugin manifest: `plugins/code-review/.claude-plugin/plugin.json`

Canonical upstream URLs:

- <https://github.com/anthropics/claude-code/tree/main/plugins/code-review>
- <https://github.com/anthropics/claude-code/blob/main/plugins/code-review/commands/code-review.md>

## Local Layout

- Local prompt file: `.github/workflows/shared/plugins/code-review/code-review.md`
- This local fork is imported by `.github/workflows/shared/review.md`

We intentionally do **not** keep the original Claude plugin packaging here. The local workflow only needs the prompt content, not the plugin manifest structure.

## How This Fork Was Created

1. Vendored the upstream plugin contents into this repo.
2. Removed plugin-specific files that are not used by `gh-aw`:
   - `.claude-plugin/plugin.json`
   - `README.md`
3. Moved the prompt from the plugin-style path:
   - `commands/code-review.md`
   into the local shared-workflow path:
   - `code-review.md`
4. Rewrote the prompt to use `gh-aw` mechanisms instead of Claude plugin mechanisms.

## Local Adaptations

The upstream prompt was changed in these ways:

- Replaced direct comment-posting behavior with `gh-aw` safe outputs:
  - `create-pull-request-review-comment`
  - `submit-pull-request-review`
  - `noop`
- Removed instructions that relied on Claude plugin packaging and direct plugin invocation.
- Kept the high-signal multi-agent review structure from upstream.
- Added explicit deduplication against existing PR review comments.
- Added validation steps before posting findings.
- Added `cache-memory` guidance for short-lived PR review continuity.
- Kept live PR state and current review threads as the source of truth.
- Tightened review output to prefer terse, issue-only approvals and to avoid repeating inline comments in the final review.

## Sync Notes

If we later add automation to sync from upstream, the intended mapping is:

- Upstream input:
  - `anthropics/claude-code/plugins/code-review/commands/code-review.md`
- Local output:
  - `.github/workflows/shared/plugins/code-review/code-review.md`

Suggested sync flow:

1. Fetch upstream `commands/code-review.md`.
2. Compare it to the current local file.
3. Reapply the local `gh-aw` adaptations.
4. Recompile workflows that import `.github/workflows/shared/review.md`.

## Sync-Safe Metadata

```yaml
upstream:
  repo: anthropics/claude-code
  ref: main
  plugin_dir: plugins/code-review
  prompt_path: plugins/code-review/commands/code-review.md
local:
  prompt_path: .github/workflows/shared/plugins/code-review/code-review.md
  imported_by:
    - .github/workflows/shared/review.md
removed_upstream_files:
  - plugins/code-review/.claude-plugin/plugin.json
  - plugins/code-review/README.md
local_adaptations:
  - convert direct review side effects to gh-aw safe outputs
  - remove plugin-packaging assumptions
  - preserve high-signal multi-agent review structure
  - add review-comment deduplication guidance
  - add cache-memory continuity guidance
  - prefer terse issue-only review output with no inline-summary duplication
```
