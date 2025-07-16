# LunaCare Cosmo Router Fork

This is LunaCare's fork of the [Wundergraph Cosmo Router](https://github.com/wundergraph/cosmo) with custom features and modifications.

## Custom Features

- **498 Status Code Support**: Detects GraphQL errors with 498 status codes and sets the HTTP response status accordingly
- **Custom Docker Build**: Specialized build process for LunaCare deployment needs
- **Custom Entry Point**: Uses `cmd/custom-luna/main.go` for LunaCare-specific configuration

## Versioning Strategy

Our versioning follows a two-part system:

### Main Version (Upstream Tracking)
- **Source**: Latest version from upstream `router/CHANGELOG.md`
- **Updates**: When syncing with upstream Cosmo releases
- **Current**: `0.192.1`

### LunaCare Version (Custom Features)
- **Source**: `router/VERSION-LUNACARE` file
- **Updates**: Only when WE add/modify custom features
- **Current**: `1.1.0`

### Docker Tag Format
```
{upstream_version}-lunacare-{lunacare_version}
```

**Example**: `0.192.1-lunacare-1.1.0`

This clearly shows:
- `0.192.1` = Upstream Cosmo version
- `1.1.0` = LunaCare custom features version

## When to Increment Versions

### Main Version (Automatic)
- ✅ Syncing with upstream
- ✅ Rebasing onto new upstream releases
- ❌ Never increment manually

### LunaCare Version (Manual)
- ✅ Adding new custom features
- ✅ Modifying existing custom features
- ✅ Custom configuration changes
- ❌ Just syncing with upstream

## CICD and Docker Builds

### GitHub Actions Workflow
Location: `.github/workflows/build-luna-router.yml`

### Manual Trigger (Recommended)
```bash
# Go to GitHub Actions tab
# Select "Build Luna Router" workflow  
# Click "Run workflow"
# Choose BUILD_AS_TEST: false for production
```

### Automatic Trigger
- Triggers on every push to `main` branch
- Builds production image (`latest` tag)

### Build Script
Location: `router/build-and-deploy.sh`

### Local Testing
```bash
# Build test image
./router/build-and-deploy.sh true

# Build production image  
./router/build-and-deploy.sh false
```

### ECR Repository
- **Repository**: `lunacare-cosmo-router`
- **Registry**: `836236105554.dkr.ecr.us-west-2.amazonaws.com`
- **Full URL**: `836236105554.dkr.ecr.us-west-2.amazonaws.com/lunacare-cosmo-router`

### Docker Tags Created
- **Version Tag**: `{upstream_version}-lunacare-{lunacare_version}`
- **Latest Tag**: `latest` (production builds only)
- **Test Tag**: `latest-test` (test builds only)

## Syncing with Upstream

### Clean Sync Process
```bash
# Fetch latest upstream
git fetch upstream

# Rebase your changes on top of upstream
git rebase upstream/main

# Force push (safe)
git push origin main --force-with-lease
```

### After Sync
1. ✅ Main version automatically updates from upstream CHANGELOG
2. ✅ LunaCare version stays the same (unless you added features)
3. ✅ Run CICD to build new Docker image
4. ✅ New Docker tag: `{new_upstream_version}-lunacare-1.1.0`

## File Structure

```
router/
├── VERSION-LUNACARE          # LunaCare version (manual)
├── CHANGELOG.md              # Upstream changelog (auto-updated)
├── custom-luna.Dockerfile    # Custom Docker build
├── build-and-deploy.sh       # Build script
└── cmd/custom-luna/main.go   # Custom entry point
```

## Development Workflow

1. **Sync with upstream** (as needed)
2. **Add custom features** (increment LunaCare version)
3. **Test locally** with build script
4. **Run CICD** to deploy to ECR
5. **Deploy** using new Docker tag

## Backup Strategy

Always create backup branches before major operations:
```bash
git checkout -b backup-main-$(date +%Y%m%d)
```