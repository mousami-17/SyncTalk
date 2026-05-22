# Database Management Scripts

This directory contains scripts for database maintenance, backup, and archival.

## Scripts Overview

### 1. backup_database.sh

Creates a compressed backup of the PostgreSQL database.

**Usage:**

```bash
chmod +x scripts/backup_database.sh
./scripts/backup_database.sh
```

**Features:**

- Creates timestamped backups
- Compresses backups with gzip
- Automatic cleanup of old backups (30 days retention)
- Shows backup size and recent backups

**Backup Location:** `./backups/`

---

### 2. restore_database.sh

Restores the database from a backup file.

**Usage:**

```bash
chmod +x scripts/restore_database.sh
./scripts/restore_database.sh ./backups/chat_db_backup_20251109_103045.sql.gz
```

**Warning:** This will replace the current database!

---

### 3. database_maintenance.sh

Performs routine database maintenance tasks.

**Usage:**

```bash
chmod +x scripts/database_maintenance.sh
./scripts/database_maintenance.sh
```

**Tasks Performed:**

- VACUUM ANALYZE (reclaim space, update statistics)
- REINDEX (rebuild indexes)
- Check for table bloat
- Show database and table sizes
- Identify unused indexes
- Update statistics

**Recommended:** Run daily via cron

---

### 4. archive_old_messages.go

Archives or deletes old messages from the database.

**Usage:**

```bash
# Dry run (see what would be archived)
go run scripts/archive_old_messages.go -days=90 -dry-run

# Actually archive messages older than 90 days
go run scripts/archive_old_messages.go -days=90
```

**Options:**

- `-days`: Number of days (default: 90)
- `-dry-run`: Preview without deleting

---

## Automation with Cron

### Daily Backup (2 AM)

```bash
0 2 * * * /path/to/scripts/backup_database.sh >> /var/log/db_backup.log 2>&1
```

### Weekly Maintenance (Sunday 3 AM)

```bash
0 3 * * 0 /path/to/scripts/database_maintenance.sh >> /var/log/db_maintenance.log 2>&1
```

### Monthly Archival (1st of month, 4 AM)

```bash
0 4 1 * * cd /path/to/project && go run scripts/archive_old_messages.go -days=90 >> /var/log/db_archive.log 2>&1
```

---

## Windows Users

For Windows, use Task Scheduler or run scripts manually:

**Backup:**

```powershell
docker exec postgres pg_dump -U postgres chat_db > backup.sql
```

**Restore:**

```powershell
Get-Content backup.sql | docker exec -i postgres psql -U postgres chat_db
```

---

## Best Practices

1. **Test Restores:** Regularly test your backup restoration process
2. **Off-site Backups:** Copy backups to a different location/server
3. **Monitor Disk Space:** Ensure sufficient space for backups
4. **Retention Policy:** Adjust retention days based on your needs
5. **Maintenance Window:** Run maintenance during low-traffic periods

---

## Troubleshooting

### Backup fails with "permission denied"

```bash
chmod +x scripts/*.sh
```

### Docker container not found

Ensure Docker containers are running:

```bash
docker-compose ps
```

### Restore fails

Check if backup file is corrupted:

```bash
gunzip -t backup.sql.gz
```

---

## Security Notes

- Backup files contain sensitive data - store securely
- Restrict access to backup directory
- Consider encrypting backups for production
- Never commit backups to version control

---

Last Updated: November 9, 2025
