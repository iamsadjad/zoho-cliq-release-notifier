# Zoho Cliq Release Notifier

[![CI](https://github.com/iamsadjad/zoho-cliq-release-notifier/actions/workflows/ci.yml/badge.svg)](https://github.com/iamsadjad/zoho-cliq-release-notifier/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

GitHub Action that sends Release and Pre-release notifications to [Zoho Cliq](https://www.zoho.com/cliq/).

## Usage

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

Secret: `ZOHO_CLIQ_WEBHOOK_URL`

```yaml
- uses: iamsadjad/zoho-cliq-release-notifier@v1
  with:
    webhook_url: ${{ secrets.ZOHO_CLIQ_WEBHOOK_URL }}
    tag: ${{ steps.release.outputs.tag_name }}
    github_token: ${{ secrets.GITHUB_TOKEN }}
```

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

## Versioning

Pushes to `main` with Conventional Commits create tags and GitHub Releases:

| Commit | Bump |
| ------ | ---- |
| `fix:` | PATCH |
| `feat:` | MINOR |
| `BREAKING CHANGE:` / `feat!:` | MAJOR |

Floating tags `v1` and `v1.x` are updated after each release.

```yaml
uses: iamsadjad/zoho-cliq-release-notifier@v1
```

Repo setting: **Actions → Workflow permissions → Read and write**.

## Development

```bash
make all
go test ./... -race
```

## License

[MIT](LICENSE) — see [SECURITY.md](SECURITY.md) for reporting.
