FROM maven:3.9-eclipse-temurin-21-alpine AS demo_java

WORKDIR /app
COPY pom.xml .
COPY src .

RUN mvn package -DskipTests

FROM openjdk:21-jdk
ARG module
ARG port
WORKDIR /app

COPY --from=demo_java /app/${module}/target/*.jar ${module}.jar
EXPOSE ${port}
ENV SERVICE_NAME=${module}
CMD java -jar ${SERVICE_NAME}.jar
