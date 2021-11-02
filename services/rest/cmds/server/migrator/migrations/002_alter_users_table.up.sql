BEGIN;

ALTER TABLE users
ADD COLUMN personal_number TEXT;

COMMIT;
