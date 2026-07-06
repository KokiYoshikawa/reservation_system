CREATE TABLE IF NOT EXISTS reservation_slots (
    id BIGSERIAL PRIMARY KEY,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    capacity INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT uk_slots_time UNIQUE (start_time, end_time)
);

CREATE INDEX IF NOT EXISTS idx_slots_start_time ON reservation_slots (start_time);
CREATE INDEX IF NOT EXISTS idx_slots_time_range ON reservation_slots (start_time, end_time);
