CREATE TABLE tags (
  id   BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  UNIQUE KEY uq_tag_name (name)
);

CREATE TABLE post_tags (
  post_id BIGINT UNSIGNED NOT NULL,
  tag_id  BIGINT UNSIGNED NOT NULL,
  PRIMARY KEY (post_id, tag_id),
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id)  REFERENCES tags(id)  ON DELETE CASCADE
);

INSERT INTO tags (name)
SELECT DISTINCT jt.tag
FROM posts p,
     JSON_TABLE(p.tags, '$[*]' COLUMNS (tag VARCHAR(100) PATH '$')) jt
WHERE p.tags IS NOT NULL AND JSON_LENGTH(p.tags) > 0
ON DUPLICATE KEY UPDATE name = name;

INSERT INTO post_tags (post_id, tag_id)
SELECT p.id, t.id
FROM posts p,
     JSON_TABLE(p.tags, '$[*]' COLUMNS (tag VARCHAR(100) PATH '$')) jt
JOIN tags t ON t.name = jt.tag
WHERE p.tags IS NOT NULL AND JSON_LENGTH(p.tags) > 0;

ALTER TABLE posts DROP COLUMN tags;
