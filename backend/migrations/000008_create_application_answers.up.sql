CREATE TABLE application_answers (
  id             BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  application_id BIGINT UNSIGNED NOT NULL,
  question       VARCHAR(255) NOT NULL,
  required       BOOLEAN NOT NULL,
  answer         VARCHAR(500) NOT NULL,
  passed         BOOLEAN NOT NULL,
  CONSTRAINT fk_answers_application FOREIGN KEY (application_id) REFERENCES applications(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
