CREATE TABLE companies (
  id               BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  user_id          BIGINT UNSIGNED NOT NULL UNIQUE,
  name             VARCHAR(150) NOT NULL,
  slug             VARCHAR(170) NOT NULL UNIQUE,
  logo_path        VARCHAR(255) NULL,
  industry         VARCHAR(100) NULL,
  size             VARCHAR(20) NULL,
  location         VARCHAR(100) NULL,
  website          VARCHAR(255) NULL,
  description      TEXT NULL,
  status           ENUM('pending','verified','rejected') NOT NULL DEFAULT 'pending',
  rejection_reason VARCHAR(500) NULL,
  verified_at      DATETIME(3) NULL,
  created_at       DATETIME(3) NOT NULL,
  updated_at       DATETIME(3) NOT NULL,
  CONSTRAINT fk_companies_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;