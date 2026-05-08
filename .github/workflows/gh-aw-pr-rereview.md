---
description: Run PR re-review on explicit maintainer slash command.
timeout-minutes: 30
strict: true
on:
  slash_command:
    name: review-again
    events: [pull_request_comment, pull_request_review_comment]
imports:
  - shared/review.md
  - shared/plugins/code-review/code-review.md
permissions:
  contents: read
  pull-requests: read
  id-token: write
---

# Internal PR Re-Review (Slash Command)
