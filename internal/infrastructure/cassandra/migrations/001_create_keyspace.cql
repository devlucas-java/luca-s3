CREATE KEYSPACE IF NOT EXISTS luca_s3 WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1} AND durable_writes = true;
