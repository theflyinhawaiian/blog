ALTER TABLE comments ADD COLUMN display_name VARCHAR(255) NOT NULL DEFAULT '';
UPDATE comments c JOIN users u ON u.id = c.user_id SET c.display_name = u.display_name;
ALTER TABLE comments MODIFY COLUMN user_id BIGINT UNSIGNED NULL
