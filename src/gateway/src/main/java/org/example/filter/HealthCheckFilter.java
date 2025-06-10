package org.example.filter;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

/**
 * 健康检查过滤器，用于响应Envoy的健康检查请求
 * Envoy会定期向/health路径发送GET请求，返回200表示服务健康
 */
public class HealthCheckFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(HealthCheckFilter.class);

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        // 记录请求详细信息
        log.info("Health check request received: Path={}, Method={}, Headers={}", 
                exchange.getRequest().getPath(), 
                exchange.getRequest().getMethod(),
                exchange.getRequest().getHeaders());
        
        // 检查请求方法是否为GET
        if (exchange.getRequest().getMethod() == HttpMethod.GET) {
            // 设置HTTP状态为200 OK
            exchange.getResponse().setStatusCode(HttpStatus.OK);
            exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
            
            // 返回简单的JSON响应
            String responseBody = "{\"status\":\"UP\"}";
            DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(responseBody.getBytes());
            
            // 记录响应内容
            log.info("Health check response: Status=200 OK, ContentType={}, Body={}", 
                    MediaType.APPLICATION_JSON, 
                    responseBody);
            
            return exchange.getResponse().writeWith(Mono.just(buffer));
        } else {
            // 对于非GET请求，返回405 Method Not Allowed
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            
            // 记录拒绝请求的日志
            log.warn("Health check rejected: non-GET method {} used", 
                    exchange.getRequest().getMethod());
            
            return exchange.getResponse().setComplete();
        }
    }
} 