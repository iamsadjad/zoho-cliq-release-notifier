# Zoho Cliq Release Notifier

[![CI](https://github.com/iamsadjad/zoho-cliq-release-notifier/actions/workflows/ci.yml/badge.svg)](https://github.com/iamsadjad/zoho-cliq-release-notifier/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

GitHub Action that sends Release and Pre-release notifications to [Zoho Cliq](https://www.zoho.com/cliq/).

## Setup

1. Create an incoming webhook in Zoho Cliq for the channel that should receive release messages.
2. Add the webhook URL as a repository secret named `ZOHO_CLIQ_WEBHOOK_URL`.
3. Add a workflow that runs this action when a release is published.

## Usage

Notify Cliq when a GitHub Release is published:

```yaml
name: Notify Zoho Cliq on Release

on:
  release:
    types: [published]

permissions:
  contents: read

jobs:
  notify:
    runs-on: ubuntu-latest
    steps:
      - uses: iamsadjad/zoho-cliq-release-notifier@v1
        with:
          webhook_url: ${{ secrets.ZOHO_CLIQ_WEBHOOK_URL }}
```

### Notify for a specific tag

Use this when the release is created outside a `release` event (for example after a custom release step):

```yaml
- uses: iamsadjad/zoho-cliq-release-notifier@v1
  with:
    webhook_url: ${{ secrets.ZOHO_CLIQ_WEBHOOK_URL }}
    tag: ${{ steps.release.outputs.tag_name }}
    github_token: ${{ secrets.GITHUB_TOKEN }}
```

### Include pre-releases

By default, pre-releases are skipped. Set `notify_prerelease` to `true` to send them:

```yaml
- uses: iamsadjad/zoho-cliq-release-notifier@v1
  with:
    webhook_url: ${{ secrets.ZOHO_CLIQ_WEBHOOK_URL }}
    notify_prerelease: 'true'
```

A full example workflow is in [`examples/notify-on-release.yml`](examples/notify-on-release.yml).

## Inputs

| Input | Required | Default | Description |
| ----- | -------- | ------- | ----------- |
| `webhook_url` | No | `''` | Zoho Cliq webhook URL. Empty skips successfully. |
| `tag` | No | `''` | Fetch release by tag; otherwise use the release event. |
| `notify_prerelease` | No | `'false'` | Notify pre-releases when `true`. |
| `github_token` | No | `${{ github.token }}` | Used when `tag` is set. |
| `repository` | No | `${{ github.repository }}` | `owner/repo` for API and card. |

## Outputs

| Output | Description |
| ------ | ----------- |
| `notified` | `true` if Cliq was notified |
| `skipped_reason` | `missing_webhook`, `prerelease_skipped`, or empty |
| `http_status` | HTTP status from Cliq |
| `release_tag` | Processed tag |
| `prerelease` | `true` / `false` when resolved |

## License

[MIT](LICENSE) — see [SECURITY.md](SECURITY.md) for reporting.
