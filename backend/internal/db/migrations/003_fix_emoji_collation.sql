DELETE FROM comment_reactions;

ALTER TABLE comment_reactions
  MODIFY COLUMN emoji VARCHAR(16) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin NOT NULL,
  DROP INDEX uq_comment_user_emoji,
  ADD UNIQUE KEY uq_comment_user_emoji (comment_id, user_id, emoji);
