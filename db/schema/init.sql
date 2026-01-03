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
        endpoint TEXT NOT NULL,
        is_active BOOLEAN NOT NULL DEFAULT true,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        UNIQUE (user_id, p256dh_key, auth_key, endpoint)
    );

-- endpoints
CREATE TABLE
    endpoints (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        user_id UUID NOT NULL REFERENCES users (id),
        name TEXT NOT NULL,
        token TEXT NOT NULL UNIQUE,
        notification_enabled BOOLEAN NOT NULL DEFAULT true,
        notification_disabled_at TIMESTAMP NULL,
        created_at TIMESTAMP NOT NULL DEFAULT now(),

        CONSTRAINT endpoints_user_name_uniq
        UNIQUE (user_id, name)
    );

-- notifications
CREATE TABLE
    notifications (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        service_id UUID NOT NULL REFERENCES endpoints (id),
        body TEXT NOT NULL,
        is_read BOOLEAN NOT NULL DEFAULT false, 
        read_at TIMESTAMP NULL,
        is_deleted BOOLEAN NOT NULL DEFAULT false,
        created_at TIMESTAMP NOT NULL DEFAULT now()
    );