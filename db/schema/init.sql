-- users
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        email TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        updated_at TIMESTAMP
    );

-- push_tokens
CREATE TABLE
    push_tokens (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        user_id UUID NOT NULL REFERENCES users (id),
        p256dh_key TEXT NOT NULL,
        auth_key TEXT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT true,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        UNIQUE (user_id, p256dh_key, auth_key)
    );