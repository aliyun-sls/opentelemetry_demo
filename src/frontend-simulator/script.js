const puppeteer = require("puppeteer");
const request = require("request");

function sleep(time) {
  return new Promise((resolve) => {
    setTimeout(() => {
      resolve();
    }, time * 1000);
  });
}

const FRONTEND_URL = process.env.FRONTEND_URL || "http://frontend-proxy:8080";

async function visitProducts(page) {
  try {
    await page.goto(
      `${FRONTEND_URL}/#hot-products`
    );
//    await page.waitForNetworkIdle();
    await sleep(5);
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
    console.error(`${new Date().toISOString()}:`, error);
  }
}

async function cartList(page) {
  try {
    await page.goto(
      `${FRONTEND_URL}/cart`
    );
    await sleep(5);
  } catch (error) {
    console.error(`${new Date().toISOString()}:`, error);
  }
}

async function order(page) {
  try {
    await page.goto(
      `${FRONTEND_URL}/#hot-products`
    );
    await sleep(5);
    const randomIndex = Math.floor(Math.random() * 10) + 1;
    await page.click(`div[data-cy="product-list"]  a:nth-of-type(${randomIndex})`);
    await sleep(5);
    await page.click(`button[data-cy="product-add-to-cart"]`);
    console.info(`${new Date().toISOString()}: RunScript product-add-to-cart `);
    await sleep(5);
    let retryCount = 0;
    const maxRetries = 3;
    let success = false;
    while (retryCount < maxRetries && !success) {
      await page.click(`button[data-cy="checkout-place-order"]`);
      await sleep(5);
      const elementExists = await page.$(`button[data-cy="checkout-place-order"]`);
      if (!elementExists) {
        const title = await page.title();
        if (title.includes("Checkout")) {
           console.info(`${new Date().toISOString()}: RunScript Checkout found in title`);
           success = true;
        }
        console.info(`${new Date().toISOString()}: RunScript checkout-place-order `);
      } else {
        retryCount++;
        console.info(`${new Date().toISOString()}: Checkout not found in title, retrying... (${retryCount}/${maxRetries})`);
      }
    }
    if (!success) {
      console.error(`${new Date().toISOString()}: Failed to find "Checkout" in title after ${maxRetries} retries.`);
    }
  } catch (error) {
    console.error(`${new Date().toISOString()}:`, error);
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
    console.info(`${new Date().toISOString()}: RunScript start `);

    const browser = await launchBrowser();
    const page = await browser.newPage();
    await page.setRequestInterception(true);
    page.on('request', request => {
        if (request.url().includes('api/products') || request.url().includes('api/shipping')
         || request.url().includes('api/recommendations') || request.url().includes('api/cart')
         || request.url().includes('api/checkout')) {
          if (request.postData()) {
             console.log('Request URL:', request.url(), 'Request Post Data:', request.postData());
          }else{
             console.log('Request URL:', request.url());
          }


        }
        request.continue();
    });
//    page.on('response', async response => {
//        console.log('Response URL:', response.url());
//        console.log('Response Status:', response.status());
//    });
    await order(page);
    await order(page);
    await order(page);
    await visitProducts(page);
    await cartList(page);
    await sleep(10);
    await page.close();
    await browser.close();
    callback();
    console.info(`${new Date().toISOString()}: RunScript end `);
  } catch (error) {
    callback();
    console.error(`${new Date().toISOString()}:`, error);
  }
};

module.exports = runScript;