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

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;

public class ListProductFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(ListProductFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private ManagedChannel productCatalogChannel;
    private ManagedChannel currencyChannel;

    public ListProductFilter(Config config) {
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


    private Demo.ListProductsResponse DoListProducts(Demo.Empty request) {
        ProductCatalogServiceGrpc.ProductCatalogServiceBlockingStub productCatalogStub = ProductCatalogServiceGrpc.newBlockingStub(productCatalogChannel)
                .withDeadlineAfter(5, TimeUnit.SECONDS); // 添加5秒超时 - 产品列表查询
        return productCatalogStub.listProducts(request);
    }

    private Demo.Money DoCurrencyConvert(Demo.CurrencyConversionRequest request) {
        CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = CurrencyServiceGrpc.newBlockingStub(currencyChannel)
                .withDeadlineAfter(3, TimeUnit.SECONDS); // 添加3秒超时 - 货币转换
        return currencyServiceBlockingStub.convert(request);
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());
        // query: currencyCode productId
        if (!exchange.getRequest().getMethod().equals(HttpMethod.GET)) {
            exchange.getResponse().setStatusCode(HttpStatus.METHOD_NOT_ALLOWED);
            return exchange.getResponse().setComplete();
        }

        MultiValueMap<String, String> queryParams = exchange.getRequest().getQueryParams();
        String currencyCode = queryParams.getFirst("currencyCode");
        if (currencyCode == null || currencyCode.isEmpty()) {
            currencyCode = "USD";
        }
        String finalCurrencyCode = currencyCode;


        return Mono.<List<Demo.Product>>create(sink-> {
            Demo.ListProductsResponse listProductsResponse = DoListProducts(Demo.Empty.newBuilder().build());
            sink.success(listProductsResponse.getProductsList());
        }).flatMap(products -> {

            try {
                List<JsonObject> rst = new ArrayList<>();

                for (Demo.Product product : products) {
                    JsonObject productJson = new Gson().fromJson(JsonFormat.printer().print(product), JsonObject.class);

                    Demo.CurrencyConversionRequest.Builder request = Demo.CurrencyConversionRequest.newBuilder();
                    request.setFrom(product.getPriceUsd());
                    request.setToCode(finalCurrencyCode);

                    Demo.Money money = DoCurrencyConvert(request.build());
                    JsonObject moneyJson = new Gson().fromJson(JsonFormat.printer().print(money), JsonObject.class);
                    if (moneyJson.get("units") != null) {
                        moneyJson.addProperty("units", Integer.valueOf(moneyJson.get("units").getAsString()));
                    }
                    productJson.add("priceUsd", moneyJson);

                    rst.add(productJson);
                }

                return Mono.just(rst);
            }catch (Exception e) {
                log.info("ListProductFilter", e);
                return Mono.error(e);
            }
        }).map(f -> new Gson().toJson(f)).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.info("ListProductFilter", e);
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
