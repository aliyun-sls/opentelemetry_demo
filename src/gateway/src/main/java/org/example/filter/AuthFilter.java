package org.example.filter;

import com.networknt.schema.utils.StringUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.annotation.Order;
import org.springframework.data.redis.core.StringRedisTemplate;
import org.springframework.cloud.gateway.filter.GlobalFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.stereotype.Component;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;

import java.util.Objects;
import java.util.concurrent.ConcurrentHashMap;


@Order(-200)
@Component
public class AuthFilter implements GlobalFilter {

    private static final Logger log = LoggerFactory.getLogger(AuthFilter.class);


    @Autowired
    private StringRedisTemplate redisTemplate;

    private final String uidKey = "uid";

    private final String sidKey = "sid";
    //cache for user info.
    private final ConcurrentHashMap<String, Integer> userCaches = new ConcurrentHashMap<>();
    
    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        String sid = exchange.getRequest().getHeaders().getFirst(sidKey);
        if (StringUtils.isNotBlank(sid)) {
            exchange.getAttributes().put(sidKey, sid);
            Integer uid = getUidBySid(sid);
            if (uid != null) {
                userCaches.put(sid, uid);
                exchange.getAttributes().put(uidKey, uid);
                exchange = exchange.mutate()
                    .request(r -> r.header(sidKey, sid).header(uidKey, uid.toString()))
                    .build();
            }
        }
        return chain.filter(exchange);
    }

    private Integer getUidBySid(String sid) {
        try {
            final Integer uid = userCaches.get(sid);
            if (Objects.nonNull(uid)) {
                return uid;
            }
            final String uidStr = redisTemplate.opsForValue().get(sid);
            if (StringUtils.isBlank(uidStr)) {
                return null;
            }
            return Integer.parseInt(uidStr);
        } catch (RuntimeException e) {
            log.error(e.getMessage(), e);
            return null;
        }
    }
}