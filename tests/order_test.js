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
  userId : "bb69958d2c61b3d8cf2003265bb88621bd6dcb0f",
  productId : "443b733e557d325b53f74c85accc4afced3036e9",
  quantity : 3
};

export default function () {
  // Request POST
  http.post('http://localhost:3000/api/orders/user/create-order', JSON.stringify(body), {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request GET
  http.get('http://localhost:3000/api/orders/user/view-order-all', {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request DELETE
  http.del('http://localhost:3000/api/orders/user/delete-order/8b311662309aa461f667a04a05e5bc28dd0eabb9', {
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

  // Request CANCEL
  http.del('http://localhost:3000/api/orders/user/cancel-order/7de907c0156ce46d4d801643ceee4d436adb92e7', {
    headers: {
      'Authorization': `Bearer ${jwtToken}`,
    },
    cookies: cookies,
  });

}
