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
        if (path.startsWith("/login") || path.startsWith("/register") || path.startsWith("/api/products")
                || path.startsWith("/_next") || path.startsWith("/images") || path.startsWith("/icons")
                || path.isEmpty() || path.equals("/")) {
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
                    // 获取用户权限
                    return getUserRole(userId)
                            .flatMap(role -> {
                                // 检查权限是否允许访问路径
                                if (!isAccessAllowed(path, role)) {
                                    ServerHttpResponse response = exchange.getResponse();
                                    response.setStatusCode(HttpStatus.FORBIDDEN);
                                    return response.setComplete();
                                }
                                // 将 userId 和 role 添加到 exchange 的属性中
                                exchange.getAttributes().put("userId", userId);
                                exchange.getAttributes().put("role", role);
                                return chain.filter(exchange);
                            });
                });
    }

    private Mono<String> isAuthenticated(ServerWebExchange exchange) {
        return exchange.getSession()
                .flatMap(session -> {
                    String sessionId = session.getId();
                    return Mono.justOrEmpty(redisTemplate.opsForValue().get(sessionId));
                });
    }

    // 获取用户角色
    private Mono<String> getUserRole(String userId) {
        return Mono.justOrEmpty(redisTemplate.opsForValue().get("role:" + userId));
    }

    // 检查路径是否允许访问
    private boolean isAccessAllowed(String path, String role) {
        // 示例：0 admin 可以访问所有路径，1 user 只能访问特定路径
        if ("0".equals(role)) {
            return true;
        } else if ("1".equals(role)) {
            return path.startsWith("/user/") || path.startsWith("/product/");
        }
        return false;
    }
}