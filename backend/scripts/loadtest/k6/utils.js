import http from 'k6/http';
import { check } from 'k6';

export const BASE_URL = 'http://localhost:8080';

export const childIds = [
  'da50329a-2ec4-4ec5-b628-00d218489dfa',
  '72ca378e-2b0a-4185-acea-2fa95860c2ad',
  'd3f3904a-6436-436f-a46b-07ca1b9bb33a',
  '1c17becb-70b3-41c0-b8c1-dcc6ef4fee58',
  '7ac3e86a-1d8d-49fd-8fdd-4232c0c81178',
];

export function loginParent(email, password) {
  const res = http.post(`${BASE_URL}/api/auth/login`, JSON.stringify({
    email: email,
    password: password,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(res, { 'Parent login success': (r) => r.status === 200 });
  
  if (res.status === 200) {
    return res.json('token');
  }
  return null;
}

export function loginChild(username, password) {
  const res = http.post(`${BASE_URL}/api/auth/child/login`, JSON.stringify({
    username: username,
    password: password,
  }), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(res, { 'Child login success': (r) => r.status === 200 });
  
  if (res.status === 200) {
    return res.json('token');
  }
  return null;
}

export function getTasks(token) {
  const res = http.get(`${BASE_URL}/api/tasks`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  
  check(res, { 'Get tasks success': (r) => r.status === 200 });
  return res;
}

export function getWishes(token, childId) {
  const res = http.get(`${BASE_URL}/api/children/${childId}/wishes`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  
  check(res, { 'Get wishes success': (r) => r.status === 200 });
  return res;
}

export function createTask(token, childId, title, reward) {
  const res = http.post(`${BASE_URL}/api/tasks`, JSON.stringify({
    child_id: childId,
    title: title,
    reward: reward,
  }), {
    headers: { 
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  });
  
  check(res, { 'Create task success': (r) => r.status === 201 });
  return res;
}
