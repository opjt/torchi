-- users 테이블에 OAuth provider 정보 추가
ALTER TABLE users ADD COLUMN provider TEXT;
ALTER TABLE users ADD COLUMN provider_id TEXT;

-- provider 기반 유니크 제약 추가
ALTER TABLE users ADD CONSTRAINT users_provider_uniq UNIQUE (provider, provider_id);

-- email 유니크 제약 제거 (provider_id가 식별자 역할)
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_email_key;
