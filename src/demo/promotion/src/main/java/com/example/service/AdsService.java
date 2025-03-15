package com.example.service;

import com.example.entity.AdEntity;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.util.List;

@Service
public class AdsService {
    @Bean
    public WebClient webClient(WebClient.Builder webClientBuilder) {
        return webClientBuilder.baseUrl("http://ads:8080").build();
    }

    @Autowired
    private WebClient webClient;

    public Mono<List> listAds() {
        return webClient.get()
                .uri("/ads/listAds")
                .retrieve()
                .bodyToMono(List.class);
    }
}
