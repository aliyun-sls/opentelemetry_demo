package com.example.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;

import java.util.List;

@Service
public class AdsService {

    private final WebClient webClient;

    @Autowired
    public AdsService(WebClient webClient) {
        this.webClient = webClient;
    }

    public Mono<List> listAds() {
        return webClient.get()
                .uri("/listAds")
                .retrieve()
                .bodyToMono(List.class);
    }
}
