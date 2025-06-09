package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.protobuf.util.JsonFormat;
import io.grpc.ManagedChannel;
import io.grpc.netty.NettyChannelBuilder;
import org.example.config.Config;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.MediaType;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import reactor.core.publisher.Mono;
import oteldemo.Demo;

import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;


public class CurrencyFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(CurrencyFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel currencyChannel;

    public CurrencyFilter(Config config) {

        currencyChannel = NettyChannelBuilder.forTarget(config.currencyAddr).usePlaintext() // 明文通信（仅限开发环境）
                .maxInboundMessageSize(1024 * 1024 * 20) // 20MB 最大消息
                .keepAliveTime(30, TimeUnit.SECONDS) // 保活间隔
                .keepAliveTimeout(10, TimeUnit.SECONDS) // 保活超时
                .keepAliveWithoutCalls(true) // 即使没有活跃调用也发送keepalive
                .enableRetry() // 启用重试
                .idleTimeout(5, TimeUnit.MINUTES) // 添加空闲超时
                .build();
    }

    private Demo.GetSupportedCurrenciesResponse DoGetSupportedCurrencies(Demo.Empty request) {
        oteldemo.CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = oteldemo.CurrencyServiceGrpc.newBlockingStub(currencyChannel)
                .withDeadlineAfter(3, TimeUnit.SECONDS); // 添加3秒超时 - 获取支持的货币列表
        return currencyServiceBlockingStub.getSupportedCurrencies(request);
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, GatewayFilterChain chain) {
        log.info("Filtering request {}", exchange.getRequest().getPath());

        return Mono.<String>create(sink -> {
            try {
                Demo.Empty.Builder builder = Demo.Empty.newBuilder();
                Demo.GetSupportedCurrenciesResponse currencyCodes = DoGetSupportedCurrencies(builder.build());
                List<String> currencyList = currencyCodes.getCurrencyCodesList();
                String result = "[" + currencyList.stream().map(code -> "\"" + code + "\"").collect(Collectors.joining(",")) + "]";                sink.success(result);
            } catch (Exception e) {
                log.error("Failed Currency", e);
                sink.error(e);
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed Currency", e);
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
