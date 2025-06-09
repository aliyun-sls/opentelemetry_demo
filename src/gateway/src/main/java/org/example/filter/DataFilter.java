package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.example.config.Config;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.MediaType;
import org.springframework.util.MultiValueMap;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.AdServiceGrpc;
import oteldemo.Demo;
import reactor.core.publisher.Mono;

import java.util.concurrent.TimeUnit;


public class DataFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(DataFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel channel;

    public DataFilter(Config config) {
        channel = ManagedChannelBuilder.forTarget(config.dataAddr)
                .usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .keepAliveWithoutCalls(true) // 即使没有活跃调用也发送keepalive
                .enableRetry() // 启用重试
                .disableServiceConfigLookUp() // 禁用服务配置查找，防止缓存
                .defaultLoadBalancingPolicy("round_robin") // 使用轮询策略
                .idleTimeout(5, TimeUnit.MINUTES) // 空闲5分钟后关闭连接
                .build();
    }

    private Demo.AdResponse DoGetAds(Demo.AdRequest request) {
        AdServiceGrpc.AdServiceBlockingStub adServiceStub = AdServiceGrpc.newBlockingStub(channel)
                .withDeadlineAfter(5, TimeUnit.SECONDS); // 添加5秒超时 - 广告查询
                
        try {
            return adServiceStub.getAds(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error getting ads, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return adServiceStub.withDeadlineAfter(10, TimeUnit.SECONDS).getAds(request);
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                    log.error("Retry interrupted: {}", ie.getMessage());
                    throw new RuntimeException(ie);
                } catch (Exception retryEx) {
                    log.error("Retry failed: {}", retryEx.getMessage());
                    throw new RuntimeException(retryEx);
                }
            }
            throw e;
        }
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());
        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String contextKeys = queryParams.getFirst("contextKeys");
        contextKeys=contextKeys==null?"":contextKeys;
        contextKeys = contextKeys.replaceAll("\\[", "");
        contextKeys = contextKeys.replaceAll("]", "");
        String[] keys = contextKeys.split(",");
        return Mono.<String>create(sink -> {
            try {
                Demo.AdRequest.Builder builder = Demo.AdRequest.newBuilder();
                for (String key : keys) {
                    if (!key.isEmpty()) {
                        builder.addContextKeys(key);
                    }
                }
                Demo.AdResponse ad = DoGetAds(builder.build());
                JsonArray jsonArray = new JsonArray();
                for (Demo.Ad adItem : ad.getAdsList()) {
                    JsonObject jsonObject = new JsonObject();
                    jsonObject.addProperty("redirectUrl", adItem.getRedirectUrl());
                    jsonObject.addProperty("text", adItem.getText());
                    jsonArray.add(jsonObject);
                }
                String json = jsonArray.toString();
                sink.success(json);
            } catch (Exception e) {
                log.error("Failed Data", e);
                sink.error(new RuntimeException(e));
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed Data", e);
                return Mono.error(new RuntimeException(e));
            }
        }).onErrorResume(ResponseStatusException.class, e -> {
            exchange.getResponse().setStatusCode(e.getStatusCode());
            exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
            String errorJson = String.format("{\"error\":\"%s\"}", e.getReason());
            DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(errorJson.getBytes());
            return exchange.getResponse().writeWith(Mono.just(buffer));
        }).then(chain.filter(exchange));
    }
}
