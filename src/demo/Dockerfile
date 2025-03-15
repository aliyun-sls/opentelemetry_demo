FROM maven:3.6.3-jdk-11-openj9 AS demo_java

WORKDIR /app
COPY pom.xml .

COPY ads/pom.xml ads/
COPY ads/src ads/src/

RUN mvn package -DskipTests

FROM openjdk:11-jdk
WORKDIR /app

COPY --from=demo_java /app/ads/target/*.jar demo.jar
EXPOSE 8080
CMD java -jar demo.jar
