CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT uk_users_email UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);

CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    duration_minutes INTEGER NOT NULL,
    price INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_services_active ON services (is_active);

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

CREATE TABLE IF NOT EXISTS reservations (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    service_id BIGINT NOT NULL,
    slot_id BIGINT NOT NULL,
    status VARCHAR(20) NOT NULL,
    note TEXT,
    reserved_at TIMESTAMP NOT NULL,
    cancelled_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_reservations_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_reservations_service
        FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE RESTRICT,
    CONSTRAINT fk_reservations_slot
        FOREIGN KEY (slot_id) REFERENCES reservation_slots (id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_reservations_user ON reservations (user_id);
CREATE INDEX IF NOT EXISTS idx_reservations_slot ON reservations (slot_id);
CREATE INDEX IF NOT EXISTS idx_reservations_status ON reservations (status);
CREATE INDEX IF NOT EXISTS idx_reservations_reserved_at ON reservations (reserved_at);
CREATE INDEX IF NOT EXISTS idx_reservations_slot_status ON reservations (slot_id, status);

INSERT INTO users (
    id,
    name,
    email,
    password_hash,
    role,
    created_at,
    updated_at
)
VALUES
    (1, '管理者ユーザー', 'admin@example.com', 'adminpass', 'admin', NOW(), NOW()),
    (2, '一般ユーザー', 'user@example.com', 'userpass', 'user', NOW(), NOW()),
    (3, 'テストユーザー', 'guest@example.com', 'guestpass', 'user', NOW(), NOW())
ON CONFLICT (email) DO NOTHING;

INSERT INTO services (
    id,
    name,
    duration_minutes,
    price,
    is_active,
    created_at,
    updated_at
)
VALUES
    (1, 'カット', 60, 5000, TRUE, NOW(), NOW()),
    (2, 'カラー', 90, 8000, TRUE, NOW(), NOW()),
    (3, 'ヘッドスパ', 45, 4500, TRUE, NOW(), NOW()),
    (4, 'パーマ', 120, 11000, FALSE, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO reservation_slots (
    id,
    start_time,
    end_time,
    capacity,
    created_at,
    updated_at
)
VALUES
    (1, '2026-07-06 10:00:00', '2026-07-06 11:00:00', 1, NOW(), NOW()),
    (2, '2026-07-06 11:30:00', '2026-07-06 12:30:00', 2, NOW(), NOW()),
    (3, '2026-07-06 14:00:00', '2026-07-06 15:00:00', 1, NOW(), NOW()),
    (4, '2026-07-07 10:00:00', '2026-07-07 11:00:00', 1, NOW(), NOW()),
    (5, '2026-07-07 13:00:00', '2026-07-07 14:00:00', 2, NOW(), NOW())
ON CONFLICT ON CONSTRAINT uk_slots_time DO NOTHING;

INSERT INTO reservations (
    id,
    user_id,
    service_id,
    slot_id,
    status,
    note,
    reserved_at,
    cancelled_at,
    created_at,
    updated_at
)
VALUES
    (1, 2, 1, 1, 'reserved', '初回予約', '2026-07-05 18:00:00', NULL, NOW(), NOW()),
    (2, 3, 2, 2, 'reserved', 'カラー希望', '2026-07-05 19:00:00', NULL, NOW(), NOW()),
    (3, 2, 3, 5, 'cancelled', '都合によりキャンセル', '2026-07-05 20:00:00', '2026-07-05 21:00:00', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

SELECT setval('users_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM users), 1));
SELECT setval('services_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM services), 1));
SELECT setval('reservation_slots_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM reservation_slots), 1));
SELECT setval('reservations_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM reservations), 1));
