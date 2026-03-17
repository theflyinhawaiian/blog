ALTER TABLE user_identities
  MODIFY COLUMN provider ENUM('google','apple','facebook','linkedin','github','twitter') NOT NULL;
