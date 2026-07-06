DROP INDEX IF EXISTS idx_reservations_slot_status;
DROP INDEX IF EXISTS idx_reservations_reserved_at;
DROP INDEX IF EXISTS idx_reservations_status;
DROP INDEX IF EXISTS idx_reservations_slot;
DROP INDEX IF EXISTS idx_reservations_user;

ALTER TABLE IF EXISTS reservations
    DROP CONSTRAINT IF EXISTS fk_reservations_slot;

ALTER TABLE IF EXISTS reservations
    DROP CONSTRAINT IF EXISTS fk_reservations_service;

ALTER TABLE IF EXISTS reservations
    DROP CONSTRAINT IF EXISTS fk_reservations_user;

DROP TABLE IF EXISTS reservations;
