CREATE SCHEMA IF NOT EXISTS auth_tokens;

CREATE TABLE IF NOT EXISTS service_go_gin.auth_tokens (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_auth_token_user_id ON service_go_gin.auth_tokens(user_id);
