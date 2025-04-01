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
import org.springframework.util.MultiValueMap;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.CurrencyServiceGrpc;
import oteldemo.Demo;
import oteldemo.ProductCatalogServiceGrpc;
import reactor.core.publisher.Mono;

import java.util.concurrent.TimeUnit;

public class GetProductFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(GetProductFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private ManagedChannel productCatalogChannel;
    private ManagedChannel currencyChannel;

    public GetProductFilter(Config config) {
        productCatalogChannel = ManagedChannelBuilder.forTarget(config.productAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry() // 启用重试
                .build();


        currencyChannel = ManagedChannelBuilder.forTarget(config.currencyAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .enableRetry() // 启用重试
                .build();
    }


    private Demo.Product DoGetProductCatalog(Demo.GetProductRequest request) {
//        String json = "{ \"id\": \"OLJCESPC7Z\", \"name\": \"National Park Foundation Explorascope\", \"description\": \"The National Park Foundation’s (NPF) Explorascope 60AZ is a manual alt-azimuth, refractor telescope perfect for celestial viewing on the go. The NPF Explorascope 60 can view the planets, moon, star clusters and brighter deep sky objects like the Orion Nebula and Andromeda Galaxy.\", \"picture\": \"NationalParkFoundationExplorascope.jpg\", \"priceUsd\": { \"currencyCode\": \"USD\", \"units\": 101, \"nanos\": 960000000 }, \"categories\": [ \"telescopes\" ] }";
//
//        try {
//            Demo.Product.Builder builder = Demo.Product.newBuilder();
//            JsonFormat.parser().ignoringUnknownFields().merge(json, builder);
//            return builder.build();
//        } catch (InvalidProtocolBufferException e) {
//            throw new RuntimeException(e);
//        }


        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel);
        return productCatalogStub.getProduct(request);
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = CurrencyServiceGrpc.newBlockingStub(currencyChannel);
        return currencyServiceBlockingStub.convert(request);
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());
        // query: currencyCode productId
        if (!exchange.getRequest().getMethod().equals(HttpMethod.POST)) {
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            return exchange.getResponse().setComplete();
        }

        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String productId = queryParams.getFirst("productId");
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
                orderJson.add("priceUsd", moneyJson);
                sink.success(orderJson);

            } catch (Exception e) {
                sink.error(e);
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
