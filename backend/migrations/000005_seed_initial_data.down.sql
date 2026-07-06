DELETE FROM reservations
WHERE id IN (1, 2, 3);

DELETE FROM reservation_slots
WHERE id IN (1, 2, 3, 4, 5);

DELETE FROM services
WHERE id IN (1, 2, 3, 4);

DELETE FROM users
WHERE email IN ('admin@example.com', 'user@example.com', 'guest@example.com');
