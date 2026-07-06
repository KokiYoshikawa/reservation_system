DROP INDEX IF EXISTS idx_slots_time_range;
DROP INDEX IF EXISTS idx_slots_start_time;

ALTER TABLE IF EXISTS reservation_slots
    DROP CONSTRAINT IF EXISTS uk_slots_time;

DROP TABLE IF EXISTS reservation_slots;
