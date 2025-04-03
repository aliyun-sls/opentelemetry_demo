package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

import javax.annotation.PreDestroy;
import java.util.Random;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

@org.springframework.stereotype.Controller
public class Controller {

    // 每次调用时重新创建线程池
    ExecutorService executorService = Executors.newFixedThreadPool(1000);

    @GetMapping("/reporting")
    public ResponseEntity<?> reporting() {
        executorService.submit(() -> {
            System.out.println("Starting CPU load simulation...");
            long startTime = System.currentTimeMillis();
            for (int i = 0; i < Integer.MAX_VALUE; i++) {
                // 进行一些计算密集型的操作
                double dummy = Math.atan(Math.random());
                // 每进行一定次数的计算后休眠，避免过高的CPU占用
                if (i % 100000 == 0) {
                    if(System.currentTimeMillis() - startTime < 100){
                        System.out.println("CPU load simulation ended.");
                        return;
                    }
                    try {
                        Thread.sleep(10); // 休眠10毫秒
                    } catch (InterruptedException e) {
                        Thread.currentThread().interrupt();
                        return; // 如果被中断，终止线程
                    }
                }
            }
        });


        return ResponseEntity.ok("CPU load started and will be shut down in 1.5 seconds!");
    }
}
