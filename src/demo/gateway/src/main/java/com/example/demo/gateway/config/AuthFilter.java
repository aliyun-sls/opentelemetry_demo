package com.example.demo.gateway.config;

import org.springframework.http.HttpStatus;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.http.server.reactive.ServerHttpResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

//@Component
public class AuthFilter implements WebFilter {

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        ServerHttpRequest request = exchange.getRequest();
        String path = request.getPath().value();

        // 排除前端路由的认证检查
        if (path.startsWith("/login") || path.startsWith("/register")) {
            return chain.filter(exchange);
        }

        return isAuthenticated(exchange)
                .flatMap(isAuthenticated -> {
                    if (!isAuthenticated) {
                        ServerHttpResponse response = exchange.getResponse();
                        response.setStatusCode(HttpStatus.FOUND);
                        response.getHeaders().add("Location", "/login");
                        return response.setComplete();
                    }
                    return chain.filter(exchange);
                });
    }

    private Mono<Boolean> isAuthenticated(ServerWebExchange exchange) {
        return exchange.getSession()
                .map(session -> session.getAttribute("user") != null)
                .defaultIfEmpty(false);
    }
}