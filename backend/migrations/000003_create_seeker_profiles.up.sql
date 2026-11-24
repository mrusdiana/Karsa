CREATE TABLE seeker_profiles (
  user_id          BIGINT UNSIGNED PRIMARY KEY,
  headline         VARCHAR(150) NULL,
  phone            VARCHAR(20) NULL,
  location         VARCHAR(100) NULL,
  education        VARCHAR(10) NULL,
  major            VARCHAR(100) NULL,
  graduation_year  SMALLINT UNSIGNED NULL,
  experience_years TINYINT UNSIGNED NOT NULL DEFAULT 0,
  skills           JSON NULL,
  bio              TEXT NULL,
  portfolio_url    VARCHAR(255) NULL,
  cv_path          VARCHAR(255) NULL,
  is_discoverable  BOOLEAN NOT NULL DEFAULT FALSE,
  created_at       DATETIME(3) NOT NULL,
  updated_at       DATETIME(3) NOT NULL,
  INDEX idx_profiles_discoverable (is_discoverable, location),
  CONSTRAINT fk_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;