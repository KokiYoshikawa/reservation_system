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

CREATE TABLE IF NOT EXISTS operation_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    reservation_id BIGINT,
    action VARCHAR(50) NOT NULL,
    detail TEXT,
    created_at TIMESTAMP NOT NULL,
    CONSTRAINT fk_operation_logs_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT fk_operation_logs_reservation
        FOREIGN KEY (reservation_id) REFERENCES reservations (id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_logs_user ON operation_logs (user_id);
CREATE INDEX IF NOT EXISTS idx_logs_reservation ON operation_logs (reservation_id);
CREATE INDEX IF NOT EXISTS idx_logs_created_at ON operation_logs (created_at);

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
    (1, date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '10 hours', date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '11 hours', 1, NOW(), NOW()),
    (2, date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '11 hours 30 minutes', date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '12 hours 30 minutes', 2, NOW(), NOW()),
    (3, date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '14 hours', date_trunc('day', NOW() + INTERVAL '6 months') + INTERVAL '15 hours', 1, NOW(), NOW()),
    (4, date_trunc('day', NOW() + INTERVAL '6 months 1 day') + INTERVAL '10 hours', date_trunc('day', NOW() + INTERVAL '6 months 1 day') + INTERVAL '11 hours', 1, NOW(), NOW()),
    (5, date_trunc('day', NOW() + INTERVAL '6 months 1 day') + INTERVAL '13 hours', date_trunc('day', NOW() + INTERVAL '6 months 1 day') + INTERVAL '14 hours', 2, NOW(), NOW()),
    (1001, date_trunc('day', NOW() + INTERVAL '6 months 2 days') + INTERVAL '10 hours', date_trunc('day', NOW() + INTERVAL '6 months 2 days') + INTERVAL '11 hours', 1, NOW(), NOW()),
    (1002, date_trunc('day', NOW() + INTERVAL '6 months 3 days') + INTERVAL '13 hours', date_trunc('day', NOW() + INTERVAL '6 months 3 days') + INTERVAL '14 hours', 1, NOW(), NOW()),
    (1003, date_trunc('day', NOW() + INTERVAL '6 months 4 days') + INTERVAL '11 hours', date_trunc('day', NOW() + INTERVAL '6 months 4 days') + INTERVAL '12 hours', 1, NOW(), NOW()),
    (1004, date_trunc('day', NOW() + INTERVAL '6 months 7 days') + INTERVAL '15 hours', date_trunc('day', NOW() + INTERVAL '6 months 7 days') + INTERVAL '16 hours', 1, NOW(), NOW())
ON CONFLICT (id) DO UPDATE
SET
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    capacity = EXCLUDED.capacity,
    updated_at = NOW();

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
    (1, 2, 1, 1, 'reserved', '初回予約です。よろしくお願いします。', NOW() - INTERVAL '3 days', NULL, NOW() - INTERVAL '3 days', NOW() - INTERVAL '3 days'),
    (2, 3, 2, 2, 'reserved', 'カラーは落ち着いた色を希望します。', NOW() - INTERVAL '2 days', NULL, NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days'),
    (3, 2, 3, 5, 'cancelled', '都合によりキャンセルしました。', NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days', NOW() - INTERVAL '10 days', NOW() - INTERVAL '8 days'),
    (1001, 2, 2, 1001, 'reserved', '敏感肌なので、施術前に相談したいです。', NOW() - INTERVAL '1 day', NULL, NOW() - INTERVAL '1 day', NOW() - INTERVAL '1 day'),
    (1002, 2, 1, 1002, 'cancelled', NULL, NOW() - INTERVAL '14 days', NOW() - INTERVAL '12 days', NOW() - INTERVAL '14 days', NOW() - INTERVAL '12 days'),
    (1003, 3, 3, 1003, 'reserved', '頭皮の乾燥が気になります。', NOW() - INTERVAL '12 hours', NULL, NOW() - INTERVAL '12 hours', NOW() - INTERVAL '12 hours'),
    (1004, 2, 1, 1004, 'reserved', '前回と同じスタイルを希望します。', NOW() - INTERVAL '2 hours', NULL, NOW() - INTERVAL '2 hours', NOW() - INTERVAL '2 hours')
ON CONFLICT (id) DO UPDATE
SET
    user_id = EXCLUDED.user_id,
    service_id = EXCLUDED.service_id,
    slot_id = EXCLUDED.slot_id,
    status = EXCLUDED.status,
    note = EXCLUDED.note,
    reserved_at = EXCLUDED.reserved_at,
    cancelled_at = EXCLUDED.cancelled_at,
    created_at = EXCLUDED.created_at,
    updated_at = NOW();

SELECT setval('users_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM users), 1));
SELECT setval('services_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM services), 1));
SELECT setval('reservation_slots_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM reservation_slots), 1));
SELECT setval('reservations_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM reservations), 1));
SELECT setval('operation_logs_id_seq', GREATEST((SELECT COALESCE(MAX(id), 1) FROM operation_logs), 1));
