package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.concurrent.*;

@org.springframework.stereotype.Controller
public class Controller {

    // 每次调用时重新创建线程池
    ExecutorService executorService =  new ThreadPoolExecutor(
        Runtime.getRuntime().availableProcessors(),  //corePoolSize:线程核心数量，及时处于idle状态，也不会回收
                20, //maximumPoolSize:线程数的上限
                60, //keepAliveTime:超过这个时间，超过corePoolSize的线程，多余的线程将会被回收
        TimeUnit.SECONDS,
        new LinkedBlockingQueue<>(10000), //任务的排队队列  不显示指定大小，会默认Integer.MAX_VALUE，设置不当容易OOM
        new ThreadPoolExecutor.DiscardPolicy()  //拒绝策略
                );

    // 从环境变量中获取默认值为1000的值
    private static final int MAX_DURATION = Integer.parseInt(System.getenv().getOrDefault("MAX_DURATION", "1000"));

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
                    if(System.currentTimeMillis() - startTime < MAX_DURATION){
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
