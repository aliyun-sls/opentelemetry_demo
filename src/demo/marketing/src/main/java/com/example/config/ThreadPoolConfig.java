package com.example.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

@Configuration
public class ThreadPoolConfig {

    @Bean("customThreadPool")
    public ExecutorService customThreadPool() {
        // nThread == 2
        return Executors.newFixedThreadPool(2);
    }
}
