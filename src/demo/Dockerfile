FROM maven:3.6.3-jdk-11-openj9 AS demo_java

WORKDIR /app
COPY pom.xml .

COPY ads/pom.xml ads/
COPY ads/src ads/src/

COPY marketing/pom.xml marketing/
COPY marketing/src marketing/src/

COPY notification/pom.xml notification/
COPY notification/src notification/src/

COPY promotion/pom.xml promotion/
COPY promotion/src promotion/src/

RUN mvn package -DskipTests
ARG module
FROM openjdk:11-jdk
WORKDIR /app

COPY --from=demo_java /app/${module}/target/*.jar demo.jar
EXPOSE 8080
CMD java -jar demo.jar
