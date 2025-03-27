FROM --platform=${BUILDPLATFORM} eclipse-temurin:21-jdk AS builder
ARG _JAVA_OPTIONS
WORKDIR /usr/src/app/

COPY ./src/gateway/gradlew* ./src/gateway/settings.gradle* ./src/gateway/build.gradle ./
COPY ./src/gateway/gradle ./gradle

RUN chmod +x ./gradlew
RUN ./gradlew

COPY ./src/gateway/ ./
COPY ./pb/ ./proto
RUN chmod +x ./gradlew
RUN ./gradlew build -PprotoSourceDir=./proto

# -----------------------------------------------------------------------------

FROM eclipse-temurin:21-jre

ARG OTEL_JAVA_AGENT_VERSION
ARG _JAVA_OPTIONS

WORKDIR /usr/src/app/

COPY --from=builder /usr/src/app/ ./
# ADD --chmod=644 https://github.com/open-telemetry/opentelemetry-java-instrumentation/releases/download/v$OTEL_JAVA_AGENT_VERSION/opentelemetry-javaagent.jar /usr/src/app/opentelemetry-javaagent.jar
# ENV JAVA_TOOL_OPTIONS=-javaagent:/usr/src/app/opentelemetry-javaagent.jar

EXPOSE ${GATEWAY_PORT}
ENTRYPOINT [ "java", "-jar" , "./build/libs/gateway-1.0.0-SNAPSHOT.jar" ]
