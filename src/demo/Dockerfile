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

COPY logistics/pom.xml logistics/
COPY logistics/src logistics/src/

COPY reporting/pom.xml reporting/
COPY reporting/src reporting/src/

COPY review/pom.xml review/
COPY review/src review/src/

RUN mvn package -DskipTests

FROM openjdk:11-jdk
ARG module
WORKDIR /app

COPY --from=demo_java /app/${module}/target/*.jar demo.jar
COPY --from=demo_java /app/${module}/src/main/resources/application.properties application.properties
EXPOSE 8080
CMD java -jar demo.jar
