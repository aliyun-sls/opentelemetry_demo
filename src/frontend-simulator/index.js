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
  console.info(`Add task -------------`);
  for (let index = 0; index < 3; index++) {
    scriptQueue.push({});
  }
}

run();

setInterval(() => {
  run();
}, 200000);
