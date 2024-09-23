# Copyright The OpenTelemetry Authors
# SPDX-License-Identifier: Apache-2.0


#FROM cr-sync-registry.cn-hangzhou.cr.aliyuncs.com/otel-demo/eclipse-temurin:21-jdk as builder
#
#WORKDIR /usr/src/app/
#
#COPY ./src/adservice/gradlew* ./src/adservice/settings.gradle* ./src/adservice/build.gradle ./
#COPY ./src/adservice/gradle ./gradle
#
#RUN ./gradlew
#RUN ./gradlew downloadRepos
#
#COPY ./src/adservice/ ./
#COPY ./pb/ ./proto
#RUN ./gradlew installDist -PprotoSourceDir=./proto

# -----------------------------------------------------------------------------

FROM arms-deploy-registry.cn-hangzhou.cr.aliyuncs.com/unify-demo/eclipse-temurin:21-jdk

ARG OTEL_JAVA_AGENT_VERSION
ARG _JAVA_OPTIONS

WORKDIR /usr/src/app/

#COPY --from=builder /usr/src/app/ ./
COPY ./ ./
#ADD --chmod=644 https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v$version/opentelemetry-javaagent.jar /usr/src/app/opentelemetry-javaagent.jar
#ENV JAVA_TOOL_OPTIONS=-javaagent:/usr/src/app/opentelemetry-javaagent.jar

EXPOSE ${AD_PORT}
ENTRYPOINT [ "./build/install/opentelemetry-demo-ad/bin/Ad" ]
