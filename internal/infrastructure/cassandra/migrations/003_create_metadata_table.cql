CREATE TABLE IF NOT EXISTS luca_s3.metadata (user_id UUID, id UUID, filename TEXT, object_key TEXT, manifest_key TEXT, video_type TEXT, size BIGINT, mime_type TEXT, duration_seconds DOUBLE, width INT, height INT, created_at TIMESTAMP, PRIMARY KEY (user_id, id)) WITH CLUSTERING ORDER BY (id ASC);
CREATE INDEX IF NOT EXISTS metadata_object_key_idx ON luca_s3.metadata (object_key);
