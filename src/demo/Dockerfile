FROM maven:3.6.3-jdk-11-openj9 AS demo_java

WORKDIR /app
COPY pom.xml .

COPY ads/pom.xml ads/
COPY ads/src ads/src/

COPY ads-error/pom.xml ads-error/
COPY ads-error/src ads-error/src/

COPY marketing/pom.xml marketing/
COPY marketing/src marketing/src/

COPY notification/pom.xml notification/
COPY notification/src notification/src/

COPY promotion/pom.xml promotion/
COPY promotion/src promotion/src/

COPY logistics/pom.xml logistics/
COPY logistics/src logistics/src/

RUN mvn package -DskipTests

FROM openjdk:11-jdk
ARG module
WORKDIR /app

COPY --from=demo_java /app/${module}/target/*.jar demo.jar
COPY --from=demo_java /app/${module}/src/main/resources/application.properties application.properties
EXPOSE 8080
CMD java --add-opens java.base/jdk.internal.misc=ALL-UNNAMED \
         -Dio.netty.tryReflectionSetAccessible=true -Xmx1536m -jar demo.jar
