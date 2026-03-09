CREATE TABLE IF NOT EXISTS users (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  display_name VARCHAR(255) NOT NULL,
  created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_identities (
  id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id          BIGINT UNSIGNED NOT NULL,
  provider         ENUM('google','apple','facebook','linkedin','github') NOT NULL,
  provider_user_id VARCHAR(255) NOT NULL,
  email            VARCHAR(255),
  UNIQUE KEY uq_provider_identity (provider, provider_user_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS posts (
  id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  title            VARCHAR(512) NOT NULL,
  slug             VARCHAR(255) NOT NULL UNIQUE,
  content          LONGTEXT NOT NULL,
  excerpt          TEXT,
  meta_description VARCHAR(160),
  canonical_url    VARCHAR(2048),
  post_image       VARCHAR(2048),
  tags             JSON,
  created_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  post_id    BIGINT UNSIGNED NOT NULL,
  user_id    BIGINT UNSIGNED NOT NULL,
  content    TEXT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_reactions (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  comment_id BIGINT UNSIGNED NOT NULL,
  emoji      VARCHAR(16) NOT NULL,
  count      INT UNSIGNED NOT NULL DEFAULT 1,
  UNIQUE KEY uq_comment_emoji (comment_id, emoji),
  FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE
);
