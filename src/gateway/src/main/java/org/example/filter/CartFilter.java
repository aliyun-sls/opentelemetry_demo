package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
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

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;


public class CartFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(CartFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel cartChannel;
    private final ManagedChannel productCatalogChannel;
    private final ManagedChannel currencyChannel;

    public CartFilter(Config config) {
        cartChannel = ManagedChannelBuilder.forTarget(config.cartAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry() // 启用重试
                .idleTimeout(5, TimeUnit.MINUTES) // 添加空闲超时
                .build();

        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry() // 启用重试
                .idleTimeout(5, TimeUnit.MINUTES) // 添加空闲超时
                .build();

        currencyChannel = ManagedChannelBuilder.forTarget(config.currencyAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry() // 启用重试
                .idleTimeout(5, TimeUnit.MINUTES) // 添加空闲超时
                .build();
    }


    private Demo.Cart DoGetCart(Demo.GetCartRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel);
        return cartServiceStub.getCart(request);
    }

    private Demo.Empty DoAddCartItem(Demo.AddItemRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel);
        return cartServiceStub.addItem(request);
    }

    private Demo.Empty DoEmptyCart(Demo.EmptyCartRequest request) {
        CartServiceGrpc.CartServiceBlockingStub cartServiceStub = CartServiceGrpc.newBlockingStub(cartChannel);
        return cartServiceStub.emptyCart(request);
    }

    private Demo.Product DoGetProductCatalog(Demo.GetProductRequest request) {
        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel);
        return productCatalogStub.getProduct(request);
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        oteldemo.CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = oteldemo.CurrencyServiceGrpc.newBlockingStub(currencyChannel);
        return currencyServiceBlockingStub.convert(request);
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());

        return Mono.<JsonObject>create(sink -> {
            try {
                if (exchange.getRequest().getMethod().equals(HttpMethod.GET)) {
                    handleGetRequest(exchange, sink);
                } else if (exchange.getRequest().getMethod().equals(HttpMethod.POST)) {
                    handlePostRequest(exchange, sink, chain);
                } else if (exchange.getRequest().getMethod().equals(HttpMethod.DELETE)) {
                    handleDeleteRequest(exchange, sink, chain);
                } else {
                    throw new ResponseStatusException(HttpStatus.METHOD_NOT_ALLOWED, "Unsupported HTTP method");
                }
            } catch (Exception e) {
                log.error("Failed Cart", e);
                sink.error(e);
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.toString().getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed Cart", e);
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

    private void handleGetRequest(ServerWebExchange exchange, MonoSink<JsonObject> sink) {
        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String sessionId = queryParams.getFirst("sessionId");
        String currencyCode = queryParams.getFirst("currencyCode");
        log.info("handleGetRequest sessionId: {} currencyCode: {}", sessionId, currencyCode);
        if (sessionId == null){
            return;
        }
        Demo.GetCartRequest.Builder builder = Demo.GetCartRequest.newBuilder();
        builder.setUserId(sessionId);
        Demo.Cart cart = DoGetCart(builder.build());
        String userId = cart.getUserId();
        log.info("handleGetRequest userId: {} cart: {}", userId, cart);
        List<JsonObject> productList = new ArrayList<>();
        cart.getItemsList().forEach(item -> {
            Demo.GetProductRequest.Builder getProductBuilder = Demo.GetProductRequest.newBuilder();
            getProductBuilder.setId(item.getProductId());
            Demo.Product product = DoGetProductCatalog(getProductBuilder.build());
            if (product != null) {
                JsonObject jsonObject = new JsonObject();
                jsonObject.addProperty("productId", item.getProductId());
                jsonObject.addProperty("quantity", item.getQuantity());
                JsonObject productObj;
                try {
                    productObj = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);
                    if (currencyCode!=null){
                        Demo.CurrencyConversionRequest.Builder currencyRequest = Demo.CurrencyConversionRequest.newBuilder();
                        currencyRequest.setFrom(product.getPriceUsd());
                        currencyRequest.setToCode(currencyCode);
                        Demo.Money money = DoCurrencyConvert(currencyRequest.build());
                        if (money != null && money.getUnits() != 0) {
                            JsonObject moneyObj = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                            moneyObj.addProperty("units", money.getUnits());
                            productObj.add("priceUsd", moneyObj);
                        }
                    }
                    jsonObject.add("product", productObj);
                    productList.add(jsonObject);
                } catch (InvalidProtocolBufferException e) {
                    throw new RuntimeException(e);
                }
            }
        });
        JsonObject resObject = new JsonObject();
        resObject.addProperty("userId", userId);
        resObject.add("items", new Gson().toJsonTree(productList));
        sink.success(resObject);
    }

    private void handlePostRequest(ServerWebExchange exchange, MonoSink<JsonObject> sink, GatewayFilterChain chain) {
        exchange.getRequest().getBody()
                .collectList()
                .flatMap(dataBuffers -> {
                    StringBuilder bodyBuilder = new StringBuilder();
                    dataBuffers.forEach(dataBuffer -> {
                        byte[] bytes = new byte[dataBuffer.readableByteCount()];
                        dataBuffer.read(bytes);
                        bodyBuilder.append(new String(bytes));
                    });
                    String body = bodyBuilder.toString();
                    try {
                        JsonObject jsonObject = new Gson().fromJson(body, JsonObject.class);
                        String userId = jsonObject.get("userId").getAsString();
                        JsonObject itemJson = jsonObject.getAsJsonObject("item");
                        Demo.CartItem.Builder cartItemBuilder = Demo.CartItem.newBuilder();
                        JsonFormat.parser().merge(itemJson.toString(), cartItemBuilder);

                        Demo.AddItemRequest.Builder builder = Demo.AddItemRequest.newBuilder();
                        builder.setUserId(userId);
                        builder.setItem(cartItemBuilder.build());
                        Demo.Empty addItemResponse = DoAddCartItem(builder.build());
                        log.info("AddItem response: {}", addItemResponse);

                        Demo.GetCartRequest.Builder getCartBuilder = Demo.GetCartRequest.newBuilder();
                        builder.setUserId(userId);
                        Demo.Cart cart = DoGetCart(getCartBuilder.build());
                        JsonObject cartObj = objectMapper.convertValue(cart, JsonObject.class);
                        sink.success(cartObj);
                        return chain.filter(exchange);
                    } catch (Exception e) {
                        log.error("Failed to parse request body", e);
                        sink.error(e);
                        return Mono.error(e);
                    }
                });
    }

    private void handleDeleteRequest(ServerWebExchange exchange, MonoSink<JsonObject> sink, GatewayFilterChain chain) {
        exchange.getRequest().getBody()
                .collectList()
                .flatMap(dataBuffers -> {
                    StringBuilder bodyBuilder = new StringBuilder();
                    dataBuffers.forEach(dataBuffer -> {
                        byte[] bytes = new byte[dataBuffer.readableByteCount()];
                        dataBuffer.read(bytes);
                        bodyBuilder.append(new String(bytes));
                    });
                    String body = bodyBuilder.toString();
                    try {
                        JsonObject jsonObject = new Gson().fromJson(body, JsonObject.class);
                        String userId = jsonObject.get("userId").getAsString();
                        Demo.EmptyCartRequest.Builder builder = Demo.EmptyCartRequest.newBuilder();
                        builder.setUserId(userId);
                        DoEmptyCart(builder.build());
                        Demo.Empty response = DoEmptyCart(builder.build());
                        log.info("EmptyCart response: {}", response);
                        exchange.getResponse().setStatusCode(HttpStatus.NO_CONTENT);

                        return chain.filter(exchange);
                    } catch (Exception e) {
                        log.error("Failed to parse request body", e);
                        sink.error(e);
                        return Mono.error(e);
                    }
                });
    }
}
