import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '10s', target: 20 },  // разогрев до 20 пользователей
        { duration: '30s', target: 50 },  // до 50 пользователей
        { duration: '30s', target: 100 }, // до 100 пользователей
        { duration: '10s', target: 0 },   // спад
    ],
    thresholds: {
        http_req_duration: ['p(95)<200'], // 95% запросов быстрее 200ms
        http_req_failed: ['rate<0.01'],   // ошибок менее 1%
    },
};

export default function () {
    // 1. Получаем токен пользователя
    const loginRes = http.post('http://localhost:8080/dummyLogin',
        JSON.stringify({ role: 'user' }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    check(loginRes, {
        'login status 200': (r) => r.status === 200,
    });

    const token = loginRes.json('token');

    if (!token) {
        console.error('Failed to get token');
        return;
    }

    // 2. Получаем комнаты
    const roomsRes = http.get('http://localhost:8080/rooms/list', {
        headers: { 'Authorization': `Bearer ${token}` }
    });

    check(roomsRes, {
        'rooms status 200': (r) => r.status === 200,
    });

    // 3. Получаем доступные слоты (самый нагруженный эндпоинт)
    const rooms = roomsRes.json('rooms');
    if (rooms && rooms.length > 0) {
        const roomId = rooms[0].id;
        const today = new Date().toISOString().split('T')[0];

        const slotsRes = http.get(
            `http://localhost:8080/rooms/${roomId}/slots/list?date=${today}`,
            { headers: { 'Authorization': `Bearer ${token}` } }
        );

        check(slotsRes, {
            'slots status 200': (r) => r.status === 200,
        });

        const slots = slotsRes.json('slots');

        // 4. Создаём бронь на первый свободный слот
        if (slots && slots.length > 0) {
            const slotId = slots[0].id;
            const bookingRes = http.post('http://localhost:8080/bookings/create',
                JSON.stringify({ slotId: slotId }),
                { headers: { 'Authorization': `Bearer ${token}`, 'Content-Type': 'application/json' } }
            );

            check(bookingRes, {
                'booking status 201': (r) => r.status === 201,
            });

            // 5. Отменяем бронь (идемпотентность)
            if (bookingRes.status === 201) {
                const bookingId = bookingRes.json('booking').id;
                const cancelRes = http.post(`http://localhost:8080/bookings/${bookingId}/cancel`,
                    null,
                    { headers: { 'Authorization': `Bearer ${token}` } }
                );

                check(cancelRes, {
                    'cancel status 200': (r) => r.status === 200,
                });
            }
        }
    }

    sleep(0.5);
}