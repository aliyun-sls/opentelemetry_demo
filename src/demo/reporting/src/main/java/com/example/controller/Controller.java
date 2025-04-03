package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

import javax.annotation.PreDestroy;
import java.util.Random;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.RejectedExecutionException;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.logging.Logger;

@org.springframework.stereotype.Controller
public class Controller {

    private static final Random random = new Random();
    private static final Logger logger = Logger.getLogger(Controller.class.getName());

    private final ExecutorService executorService = Executors.newFixedThreadPool(1000); // 创建线程池
    private final ScheduledExecutorService scheduledExecutorService = Executors.newScheduledThreadPool(1); // 创建调度线程池

    @GetMapping("/reporting")
    public ResponseEntity<?> reporting() {
        try {
            executorService.submit(() -> {
                while (true) {
                    for (int i = 0; i < Integer.MAX_VALUE; i++) {
                        // 进行一些计算密集型的操作
                        double dummy = Math.atan(Math.random());

                        // 每进行一定次数的计算后休眠，避免过高的CPU占用
                        if (i % 1000000 == 0) {
                            try {
                                Thread.sleep(10); // 休眠10毫秒
                            } catch (InterruptedException e) {
                                Thread.currentThread().interrupt();
                                return; // 如果被中断，终止线程
                            }
                        }
                    }
                }
            });
        } catch (RejectedExecutionException e) {
            logger.severe("Task submission rejected: " + e.getMessage());
        }

        // 随机生成一个 0 到 1000 毫秒之间的延迟时间
        int delay = random.nextInt(1500);

        // 调度一个任务在延迟时间后关闭线程池
        scheduledExecutorService.schedule(() -> {
            executorService.shutdownNow();
        }, delay, TimeUnit.MILLISECONDS);

        return ResponseEntity.ok("CPU load started and will be shut down in " + delay + " milliseconds!");
    }

    @PreDestroy
    public void shutdown() {
        // 确保在应用关闭时关闭所有线程池
        if (!executorService.isShutdown()) {
            executorService.shutdownNow();
            logger.info("ExecutorService has been shut down.");
        }
        if (!scheduledExecutorService.isShutdown()) {
            scheduledExecutorService.shutdownNow();
            logger.info("ScheduledExecutorService has been shut down.");
        }
    }
}