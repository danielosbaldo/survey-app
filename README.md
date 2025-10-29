# Survey App – Go + Gin + HTMX + GORM

A modern, configurable survey application with server-side rendering, HTMX for dynamic updates, and a full admin panel.

## Quick Start

### Development with Docker
```bash
docker compose up --build
```
Then open http://localhost:8080/form and http://localhost:8080/admin

### Local Development
```bash
make run
```

## 🎨 Customization

### Branding Your Application

The application name is **fully configurable** via environment variables, allowing you to:
- Use generic "Survey App" in the public repository
- Deploy with your actual brand name in production

**How it works:**

1. **Public Repository** - Default is "Survey App"
2. **Your Production** - Set `APP_NAME` in your `.env` file

#### For Your Production Deployment

Create a `.env` file with your branding:

```bash
# .env (on your server - NEVER commit this!)
APP_NAME=LaTienditaShop
DB_PASS=your_secure_password
DB_NAME=heladeria
# ... other settings
```

The app will automatically display "LaTienditaShop" everywhere instead of "Survey App":
- Page titles
- Headers
- Footers
- Admin panel

#### Configuration Files

- **`.env.example`** - Generic defaults for public repo
- **`.env.production.example`** - Template with your actual values (git-ignored)
- **`.env`** - Your actual configuration (git-ignored, create from template)

**Quick Setup:**
```bash
# On your server
cp .env.production.example .env
nano .env  # Edit with your values
```

## Tech Stack
- Router: Gin
- ORM: GORM (PostgreSQL)
- Frontend: Server-side rendered HTML + HTMX
- CSS: Tailwind CSS
- Admin: Full CRUD interface

## Production Deployment

See [SECURITY.md](SECURITY.md) for important security information before deploying.

### Prerequisites
- Git installed on server
- Docker and Docker Compose
- SSH access to server
- Repository on GitHub/GitLab/Bitbucket

### Deployment via Git (Recommended)

This project uses Git-based deployment for security (no file upload vulnerabilities).

1. Push code to your repository
2. SSH to your server
3. Create `/opt/myapp/.env` with your production values
4. Run the deployment script:

```bash
curl -o deploy.sh https://raw.githubusercontent.com/yourusername/survey-app/main/deploy.sh
chmod +x deploy.sh
./deploy.sh
```

## Security
- No hardcoded credentials
- Environment-based configuration
- Configurable branding (no business info in code)
- Comprehensive .gitignore
- See [SECURITY.md](SECURITY.md) for details

## Development Commands
```bash
make run              # Run locally
make compose-up       # Start dev environment
make compose-down     # Stop dev environment
make deploy           # Deploy to production
make help            # Show all commands
```

## Project Structure
```
.
├── cmd/server/          # Application entry point
├── internal/            # Internal packages
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   └── db/            # Database
├── assets/web/         # Static assets
├── .env.example        # Generic config (public)
├── .env.production.example  # Your config template (git-ignored)
├── deploy.sh          # Deployment script
└── SECURITY.md        # Security guidelines
```

## Environment Variables

### Application
- `APP_NAME` - Application display name (default: "Survey App")
- `PORT` - Internal port (default: 8080)
- `APP_PORT` - External port (default: 8080)

### Database
- `DB_HOST` - Database host (default: "db")
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: "postgres")
- `DB_PASS` - Database password (REQUIRED in production)
- `DB_NAME` - Database name (default: "myapp")
- `DB_SSLMODE` - SSL mode (default: "disable")

See `.env.example` for complete list.
