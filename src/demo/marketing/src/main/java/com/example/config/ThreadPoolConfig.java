package com.example.config;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.util.concurrent.*;

@Configuration
public class ThreadPoolConfig {
    @Bean("customThreadPool")
    public ExecutorService customThreadPool() {
        return new ThreadPoolExecutor(
                2,
                2,
                60L, TimeUnit.SECONDS,
                new LinkedBlockingQueue<>(1),  //fast fail
                new ThreadPoolExecutor.AbortPolicy()
        );
    }
}
