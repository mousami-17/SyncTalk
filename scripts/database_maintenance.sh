#!/bin/bash

# Database Maintenance Script
# Run this script periodically (e.g., daily via cron) to maintain database health

echo "=== Database Maintenance Started at $(date) ==="

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# 1. Run VACUUM ANALYZE to optimize database
echo "Running VACUUM ANALYZE..."
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "VACUUM ANALYZE;"

if [ $? -eq 0 ]; then
    echo "✓ VACUUM ANALYZE completed"
else
    echo "✗ VACUUM ANALYZE failed"
fi

# 2. Reindex database to rebuild indexes
echo "Reindexing database..."
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "REINDEX DATABASE $POSTGRES_DATABASE;"

if [ $? -eq 0 ]; then
    echo "✓ Reindex completed"
else
    echo "✗ Reindex failed"
fi

# 3. Check for bloat and dead tuples
echo "Checking for table bloat..."
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "
SELECT 
    schemaname,
    tablename,
    n_dead_tup,
    n_live_tup,
    ROUND(n_dead_tup * 100.0 / NULLIF(n_live_tup + n_dead_tup, 0), 2) as dead_ratio
FROM pg_stat_user_tables
WHERE n_dead_tup > 1000
ORDER BY n_dead_tup DESC;
"

# 4. Check database size
echo "Database size:"
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "
SELECT 
    pg_size_pretty(pg_database_size('$POSTGRES_DATABASE')) as database_size;
"

# 5. Check table sizes
echo "Top 5 largest tables:"
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "
SELECT 
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC
LIMIT 5;
"

# 6. Check index usage
echo "Checking unused indexes..."
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan
FROM pg_stat_user_indexes
WHERE idx_scan = 0
AND indexname NOT LIKE '%_pkey';
"

# 7. Update statistics
echo "Updating statistics..."
docker exec postgres psql -U "$POSTGRES_USER" -d "$POSTGRES_DATABASE" -c "ANALYZE;"

echo "=== Database Maintenance Completed at $(date) ==="
echo ""
