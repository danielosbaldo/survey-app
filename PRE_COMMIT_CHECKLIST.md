# Pre-Commit Security Checklist

Use this checklist before committing code to your public repository.

## ✅ Quick Verification

Run these commands to verify your repository is safe:

```bash
# 1. Check for .env files (should return nothing)
git status | grep "\.env$"

# 2. Verify .env is ignored
git check-ignore .env
# Should output: .env

# 3. Check for hardcoded passwords
grep -r "password.*postgres" --include="*.yml" --include="*.go" . | grep -v "changeme\|your_secure\|:-"
# Should return nothing or only safe defaults

# 4. Check for business/personal names
grep -r "your-company-name-here" --include="*.go" --include="*.yml" .
# Replace "your-company-name-here" with your actual business name

# 5. Review files to be committed
git add -A
git status
# Verify no .env files are listed
```

## 📋 Manual Checklist

Before pushing to GitHub/GitLab/Bitbucket:

### Environment & Credentials
- [ ] No `.env` files in git status
- [ ] `.env.example` contains only placeholder values
- [ ] No hardcoded passwords in code
- [ ] No API keys or tokens in code
- [ ] Database credentials use environment variables

### Personal Information
- [ ] No real email addresses
- [ ] No phone numbers
- [ ] No physical addresses
- [ ] No personal names (except in license/credits)

### Business Information
- [ ] No company-specific names in code (unless intended)
- [ ] No internal URLs or domains
- [ ] No customer/client names
- [ ] No proprietary business logic (if applicable)

### Server Information
- [ ] No production server IPs
- [ ] No internal hostnames
- [ ] Only localhost/example.com in examples
- [ ] Deployment paths are generic

### Code Quality
- [ ] `go.mod` has correct module path
- [ ] Import paths match module path
- [ ] Code compiles: `go build ./...`
- [ ] Tests pass: `go test ./...`
- [ ] No debug print statements with sensitive data

### Documentation
- [ ] README has no sensitive info
- [ ] SECURITY.md is up to date
- [ ] Comments don't reveal internal details
- [ ] Example values are clearly placeholders

## 🔧 Fix Common Issues

### If .env is showing in git status:

```bash
# Remove from staging
git rm --cached .env

# Verify it's in .gitignore
echo ".env" >> .gitignore

# Commit the .gitignore change
git add .gitignore
git commit -m "Add .env to gitignore"
```

### If you committed secrets by accident:

```bash
# Remove file from all history (USE WITH CAUTION)
git filter-branch --force --index-filter \
  "git rm --cached --ignore-unmatch .env" \
  --prune-empty --tag-name-filter cat -- --all

# Force push (WARNING: Rewrites history)
git push origin --force --all

# IMPORTANT: Rotate all exposed credentials immediately!
```

### If module path needs updating:

```bash
# Update go.mod
sed -i 's|old-module-path|new-module-path|g' go.mod

# Update all imports
find . -type f -name "*.go" -exec sed -i 's|old-module-path|new-module-path|g' {} +

# Tidy dependencies
go mod tidy
```

## 🎯 Best Practices

1. **Review diffs before committing:**
   ```bash
   git diff
   git diff --staged
   ```

2. **Use .gitignore templates:**
   - Start with a good .gitignore (already included)
   - Add project-specific patterns as needed

3. **Scan for secrets (optional):**
   ```bash
   # Install gitleaks
   brew install gitleaks  # macOS

   # Scan repository
   gitleaks detect --source . --verbose
   ```

4. **Use environment variables:**
   - Never hardcode credentials
   - Always use `os.Getenv()` with safe defaults
   - Document required env vars in README

5. **Regular audits:**
   - Review this checklist monthly
   - Update SECURITY.md as needed
   - Check for new sensitive data patterns

## 🚨 If You Find Sensitive Data

### In Your Own Code:
1. Remove it immediately
2. Use git filter-branch if already committed
3. Rotate any exposed credentials
4. Update .gitignore to prevent recurrence

### In Someone Else's Code:
1. DO NOT open a public issue
2. Contact maintainers privately
3. Provide file and line numbers
4. Suggest appropriate fixes

## 📚 Additional Resources

- [SECURITY.md](./SECURITY.md) - Comprehensive security guidelines
- [.gitignore](./.gitignore) - Files excluded from git
- [Gitleaks](https://github.com/gitleaks/gitleaks) - Secret scanning tool

---

**Remember:** Once sensitive data is pushed to a public repository, assume it's compromised. Prevention is key!
