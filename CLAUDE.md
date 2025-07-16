# Claude Context for Cosmo Router

This file provides context for AI assistants working on the LunaCare Cosmo Router fork.

## Repository Context

This is LunaCare's fork of the Wundergraph Cosmo Router with custom features and modifications. The full documentation is available in `LUNACARE_README.md`.

## Key Information

### Repository Structure
- **Upstream**: `git@github.com:wundergraph/cosmo.git` (original)
- **Fork**: `git@github.com:lunacare/cosmo.git` (our version)
- **Main Branch**: `main`

### Custom Features
- **498 Status Code Detection**: Custom HTTP response writer in `router/core/http_498_header.go`
- **Custom Docker Build**: Specialized build process for ECR deployment
- **Custom Entry Point**: `router/cmd/custom-luna/main.go`

### Versioning Strategy
```
Docker Tag Format: {upstream_version}-lunacare-{lunacare_version}
Example: 0.192.1-lunacare-1.1.0
```

**Two-part versioning system:**
1. **Main Version**: Tracks upstream releases (auto-updated from `router/CHANGELOG.md`)
2. **LunaCare Version**: Custom features version (manual in `router/VERSION-LUNACARE`)

**When to increment:**
- **Main Version**: Only when syncing with upstream (automatic)
- **LunaCare Version**: Only when adding/modifying custom features (manual)

### CICD Process

**GitHub Actions**: `.github/workflows/build-luna-router.yml`

**Manual Trigger (Recommended):**
1. Go to GitHub Actions tab
2. Select "Build Luna Router" workflow
3. Click "Run workflow"
4. Set `BUILD_AS_TEST: false` for production

**ECR Details:**
- Repository: `lunacare-cosmo-router`
- Registry: `836236105554.dkr.ecr.us-west-2.amazonaws.com`

### Common Operations

**Sync with upstream:**
```bash
git fetch upstream
git rebase upstream/main
git push origin main --force-with-lease
```

**Local Docker build:**
```bash
# Test build
./router/build-and-deploy.sh true

# Production build
./router/build-and-deploy.sh false
```

**After upstream sync:**
1. Main version auto-updates from upstream CHANGELOG
2. LunaCare version stays the same (unless custom features added)
3. Run CICD to build new Docker image with updated upstream version

### Important Files
- `router/VERSION-LUNACARE` - LunaCare version (manual updates only)
- `router/CHANGELOG.md` - Upstream changelog (auto-updated)
- `router/custom-luna.Dockerfile` - Custom Docker build
- `router/build-and-deploy.sh` - Build and deploy script
- `router/core/http_498_header.go` - 498 status code feature
- `router/cmd/custom-luna/main.go` - Custom entry point

### Development Guidelines

1. **Never increment LunaCare version** just for upstream syncs
2. **Always create backup branches** before major operations
3. **Use rebase instead of merge** for upstream syncs
4. **Test locally** before running CICD
5. **Use manual CICD trigger** for production deployments

### Backup Strategy
```bash
# Create backup before major operations
git checkout -b backup-main-$(date +%Y%m%d)
```

## For AI Assistants

When working with this repository:
- Understand the dual versioning system
- Don't increment LunaCare version unless custom features are added
- Preserve custom features (like 498 status code) during upstream syncs
- Use the rebase strategy for cleaner commit history
- Reference `LUNACARE_README.md` for complete documentation