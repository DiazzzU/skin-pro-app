CREATE TABLE IF NOT EXISTS users (
    id        BIGSERIAL PRIMARY KEY,
    login     VARCHAR(255) UNIQUE NOT NULL,
    name      VARCHAR(255) NOT NULL,
    password  VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_tokens (
    id            BIGSERIAL PRIMARY KEY,
    refresh_token VARCHAR(512) UNIQUE NOT NULL,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_user_tokens_user_id ON user_tokens(user_id);
CREATE INDEX idx_user_tokens_refresh_token ON user_tokens(refresh_token);
