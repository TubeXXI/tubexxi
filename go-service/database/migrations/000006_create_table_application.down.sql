BEGIN;

DROP TRIGGER IF EXISTS update_applications_modtime ON applications;

DROP FUNCTION IF EXISTS update_applications_modtime();

DROP INDEX IF EXISTS idx_applications_package_name;
DROP INDEX IF EXISTS idx_applications_package_name_key;

DROP TABLE IF EXISTS applications;

COMMIT;