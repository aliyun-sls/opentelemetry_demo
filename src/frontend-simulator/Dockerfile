#FROM sls-registry.cn-beijing.cr.aliyuncs.com/sls-microservices/node-puppeteer:1.0.7
FROM ghcr.io/puppeteer/puppeteer:latest
ENV PUPPETEER_SKIP_CHROMIUM_DOWNLOAD true
ENV USE_DOCKER_CHROME true

ENV PORT 5174
EXPOSE 5174


# Prepare app directory
WORKDIR /usr/src/app
COPY . /usr/src/app/

RUN npm install

# Start the app
CMD ./replace.sh ;npm run start