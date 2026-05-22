#!/bin/bash

# Database Restore Script
# This script restores a PostgreSQL database from a backup

# Check if backup file is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file.sql.gz>"
    echo ""
    echo "Available backups:"
    ls -lh ./backups/*.sql.gz 2>/dev/null || echo "No backups found"
    exit 1
fi

BACKUP_FILE="$1"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo "Error: Backup file not found: $BACKUP_FILE"
    exit 1
fi

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

echo "WARNING: This will replace the current database with the backup!"
echo "Backup file: $BACKUP_FILE"
echo ""
read -p "Are you sure you want to continue? (yes/no): " CONFIRM

if [ "$CONFIRM" != "yes" ]; then
    echo "Restore cancelled"
    exit 0
fi

echo "Starting database restore..."

# Decompress if needed
if [[ "$BACKUP_FILE" == *.gz ]]; then
    echo "Decompressing backup..."
    gunzip -c "$BACKUP_FILE" > /tmp/restore_temp.sql
    RESTORE_FILE="/tmp/restore_temp.sql"
else
    RESTORE_FILE="$BACKUP_FILE"
fi

# Drop existing database and recreate
echo "Dropping existing database..."
docker exec postgres psql -U "$POSTGRES_USER" -c "DROP DATABASE IF EXISTS $POSTGRES_DATABASE;"
docker exec postgres psql -U "$POSTGRES_USER" -c "CREATE DATABASE $POSTGRES_DATABASE;"

# Restore from backup
echo "Restoring database..."
docker exec -i postgres psql -U "$POSTGRES_USER" "$POSTGRES_DATABASE" < "$RESTORE_FILE"

if [ $? -eq 0 ]; then
    echo "Database restored successfully!"
    
    # Clean up temp file
    if [ -f /tmp/restore_temp.sql ]; then
        rm /tmp/restore_temp.sql
    fi
else
    echo "Restore failed!"
    exit 1
fi
