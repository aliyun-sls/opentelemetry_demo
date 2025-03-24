package com.example.demo.gateway.config;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.server.reactive.ServerHttpRequest;
import org.springframework.http.server.reactive.ServerHttpResponse;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;
import reactor.core.publisher.Mono;

@Component
public class AuthFilter implements WebFilter {

    @Autowired
    private StringRedisTemplate redisTemplate;

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        ServerHttpRequest request = exchange.getRequest();
        String path = request.getPath().value();

        // 排除前端路由的认证检查
        if (path.startsWith("/login") || path.startsWith("/register")) {
            return chain.filter(exchange);
        }

        return isAuthenticated(exchange)
                .flatMap(userId -> {
                    if (userId == null || userId.isEmpty()) {
                        ServerHttpResponse response = exchange.getResponse();
                        response.setStatusCode(HttpStatus.FOUND);
                        response.getHeaders().add("Location", "/login");
                        return response.setComplete();
                    }
                    // 将 userId 添加到 exchange 的属性中
                    exchange.getAttributes().put("userId", userId);
                    return chain.filter(exchange);
                });
    }

    private Mono<String> isAuthenticated(ServerWebExchange exchange) {
        return exchange.getSession()
                .flatMap(session -> {
                    String sessionId = session.getId();
                    return Mono.justOrEmpty(redisTemplate.opsForValue().get(sessionId));
                });
    }
}