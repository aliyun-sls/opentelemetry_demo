package org.example.filter;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

public class NoopFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(NoopFilter.class);

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("NoopFilter: {}", exchange.getRequest().getURI());
        return chain.filter(exchange);
    }
}
