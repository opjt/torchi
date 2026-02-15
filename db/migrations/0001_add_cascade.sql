ALTER TABLE
    push_tokens DROP CONSTRAINT push_tokens_user_id_fkey,
ADD
    CONSTRAINT push_tokens_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE
    endpoints DROP CONSTRAINT endpoints_user_id_fkey,
ADD
    CONSTRAINT endpoints_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE
    notifications DROP CONSTRAINT notifications_user_id_fkey,
ADD
    CONSTRAINT notifications_user_id_fkey FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;