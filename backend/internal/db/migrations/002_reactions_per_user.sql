ALTER TABLE comment_reactions
  DROP INDEX uq_comment_emoji,
  DROP COLUMN count,
  ADD COLUMN user_id BIGINT UNSIGNED NOT NULL AFTER comment_id,
  ADD UNIQUE KEY uq_comment_user_emoji (comment_id, user_id, emoji),
  ADD FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
