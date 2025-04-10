package org.example.filter;

import org.example.config.Config;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.factory.AbstractGatewayFilterFactory;
import org.springframework.stereotype.Component;

@Component
public class RecommendationGatewayFilterFactory extends AbstractGatewayFilterFactory<Config> {

    private static final Logger log = LoggerFactory.getLogger(RecommendationGatewayFilterFactory.class);

    public RecommendationGatewayFilterFactory() {
        super(Config.class);
    }

    @Override
    public GatewayFilter apply(Config config) {
        log.info("RecommendationGatewayFilterFactory apply");
        return new RecommendationFilter(config);
    }
}
