const express = require("express");
const queue = require("async").queue;

const runScript = require("./script");

const app = express();

const scriptQueue = queue((task, callback) => {
  runScript(callback);
}, 5);

const port = process.env.PORT || 5174;

const server = app.listen(port, () => {
  console.info(`Listened at http://localhost:${port}`);
});

// 外部调用检测心跳
app.get("/heartbeat", function (request, response) {
  response.status(200).json({
    code: 200,
    result: "heartbeat normally",
  });
});

function run() {
  const taskCount = process.env.TASK_COUNT || 6; // 从环境变量中获取任务数量，默认值为6
  console.info(`Add task ------------- Task count: ${taskCount}`);
  for (let index = 0; index < taskCount; index++) {
    scriptQueue.push({});
  }
}

run();

setInterval(() => {
  run();
}, 150000);
