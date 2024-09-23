FROM arms-deploy-registry.cn-hangzhou.cr.aliyuncs.com/unify-demo/openjdk:8

# RUN apt-get update --allow-unauthenticated --allow-insecure-repositories \
#     && apt-get install -y --no-install-recommends procps net-tools wget curl vim \
#     && rm -rf /var/lib/apt/lists/*  # 清除缓存以减小镜像大小
#
# RUN apt-get update --allow-unauthenticated --allow-insecure-repositories \
#     && apt-get install -y tcpdump iputils-ping \
#     && rm -rf /var/lib/apt/lists/*  # 清除缓存以减小镜像大小
#
# RUN apt-get update --allow-unauthenticated --allow-insecure-repositories \
#     && apt-get install -y iproute2 && \
#     rm -rf /var/lib/apt/lists/*  # 清除缓存以减小镜像大小

RUN groupadd admin && useradd -m -g admin admin \
    && echo 'admin:admin' | chpasswd && usermod -u 1000 admin && groupmod -g 1000 admin

RUN mkdir -p /home/admin/sanmu-aiops-demo-server/kafka-conf
ADD ./docker/resources/mix.4096.client.truststore.jks /home/admin/sanmu-aiops-demo-server/kafka-conf/mix.4096.client.truststore.jks

RUN chown admin:admin /home/admin/ -R && chmod 755 /home/admin/ -R

WORKDIR /home/admin/sanmu-aiops-demo-server
USER admin:admin
COPY --chown=admin:admin ../arthas  /home/admin/arthas/
COPY --chown=admin:admin ../mock-server/target/sanmu-aiops-demo-server  /home/admin/sanmu-aiops-demo-server

ENV TZ=Asia/Shanghai
# COPY --chown=admin:admin ../ArmsAgent/prod-cn-hangzhou  /home/admin/ArmsAgent/
# official exp
ENV CURRENT_ENV prod
ENTRYPOINT ["sh", "./bin/bootstrap.sh"]
