package org.example.filter;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.example.config.Config;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.core.io.buffer.DataBufferUtils;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.CheckoutServiceGrpc;
import oteldemo.Demo;
import oteldemo.ProductCatalogServiceGrpc;
import reactor.core.publisher.Mono;

import java.io.IOException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class CheckFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(CheckFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private ManagedChannel checkOutChannel;
    private ManagedChannel productCatalogChannel;

    public CheckFilter(Config config) {
        checkOutChannel = ManagedChannelBuilder.forTarget(config.checkoutAddr)
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

        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr).usePlaintext() // 明文通信（仅限开发环境）
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


    private Demo.PlaceOrderResponse DoPlaceOrder(Demo.PlaceOrderRequest request) {
        CheckoutServiceGrpc.CheckoutServiceBlockingStub checkoutServiceBlockingStub = CheckoutServiceGrpc.newBlockingStub(checkOutChannel)
                .withDeadlineAfter(10, TimeUnit.SECONDS); // 添加10秒超时 - 结账复杂操作
                
        try {
            return checkoutServiceBlockingStub.placeOrder(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error placing order, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return checkoutServiceBlockingStub.withDeadlineAfter(15, TimeUnit.SECONDS).placeOrder(request);
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

    private Demo.Product DoGetProductCatalog(Demo.GetProductRequest request) {
//        String json = "{ \"id\": \"OLJCESPC7Z\", \"name\": \"National Park Foundation Explorascope\", \"description\": \"The National Park Foundation's (NPF) Explorascope 60AZ is a manual alt-azimuth, refractor telescope perfect for celestial viewing on the go. The NPF Explorascope 60 can view the planets, moon, star clusters and brighter deep sky objects like the Orion Nebula and Andromeda Galaxy.\", \"picture\": \"NationalParkFoundationExplorascope.jpg\", \"priceUsd\": { \"currencyCode\": \"USD\", \"units\": 101, \"nanos\": 960000000 }, \"categories\": [ \"telescopes\" ] }";
//
//        try {
//            Demo.Product.Builder builder = Demo.Product.newBuilder();
//            JsonFormat.parser().ignoringUnknownFields().merge(json, builder);
//            return builder.build();
//        } catch (InvalidProtocolBufferException e) {
//            throw new RuntimeException(e);
//        }

        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel)
                .withDeadlineAfter(5, TimeUnit.SECONDS); // 添加5秒超时 - 产品查询相对简单
        return productCatalogStub.getProduct(request);
    }


    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {

        log.info("Filtering request {}", exchange.getRequest().getPath());
        // 验证1：仅允许POST请求
        if (!exchange.getRequest().getMethod().equals(HttpMethod.POST)) {
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            return exchange.getResponse().setComplete();
        }

        // 验证2：检查请求体是否存在
        if (exchange.getRequest().getHeaders().getContentLength() <= 0) {
            return Mono.error(new ResponseStatusException(HttpStatus.BAD_REQUEST, "Request body is required"));
        }

        return DataBufferUtils.join(exchange.getRequest().getBody()).switchIfEmpty(Mono.error(new ResponseStatusException(HttpStatus.BAD_REQUEST, "Request body cannot be empty"))).flatMap(dataBuffer -> {
            try {
                JsonNode json = objectMapper.readTree(dataBuffer.asInputStream());
                Demo.PlaceOrderRequest.Builder requestBuilder = Demo.PlaceOrderRequest.newBuilder();
                JsonFormat.parser().ignoringUnknownFields().merge(json.toString(), requestBuilder);


                return Mono.<Demo.OrderResult>create(sink -> {
                    sink.success(DoPlaceOrder(requestBuilder.build()).getOrder());
                }).flatMap(order -> {
                    try {
                        List<JsonObject> orderItems = new ArrayList<>();

                        for (Demo.OrderItem item : order.getItemsList()) {
                            JsonObject baseItem = new JsonObject();

                            baseItem.addProperty("productId", item.getItem().getProductId());
                            baseItem.addProperty("quantity", item.getItem().getQuantity());

                            // 转换货币信息
                            JsonObject costJson = new Gson().fromJson(JsonFormat.printer().print(item.getCost()), JsonObject.class);
                            baseItem.add("cost", costJson);

                            // 异步查询商品详情
                            Demo.GetProductRequest request = Demo.GetProductRequest.newBuilder().setId(item.getItem().getProductId()).build();

                            Demo.Product product = DoGetProductCatalog(request);
                            JsonObject productJson = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);
                            baseItem.add("product", productJson);

                            orderItems.add(baseItem);
                        }

                        JsonObject orderJson = new Gson().fromJson(JsonFormat.printer().print(order), JsonObject.class);
                        orderJson.add("items", new JsonParser().parse(new Gson().toJson(orderItems)));
                        return Mono.just(orderJson);
                    } catch (InvalidProtocolBufferException e) {
                        return Mono.error(new RuntimeException(e));
                    }
                });
            } catch (IOException e) {
                return Mono.error(new RuntimeException(e));
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = objectMapper.writeValueAsBytes(responseBody.toString());
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed check", e);
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
