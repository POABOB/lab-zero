const events = require('events'), fs = require('fs'), readline = require('readline'), path = require('node:path');
const argv = require('minimist')(process.argv.slice(2));
const filename = argv._[0] ?? '../raw_data/case-E-resouce-breakpoint-3.json';
const title = argv['t'] ?? `Stress Test Report - ${new Date().toLocaleString('en-US')}`;
const chartsDir = argv['d'] ?? 'charts';
if (!fs.existsSync(chartsDir)) fs.mkdirSync(chartsDir);
const htmlPath = path.resolve(chartsDir, argv['f'] ?? `resource_chart-${new Date().toISOString().replace(/[-T:.Z]/g, '')}.html`);

const chartOpt = { 
  labels: [], datasets: [],
  addDataSet(label, data, color, fill = false, yAxisID = 'y1', backgroundColor = 'rgba(0,0,0,0.1)') {
    this.datasets.push({
      label, data, borderColor: color, yAxisID, backgroundColor,
      pointRadius: 0, fill, borderWidth: 2, lineType: 'line'
    });
  }
};

// TODO 要自己定義
let baseTime = "2023-09-09T10:35:30.508Z";
let pointsArr = [], errors = [];
const toHHmmss = d => d.toLocaleTimeString('en-US', { hour12: false });
(async function processLineByLine() {
  try {
    const rl = readline.createInterface({
      input: fs.createReadStream(filename), crlfDelay: Infinity
    });
    
    let lineCounts = 0;
    rl.on('line', (line) => {
      lineCounts++;

      // 每行都有不同的pods
      const lineArr = JSON.parse(line)
      let points = [];
      let periodTime;
      lineArr.forEach(element => {
        const { sentTime, POD, CPU, MEMORY } = element
        periodTime = sentTime
        points.push({
          POD: POD,
          CPU: CPU,
          MEMORY: MEMORY,
        })
      });
      pointsArr.push({sentTime: periodTime, data: points})
      // console.log(pointsArr);

      const startTimeNative = new Date(periodTime);
      startTimeNative.setMilliseconds(0);
      const startTime = toHHmmss(startTimeNative);

      readline.cursorTo(process.stdout, 0);
      process.stdout.write(`${lineCounts.toLocaleString('en-US')} lines processed`);
    });

    process.stdout.write("\x1B[?25l"); //hide cursor when processing
    await events.once(rl, 'close');
    process.stdout.write("\x1B[?25h"); //show cursor

    baseTime = new Date(baseTime);
    baseTime.setMilliseconds(0);

    let labels = [], podsNums = [], pods = [[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[],[]], podLabel = [], podColors = ['orange','green','red','blue','#93FF93','#9393FF','#4DFFFF','#467500','#737300','#BB5E00','#B87070','#6FB7B7','#3A006F','#EBD3E8','#D1E9E9','#EBD6D6','#C4E1E1','#E6E6F2','#FFF4C1','#FFE6D9','#FFD9EC','#F1E1FF','#DDDDFF','#E8FFF5']
    pointsArr.forEach(element => {
      // TODO 低於 baseTime不紀錄
      const timespan = remain_time_string(new Date(element.sentTime), baseTime)
      if(timespan.slice(0,1) === "-") return;
      labels.push(timespan)
      podsNums.push(element.data.length)
      element.data.forEach(data => {
        let podName = data.POD.split("-")[0] + "-" + data.POD.split("-")[3];
        if(podName.slice(0, 1) === "l") podName = data.POD.split("-")[4] + "-" + data.POD.split("-")[2];
        let index = podLabel.findIndex(d => d === podName)
        if(index === -1) {
          podLabel.push(podName)
          index= podLabel.length - 1
          for(let i = 0; i < labels.length - 1; i++) pods[index].push(null)
        }
        pods[index].push(data.CPU)

        // podName = data.POD + "-MEMORY"
        // index = podLabel.findIndex(d => d === podName)
        // if(index === -1) {
        //   podLabel.push(podName)
        //   index= podLabel.length - 1
        //   for(let i = 0; i < labels.length - 1; i++) pods[index].push(0)
        // }
        // pods[index].push(data.MEMORY)
      })
    })

    chartOpt.labels = labels;
    chartOpt.addDataSet('Pod nums', podsNums,'lightgray', true, 'y3');
    
    pods.forEach((pod, index) => {
      if(pod.length > 0) {
        if(podLabel[index].slice(0, 1) === "u") { 
          chartOpt.addDataSet(`${podLabel[index]}`, pod, podColors[index], false, 'y1');
        } else if(podLabel[index].slice(0, 1) === "o") {
          chartOpt.addDataSet(`${podLabel[index]}`, pod, podColors[index], false, 'y2');
        } else if(podLabel[index].slice(0, 1) === "w") {
          chartOpt.addDataSet(`${podLabel[index]}`, pod, podColors[index], false, 'y4');
        }
      }
    })

    let html = fs.readFileSync(path.resolve('.', 'resource_chart.html'), 'utf-8', 'r');
    html = html.replace(
      `<script></script>`,`<script>setChartTitle(${JSON.stringify(title)}, "${baseTime.toLocaleString('en-US')}");
let data=${JSON.stringify(chartOpt)};drawChart(data);
let errors=${JSON.stringify(errors)};listErrors(errors);
</script>`);
    fs.writeFileSync(htmlPath, html);
    console.log(`\n${pointsArr.length} points saved`);
    var url = `file:///${htmlPath}`;
    var start = (process.platform == 'darwin' ? 'open' : process.platform == 'win32' ? 'start' : 'xdg-open');
    require('child_process').exec(start + ' ' + url);
  } catch (err) {
    console.error(`Error parsing JSON: ${err}`);
    return;
  }
})();

function remain_time_string (startDate, endDate) {
  var diff = (startDate.getTime() - endDate.getTime());
  var hours = Math.floor(diff / 1000 / 60 / 60);
  diff -= hours * 1000 * 60 * 60;
  var minutes = Math.floor(diff / 1000 / 60);
  diff -= minutes * 1000 * 60;
  var seconds = Math.floor(diff / 1000);
  resultTime = hours.toString().padStart(2, '0') + ":" + minutes.toString().padStart(2, '0') + ":" + seconds.toString().padStart(2, '0');
  return resultTime
}

function getRandomColor() {
    var letters = '0123456789ABCDEF'.split('');
    var color = '#';
    for (var i = 0; i < 6; i++ ) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
}