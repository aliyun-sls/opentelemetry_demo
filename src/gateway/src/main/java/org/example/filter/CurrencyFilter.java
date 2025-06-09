package org.example.filter;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.google.gson.Gson;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.example.config.Config;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.cloud.gateway.filter.GatewayFilter;
import org.springframework.cloud.gateway.filter.GatewayFilterChain;
import org.springframework.core.io.buffer.DataBuffer;
import org.springframework.http.MediaType;
import org.springframework.web.server.ResponseStatusException;
import org.springframework.web.server.ServerWebExchange;
import oteldemo.Demo;
import oteldemo.CurrencyServiceGrpc;
import reactor.core.publisher.Mono;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.stream.Collectors;


public class CurrencyFilter implements GatewayFilter {
    private static final Logger log = LoggerFactory.getLogger(CurrencyFilter.class);
    public static ObjectMapper objectMapper = new ObjectMapper();
    private final ManagedChannel currencyChannel;

    public CurrencyFilter(Config config) {
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

    private List<Demo.GetSupportedCurrenciesResponse> DoGetSupportedCurrencies() {
        CurrencyServiceGrpc.CurrencyServiceBlockingStub currencyServiceBlockingStub = CurrencyServiceGrpc.newBlockingStub(currencyChannel)
                .withDeadlineAfter(3, TimeUnit.SECONDS); // 添加3秒超时 - 货币服务
        
        try {
            List<Demo.GetSupportedCurrenciesResponse> currenciesList = new ArrayList<>();
            currenciesList.add(currencyServiceBlockingStub.getSupportedCurrencies(Demo.Empty.newBuilder().build()));
            return currenciesList;
        } catch (io.grpc.StatusRuntimeException e) {
            log.error("gRPC error getting supported currencies, retrying once: {}", e.getMessage());
            // 连接错误时重试一次
            if (e.getStatus().getCode() == io.grpc.Status.Code.INTERNAL || 
                e.getStatus().getCode() == io.grpc.Status.Code.UNAVAILABLE) {
                try {
                    Thread.sleep(500); // 短暂延迟后重试
                    List<Demo.GetSupportedCurrenciesResponse> currenciesList = new ArrayList<>();
                    currenciesList.add(currencyServiceBlockingStub.withDeadlineAfter(6, TimeUnit.SECONDS)
                            .getSupportedCurrencies(Demo.Empty.newBuilder().build()));
                    return currenciesList;
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
        return Mono.<String>create(sink -> {
            try {
                List<Demo.GetSupportedCurrenciesResponse> currencyCodes = DoGetSupportedCurrencies();
                List<String> currencyList = currencyCodes.stream()
                        .flatMap(response -> response.getCurrencyCodesList().stream())
                        .collect(Collectors.toList());
                String result = "[" + currencyList.stream().map(code -> "\"" + code + "\"").collect(Collectors.joining(",")) + "]";
                sink.success(result);
            } catch (Exception e) {
                log.error("Failed Currency", e);
                sink.error(new RuntimeException(e));
            }
        }).flatMap(responseBody -> {
            try {
                exchange.getResponse().getHeaders().setContentType(MediaType.APPLICATION_JSON);
                byte[] bytes = responseBody.getBytes();
                DataBuffer buffer = exchange.getResponse().bufferFactory().wrap(bytes);
                return exchange.getResponse().writeWith(Mono.just(buffer));
            } catch (Exception e) {
                log.error("Failed Currency", e);
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
