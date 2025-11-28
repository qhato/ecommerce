#!/bin/bash

# Database Migration Script for E-Commerce Platform

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-ecommerce}

# Script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
SQL_FILE="$PROJECT_ROOT/database.sql"

echo -e "${GREEN}E-Commerce Database Migration Tool${NC}"
echo "===================================="
echo ""

# Check if SQL file exists
if [ ! -f "$SQL_FILE" ]; then
    echo -e "${RED}Error: database.sql not found at $SQL_FILE${NC}"
    exit 1
fi

# Function to check if database exists
database_exists() {
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -lqt | cut -d \| -f 1 | grep -qw $DB_NAME
}

# Function to create database
create_database() {
    echo -e "${YELLOW}Creating database '$DB_NAME'...${NC}"
    PGPASSWORD=$DB_PASSWORD createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    echo -e "${GREEN}Database created successfully!${NC}"
}

# Function to drop database
drop_database() {
    echo -e "${YELLOW}Dropping database '$DB_NAME'...${NC}"
    PGPASSWORD=$DB_PASSWORD dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    echo -e "${GREEN}Database dropped successfully!${NC}"
}

# Function to run migrations
run_migrations() {
    echo -e "${YELLOW}Running migrations...${NC}"
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $SQL_FILE
    echo -e "${GREEN}Migrations completed successfully!${NC}"
}

# Main script
case "$1" in
    create)
        if database_exists; then
            echo -e "${YELLOW}Database '$DB_NAME' already exists.${NC}"
        else
            create_database
        fi
        ;;
    drop)
        if database_exists; then
            read -p "Are you sure you want to drop database '$DB_NAME'? (yes/no): " confirm
            if [ "$confirm" == "yes" ]; then
                drop_database
            else
                echo "Operation cancelled."
            fi
        else
            echo -e "${YELLOW}Database '$DB_NAME' does not exist.${NC}"
        fi
        ;;
    migrate)
        if ! database_exists; then
            echo -e "${RED}Error: Database '$DB_NAME' does not exist. Create it first with: $0 create${NC}"
            exit 1
        fi
        run_migrations
        ;;
    reset)
        read -p "Are you sure you want to reset database '$DB_NAME'? All data will be lost! (yes/no): " confirm
        if [ "$confirm" == "yes" ]; then
            if database_exists; then
                drop_database
            fi
            create_database
            run_migrations
        else
            echo "Operation cancelled."
        fi
        ;;
    *)
        echo "Usage: $0 {create|drop|migrate|reset}"
        echo ""
        echo "Commands:"
        echo "  create  - Create database"
        echo "  drop    - Drop database"
        echo "  migrate - Run migrations"
        echo "  reset   - Drop, create, and migrate (fresh start)"
        echo ""
        echo "Environment variables:"
        echo "  DB_HOST     (default: localhost)"
        echo "  DB_PORT     (default: 5432)"
        echo "  DB_USER     (default: postgres)"
        echo "  DB_PASSWORD (default: postgres)"
        echo "  DB_NAME     (default: ecommerce)"
        exit 1
        ;;
esac

echo ""
echo -e "${GREEN}Done!${NC}"
