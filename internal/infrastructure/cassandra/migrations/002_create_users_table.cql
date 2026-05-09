CREATE TABLE IF NOT EXISTS luca_s3.users (id UUID, name TEXT, email TEXT, username TEXT, password TEXT, roles LIST<TEXT>, created_at TIMESTAMP, updated_at TIMESTAMP, PRIMARY KEY (id));
CREATE INDEX IF NOT EXISTS users_email_idx ON luca_s3.users (email);
CREATE INDEX IF NOT EXISTS users_username_idx ON luca_s3.users (username);
