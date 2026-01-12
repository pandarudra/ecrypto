# NPM Publishing Guide for ecrypto-cli

## Prerequisites

1. **npm account**: Create one at [npmjs.com](https://www.npmjs.com/signup)
2. **npm login**: Login to your account

```bash
npm login
```

## Publishing Steps

### 1. Navigate to npm directory

```bash
cd npm
```

### 2. Test the package locally (Optional)

```bash
# Install dependencies and test
npm pack

# This creates ecrypto-cli-1.0.7.tgz
# Test install it locally
npm install -g ./ecrypto-cli-1.0.7.tgz

# Test the command
ecrypto --version
```

### 3. Publish to npm

```bash
# Dry run first (optional - shows what will be published)
npm publish --dry-run

# Actually publish
npm publish
```

**If this is your first time publishing:**

- You might need to verify your email first
- Visit: https://www.npmjs.com/settings/YOUR_USERNAME/profile

### 4. Verify the publication

```bash
# Check on npm
npm view ecrypto-cli

# Install from npm to test
npm install -g ecrypto-cli

# Test it
ecrypto --version
```

## Updating to a New Version

When you want to release a new version:

### Option 1: Manual version bump

1. Edit `package.json` and increment the version:

   ```json
   "version": "1.0.8"
   ```

2. Edit `postinstall.js` and update the version:

   ```javascript
   const version = "v1.0.8";
   ```

3. Publish:
   ```bash
   npm publish
   ```

### Option 2: Using npm version command

```bash
# For patch release (1.0.7 -> 1.0.8)
npm version patch

# For minor release (1.0.7 -> 1.1.0)
npm version minor

# For major release (1.0.7 -> 2.0.0)
npm version major

# Then publish
npm publish
```

**Note**: After using `npm version`, don't forget to also update `postinstall.js` version manually!

## Troubleshooting

### Error: Package already exists

If you get an error that the version already exists:

```bash
# Increment version
npm version patch

# Try publishing again
npm publish
```

### Error: Need to authenticate

```bash
# Login again
npm login

# Try publishing
npm publish
```

### Error: Package name already taken

If "ecrypto-cli" is taken (unlikely, but possible):

1. Update package name in `package.json`:

   ```json
   "name": "@YOUR_USERNAME/ecrypto-cli"
   ```

2. Publish with public access:
   ```bash
   npm publish --access public
   ```

## Post-Publication Checklist

✅ Verify package on npm: https://www.npmjs.com/package/ecrypto-cli  
✅ Test installation: `npm install -g ecrypto-cli`  
✅ Test command: `ecrypto --version`  
✅ Update main README.md with new version badge  
✅ Create GitHub release with same version  
✅ Tweet/share the release! 🎉

## Package Info

- **Package name**: ecrypto-cli
- **Current version**: 1.0.7
- **Registry**: https://www.npmjs.com/
- **Installation**: `npm install -g ecrypto-cli`
- **Repository**: https://github.com/pandarudra/ecrypto

## Useful Commands

```bash
# View package info
npm view ecrypto-cli

# View all versions
npm view ecrypto-cli versions

# View downloads stats
npm view ecrypto-cli

# Unpublish a version (within 72 hours)
npm unpublish ecrypto-cli@1.0.7

# Deprecate a version
npm deprecate ecrypto-cli@1.0.7 "Use version 1.0.8 instead"
```

## GitHub Release Integration

After publishing to npm, create a GitHub release:

```bash
# Tag the release
git tag v1.0.7
git push origin v1.0.7

# Or create release on GitHub UI
# Go to: https://github.com/pandarudra/ecrypto/releases/new
```

---

## Quick Publish Checklist

- [ ] Tested locally
- [ ] Version updated in `package.json`
- [ ] Version updated in `postinstall.js`
- [ ] README.md updated
- [ ] Logged in to npm (`npm login`)
- [ ] Run `npm publish`
- [ ] Verify on npmjs.com
- [ ] Test global install
- [ ] Create GitHub release
- [ ] Update main README if needed

---

**Ready to publish?**

```bash
cd npm
npm publish
```

🎉 Your package will be live at: https://www.npmjs.com/package/ecrypto-cli
