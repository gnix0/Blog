CREATE TABLE IF NOT EXISTS posts (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    content         TEXT         NOT NULL,
    user_id         BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tags            TEXT[]       NOT NULL DEFAULT '{}',
    cover_image_url TEXT,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
