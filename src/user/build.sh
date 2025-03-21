#!/bin/bash


DOCKER_REPO=o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo
DOCKER_TAG=latest

module=user
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}

DOCKER_REPO=ghcr.io/aliyun-sls/demo
DOCKER_TAG=1.0

module=user
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module} .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}