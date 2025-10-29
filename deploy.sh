#!/bin/bash
set -e

# Configuration
REPO_URL="${REPO_URL:-https://github.com/yourusername/go_htmx_gorm_compose.git}"
BRANCH="${BRANCH:-main}"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/heladeria}"
ENV_FILE="${ENV_FILE:-${DEPLOY_DIR}/.env}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root or with sudo
if [ "$EUID" -eq 0 ]; then
    log_warn "Running as root. This is generally not recommended."
fi

# Check if git is installed
if ! command -v git &> /dev/null; then
    log_error "Git is not installed. Please install git first."
    exit 1
fi

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    log_error "Docker is not installed. Please install docker first."
    exit 1
fi

# Check if docker compose is available
if ! docker compose version &> /dev/null; then
    log_error "Docker Compose is not available. Please install docker compose plugin."
    exit 1
fi

log_info "Starting deployment process..."

# Create deploy directory if it doesn't exist
if [ ! -d "$DEPLOY_DIR" ]; then
    log_info "Creating deployment directory: $DEPLOY_DIR"
    mkdir -p "$DEPLOY_DIR"
fi

cd "$DEPLOY_DIR"

# Check if this is the first deployment or an update
if [ -d ".git" ]; then
    log_info "Repository exists. Pulling latest changes..."

    # Stash any local changes (shouldn't be any, but just in case)
    git stash

    # Fetch and checkout the specified branch
    git fetch origin
    git checkout "$BRANCH"
    git pull origin "$BRANCH"

    log_info "Repository updated to latest commit: $(git rev-parse --short HEAD)"
else
    log_info "First deployment. Cloning repository..."

    # Remove any existing files in the directory
    rm -rf ./*

    # Clone the repository
    git clone "$REPO_URL" temp_clone

    # Move contents to current directory
    mv temp_clone/.git .
    mv temp_clone/* .
    mv temp_clone/.* . 2>/dev/null || true
    rm -rf temp_clone

    git checkout "$BRANCH"

    log_info "Repository cloned successfully"
fi

# Check if .env file exists
if [ ! -f "$ENV_FILE" ]; then
    log_warn ".env file not found at $ENV_FILE"

    if [ -f ".env.example" ]; then
        log_info "Creating .env from .env.example"
        cp .env.example .env
        log_warn "Please edit .env file with your production values before continuing"
        log_warn "Run: nano $DEPLOY_DIR/.env"
        exit 1
    else
        log_error "No .env.example found. Please create .env file manually."
        exit 1
    fi
fi

# Load environment variables
log_info "Loading environment variables..."
export $(grep -v '^#' "$ENV_FILE" | xargs)

# Build and deploy with docker compose
log_info "Building and starting containers..."

# Stop existing containers
docker compose -f docker-compose.prod.yml down

# Build new images
docker compose -f docker-compose.prod.yml build --no-cache

# Start containers
docker compose -f docker-compose.prod.yml up -d

# Wait for services to be healthy
log_info "Waiting for services to be healthy..."
sleep 10

# Check container status
if docker compose -f docker-compose.prod.yml ps | grep -q "Up"; then
    log_info "Deployment successful!"
    log_info "Application is running at: http://localhost:${APP_PORT:-8080}"

    # Show container status
    docker compose -f docker-compose.prod.yml ps

    # Clean up old images
    log_info "Cleaning up old Docker images..."
    docker image prune -f
else
    log_error "Deployment failed. Checking logs..."
    docker compose -f docker-compose.prod.yml logs --tail=50
    exit 1
fi

log_info "Deployment completed at $(date)"
