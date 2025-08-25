FROM ghcr.io/open-telemetry/opentelemetry-operator/autoinstrumentation-java:2.18.1 as java_agent

# Stage 1: Build the application using Maven
FROM maven:3.9-eclipse-temurin-17 AS build
WORKDIR /app
COPY pom.xml .
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Create the final image
FROM eclipse-temurin:17-jre-jammy

# The NetworkFault requires the 'tc' command, which is in the 'iproute2' package.
# The 'jammy' base image (Ubuntu 22.04) includes this by default, but we ensure it's there.
RUN apt-get update && apt-get install -y iproute2 && rm -rf /var/lib/apt/lists/*

WORKDIR /app
# Copy the executable JAR from the build stage
COPY --from=build /app/target/*.jar app.jar

EXPOSE 3500
RUN mkdir -p /opt/otel/
COPY --from=java_agent javaagent.jar /opt/otel/opentelemetry-javaagent.jar

# Set entrypoint to run the application
CMD java -javaagent:/opt/otel/opentelemetry-javaagent.jar -jar app.jar --server.port=3500