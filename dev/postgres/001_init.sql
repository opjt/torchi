-- users
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        email TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        updated_at TIMESTAMP
    );

-- services
CREATE TABLE
    services (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        name TEXT NOT NULL,
        service_key TEXT NOT NULL UNIQUE,
        endpoint TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- push_tokens
CREATE TABLE
    push_tokens (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        p256dh_key TEXT NOT NULL,
        auth_key TEXT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT true,
        created_at TIMESTAMP NOT NULL DEFAULT now()
    );

-- service_access
CREATE TABLE
    service_access (
        id UUID PRIMARY KEY,
        service_id UUID NOT NULL REFERENCES services (id) ON DELETE CASCADE,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        UNIQUE (service_id, user_id)
    );

-- noti_history
CREATE TABLE
    noti_history (
        id UUID PRIMARY KEY,
        service_id UUID NOT NULL REFERENCES services (id) ON DELETE CASCADE,
        body TEXT NOT NULL,
        success_count INTEGER NOT NULL,
        status TEXT NOT NULL,
        sent_at TIMESTAMP NOT NULL DEFAULT now()
    );