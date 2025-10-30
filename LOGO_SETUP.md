# Logo Setup Instructions

## ⚠️ Logo File Required

This application requires a logo image to display properly. The logo is **NOT included in the repository** to protect your business branding.

## 📍 Logo Location

Place your logo file at:
```
assets/web/images/logo.png
```

## 📐 Recommended Specifications

- **Format**: PNG with transparency (recommended) or JPG
- **Dimensions**: 500x500 pixels (square)
- **File size**: < 500KB
- **Background**: Transparent PNG preferred

## 🚀 Setup for Development

1. **Copy your logo**:
   ```bash
   cp /path/to/your/logo.png assets/web/images/logo.png
   ```

2. **Or create a placeholder** (if you don't have a logo yet):
   ```bash
   # Create a simple colored square placeholder
   convert -size 500x500 xc:#3B82F6 -gravity center \
     -pointsize 48 -fill white -annotate +0+0 "LOGO" \
     assets/web/images/logo.png
   ```

## 🎨 Using a Different Logo File

If your logo has a different name or format:

1. Update the template references in:
   - `assets/web/templates/admin.gohtml`
   - `assets/web/templates/form.gohtml`

2. Change `logo.png` to your filename (e.g., `company-logo.svg`)

## 🔒 Security Note

The logo file is **git-ignored** to prevent:
- Exposing your business branding in a public repository
- Copyright/trademark issues
- Unwanted use of your brand assets

## 📦 Deployment

When deploying:

1. **Copy your logo to the server**:
   ```bash
   scp /path/to/your/logo.png user@server:/opt/myapp/assets/web/images/logo.png
   ```

2. **Or include in your private deployment**:
   - Keep logo in a private folder
   - Copy during deployment process
   - Add to your deployment script

## ✅ Verify Logo is Working

After adding the logo:

1. Start the application:
   ```bash
   docker compose up
   ```

2. Visit:
   - http://localhost:8080/form (should show logo)
   - http://localhost:8080/admin (should show logo)

3. Check browser console for image loading errors

## 🆘 Troubleshooting

### Logo not showing
- Check file exists: `ls -la assets/web/images/logo.png`
- Check file permissions: `chmod 644 assets/web/images/logo.png`
- Check Docker volume mounts in `docker-compose.yml`
- Clear browser cache (Ctrl+Shift+R)

### Wrong logo appears
- Ensure correct file path
- Rebuild Docker image: `docker compose build --no-cache`
- Restart containers: `docker compose restart`
