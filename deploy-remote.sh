#!/bin/bash
set -e

# Remote Deployment Script
# Run this from your LOCAL machine to deploy to a REMOTE server
# Usage: ./deploy-remote.sh user@server-ip

# Configuration
SERVER="$1"
DEPLOY_DIR="${DEPLOY_DIR:-/opt/myapp}"
REPO_URL="${REPO_URL:-https://github.com/yourusername/survey-app.git}"
BRANCH="${BRANCH:-main}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if server argument is provided
if [ -z "$SERVER" ]; then
    log_error "Usage: $0 user@server-ip"
    echo ""
    echo "Example:"
    echo "  $0 daniel@192.168.1.100"
    echo "  $0 root@myserver.com"
    echo ""
    echo "Optional environment variables:"
    echo "  DEPLOY_DIR - Directory on server (default: /opt/myapp)"
    echo "  REPO_URL   - Git repository URL"
    echo "  BRANCH     - Git branch to deploy (default: main)"
    exit 1
fi

log_info "Remote Deployment to: $SERVER"
log_info "Deploy directory: $DEPLOY_DIR"
log_info "Repository: $REPO_URL"
log_info "Branch: $BRANCH"
echo ""

# Test SSH connection
log_step "Testing SSH connection..."
if ! ssh -o ConnectTimeout=5 "$SERVER" "echo 'SSH connection successful'" > /dev/null 2>&1; then
    log_error "Cannot connect to $SERVER via SSH"
    echo ""
    echo "Please ensure:"
    echo "  1. The server is reachable"
    echo "  2. SSH is running on the server"
    echo "  3. You have SSH access (ssh-key or password)"
    echo ""
    echo "Test manually: ssh $SERVER"
    exit 1
fi
log_info "SSH connection successful"

# Check if deploy.sh exists locally
if [ ! -f "deploy.sh" ]; then
    log_error "deploy.sh not found in current directory"
    echo "Please run this script from the project root directory"
    exit 1
fi

# Step 1: Copy deployment script to server
log_step "Copying deploy.sh to server..."
ssh "$SERVER" "mkdir -p $DEPLOY_DIR"
scp deploy.sh "$SERVER:$DEPLOY_DIR/deploy.sh"
ssh "$SERVER" "chmod +x $DEPLOY_DIR/deploy.sh"
log_info "Deployment script copied"

# Step 2: Check if .env exists on server
log_step "Checking server configuration..."
if ! ssh "$SERVER" "[ -f $DEPLOY_DIR/.env ]"; then
    log_warn ".env file not found on server"

    # Ask if user wants to copy .env.production.example
    if [ -f ".env.production.example" ]; then
        echo ""
        read -p "Copy .env.production.example to server? (y/n) " -n 1 -r
        echo ""
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            scp .env.production.example "$SERVER:$DEPLOY_DIR/.env"
            log_warn "IMPORTANT: Edit .env on the server with actual values!"
            echo ""
            echo "Run this command to edit:"
            echo "  ssh $SERVER 'nano $DEPLOY_DIR/.env'"
            echo ""
            read -p "Press Enter to continue or Ctrl+C to abort..."
        else
            log_error "Deployment cannot continue without .env file on server"
            echo ""
            echo "Create .env on server manually:"
            echo "  ssh $SERVER"
            echo "  cd $DEPLOY_DIR"
            echo "  nano .env"
            exit 1
        fi
    else
        log_error "No .env.production.example found locally"
        exit 1
    fi
fi

# Step 3: Set environment variables on server
log_step "Configuring deployment environment..."
ssh "$SERVER" "cat > $DEPLOY_DIR/.deploy_env << 'ENVEOF'
export REPO_URL='$REPO_URL'
export BRANCH='$BRANCH'
export DEPLOY_DIR='$DEPLOY_DIR'
ENVEOF"

# Step 4: Run deployment on server
log_step "Running deployment on server..."
echo ""
ssh -t "$SERVER" "cd $DEPLOY_DIR && source .deploy_env && ./deploy.sh"

# Step 5: Cleanup temporary files
ssh "$SERVER" "rm -f $DEPLOY_DIR/.deploy_env"

echo ""
log_info "Remote deployment completed!"
echo ""
echo "Your application should now be running on:"
echo "  http://$SERVER:8080"
echo ""
echo "To view logs:"
echo "  ssh $SERVER 'cd $DEPLOY_DIR && docker compose -f docker-compose.prod.yml logs -f'"
echo ""
echo "To check status:"
echo "  ssh $SERVER 'cd $DEPLOY_DIR && docker compose -f docker-compose.prod.yml ps'"
