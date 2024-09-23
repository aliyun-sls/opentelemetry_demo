#!/bin/bash -e
if [ -z "$1" ]
then
  echo "tag is empty,use default one"
  TAG=helm-unify-demo
else
  TAG=$1
fi
echo "TAG: $TAG"

gitMessage=`git show -s --format=%s`
gitBranch=`git branch --show-current`
gitCommit=`git rev-parse --short HEAD`
commitId=${gitBranch}_${gitCommit}
tag=`whoami`_${commitId}

DOCKER_IMAGE=arms-deploy-registry.cn-hangzhou.cr.aliyuncs.com/unify-demo/arms-mock-server
echo imageName=${DOCKER_IMAGE}
mvn clean package -am -Dmaven.test.skip=true
docker build --platform linux/amd64 -t ${DOCKER_IMAGE}:${tag} --build-arg ENV="prod" .
docker push ${DOCKER_IMAGE}:${tag}
