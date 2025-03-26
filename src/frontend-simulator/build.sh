#!/bin/bash

DOCKER_REPO=o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo
DOCKER_TAG=latest


module=frontend-simulator
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

