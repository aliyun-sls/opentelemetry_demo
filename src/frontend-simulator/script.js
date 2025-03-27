const puppeteer = require("puppeteer");
const request = require("request");

function sleep(time) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, time * 1000);
  });
}

async function visitProducts(page) {
  try {
    await page.goto(
      "http://frontend-proxy:8080/#hot-products"
    );
    await page.waitForNetworkIdle();
    await sleep(10);
    const array = new Array(3);
    let j = 0;
    for (let i of array) {
      j += 1;
      await page.click(`div[data-cy="product-list"] a:nth-of-type(${j})`);
      await sleep(5);
      await page.goBack();
      await sleep(5);
    }
  } catch (error) {
    console.error(error);
  }
}

async function cartList(page) {
  try {
    await page.goto(
      "http://frontend-proxy:8080/cart"
    );
    await sleep(5);
  } catch (error) {
    console.error(error);
  }
}

async function order(page) {
  try {
    await page.goto(
      "http://frontend-proxy:8080/#hot-products"
    );
    await sleep(5);
    const randomIndex = Math.floor(Math.random() * 10) + 1;
    await page.click(`div[data-cy="product-list"]  a:nth-of-type(${randomIndex})`);
    await sleep(5);
    await page.click(`button[data-cy="product-add-to-cart"]`);
    await sleep(5);
    await page.click(`button[data-cy="checkout-place-order"]`);

  } catch (error) {
    console.error(error);
  }
}

async function launchBrowser() {
  if (process.env.USE_DOCKER_CHROME) {
    return puppeteer.launch({
      executablePath: "google-chrome-stable",
      headless: "new",
      protocolTimeout: 240000,
      defaultViewport: {
        width: 1440,
        height: 900,
      },
      args: [
        "--no-sandbox",
        "--disable-setuid-sandbox",
        "--disable-dev-shm-usage",
        "--no-first-run",
        "--no-zygote",
      ],
    });
  }
  return puppeteer.launch({
    headless: false, //"new",
    defaultViewport: {
      width: 1440,
      height: 900,
    },
    args: [
      "--no-sandbox",
      "--disable-setuid-sandbox",
      "--disable-dev-shm-usage",
      "--no-first-run",
      "--no-zygote",
      "--ignore-gpu-blocklist",
      "--enable-gpu-rasterization",
    ],
    protocolTimeout: 240000,
  });
}

const runScript = async (callback) => {
  try {
    console.info(`RunScript start -------------`);

    const browser = await launchBrowser();
    const page = await browser.newPage();

    await visitProducts(page);
    await order(page);
    await order(page);
    await order(page);
    await cartList(page);
    await sleep(10);
    await page.close();
    await browser.close();
    callback();
    console.info(`RunScript end -------------`);
  } catch (error) {
    callback();
    console.error(error);
  }
};

module.exports = runScript;