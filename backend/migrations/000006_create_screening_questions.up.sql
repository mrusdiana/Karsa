CREATE TABLE screening_questions (
  id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  job_id          BIGINT UNSIGNED NOT NULL,
  question        VARCHAR(255) NOT NULL,
  type            ENUM('yes_no','number','text') NOT NULL,
  required        BOOLEAN NOT NULL DEFAULT FALSE,
  expected_answer VARCHAR(20) NULL,
  position        TINYINT UNSIGNED NOT NULL,
  CONSTRAINT fk_questions_job FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
