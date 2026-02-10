#!/bin/bash
# uninstall.sh - Complete removal script with metrics

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Installation directory
INSTALL_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
CONFIG_DIR="$INSTALL_DIR/config"
LOGS_DIR="$INSTALL_DIR/logs"
BACKUP_DIR="$HOME/telegram_osint_backup_$(date +%Y%m%d_%H%M%S)"

# Metrics file
METRICS_FILE="$CONFIG_DIR/installation_metrics.json"

print_banner() {
    echo -e "${RED}"
    cat << "EOF"
    ███████╗██╗   ██╗██████╗ ███████╗██████╗ ██╗██╗     ███████╗
    ██╔════╝██║   ██║██╔══██╗██╔════╝██╔══██╗██║██║     ██╔════╝
    █████╗  ██║   ██║██████╔╝█████╗  ██║  ██║██║██║     █████╗  
    ██╔══╝  ██║   ██║██╔══██╗██╔══╝  ██║  ██║██║██║     ██╔══╝  
    ██║     ╚██████╔╝██████╔╝███████╗██████╔╝██║███████╗███████╗
    ╚═╝      ╚═════╝ ╚═════╝ ╚══════╝╚═════╝ ╚═╝╚══════╝╚══════╝
                          UNINSTALLER
EOF
    echo -e "${NC}"
}

load_metrics() {
    if [[ -f "$METRICS_FILE" ]]; then
        INSTALL_DATE=$(jq -r '.installation_date' "$METRICS_FILE" 2>/dev/null || echo "Unknown")
        TOTAL_RUNTIME=$(jq -r '.user_metrics.total_runtime_seconds' "$METRICS_FILE" 2>/dev/null || echo "0")
        MESSAGES_PROCESSED=$(jq -r '.user_metrics.total_messages_processed' "$METRICS_FILE" 2>/dev/null || echo "0")
        USERS_ANALYZED=$(jq -r '.user_metrics.total_users_analyzed' "$METRICS_FILE" 2>/dev/null || echo "0")
    else
        INSTALL_DATE="Unknown"
        TOTAL_RUNTIME="0"
        MESSAGES_PROCESSED="0"
        USERS_ANALYZED="0"
    fi
}

