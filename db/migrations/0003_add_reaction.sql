ALTER TABLE notifications ADD COLUMN actions TEXT[] NULL;
ALTER TABLE notifications ADD COLUMN reaction TEXT NULL;
ALTER TABLE notifications ADD COLUMN reaction_at TIMESTAMP NULL;