import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 50,
  duration: '1m',
  thresholds: {
    http_req_duration: ['p(95)<2000'],
  },
};

export default function () {
  const res = http.get('http://localhost:8080/api/productos', {
    headers: { 'X-User-Role': 'Cliente' },
  });
  check(res, { 'status es 200': (r) => r.status === 200 });
  sleep(1);
}