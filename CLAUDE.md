# Claude Context for Cosmo Router

This file provides context for AI assistants working on the LunaCare Cosmo Router fork.

## Repository Overview

This is LunaCare's fork of the Wundergraph Cosmo Router with custom features. **See `LUNACARE_README.md` for complete documentation.**
@LUNACARE_README.md

### Repository Structure
- **Upstream**: `git@github.com:wundergraph/cosmo.git` (original)
- **Fork**: `git@github.com:lunacare/cosmo.git` (our version)
- **Main Branch**: `main`

### Key Files
- `LUNACARE_README.md` - Complete documentation
- `router/VERSION-LUNACARE` - LunaCare version (manual updates only)
- `router/core/http_498_header.go` - 498 status code feature
- `router/cmd/custom-luna/main.go` - Custom entry point
- `router/custom-luna.Dockerfile` - Custom Docker build
- `router/build-and-deploy.sh` - Build and deploy script

### Important Guidelines for AI Assistants

1. **Don't increment LunaCare version** just for upstream syncs
2. **Preserve custom features** (like 498 status code) during upstream syncs
3. **Use rebase instead of merge** for upstream syncs for cleaner history
4. **Create backup branches** before major operations: `git checkout -b backup-main-$(date +%Y%m%d)`

### Quick Reference

**Sync with upstream:**
```bash
git fetch upstream
git rebase upstream/main
git push origin main --force-with-lease
```

**Build Docker image:**
Manual trigger via GitHub Actions → "Build Luna Router" workflow

**Versioning:**
- Upstream version: Auto-pulled from `router/CHANGELOG.md`
- LunaCare version: Manual in `router/VERSION-LUNACARE`
- Docker tag: `{upstream_version}-lunacare-{lunacare_version}`
