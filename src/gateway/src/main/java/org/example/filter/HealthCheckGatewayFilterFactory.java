package org.example.filter;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.factory.AbstractGatewayFilterFactory;
import org.springframework.stereotype.Component;

/**
 * 健康检查过滤器工厂类
 * 用于创建健康检查过滤器
 */
@Component
public class HealthCheckGatewayFilterFactory extends AbstractGatewayFilterFactory<Object> {

    private static final Logger log = LoggerFactory.getLogger(HealthCheckGatewayFilterFactory.class);

    public HealthCheckGatewayFilterFactory() {
        super(Object.class);
    }

    @Override
    public GatewayFilter apply(Object config) {
        log.info("HealthCheckGatewayFilterFactory apply");
        return new HealthCheckFilter();
    }
}