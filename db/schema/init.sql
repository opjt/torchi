-- users
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        email TEXT NOT NULL UNIQUE,
        created_at TIMESTAMP NOT NULL DEFAULT now(),
        updated_at TIMESTAMP
    );