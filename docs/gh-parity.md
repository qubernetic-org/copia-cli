# gh CLI Feature Parity Tracker

Mapping of GitHub CLI (`gh`) commands to Copia CLI equivalents. Tracks what is achievable via Gitea API and implementation status.

## Legend

- **Full** — Gitea API provides complete equivalent
- **Partial** — Possible but with limitations
- **None** — No Gitea API equivalent
- [ ] — Not implemented
- [x] — Implemented

## auth

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh auth login` | Full | `copia auth login` | [ ] |
| `gh auth logout` | Full | `copia auth logout` | [ ] |
| `gh auth status` | Full | `copia auth status` | [ ] |
| `gh auth token` | Full | `copia auth token` | [ ] |
| `gh auth setup-git` | None | n/a | — |

## repo

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh repo list` | Full | `copia repo list` | [ ] |
| `gh repo view` | Full | `copia repo view` | [ ] |
| `gh repo create` | Full | `copia repo create` | [ ] |
| `gh repo clone` | Full | `copia repo clone` | [ ] |
| `gh repo fork` | Full | `copia repo fork` | [ ] |
| `gh repo delete` | Full | `copia repo delete` | [ ] |
| `gh repo archive` | Partial | `copia repo archive` | [ ] |
| `gh repo rename` | Full | `copia repo rename` | [ ] |
| `gh repo edit` | Full | `copia repo edit` | [ ] |
| `gh repo set-default` | n/a | `copia repo set-default` | [ ] |

## issue

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh issue list` | Full | `copia issue list` | [ ] |
| `gh issue create` | Full | `copia issue create` | [ ] |
| `gh issue view` | Full | `copia issue view` | [ ] |
| `gh issue close` | Full | `copia issue close` | [ ] |
| `gh issue reopen` | Full | `copia issue reopen` | [ ] |
| `gh issue comment` | Full | `copia issue comment` | [ ] |
| `gh issue edit` | Full | `copia issue edit` | [ ] |
| `gh issue delete` | Full | `copia issue delete` | [ ] |
| `gh issue pin` | Full | `copia issue pin` | [ ] |
| `gh issue lock` | None | n/a | — |
| `gh issue transfer` | None | n/a | — |

## pr

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh pr list` | Full | `copia pr list` | [ ] |
| `gh pr create` | Full | `copia pr create` | [ ] |
| `gh pr view` | Full | `copia pr view` | [ ] |
| `gh pr merge` | Full | `copia pr merge` | [ ] |
| `gh pr close` | Full | `copia pr close` | [ ] |
| `gh pr reopen` | Full | `copia pr reopen` | [ ] |
| `gh pr diff` | Full | `copia pr diff` | [ ] |
| `gh pr checkout` | Partial | `copia pr checkout` | [ ] |
| `gh pr review` | Full | `copia pr review` | [ ] |
| `gh pr checks` | Partial | `copia pr checks` | [ ] |
| `gh pr comment` | Full | `copia pr comment` | [ ] |
| `gh pr edit` | Full | `copia pr edit` | [ ] |
| `gh pr ready` | Full | `copia pr ready` | [ ] |

## release

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh release list` | Full | `copia release list` | [ ] |
| `gh release create` | Full | `copia release create` | [ ] |
| `gh release view` | Full | `copia release view` | [ ] |
| `gh release delete` | Full | `copia release delete` | [ ] |
| `gh release upload` | Full | `copia release upload` | [ ] |
| `gh release download` | Full | `copia release download` | [ ] |
| `gh release edit` | Full | `copia release edit` | [ ] |

## label

| gh command | Parity | Copia CLI | Status |
|-----------|--------|-----------|--------|
| `gh label list` | Full | `copia label list` | [ ] |
| `gh label create` | Full | `copia label create` | [ ] |
| `gh label edit` | Full | `copia label edit` | [ ] |
| `gh label delete` | Full | `copia label delete` | [ ] |

## Other gh commands

| gh command group | Parity | Notes |
|-----------------|--------|-------|
| `gh api` | Full | Generic REST client — high priority |
| `gh search` | Partial | Basic `/repos/search` only, no advanced query syntax |
| `gh ssh-key` | Full | `GET/POST /user/keys` |
| `gh gpg-key` | Full | `GET/POST /user/gpg_keys` |
| `gh status` | Partial | Compose from multiple endpoints |
| `gh browse` | Full | Construct URL client-side |
| `gh gist` | None | No Gitea equivalent |
| `gh workflow` / `gh run` | None | Gitea Actions exist but Copia may not expose them |
| `gh codespace` | None | GitHub-specific |
| `gh copilot` | None | GitHub-specific |
| `gh project` | None | GitHub Projects, no Gitea equivalent |
| `gh cache` | None | GitHub Actions cache |
| `gh attestation` | None | GitHub-specific |
| `gh ruleset` | None | Gitea has branch protection via different API |
| `gh extension` | None | gh plugin system |

## Summary

| Status | Count |
|--------|-------|
| Full parity | 42 |
| Partial parity | 6 |
| No equivalent | 14 |
| **Implementable** | **48** |
