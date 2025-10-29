# Security Policy

## Overview

This document outlines security best practices for deploying and maintaining this application in production.

## Environment Variables

### Critical Rules

1. **NEVER commit `.env` files to version control**
   - The `.gitignore` file is configured to exclude all `.env` files
   - Only `.env.example` should be committed (with placeholder values)

2. **Use strong passwords in production**
   - Generate random passwords: `openssl rand -base64 32`
   - Never use default passwords like "postgres" or "changeme"

3. **Protect your `.env` file on the server**
   ```bash
   chmod 600 /opt/yourapp/.env
   chown youruser:youruser /opt/yourapp/.env
   ```

## Sensitive Information Checklist

Before pushing code to a public repository, ensure:

- [ ] No hardcoded passwords or API keys
- [ ] No real database credentials
- [ ] No personal or business identifying information
- [ ] No server IP addresses or hostnames (except localhost/examples)
- [ ] No private email addresses
- [ ] No internal URLs or domains
- [ ] `.env` files are in `.gitignore`
- [ ] All example values in `.env.example` are placeholders

## Git Security

### Before First Commit

1. **Initialize repository correctly**
   ```bash
   # Check what will be committed
   git status

   # Verify .env is ignored
   git check-ignore .env

   # Should output: .env (if properly ignored)
   ```

2. **Scan for secrets** (optional but recommended)
   ```bash
   # Install gitleaks
   brew install gitleaks  # macOS
   # or download from https://github.com/gitleaks/gitleaks

   # Scan repository
   gitleaks detect --source . --verbose
   ```

### If You Accidentally Committed Secrets

If you've already committed sensitive data:

1. **Remove from Git history** (use with caution)
   ```bash
   # Remove a specific file from all commits
   git filter-branch --force --index-filter \
     "git rm --cached --ignore-unmatch .env" \
     --prune-empty --tag-name-filter cat -- --all

   # Force push (WARNING: rewrites history)
   git push origin --force --all
   ```

2. **Immediately rotate all exposed credentials**
   - Change all passwords
   - Regenerate API keys
   - Update secrets on the server

3. **Consider the repository compromised**
   - If the repository was public, assume the secrets are known
   - Delete and recreate the repository if necessary

## Docker Security

1. **Use environment variables, not hardcoded values**
   ```yaml
   # ✅ Good
   POSTGRES_PASSWORD: ${DB_PASS}

   # ❌ Bad
   POSTGRES_PASSWORD: postgres
   ```

2. **Don't run containers as root** (future improvement)
   ```dockerfile
   RUN addgroup -g 1001 -S appuser && \
       adduser -S -u 1001 -G appuser appuser
   USER appuser
   ```

3. **Keep images updated**
   ```bash
   docker pull postgres:16-alpine
   docker compose build --no-cache
   ```

## Database Security

1. **Use strong passwords**
   - Minimum 16 characters
   - Mix of letters, numbers, and symbols
   - Unique for each environment

2. **Network isolation**
   - Database should not be directly accessible from internet
   - Use `DB_SSLMODE=require` in production if possible

3. **Regular backups**
   ```bash
   # Backup
   docker compose exec db pg_dump -U postgres myapp > backup.sql

   # Restore
   docker compose exec -T db psql -U postgres myapp < backup.sql
   ```

## Production Deployment Security

1. **Use HTTPS with valid certificates**
   - Use Let's Encrypt with nginx/Caddy
   - Never send credentials over HTTP

2. **Configure firewall**
   ```bash
   # Example with ufw
   sudo ufw allow 22/tcp    # SSH
   sudo ufw allow 80/tcp    # HTTP
   sudo ufw allow 443/tcp   # HTTPS
   sudo ufw enable
   ```

3. **Use SSH keys, not passwords**
   ```bash
   # Generate SSH key
   ssh-keygen -t ed25519 -C "your_email@example.com"

   # Copy to server
   ssh-copy-id user@server

   # Disable password authentication
   # Edit /etc/ssh/sshd_config:
   # PasswordAuthentication no
   ```

4. **Keep system updated**
   ```bash
   sudo apt update && sudo apt upgrade -y
   ```

## Application Security

1. **Input validation**
   - All user input is validated
   - Use prepared statements (GORM does this)
   - Sanitize HTML output

2. **CORS configuration**
   - Only allow trusted origins in production
   - Don't use `*` for `Access-Control-Allow-Origin`

3. **Rate limiting** (future improvement)
   - Limit API requests per IP
   - Prevent brute force attacks

## Monitoring and Logging

1. **Monitor logs for suspicious activity**
   ```bash
   docker compose logs -f app
   ```

2. **Don't log sensitive information**
   - Never log passwords or tokens
   - Redact sensitive data in logs

3. **Set up alerts**
   - Failed login attempts
   - Unusual traffic patterns
   - Resource usage spikes

## Incident Response

If you discover a security issue:

1. **Do not post publicly** until patched
2. Change all affected credentials immediately
3. Review logs for unauthorized access
4. Patch the vulnerability
5. Document what happened and how it was fixed

## Regular Security Maintenance

- [ ] Review and rotate credentials quarterly
- [ ] Update dependencies monthly: `go get -u ./...`
- [ ] Update Docker images weekly
- [ ] Review access logs monthly
- [ ] Test backups monthly
- [ ] Review `.gitignore` when adding new files

## Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Docker Security Best Practices](https://docs.docker.com/develop/security-best-practices/)
- [Go Security Guidelines](https://golang.org/doc/security)
- [Git Secrets Detection](https://github.com/gitleaks/gitleaks)

## Reporting Security Issues

If you find a security vulnerability:

1. **Do NOT open a public issue**
2. Contact the maintainer privately
3. Provide detailed information about the vulnerability
4. Allow time for the issue to be patched before disclosure

## License

This security policy is part of the project and follows the same license.
