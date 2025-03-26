#!/bin/bash
docker build -f src/demo/goApp/Dockerfile -t o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo:latest-goapp .

docker push o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo:latest-goapp


docker build -f src/demo/product/Dockerfile -t o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo:latest-product .
docker push o11y-demo-cn-heyuan-registry.cn-heyuan.cr.aliyuncs.com/o11y-demo-cn-heyuan/demo:latest-product