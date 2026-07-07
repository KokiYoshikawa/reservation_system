DELETE FROM reservations
WHERE id IN (1, 2, 3);

DELETE FROM reservation_slots
WHERE id IN (1, 2, 3, 4, 5);

DELETE FROM services
WHERE id IN (1, 2, 3, 4);

DELETE FROM users
WHERE email IN ('admin@example.com', 'user@example.com', 'guest@example.com');

DROP INDEX IF EXISTS idx_logs_created_at;
DROP INDEX IF EXISTS idx_logs_reservation;
DROP INDEX IF EXISTS idx_logs_user;

ALTER TABLE IF EXISTS operation_logs
    DROP CONSTRAINT IF EXISTS fk_operation_logs_reservation;

ALTER TABLE IF EXISTS operation_logs
    DROP CONSTRAINT IF EXISTS fk_operation_logs_user;

DROP TABLE IF EXISTS operation_logs;

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

DROP INDEX IF EXISTS idx_slots_time_range;
DROP INDEX IF EXISTS idx_slots_start_time;

ALTER TABLE IF EXISTS reservation_slots
    DROP CONSTRAINT IF EXISTS uk_slots_time;

DROP TABLE IF EXISTS reservation_slots;

DROP INDEX IF EXISTS idx_services_active;

DROP TABLE IF EXISTS services;

DROP INDEX IF EXISTS idx_users_role;

ALTER TABLE IF EXISTS users
    DROP CONSTRAINT IF EXISTS uk_users_email;

DROP TABLE IF EXISTS users;
