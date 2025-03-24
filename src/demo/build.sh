#!/bin/bash


DOCKER_REPO=o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo
DOCKER_TAG=latest

module=ads
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

module=notification
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

module=promotion
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

module=marketing
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

module=gateway
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}