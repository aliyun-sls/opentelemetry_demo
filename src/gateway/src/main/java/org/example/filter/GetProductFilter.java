package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
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
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.server.PathContainer;
import org.springframework.util.MultiValueMap;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.CurrencyServiceGrpc;
import oteldemo.Demo;
import oteldemo.ProductCatalogServiceGrpc;
import reactor.core.publisher.Mono;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

import static org.springframework.cloud.gateway.support.ServerWebExchangeUtils.URI_TEMPLATE_VARIABLES_ATTRIBUTE;

public class GetProductFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(GetProductFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private ManagedChannel productCatalogChannel;
    private ManagedChannel currencyChannel;

    public GetProductFilter(Config config) {
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
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                    log.error("Retry interrupted: {}", ie.getMessage());
                    throw new RuntimeException(ie);
                } catch (Exception retryEx) {
                    log.error("Retry failed: {}", retryEx.getMessage());
                    throw retryEx;
                }
            }
            throw e;
        }
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = CurrencyServiceGrpc.newBlockingStub(currencyChannel)
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
                } catch (InterruptedException ie) {
                    Thread.currentThread().interrupt();
                    log.error("Retry interrupted: {}", ie.getMessage());
                    throw new RuntimeException(ie);
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
        // query: currencyCode productId
        if (!exchange.getRequest().getMethod().equals(HttpMethod.GET)) {
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            return exchange.getResponse().setComplete();
        }

        Map<String, String> uriVariables = exchange.getAttribute(URI_TEMPLATE_VARIABLES_ATTRIBUTE);
        String productId = uriVariables.get("productId");
        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String currencyCode = queryParams.getFirst("currencyCode");

        if (productId == null || productId.isEmpty()) {
            // status code 405
            return chain.filter(exchange);
        }


        if (currencyCode == null || currencyCode.isEmpty()) {
            currencyCode = "USD";
        }

        String finalCurrencyCode = currencyCode;

        return Mono.<JsonObject>create(sink -> {
            try {
                Demo.GetProductRequest.Builder builder = Demo.GetProductRequest.newBuilder();
                builder.setId(productId);

                Demo.Product product = DoGetProductCatalog(builder.build());
                JsonObject orderJson = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);


                Demo.CurrencyConversionRequest.Builder request = Demo.CurrencyConversionRequest.newBuilder();
                request.setFrom(product.getPriceUsd());
                request.setToCode(finalCurrencyCode);

                Demo.Money money = DoCurrencyConvert(request.build());

                JsonObject moneyJson = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                if (moneyJson.get("units") != null) {
                    moneyJson.addProperty("units", Integer.valueOf(moneyJson.get("units").getAsString()));
                }
                orderJson.add("priceUsd", moneyJson);
                sink.success(orderJson);

            } catch (Exception e) {
                log.error("Failed to get product", e);
                sink.error(e);
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.toString().getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed to get product", e);
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
