const k8s = require('@kubernetes/client-node');
const fs = require('fs');
const kubeSystemNamespcae = 'default';
const kubeSystemNamespcae2 = 'istio-system';
const kc = new k8s.KubeConfig();
kc.loadFromDefault();
const k8sApi = kc.makeApiClient(k8s.CoreV1Api);
const metricsClient = new k8s.Metrics(kc);

const baseTime = new Date().toISOString();
let line = 0
async function monitor() {
    try {
        const sentTime = new Date().toISOString();
        const topPodsRes1 = await k8s.topPods(k8sApi, metricsClient, kubeSystemNamespcae);
        const topPodsRes2 = await k8s.topPods(k8sApi, metricsClient, kubeSystemNamespcae2);
        const podsColumns = topPodsRes1.map((pod) => {
            const mem = (pod.Memory.CurrentUsage.toString()).slice(0, -1)
            return {
                sentTime: sentTime,
                POD: pod.Pod.metadata.name,
                'CPU': parseFloat((pod.CPU.CurrentUsage * 1024).toFixed(2)),
                'MEMORY': parseFloat((mem / (1024 * 1024) * 10).toFixed(2)),
            };
        });
        const podsColumns2 = topPodsRes2.map((pod) => {
            const mem = (pod.Memory.CurrentUsage.toString()).slice(0, -1)
            return {
                sentTime: sentTime,
                POD: pod.Pod.metadata.name,
                'CPU': parseFloat((pod.CPU.CurrentUsage * 1024).toFixed(2)),
                'MEMORY': parseFloat((mem / (1024 * 1024) * 10).toFixed(2)),
            };
        });
        // console.table(podsColumns);
        const data = JSON.stringify(podsColumns, (key, value) =>
            typeof value === 'bigint'
                ? value.toString()
                : value
        )
        const data2 = JSON.stringify(podsColumns2, (key, value) =>
            typeof value === 'bigint'
                ? value.toString()
                : value
        )
        // append JSON string to a file
        fs.appendFile(`case-D-resouce-breakpoint-1.json`, data + "\r\n" , function (err) {
            if (err) throw err;
            console.log(`The "data to append" was appended to file! line: ${line}, datetime: ${sentTime}.`);
        });
        line++

        fs.appendFile(`case-D-resouce-breakpoint-1.json`, data2 + "\r\n" , function (err) {
            if (err) throw err;
            console.log(`The "data to append" was appended to file! line: ${line}, datetime: ${sentTime}.`);
        });
        line++
    } catch (err) {
        console.error(err);
    }
}

setInterval(monitor, 1000);