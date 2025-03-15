FROM maven:3.6.3-jdk-11-openj9 AS demo_java

WORKDIR /app
COPY pom.xml .
COPY src .

RUN mvn package -DskipTests

FROM openjdk:11-jdk
ARG module
ARG port
WORKDIR /app

COPY --from=demo_java /app/${module}/target/*.jar ${module}.jar
EXPOSE ${port}
ENV SERVICE_NAME=${module}
CMD java -jar ${SERVICE_NAME}.jar
