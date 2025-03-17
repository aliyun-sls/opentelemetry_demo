#!/bin/bash


DOCKER_REPO=o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan
DOCKER_TAG=1.0

module=ads
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}/${module}:${DOCKER_TAG} .
sudo docker push ${DOCKER_REPO}/${module}:${DOCKER_TAG}

module=notification
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}/${module}:${DOCKER_TAG} .
sudo docker push ${DOCKER_REPO}/${module}:${DOCKER_TAG}

module=promotion
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}/${module}:${DOCKER_TAG} .
sudo docker push ${DOCKER_REPO}/${module}:${DOCKER_TAG}

module=marketing
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}/${module}:${DOCKER_TAG} .
sudo docker push ${DOCKER_REPO}/${module}:${DOCKER_TAG}