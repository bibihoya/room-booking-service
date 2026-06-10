INSERT INTO users (id, email, role) VALUES
                                        ('11111111-1111-1111-1111-111111111111', 'admin@example.com', 'admin'),
                                        ('22222222-2222-2222-2222-222222222222', 'user@example.com', 'user')
    ON CONFLICT (id) DO NOTHING;

INSERT INTO rooms (id, name, description, capacity) VALUES
                                                        ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Конференц-зал', 'Большой зал', 20),
                                                        ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Переговорка', 'Маленькая', 5)
    ON CONFLICT (id) DO NOTHING;

INSERT INTO schedules (id, room_id, days_of_week, start_time, end_time) VALUES
                                                                            ('11111111-1111-1111-1111-111111111111', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', ARRAY[1,2,3,4,5], '09:00', '18:00'),
                                                                            ('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', ARRAY[1,2,3,4,5], '10:00', '19:00')
    ON CONFLICT (id) DO NOTHING;