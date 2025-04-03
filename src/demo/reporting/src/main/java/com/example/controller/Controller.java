package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.Random;

@org.springframework.stereotype.Controller
public class Controller {

    private static final int ARRAY_SIZE = Integer.parseInt(System.getenv().getOrDefault("ARRAY_SIZE", "10000000")); // 默认值为 1000000    private static final int[] largeArray = new int[ARRAY_SIZE];
    private static final int[] largeArray = new int[ARRAY_SIZE];
    private static final Random random = new Random();

    static {
        // 初始化数组
        for (int i = 0; i < ARRAY_SIZE; i++) {
            largeArray[i] = random.nextInt(ARRAY_SIZE);
        }
    }

    @GetMapping("/reporting")
    public ResponseEntity<?> reporting() {
        int targetValue = random.nextInt(ARRAY_SIZE); // 随机选择一个目标值
        int index = findInArray(String.valueOf(targetValue)); // 在数组中查找目标值
        return ResponseEntity.ok("Found at index: " + index);
    }

    private int findInArray(String target) {
        for (int i = 0; i < ARRAY_SIZE; i++) {
            if (largeArray[i] == Integer.parseInt(target)) {
                return i; // 找到目标值，返回索引
            }
        }
        return -1; // 未找到目标值，返回-1
    }
}
