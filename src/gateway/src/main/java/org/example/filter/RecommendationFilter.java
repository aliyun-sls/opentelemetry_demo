package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.google.protobuf.util.JsonFormat;
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
import oteldemo.RecommendationServiceGrpc;
import oteldemo.Demo;
import oteldemo.ProductCatalogServiceGrpc;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;


public class RecommendationFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(RecommendationFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel recommendationChannel;
    private final ManagedChannel productCatalogChannel;
    private final ManagedChannel currencyChannel;

    public RecommendationFilter(Config config) {
        recommendationChannel = ManagedChannelBuilder.forTarget(config.recommendationAddr)
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

        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr)
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

        currencyChannel = ManagedChannelBuilder.forTarget(config.currencyAddr)
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


    private Demo.ListRecommendationsResponse DoListRecommendations(Demo.ListRecommendationsRequest request) {
        RecommendationServiceGrpc.RecommendationServiceBlockingStub recommendationServiceBlockingStub = RecommendationServiceGrpc.newBlockingStub(recommendationChannel)
                .withDeadlineAfter(10, TimeUnit.SECONDS); // 添加10秒超时 - 推荐算法可能较慢
                
        try {
            return recommendationServiceBlockingStub.listRecommendations(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error getting recommendations, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return recommendationServiceBlockingStub.withDeadlineAfter(15, TimeUnit.SECONDS).listRecommendations(request);
                } catch (Exception retryEx) {
                    log.error("Retry failed: {}", retryEx.getMessage());
                    throw retryEx;
                }
            }
            throw e;
        }
    }

    private Demo.Product DoGetProductCatalog(Demo.GetProductRequest request) {
        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel)
                .withDeadlineAfter(5, TimeUnit.SECONDS); // 添加5秒超时 - 产品查询
                
        try {
            return productCatalogStub.getProduct(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error getting product, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return productCatalogStub.withDeadlineAfter(10, TimeUnit.SECONDS).getProduct(request);
                } catch (Exception retryEx) {
                    log.error("Retry failed: {}", retryEx.getMessage());
                    throw retryEx;
                }
            }
            throw e;
        }
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        oteldemo.CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = oteldemo.CurrencyServiceGrpc.newBlockingStub(currencyChannel)
                .withDeadlineAfter(3, TimeUnit.SECONDS); // 添加3秒超时 - 货币转换
                
        try {
            return currencyServiceBlockingStub.convert(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error converting currency, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return currencyServiceBlockingStub.withDeadlineAfter(6, TimeUnit.SECONDS).convert(request);
                } catch (Exception retryEx) {
                    log.error("Retry failed: {}", retryEx.getMessage());
                    throw retryEx;
                }
            }
            throw e;
        }
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());
        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String productIds = queryParams.getFirst("productIds");
        String sessionId = queryParams.getFirst("sessionId");
        String currencyCode = queryParams.getFirst("currencyCode");
        log.info("handleGetRequest sessionId: {} currencyCode: {} productIds:{}", sessionId, currencyCode, productIds);
        return Mono.<JsonArray>create(sink -> {
            try {
                Demo.ListRecommendationsRequest.Builder listRecommendationsBuilder = Demo.ListRecommendationsRequest.newBuilder();
                if(sessionId != null){
                    listRecommendationsBuilder.setUserId(sessionId);
                }
                if (productIds != null) {
                    String[] ids = productIds.split(",");
                    for (String id : ids) {
                        if (!id.isEmpty()) {
                            listRecommendationsBuilder.addProductIds(id);
                        }
                    }
                }
                Demo.ListRecommendationsResponse recommendationsResponse = DoListRecommendations(listRecommendationsBuilder.build());
                List<String> productIdsList = recommendationsResponse.getProductIdsList();
                List<JsonObject> productList = new ArrayList<>();
                for (String productId : productIdsList) {
                    Demo.GetProductRequest.Builder getProductBuilder = Demo.GetProductRequest.newBuilder();
                    getProductBuilder.setId(productId);
                    Demo.Product product = DoGetProductCatalog(getProductBuilder.build());
                    if (product != null) {
                        JsonObject productObj = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);
                        if (currencyCode!=null){
                            Demo.CurrencyConversionRequest.Builder currencyRequest = Demo.CurrencyConversionRequest.newBuilder();
                            currencyRequest.setFrom(product.getPriceUsd());
                            currencyRequest.setToCode(currencyCode);
                            Demo.Money money = DoCurrencyConvert(currencyRequest.build());
                            if (money != null && money.getUnits() != 0) {
                                JsonObject moneyObj = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                                moneyObj.addProperty("units", money.getUnits());
                                moneyObj.addProperty("nanos", money.getNanos());
                                productObj.add("priceUsd", moneyObj);
                            }
                        }
                        productList.add(productObj);
                    }
                }
                sink.success(new Gson().toJsonTree(productList).getAsJsonArray());
            } catch (Exception e) {
                log.error("Failed Recommendation", e);
                sink.error(e);
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.toString().getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed Recommendation", e);
                return Mono.error(e);
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
