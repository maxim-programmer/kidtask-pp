import { sleep } from 'k6';
import { childIds, loginParent, loginChild, getTasks, getWishes, createTask } from './utils.js';

export const options = {
  stages: [
    { duration: '2m', target: 50 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 0 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<200'],
    http_req_failed: ['rate<0.001'],
  },
};

export default function () {
  const vuId = __VU % childIds.length;
  const childId = childIds[vuId];
  const parentEmail = `parent_${vuId}@example.com`;
  const parentPassword = 'password';
  const childUsername = `child_user_${vuId}`;
  const childPassword = 'password';
  
  const parentToken = loginParent(parentEmail, parentPassword);
  
  if (parentToken) {
    if (Math.random() < 0.7) {
      getTasks(parentToken);
      sleep(1);
    }
    
    if (Math.random() < 0.1) {
      createTask(parentToken, childId, `Load Test Task ${__ITER}`, 100);
      sleep(1);
    }
  }
  
  const childToken = loginChild(childUsername, childPassword);
  
  if (childToken) {
    if (Math.random() < 0.7) {
      getWishes(childToken, childId);
      sleep(1);
    }
  }
}
