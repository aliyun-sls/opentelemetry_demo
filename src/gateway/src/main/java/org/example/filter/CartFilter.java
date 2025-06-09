package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
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
import org.springframework.core.io.buffer.DataBufferUtils;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.util.MultiValueMap;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.ProductCatalogServiceGrpc;
import oteldemo.Demo;
import oteldemo.CartServiceGrpc;
import reactor.core.publisher.Mono;
import reactor.core.publisher.MonoSink;

import java.nio.charset.StandardCharsets;
import java.util.concurrent.TimeUnit;

public class CartFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(CartFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel cartChannel;
    private final ManagedChannel productCatalogChannel;
    private final ManagedChannel currencyChannel;

    public CartFilter(Config config) {
        cartChannel = ManagedChannelBuilder.forTarget(config.cartAddr)
                .usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .keepAliveWithoutCalls(true) // 即使没有活跃调用也发送keepalive
                .enableRetry() // 启用重试
                .defaultLoadBalancingPolicy("round_robin") // 使用轮询策略
                .build();

        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr)
                .usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .keepAliveWithoutCalls(true) // 即使没有活跃调用也发送keepalive
                .enableRetry() // 启用重试
                .defaultLoadBalancingPolicy("round_robin") // 使用轮询策略
                .build();

        currencyChannel = ManagedChannelBuilder.forTarget(config.currencyAddr)
                .usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .keepAliveWithoutCalls(true) // 即使没有活跃调用也发送keepalive
                .enableRetry() // 启用重试
                .defaultLoadBalancingPolicy("round_robin") // 使用轮询策略
                .build();
    }


    private Demo.Cart DoGetCart(Demo.GetCartRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel)
                .withDeadlineAfter(8, TimeUnit.SECONDS); // 添加8秒超时 - cart操作中等复杂度
                
        try {
            return cartServiceStub.getCart(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error getting cart, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return cartServiceStub.withDeadlineAfter(12, TimeUnit.SECONDS).getCart(request);
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

    private Demo.Empty DoAddCartItem(Demo.AddItemRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel)
                .withDeadlineAfter(8, TimeUnit.SECONDS); // 添加8秒超时
                
        try {
            return cartServiceStub.addItem(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error adding cart item, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return cartServiceStub.withDeadlineAfter(12, TimeUnit.SECONDS).addItem(request);
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

    private Demo.Empty DoEmptyCart(Demo.EmptyCartRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel)
                .withDeadlineAfter(8, TimeUnit.SECONDS); // 添加8秒超时
                
        try {
            return cartServiceStub.emptyCart(request);
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error emptying cart, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    return cartServiceStub.withDeadlineAfter(12, TimeUnit.SECONDS).emptyCart(request);
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
                    throw new RuntimeException(retryEx);
                }
            }
            throw e;
        }
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        oteldemo.CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = oteldemo.CurrencyServiceGrpc.newBlockingStub(currencyChannel)
                .withDeadlineAfter(3, TimeUnit.SECONDS); // 添加3秒超时 - 货币转换简单快速
                
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
                    throw new RuntimeException(retryEx);
                }
            }
            throw e;
        }
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());
        if(exchange.getRequest().getMethod() == HttpMethod.GET){
            return Mono.<JsonObject>create(sink -> handleGetRequest(exchange, sink))
                    .flatMap(responseBody -> {
                        try {
                            exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                            byte[] bytes = responseBody.toString().getBytes();
                            DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                            return exchange.getResponse().writeWith(Mono.just(buffer));
                        } catch (Exception e) {
                            log.error("Failed Cart", e);
                            return Mono.error(new RuntimeException(e));
                        }
                    }).onErrorResume(ResponseStatusException.class, e -> {
                        exchange.getResponse().setStatusCode(e.getStatusCode());
                        exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                        String errorJson = String.format("{\"error\":\"%s\"}", e.getReason());
                        DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(errorJson.getBytes());
                        return exchange.getResponse().writeWith(Mono.just(buffer));
                    }).then(chain.filter(exchange));
        }else if(exchange.getRequest().getMethod() == HttpMethod.POST){
            return handlePostRequest(exchange, chain);
        }else if(exchange.getRequest().getMethod() == HttpMethod.DELETE){
            return handleDeleteRequest(exchange, chain);
        }else{
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            return exchange.getResponse().setComplete().then(chain.filter(exchange));
        }
    }

    private void handleGetRequest(ServerWebExchange exchange, MonoSink<JsonObject> sink) {
        try{
            MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
            String sessionId = queryParams.getFirst("sessionId");
            String currencyCode = queryParams.getFirst("currencyCode");
            log.info("handleGetRequest sessionId: {} currencyCode: {}", sessionId, currencyCode);
            if(sessionId == null || sessionId.isEmpty()){
                JsonObject jsonObject = new JsonObject();
                jsonObject.add("items", new JsonArray());
                sink.success(jsonObject);
                return;
            }
            Demo.GetCartRequest.Builder builder = Demo.GetCartRequest.newBuilder();
            builder.setUserId(sessionId);
            Demo.Cart cart = DoGetCart(builder.build());
            JsonObject cartJson = new Gson().fromJson(JsonFormat.printer().print(cart), JsonObject.class);
            
            // 检查items字段是否存在，如果不存在则添加一个空数组
            if (!cartJson.has("items") || cartJson.get("items").isJsonNull()) {
                cartJson.add("items", new JsonArray());
                log.info("Returning empty cart for sessionId: {}, currencyCode: {}, cartJson: {}", sessionId, currencyCode, cartJson);
                sink.success(cartJson);
                return; // 如果没有items，直接返回空购物车
            }
            
            JsonArray itemsArray = cartJson.getAsJsonArray("items");
            if (itemsArray == null || itemsArray.size() == 0) {
                sink.success(cartJson);
                return; // 如果items是空数组，直接返回
            }
            
            for (int i = 0; i < itemsArray.size(); i++) {
                JsonElement itemElement = itemsArray.get(i);
                if (itemElement == null || !itemElement.isJsonObject()) {
                    continue; // 跳过非对象元素
                }
                
                JsonObject itemObj = itemElement.getAsJsonObject();
                if (!itemObj.has("product") || itemObj.get("product").isJsonNull() || 
                    !itemObj.get("product").isJsonObject()) {
                    continue; // 跳过没有product字段或product不是对象的元素
                }
                
                JsonObject productObj = itemObj.getAsJsonObject("product");
                if (!productObj.has("id") || productObj.get("id").isJsonNull()) {
                    continue; // 跳过没有id字段的product
                }
                
                String productId = productObj.get("id").getAsString();
                if (productId == null || productId.isEmpty()) {
                    continue; // 跳过id为空的product
                }
                
                try {
                    Demo.GetProductRequest.Builder productRequestBuilder = Demo.GetProductRequest.newBuilder();
                    productRequestBuilder.setId(productId);
                    Demo.Product product = DoGetProductCatalog(productRequestBuilder.build());
                    JsonObject productJson = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);
                    itemObj.add("product", productJson);
                    
                    if (currencyCode != null && !currencyCode.isEmpty() && 
                        product.hasPriceUsd()) { // 确保product有priceUsd字段
                        try {
                            Demo.CurrencyConversionRequest.Builder request = Demo.CurrencyConversionRequest.newBuilder();
                            request.setFrom(product.getPriceUsd());
                            request.setToCode(currencyCode);
                            Demo.Money money = DoCurrencyConvert(request.build());
                            JsonObject moneyJson = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                            productJson.add("priceUsd", moneyJson);
                        } catch (Exception e) {
                            log.error("Failed to convert currency for product {}: {}", productId, e.getMessage());
                            // 货币转换失败不应该影响整个购物车的获取
                        }
                    }
                } catch (Exception e) {
                    log.error("Failed to get product details for {}: {}", productId, e.getMessage());
                    // 获取产品详情失败不应该影响整个购物车的获取
                }
            }
            sink.success(cartJson);
        } catch (Exception e) {
            log.error("Failed to get cart", e);
            sink.error(new RuntimeException(e));
        }
    }

    private Mono<Void> handlePostRequest(ServerWebExchange exchange, GatewayFilterChain chain) {
        return DataBufferUtils.join(exchange.getRequest().getBody()).flatMap(dataBuffer -> {
            try {
                byte[] bytes = new byte[dataBuffer.readableByteCount()];
                dataBuffer.read(bytes);
                DataBufferUtils.release(dataBuffer);
                String request = new String(bytes, StandardCharsets.UTF_8);
                JsonObject requestJson = new Gson().fromJson(request, JsonObject.class);
                String userId = requestJson.get("userId").getAsString();
                String productId = requestJson.get("productId").getAsString();
                int quantity = requestJson.get("quantity").getAsInt();
                log.info("handlePostRequest userId: {} productId: {} quantity: {}", userId, productId, quantity);
                Demo.AddItemRequest.Builder builder = Demo.AddItemRequest.newBuilder();
                builder.setUserId(userId);
                builder.setItem(Demo.CartItem.newBuilder().setProductId(productId).setQuantity(quantity).build());
                DoAddCartItem(builder.build());
                exchange.getResponse().setStatusCode(HttpStatus.OK);
                return exchange.getResponse().setComplete();
            } catch (Exception e) {
                log.error("Failed to add cart item", e);
                return Mono.error(new RuntimeException(e));
            }
        }).then(chain.filter(exchange));
    }

    private Mono<Void> handleDeleteRequest(ServerWebExchange exchange, GatewayFilterChain chain) {
        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String sessionId = queryParams.getFirst("sessionId");
        log.info("handleDeleteRequest sessionId: {}", sessionId);
        try {
            if (sessionId != null) {
                Demo.EmptyCartRequest.Builder builder = Demo.EmptyCartRequest.newBuilder();
                builder.setUserId(sessionId);
                DoEmptyCart(builder.build());
            }
            exchange.getResponse().setStatusCode(HttpStatus.OK);
            return exchange.getResponse().setComplete();
        } catch (Exception e) {
            log.error("Failed to empty cart", e);
            return Mono.error(new RuntimeException(e));
        }
    }
}
