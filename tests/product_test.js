import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: 30, // Jumlah VU (Virtual Users)
  stages: [
    { duration: '10s', target: 10 },
    { duration: '10s', target: 10 },
    { duration: '10s', target: 0 },
  ],
};

// JWT token yang didapatkan dari proses login atau yang sudah Anda miliki
const jwtToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjM5MzMxMTEsInVzZXJfaWQiOiIxNzU5MGViMTE2MDZhYzgyNzgwYTZiM2RiZGMxMjhhMGE2OGEzNjlmIiwidXNlcm5hbWUiOiJhamNhanYifQ.yhpGbR1V-Ob8Y0wrMbidVrjURqfh0M2_sClfC_4duWk';

// Cookie yang berisi JWT token
const cookies = {
  cookies: {
    'auth-token': jwtToken, 
  },
};

const body = {
    Name: "Noise Headphones",
    Description: "Noise-cancelling wireless headphones with comfortable ear cups.",
    Price: 2000000,
    Stock: 50
};

export default function () {
  // Request GET
  http.get('http://localhost:3000/api/products/user/find-all', {
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request POST
  http.post('http://localhost:3000/api/products/admin/insert-product', JSON.stringify(body), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request PUT
  http.put('http://localhost:3000/api/products/admin/update-product/ee2fcd850358f6173e2be06b55f9b00dcb2b05b8', JSON.stringify(body), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request GETBYID
  http.get('http://localhost:3000/api/products/user/find-by-id/ee2fcd850358f6173e2be06b55f9b00dcb2b05b8', {
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });


    // Request DELETE
  http.del('http://localhost:3000/api/products/admin/delete-product/ee2fcd850358f6173e2be06b55f9b00dcb2b05b8', {
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

}
