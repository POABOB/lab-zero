import http from 'k6/http';
import { check } from 'k6';
import { getCurrentStageIndex } from 'https://jslib.k6.io/k6-utils/1.3.0/index.js';
// import k8s from '@kubernetes/client-node';

// 每 10 秒增加 50 個 User，直到 RPS 到達 1000
const stages = [];
for (let t = 50; t <= 2000; t += 100) {
  stages.push({ duration: '5s', target: t }, { duration: '5s', target: t });
}
stages.push({ duration: '30s', target: 0 })

export const options = {
    systemTags: ['status','error'],                 // 標記資源
    scenarios: {
      breakpoint : {
        executor: 'ramping-arrival-rate',           // 指定時間內，執行次數可以變動，根據 stages
        preAllocatedVUs: 300,                     // 執行前，預先準備 VUs 保留運行時的資源
        timeUnit: "1s",                             // 時間單位
        stages: stages                              // 執行的每次的目標 (持續時間、VU 數量)
      },
      // monitor: {
      //   executor: 'constant-arrival-rate',          // 指定時間內，執行固定次數
      //   preAllocatedVUs: 1,                         // 執行前，預先準備 VUs 保留運行時的資源
      //   rate: 1,                                    // 每秒執行 1 次
      //   duration: 10 + stages.length * 5 + 's',    // 每個 stage * 秒數 + 20 秒 
      //   timeUnit: "1s",                             // 時間單位
      //   exec: 'monitor'                             // 執行的 func
      // }
    }
};

// const jsonPayload = `{"user_id":1}`;
export default function() {
  let res = http.get(`https://lab.local.com/v1/users/info?user_id=1`, {
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
}



// // http://localhost:8001/apis/metrics.k8s.io/v1beta1/pods
// // http://localhost:8001/apis/metrics.k8s.io/v1beta1/namespaces/default/pods
// // 在執行斷點測試之前，請下 kubectl proxy 指令，才能夠透過 http://127.0.0.1:8001 存取 apiserver
// export async function monitor() {
//     let res = http.get(`http://localhost:8001/apis/metrics.k8s.io/v1beta1/pods`);
//     check(res, {
//         'status is 200': (r) => r.status === 200
//     }, { resources: res.body });
// }