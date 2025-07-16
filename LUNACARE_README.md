# LunaCare Cosmo Router Fork

LunaCare's fork of the [Wundergraph Cosmo Router](https://github.com/wundergraph/cosmo) with custom features.

## Custom Features

- **498 Status Code Support**: Detects GraphQL errors with 498 status codes and sets HTTP response status accordingly
- **Custom Docker Build**: Specialized build process for ECR deployment

## Versioning Strategy

**Docker Tag Format**: `{upstream_version}-lunacare-{lunacare_version}`

**Example**: `0.192.1-lunacare-1.1.0`

### Two-Part System:
- **Upstream Version**: Auto-pulled from `router/CHANGELOG.md` by Build Luna Router workflow
- **LunaCare Version**: Manual in `router/VERSION-LUNACARE` - only increment when adding/modifying custom features

## CICD Process

**Build Luna Router workflow** automatically pulls router version from changelog.

### Manual Trigger (Recommended)
1. Go to GitHub Actions → "Build Luna Router" workflow
2. Click "Run workflow" 
3. Set `BUILD_AS_TEST: false` for production
4. Because our workers do not trust the fork, you must change their settings here https://github.com/organizations/lunacare/settings/actions/runner-groups/8
    - Change to "Allow public repositories"
    - Immediately after the workflow is picked up by a runner, change it back to false. The current build will still succeed.

### ECR Repository
`836236105554.dkr.ecr.us-west-2.amazonaws.com/lunacare-cosmo-router`

## Syncing with Upstream

```bash
git checkout main
git pull
git branch checkout -b main-backup-$(date +%Y%m%d%H%M%S)
git push origin main-backup-$(date +%Y%m%d%H%M%S)
git fetch upstream
git rebase upstream/main
git push origin main --force-with-lease
```

After sync, run CICD to build with new upstream version.
