#!/bin/bash


DOCKER_REPO=ghcr.io/aliyun-sls/demo
DOCKER_TAG=1.0

module=ads
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks

module=notification
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks

module=promotion
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks

module=marketing
sudo docker build -f Dockerfile --build-arg module=${module} --tag ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks .
sudo docker push ${DOCKER_REPO}:${DOCKER_TAG}-${module}-eks