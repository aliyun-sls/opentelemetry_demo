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
                    
                    // 记录原始产品数据
                    log.info("Product from catalog service - ID: {}, Raw proto: {}", productId, product);
                    
                    JsonObject productJson = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);
                    
                    // 记录转换后的JSON结构
                    log.info("Product JSON after conversion - ID: {}, JSON: {}", productId, productJson);
                    
                    // 处理priceUsd，确保格式正确
                    if (productJson.has("priceUsd")) {
                        JsonObject priceUsd = productJson.getAsJsonObject("priceUsd");
                        
                        // 确保units字段是数字
                        if (priceUsd.has("units")) {
                            try {
                                if (priceUsd.get("units").isJsonPrimitive() && priceUsd.get("units").getAsJsonPrimitive().isString()) {
                                    int units = Integer.parseInt(priceUsd.get("units").getAsString());
                                    priceUsd.addProperty("units", units);
                                    log.info("Converted priceUsd.units from string to number: {}", units);
                                }
                            } catch (Exception e) {
                                log.warn("Error processing priceUsd.units: {}", e.getMessage());
                                priceUsd.addProperty("units", 0); // 默认值
                            }
                        } else {
                            priceUsd.addProperty("units", 0);
                            log.info("Added missing priceUsd.units with default value 0");
                        }
                        
                        // 确保nanos字段是数字
                        if (priceUsd.has("nanos")) {
                            try {
                                if (priceUsd.get("nanos").isJsonPrimitive() && priceUsd.get("nanos").getAsJsonPrimitive().isString()) {
                                    int nanos = Integer.parseInt(priceUsd.get("nanos").getAsString());
                                    priceUsd.addProperty("nanos", nanos);
                                    log.info("Converted priceUsd.nanos from string to number: {}", nanos);
                                }
                            } catch (Exception e) {
                                log.warn("Error processing priceUsd.nanos: {}", e.getMessage());
                                priceUsd.addProperty("nanos", 0); // 默认值
                            }
                        } else {
                            priceUsd.addProperty("nanos", 0);
                            log.info("Added missing priceUsd.nanos with default value 0");
                        }
                        
                        // 确保currencyCode字段存在
                        if (!priceUsd.has("currencyCode") || priceUsd.get("currencyCode").isJsonNull()) {
                            priceUsd.addProperty("currencyCode", "USD");
                            log.info("Added missing priceUsd.currencyCode with default value USD");
                        }
                        
                        // 更新priceUsd
                        productJson.add("priceUsd", priceUsd);
                        
                        // 同时添加price字段，复制priceUsd的内容
                        productJson.add("price", priceUsd.deepCopy());
                        log.info("Added price field (copy of priceUsd) to product JSON");
                    } else {
                        // 如果没有priceUsd字段，创建默认的price和priceUsd
                        JsonObject defaultPrice = new JsonObject();
                        defaultPrice.addProperty("units", 0);
                        defaultPrice.addProperty("nanos", 0);
                        defaultPrice.addProperty("currencyCode", "USD");
                        
                        productJson.add("priceUsd", defaultPrice);
                        productJson.add("price", defaultPrice.deepCopy());
                        log.info("Created default price and priceUsd fields for product {}", productId);
                    }
                    
                    itemObj.add("product", productJson);
                    
                    if (currencyCode != null && !currencyCode.isEmpty() && 
                        product.hasPriceUsd()) { // 确保product有priceUsd字段
                        try {
                            Demo.CurrencyConversionRequest.Builder request = Demo.CurrencyConversionRequest.newBuilder();
                            request.setFrom(product.getPriceUsd());
                            request.setToCode(currencyCode);
                            
                            // 记录货币转换请求
                            log.info("Currency conversion request - From: {}, To: {}", 
                                    JsonFormat.printer().print(product.getPriceUsd()), currencyCode);
                            
                            Demo.Money money = DoCurrencyConvert(request.build());
                            
                            // 记录货币转换结果
                            log.info("Currency conversion result - Money: {}", JsonFormat.printer().print(money));
                            
                            JsonObject moneyJson = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                            
                            // 确保price字段的数据类型正确
                            if (moneyJson.has("units") && moneyJson.get("units").isJsonPrimitive()) {
                                // 检查units是否是字符串，如果是则转换为数字
                                if (moneyJson.get("units").getAsJsonPrimitive().isString()) {
                                    try {
                                        int units = Integer.parseInt(moneyJson.get("units").getAsString());
                                        moneyJson.addProperty("units", units);
                                        log.info("Converted units from string to number: {}", units);
                                    } catch (NumberFormatException e) {
                                        log.warn("Failed to convert units to number: {}", e.getMessage());
                                        moneyJson.addProperty("units", 0);
                                    }
                                }
                            } else {
                                moneyJson.addProperty("units", 0);
                                log.warn("Added missing units field with default value 0");
                            }
                            
                            if (moneyJson.has("nanos") && moneyJson.get("nanos").isJsonPrimitive()) {
                                // 检查nanos是否是字符串，如果是则转换为数字
                                if (moneyJson.get("nanos").getAsJsonPrimitive().isString()) {
                                    try {
                                        int nanos = Integer.parseInt(moneyJson.get("nanos").getAsString());
                                        moneyJson.addProperty("nanos", nanos);
                                        log.info("Converted nanos from string to number: {}", nanos);
                                    } catch (NumberFormatException e) {
                                        log.warn("Failed to convert nanos to number: {}", e.getMessage());
                                        moneyJson.addProperty("nanos", 0);
                                    }
                                }
                            } else {
                                moneyJson.addProperty("nanos", 0);
                                log.warn("Added missing nanos field with default value 0");
                            }
                            
                            // 确保有currencyCode字段
                            if (!moneyJson.has("currencyCode") || moneyJson.get("currencyCode").isJsonNull()) {
                                moneyJson.addProperty("currencyCode", currencyCode);
                                log.info("Added missing currencyCode: {}", currencyCode);
                            }
                            
                            // 记录最终的价格对象
                            log.info("Final money JSON: {}", moneyJson);
                            
                            productJson.add("priceUsd", moneyJson);
                            
                            // 同时添加price字段，以适配前端期望
                            productJson.add("price", moneyJson.deepCopy());
                            log.info("Added both priceUsd and price fields to product JSON");
                        } catch (Exception e) {
                            log.error("Failed to convert currency for product {}: {}", productId, e.getMessage());
                            // 货币转换失败不应该影响整个购物车的获取
                        }
                    } else {
                        log.warn("Currency conversion skipped - currencyCode: {}, product has priceUsd: {}", 
                                currencyCode, product.hasPriceUsd());
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
                
                // 打印原始请求内容
                log.info("Cart POST raw request: {}", request);
                
                JsonObject requestJson = new Gson().fromJson(request, JsonObject.class);
                
                // 打印解析后的JSON对象
                log.info("Cart POST parsed JSON: {}", requestJson);
                
                // 添加空值检查，避免NullPointerException
                if (!requestJson.has("userId") || requestJson.get("userId").isJsonNull()) {
                    log.error("Missing or null userId in request: {}", requestJson);
                    exchange.getResponse().setStatusCode(HttpStatus.OK);
                    return exchange.getResponse().setComplete();
                }
                
                // 检查item对象及其内部的字段
                if (!requestJson.has("item") || requestJson.get("item").isJsonNull() || !requestJson.get("item").isJsonObject()) {
                    log.error("Missing or null item object in request: {}", requestJson);
                    exchange.getResponse().setStatusCode(HttpStatus.OK);
                    return exchange.getResponse().setComplete();
                }
                
                JsonObject itemObject = requestJson.getAsJsonObject("item");
                
                // 检查item对象中的productId和quantity
                if (!itemObject.has("productId") || itemObject.get("productId").isJsonNull()) {
                    log.error("Missing or null productId in item object: {}", itemObject);
                    exchange.getResponse().setStatusCode(HttpStatus.OK);
                    return exchange.getResponse().setComplete();
                }
                
                if (!itemObject.has("quantity") || itemObject.get("quantity").isJsonNull()) {
                    log.error("Missing or null quantity in item object: {}", itemObject);
                    exchange.getResponse().setStatusCode(HttpStatus.OK);
                    return exchange.getResponse().setComplete();
                }
                
                String userId = requestJson.get("userId").getAsString();
                String productId = itemObject.get("productId").getAsString();
                int quantity = itemObject.get("quantity").getAsInt();
                
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
