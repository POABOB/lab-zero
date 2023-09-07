import http from 'k6/http';
import {sleep} from 'k6';

export const options = {
  // Key configurations for Soak test in this section
  stages: [
    { duration: '5m', target: 500 }, // traffic ramp-up from 1 to 100 users over 5 minutes.
    { duration: '4h', target: 500 }, // stay at 100 users for 8 hours!!!
    { duration: '5m', target: 0 }, // ramp-down to 0 users
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