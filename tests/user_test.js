import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 5, // Jumlah VU (Virtual Users)
  stages: [
    { duration: '10s', target: 10 },
    { duration: '10s', target: 10 },
    { duration: '10s', target: 0 },
  ],
};


const body = {
  username: 'admin',
  password: '@1Bkabck',
};

export default function () {
// Request GET
http.post('http://localhost:3000/api/users/register',JSON.stringify(body), {
  headers: {
    'Content-Type': 'application/json',
  }
});


// Request GET
http.post('http://localhost:3000/api/users/login',JSON.stringify(body, {
  Headers: {
    'Content-Type': 'application/json',
  }
}));

}