CREATE TABLE invitations (
  id         BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  company_id BIGINT UNSIGNED NOT NULL,
  job_id     BIGINT UNSIGNED NOT NULL,
  user_id    BIGINT UNSIGNED NOT NULL,
  message    VARCHAR(500) NULL,
  status     ENUM('pending','accepted','declined','expired') NOT NULL DEFAULT 'pending',
  created_at DATETIME(3) NOT NULL,
  updated_at DATETIME(3) NOT NULL,
  UNIQUE INDEX uq_invitations_job_user (job_id, user_id),
  INDEX idx_invitations_user (user_id, status),
  INDEX idx_invitations_company (company_id, created_at),
  CONSTRAINT fk_invitations_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
  CONSTRAINT fk_invitations_job     FOREIGN KEY (job_id)     REFERENCES jobs(id)      ON DELETE CASCADE,
  CONSTRAINT fk_invitations_user    FOREIGN KEY (user_id)    REFERENCES users(id)     ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
