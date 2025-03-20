package com.example.demo.gateway.config;

import org.springframework.cloud.gateway.route.RouteLocator;
import org.springframework.cloud.gateway.route.builder.RouteLocatorBuilder;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RouteConfig {

    @Bean
    public RouteLocator customRouteLocator(RouteLocatorBuilder builder) {
        return builder.routes()
                .route("ads", r -> r.path("/ads/**").uri("http://ads:8080"))
                .route("marketing", r -> r.path("/marketing/**").uri("http://marketing:8080"))
                .route("notification", r -> r.path("/notification/**").uri("http://notification:8080"))
                .route("promotion", r -> r.path("/promotion/**").uri("http://promotion:8080"))
                .route("pms", r -> r.path("/pms/**").uri("http://pms:8080"))
                .route("default_route", r -> r.path("/**").uri("http://frontend-proxy:8080"))
                .build();
    }
}