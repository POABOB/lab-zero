import http from 'k6/http';
import {sleep} from 'k6';

export const options = {
  // Key configurations for spike in this section
  stages: [
    { duration: '2m', target: 2000 }, // fast ramp-up to a high point
    // No plateau
    { duration: '1m', target: 0 }, // quick ramp-down to 0 users
  ],
};

export default function() {
  let res = http.get(`http://lab.local.com/v1/users/info?user_id=1`, {
    tags: { 
      sentTime: new Date().toISOString(), 
      target: stages[getCurrentStageIndex()].target
    },
    headers: { 'Content-Type': 'application/json' }
  });
  // 蒐集 status code 200
  check(res, {
    'status is 200': (r) => r.status === 200
  });
  sleep(1);
}