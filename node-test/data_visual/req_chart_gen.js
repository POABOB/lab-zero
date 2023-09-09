const events = require('events'), fs = require('fs'), readline = require('readline'), path = require('node:path');
const argv = require('minimist')(process.argv.slice(2));
const filename = argv._[0] ?? path.resolve('../raw_data/case-E-request-breakpoint-3.json');
const title = argv['t'] ?? `Stress Test Report - ${new Date().toLocaleString('en-US')}`;
const chartsDir = argv['d'] ?? '../charts';
if (!fs.existsSync(chartsDir)) fs.mkdirSync(chartsDir);
const htmlPath = path.resolve(chartsDir, argv['f'] ?? `req_chart-${new Date().toISOString().replace(/[-T:.Z]/g, '')}.html`);

const chartOpt = { 
  labels: [], datasets: [],
  addDataSet(label, data, color, fill = false, yAxisID = 'y1', backgroundColor = 'rgba(0,0,0,0.1)') {
    this.datasets.push({
      label, data, borderColor: color, yAxisID, backgroundColor,
      pointRadius: 0, fill, borderWidth: 2, lineType: 'line'
    });
  }
};

let baseTime = undefined;
const perfData = {};
const toHHmmss = d => d.toLocaleTimeString('en-US', { hour12: false });
(async function processLineByLine() {
  try {
    const rl = readline.createInterface({
      input: fs.createReadStream(filename), crlfDelay: Infinity
    });
    const points = {}, errors = [];
    let lineCounts = 0;
    rl.on('line', (line) => {
      lineCounts++;
      if (!line.startsWith(`{"metric":"http_req_duration"`) && !line.startsWith(`{"metric":"checks"`)) return;
      const { tags, time, value: duration } = JSON.parse(line).data;
      const { sentTime, status, target, cpu, time: perfTime } = tags;
      // count users pods orders pods
      // each cpu & memory usage
      // if (cpu !== undefined) { // cpu data
      //   perfData[perfTime.split('.')[0]] = { cpu: cpu };
      //   return;
      // }
      if (!sentTime) return;
      const logTimeNative = new Date(time), startTimeNative = new Date(sentTime);
      logTimeNative.setMilliseconds(0);
      startTimeNative.setMilliseconds(0);
      if (!baseTime) baseTime = startTimeNative;
      const startTime = toHHmmss(startTimeNative), logTime = toHHmmss(logTimeNative);
      if (!points[startTime]) points[startTime] = new Point(startTimeNative);
      if (!points[logTime]) points[logTime] = new Point(logTimeNative);
      const startTimePoint = points[startTime], logTimePoint = points[logTime];
      startTimePoint.sentCount++;
      if (!startTimePoint.target) startTimePoint.target = target;
      if (status === '200') {
        logTimePoint.succCount++;
        logTimePoint.totalSuccDura += duration;
      } else {
        logTimePoint.failCount++;
        if (!logTimePoint.errCodes[status]) logTimePoint.errCodes[status] = 1;
        logTimePoint.errCodes[status]++;
        if (tags.error) errors.push(`${logTimePoint.timespan}\t${tags.error}`);
      }
      readline.cursorTo(process.stdout, 0);
      process.stdout.write(`${lineCounts.toLocaleString('en-US')} lines processed`);
    });

    process.stdout.write("\x1B[?25l"); //hide cursor when processing
    await events.once(rl, 'close');
    process.stdout.write("\x1B[?25h"); //show cursor

    const sortedPoints = Object.keys(points).sort().map((time) => points[time]);
    const labels = [], sent = [], succ = [], fail = [], targets = [], avgDura = [], cpu = [];
    let lastCpu = 0;
    for (const { timespan, sentCount, succCount, failCount, target, time, avgSuccDura } of sortedPoints) {
      labels.push(timespan);
      sent.push(sentCount);
      succ.push(succCount);
      fail.push(failCount);
      targets.push(target);
      // lastCpu = perfData[time]?.cpu ?? lastCpu;
      // cpu.push(lastCpu);
      avgDura.push(avgSuccDura);
    }
    console.log(sortedPoints);

    chartOpt.labels = labels;
    chartOpt.addDataSet('Target (rps)', targets,'lightgray', true);
    chartOpt.addDataSet('Sent (rps)', sent, 'orange');
    chartOpt.addDataSet('Succ (rps)', succ, 'green');
    chartOpt.addDataSet('Fail (rps)', fail, 'red');
    chartOpt.addDataSet('Avg Dura (ms)', avgDura, 'blue', false, 'y2');
    // chartOpt.addDataSet('cpu (%)', cpu, 'rgba(61, 93, 122, 1)', true, 'y3', 'rgba(61, 93, 122, 0.5)');
    let html = fs.readFileSync(path.resolve('.', 'req_chart.html'), 'utf-8', 'r');
    html = html.replace(
      `<script></script>`,`<script>setChartTitle(${JSON.stringify(title)}, "${baseTime.toLocaleString('en-US')}");
let data=${JSON.stringify(chartOpt)};drawChart(data);
let errors=${JSON.stringify(errors)};listErrors(errors);
</script>`);
    fs.writeFileSync(htmlPath, html);
    console.log(`\n${sortedPoints.length} points saved`);
    var url = `file:///${htmlPath}`;
    var start = (process.platform == 'darwin' ? 'open' : process.platform == 'win32' ? 'start' : 'xdg-open');
    require('child_process').exec(start + ' ' + url);
  } catch (err) {
    console.error(`Error parsing JSON: ${err}`);
    return;
  }
})();

class Point {
  sentCount = 0;
  succCount = 0;
  failCount = 0;
  target = 0;
  totalSuccDura = 0;
  get avgSuccDura() {
    return this.succCount ? Math.round(this.totalSuccDura / this.succCount) : 0;
  }
  errCodes = {};
  constructor(time) {
    this.time = toHHmmss(time);
    this.timespan = new Date(time - baseTime).toUTCString().match(/(\d\d:\d\d:\d\d)/)[1];
  }
}