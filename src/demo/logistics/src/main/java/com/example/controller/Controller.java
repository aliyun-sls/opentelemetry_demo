package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

@org.springframework.stereotype.Controller
public class Controller {

    @GetMapping("/logistics")
    public ResponseEntity<?> logistics() {
        // 模拟CPU高的计算任务
        performCpuIntensiveTask();
        return ResponseEntity.ok().build();
    }

    private void performCpuIntensiveTask() {
        int sum = 0;
        for (int i = 0; i < 100000000; i++) { // 大量循环
            sum += i; // 耗时计算
        }
    }
}
