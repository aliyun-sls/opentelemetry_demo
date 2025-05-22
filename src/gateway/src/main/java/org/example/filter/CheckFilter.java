package org.example.filter;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.util.JsonFormat;
import io.grpc.*;
import io.grpc.stub.MetadataUtils;
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
import java.util.concurrent.Executor;
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
                .enableRetry() // 启用重试
                .maxRetryAttempts(3) // 最大重试次数
                .build();

        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr)
                .usePlaintext()
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry()
                .maxRetryAttempts(3) // 最大重试次数
                .build();
    }


    private Demo.OrderResult DoPlaceOrder(Demo.PlaceOrderRequest request, ServerWebExchange exchange) {
        final String uid = exchange.getRequest().getHeaders().getFirst("uid");
        final String sid = exchange.getRequest().getHeaders().getFirst("sid");

        log.info("DoPlaceOrder - Starting with uid: {}, sid: {}", uid, sid);

        Metadata headers = new Metadata();
        if (uid != null) {
            Metadata.Key<String> uidKey = Metadata.Key.of("uid", Metadata.ASCII_STRING_MARSHALLER);
            headers.put(uidKey, uid);
            log.info("Added uid to metadata with key: {} value: {}", uidKey.name(), uid);
        }
        if (sid != null) {
            Metadata.Key<String> sidKey = Metadata.Key.of("sid", Metadata.ASCII_STRING_MARSHALLER);
            headers.put(sidKey, sid);
            log.info("Added sid to metadata with key: {} value: {}", sidKey.name(), sid);
        }

        log.info("Created metadata with keys: {}", headers.keys());

        ClientInterceptor headerInterceptor = MetadataUtils.newAttachHeadersInterceptor(headers);

        log.info("Creating checkout service stub with headers");
        CheckoutServiceGrpc.CheckoutServiceBlockingStub checkoutServiceStub = 
            CheckoutServiceGrpc.newBlockingStub(checkOutChannel)
            .withDeadlineAfter(10, TimeUnit.SECONDS)
            .withInterceptors(headerInterceptor);

        try {
            log.info("Calling placeOrder with request: {}", request);
            Demo.PlaceOrderResponse response = checkoutServiceStub.placeOrder(request);
            log.info("Received placeOrder response");
            return response.getOrder();
        } catch (StatusRuntimeException e) {
            log.error("gRPC call failed with status: {}, description: {}", e.getStatus(), e.getMessage());
            throw e;
        }
    }

    private Demo.Product DoGetProductCatalog(Demo.GetProductRequest request) {
        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel);
        return productCatalogStub.getProduct(request);
    }


    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {

        log.info("Filtering request path {}, headers: {}", exchange.getRequest().getPath(), exchange.getRequest().getHeaders());
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
                    sink.success(DoPlaceOrder(requestBuilder.build(), exchange));
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
                        return Mono.error(e);
                    }
                });
            } catch (IOException e) {
                return Mono.error(e);
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = objectMapper.writeValueAsBytes(responseBody.toString());
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
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