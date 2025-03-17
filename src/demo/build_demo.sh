#!/bin/bash


DOCKER_REPO=ghcr.io/aliyun-sls/demo
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