display_metrics() {
    load_metrics
    
    # Calculate human-readable runtime
    if [[ "$TOTAL_RUNTIME" != "0" && "$TOTAL_RUNTIME" != "Unknown" ]]; then
        RUNTIME_DAYS=$((TOTAL_RUNTIME / 86400))
        RUNTIME_HOURS=$(( (TOTAL_RUNTIME % 86400) / 3600 ))
        RUNTIME_MINUTES=$(( (TOTAL_RUNTIME % 3600) / 60 ))
        RUNTIME_SECONDS=$((TOTAL_RUNTIME % 60))
        RUNTIME_STR="${RUNTIME_DAYS}d ${RUNTIME_HOURS}h ${RUNTIME_MINUTES}m ${RUNTIME_SECONDS}s"
    else
        RUNTIME_STR="No recorded usage"
    fi
    
    echo -e "${CYAN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                     UNINSTALLATION SUMMARY                   ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║                                                              ║"
    echo -e "║  ${YELLOW}📅 Installation Date:${NC} $INSTALL_DATE"
    echo -e "║  ${YELLOW}⏱️  Total Usage Time:${NC}  $RUNTIME_STR"
    echo -e "║  ${YELLOW}📨 Messages Processed:${NC} $MESSAGES_PROCESSED"
    echo -e "║  ${YELLOW}👥 Users Analyzed:${NC}     $USERS_ANALYZED"
    echo "║                                                              ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

create_backup() {
    echo -e "${YELLOW}📦 Creating backup of important data...${NC}"
    
    mkdir -p "$BACKUP_DIR"
    
    # Backup database if exists
    if [[ -f "$INSTALL_DIR/telegram_osint.db" ]]; then
        cp "$INSTALL_DIR/telegram_osint.db" "$BACKUP_DIR/"
        echo -e "${GREEN}  ✓ Database backed up${NC}"
    fi
    
    # Backup configuration
    if [[ -d "$CONFIG_DIR" ]]; then
        cp -r "$CONFIG_DIR" "$BACKUP_DIR/"
        echo -e "${GREEN}  ✓ Configuration backed up${NC}"
    fi
    
    # Backup logs
    if [[ -d "$LOGS_DIR" ]]; then
        cp -r "$LOGS_DIR" "$BACKUP_DIR/"
        echo -e "${GREEN}  ✓ Logs backed up${NC}"
    fi
    
    echo -e "${GREEN}✅ Backup created at: $BACKUP_DIR${NC}"
}

remove_installation() {
    echo -e "${RED}🗑️  Removing installation...${NC}"
    
    # List of directories and files to remove
    ITEMS_TO_REMOVE=(
        "$INSTALL_DIR/ai_agents"
        "$INSTALL_DIR/core"
        "$INSTALL_DIR/ml_models"
        "$INSTALL_DIR/cli"
        "$INSTALL_DIR/installation"
        "$INSTALL_DIR/security"
        "$INSTALL_DIR/utils"
        "$INSTALL_DIR/data"
        "$INSTALL_DIR/telegram_osint.db"
        "$INSTALL_DIR/telegram_osint.log"
        "$INSTALL_DIR/requirements.txt"
        "$INSTALL_DIR/pyproject.toml"
        "$INSTALL_DIR/Makefile"
    )
    
    for item in "${ITEMS_TO_REMOVE[@]}"; do
        if [[ -e "$item" ]]; then
            if [[ -d "$item" ]]; then
                rm -rf "$item"
                echo -e "${GREEN}  ✓ Removed directory: $(basename "$item")${NC}"
            else
                rm -f "$item"
                echo -e "${GREEN}  ✓ Removed file: $(basename "$item")${NC}"
            fi
        fi
    done
    
    # Remove Python cache files
    find "$INSTALL_DIR" -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
    find "$INSTALL_DIR" -name "*.pyc" -delete 2>/dev/null || true
    find "$INSTALL_DIR" -name "*.pyo" -delete 2>/dev/null || true
    find "$INSTALL_DIR" -name ".pytest_cache" -type d -exec rm -rf {} + 2>/dev/null || true
    
    echo -e "${GREEN}✅ Installation files removed${NC}"
}

clean_system_dependencies() {
    echo -e "${YELLOW}🧹 Cleaning system dependencies...${NC}"
    
    # Check if we should remove Python packages
    read -p "$(echo -e '${YELLOW}Remove installed Python packages? (y/N): ${NC}')" -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Removing Python packages...${NC}"
        pip uninstall -y -r "$INSTALL_DIR/requirements.txt" 2>/dev/null || true
        echo -e "${GREEN}✓ Python packages removed${NC}"
    fi
    
    # Remove desktop entries (Linux)
    if [[ -f "$HOME/.local/share/applications/telegram-osint.desktop" ]]; then
        rm -f "$HOME/.local/share/applications/telegram-osint.desktop"
        echo -e "${GREEN}✓ Desktop entry removed${NC}"
    fi
    
    # Remove from PATH (if added during installation)
    if [[ -f "$HOME/.bashrc" ]]; then
        sed -i '/telegram-osint/d' "$HOME/.bashrc" 2>/dev/null || true
    fi
    if [[ -f "$HOME/.zshrc" ]]; then
        sed -i '/telegram-osint/d' "$HOME/.zshrc" 2>/dev/null || true
    fi
    
    echo -e "${GREEN}✅ System cleanup completed${NC}"
}

record_uninstallation() {
    echo -e "${YELLOW}📝 Recording uninstallation...${NC}"
    
    if [[ -f "$METRICS_FILE" ]]; then
        # Update metrics with uninstallation info
        jq '. + {
            "uninstallation_date": "'$(date -Iseconds)'",
            "uninstallation_reason": "user_initiated",
            "final_metrics": .user_metrics
        }' "$METRICS_FILE" > "$METRICS_FILE.tmp" && mv "$METRICS_FILE.tmp" "$METRICS_FILE"
        
        # Move metrics to backup
        if [[ -d "$BACKUP_DIR" ]]; then
            cp "$METRICS_FILE" "$BACKUP_DIR/"
        fi
    fi
    
    echo -e "${GREEN}✅ Uninstallation recorded${NC}"
}

main() {
    print_banner
    
    echo -e "${RED}🚨 WARNING: This will completely remove Telegram OSINT Platform${NC}"
    echo -e "${YELLOW}This action cannot be undone!${NC}"
    echo
    
    display_metrics
    echo
    
    # Confirmation
    read -p "$(echo -e '${RED}Are you sure you want to continue? (type YES to confirm): ${NC}')" -r
    if [[ ! $REPLY == "YES" ]]; then
        echo -e "${YELLOW}Uninstallation cancelled.${NC}"
        exit 0
    fi
    
    echo
    
    # Backup step
    read -p "$(echo -e '${YELLOW}Create backup before uninstallation? (Y/n): ${NC}')" -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Nn]$ ]]; then
        create_backup
    else
        echo -e "${YELLOW}⚠️  Skipping backup${NC}"
    fi
    
    echo
    
    # Removal steps
    remove_installation
    echo
    
    # System cleanup
    clean_system_dependencies
    echo
    
    # Record uninstallation
    record_uninstallation
    echo
    
    # Final message
    echo -e "${GREEN}"
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║                    UNINSTALLATION COMPLETE                              ║"
    echo "╠══════════════════════════════════════════════════════════════╣"
    echo "║                                                                         ║"
    echo -e "║   ✅ Telegram OSINT Platform has been completely removed             ║"
    if [[ -d "$BACKUP_DIR" ]]; then
        echo -e "║   📦 Backup available at: $BACKUP_DIR                            ║"
    fi
    echo "║   🗑️  All associated files and data have been deleted                   ║"
    echo "║   🧹 System has been cleaned of dependencies                            ║"
    echo "║                                                                         ║"
    echo -e "║   ${YELLOW}Thank you for using Telegram OSINT Platform!${GREEN}      ║"
    echo "║                                                                         ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

# Run main function
main "$@"
