FROM maven:3.6.3-jdk-11-openj9 AS mall_java
ARG module
ARG port
WORKDIR /app
COPY src/${module}/pom.xml ad/
COPY src/${module}/src ad/

RUN mvn package -DskipTests

FROM openjdk:11-jdk
ARG module
ARG port
WORKDIR /app

COPY --from=mall_java /app/${module}/target/*.jar ${module}.jar
EXPOSE ${port}
ENV SERVICE_NAME ${module}
CMD java -jar ${SERVICE_NAME}.jar
