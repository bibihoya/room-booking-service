import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    vus: 50,
    duration: '30s',
};

export default function () {
    const loginRes = http.post('http://localhost:8080/dummyLogin',
        JSON.stringify({ role: 'user' }),
        { headers: { 'Content-Type': 'application/json' } }
    );

    const token = loginRes.json('token');

    if (token) {
        const roomId = 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa';
        const today = new Date().toISOString().split('T')[0];

        const res = http.get(
            `http://localhost:8080/rooms/${roomId}/slots/list?date=${today}`,
            { headers: { 'Authorization': `Bearer ${token}` } }
        );

        check(res, {
            'status is 200': (r) => r.status === 200,
        });
    }

    sleep(0.5);
}