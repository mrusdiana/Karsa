CREATE TABLE email_jobs (
  id           BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  to_email     VARCHAR(191) NOT NULL,
  subject      VARCHAR(500) NOT NULL,
  template     VARCHAR(50) NOT NULL,
  payload      JSON NOT NULL,
  status       ENUM('pending','sent','failed') NOT NULL DEFAULT 'pending',
  attempts     TINYINT UNSIGNED NOT NULL DEFAULT 0,
  available_at DATETIME(3) NOT NULL,
  sent_at      DATETIME(3) NULL,
  last_error   VARCHAR(500) NULL,
  created_at   DATETIME(3) NOT NULL,
  INDEX idx_email_jobs_pick (status, available_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
