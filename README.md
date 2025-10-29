# LaMichoacana – Go + Gin + HTMX + Gorm (Docker Compose)# LaMichoacana – Go + Gin + HTMX + Gorm (Docker Compose)



## Quick Start## Run with Docker

```bash

### Development with Dockerdocker compose up --build

```bash# then open http://localhost:8080/form and http://localhost:8080/admin

docker compose up --build```

# then open http://localhost:8080/form and http://localhost:8080/admin

```## Stack

- Router: **Gin**

### Local Development- ORM: **GORM** (Postgres)

```bash- Views: **SSR HTML** + **HTMX** for partial updates

make run- Admin: CRUD rápido para **Sucursales**, **Empleados**, **Preguntas** y **Opciones**

# or

go run ./cmd/server
```

## Stack
- Router: **Gin**
- ORM: **GORM** (Postgres)
- Views: **SSR HTML** + **HTMX** for partial updates
- Admin: CRUD rápido para **Sucursales**, **Empleados**, **Preguntas** y **Opciones**

---

## 🚀 Production Deployment

### Prerequisites
- Git installed on server
- Docker and Docker Compose installed
- Access to the server via SSH
- Repository hosted on GitHub/GitLab/Bitbucket

### Deployment Strategy

This project uses **Git-based deployment** instead of file copying/unzipping for security and reliability:

#### ✅ Benefits
- **Secure**: No file uploads or extraction vulnerabilities
- **Version Control**: Track deployments with Git commits
- **Rollback**: Easy to revert to previous versions
- **Atomic**: All-or-nothing deployments
- **Audit Trail**: Full history of changes

### First-Time Setup

1. **Push your code to a Git repository**
   ```bash
   git init
   git add .
   git commit -m "Initial commit"
   git remote add origin https://github.com/yourusername/go_htmx_gorm_compose.git
   git push -u origin main
   ```

2. **SSH into your server**
   ```bash
   ssh user@your-server.com
   ```

3. **Install dependencies (if not already installed)**
   ```bash
   # Install Git
   sudo apt update
   sudo apt install git -y

   # Install Docker
   curl -fsSL https://get.docker.com -o get-docker.sh
   sudo sh get-docker.sh
   sudo usermod -aG docker $USER

   # Log out and back in for group changes to take effect
   ```

4. **Create deployment directory and environment file**
   ```bash
   sudo mkdir -p /opt/heladeria
   sudo chown $USER:$USER /opt/heladeria
   cd /opt/heladeria

   # Create .env file with production values
   nano .env
   ```

   Example `.env` content:
   ```bash
   PORT=8080
   APP_PORT=8080
   DB_USER=postgres
   DB_PASS=your_secure_password_here
   DB_NAME=heladeria
   DB_SSLMODE=disable
   REPO_URL=https://github.com/yourusername/go_htmx_gorm_compose.git
   BRANCH=main
   DEPLOY_DIR=/opt/heladeria
   ```

5. **Run the deployment script**
   ```bash
   # Download and run the deployment script
   curl -o deploy.sh https://raw.githubusercontent.com/yourusername/go_htmx_gorm_compose/main/deploy.sh
   chmod +x deploy.sh
   ./deploy.sh
   ```

### Subsequent Deployments

For updates, simply run the deployment script again:

```bash
cd /opt/heladeria
./deploy.sh
```

The script will:
1. Pull the latest changes from Git
2. Build new Docker images
3. Stop old containers
4. Start new containers
5. Clean up old images

### Manual Deployment Steps

If you prefer manual control:

```bash
# Navigate to deployment directory
cd /opt/heladeria

# Pull latest changes
git pull origin main

# Build and restart containers
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml build --no-cache
docker compose -f docker-compose.prod.yml up -d

# View logs
docker compose -f docker-compose.prod.yml logs -f
```

### Using Makefile Commands

```bash
# Deploy to production
make deploy

# Test deployment locally
make deploy-local

# Start production containers
make compose-prod

# View production logs
make compose-prod-logs

# Stop production containers
make compose-prod-down

# See all available commands
make help
```

### Rollback to Previous Version

```bash
cd /opt/heladeria
git log --oneline  # Find the commit hash to rollback to
git checkout <commit-hash>
docker compose -f docker-compose.prod.yml up -d --build
```

### Health Checks

The production setup includes health checks for both database and application:

```bash
# Check container status
docker compose -f docker-compose.prod.yml ps

# Check application health
curl http://localhost:8080/form

# View detailed logs
docker compose -f docker-compose.prod.yml logs app
docker compose -f docker-compose.prod.yml logs db
```

### Security Best Practices

1. **Use strong passwords** in your `.env` file
2. **Never commit `.env`** to Git (it's in `.gitignore`)
3. **Use SSH keys** for Git authentication
4. **Configure firewall** to only allow necessary ports
5. **Use HTTPS** with a reverse proxy (nginx/caddy)
6. **Regular updates**: Keep Docker and OS updated

### Reverse Proxy Setup (Optional)

For HTTPS and domain access, use nginx or Caddy:

```nginx
# nginx example
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 📝 Development Commands

```bash
# Run locally
make run

# Start development environment
make compose-up

# Stop development environment
make compose-down

# Clean up Docker resources
make clean
```

## 🗂️ Project Structure

```
.
├── cmd/server/          # Application entry point
├── internal/
│   ├── db/             # Database connection
│   ├── handlers/       # HTTP handlers
│   ├── models/         # Data models
│   └── server/         # Server setup
├── assets/             # Embedded static files
├── deploy.sh           # Deployment script
├── docker-compose.yml       # Development config
├── docker-compose.prod.yml  # Production config
└── Dockerfile          # Production image
```
