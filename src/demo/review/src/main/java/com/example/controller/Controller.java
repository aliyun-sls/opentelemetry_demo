package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.HashMap;
import java.util.Map;
import java.util.Random;

@org.springframework.stereotype.Controller
public class Controller {

    @GetMapping("/review")
    public ResponseEntity<?> review() {
        // 构建一个大的Map
        Map<Integer, String> largeMap = new HashMap<>();
        Random random = new Random();
        for (int i = 0; i < 1000000; i++) {
            largeMap.put(i, "value" + random.nextInt(100000));
        }

        // 在Map中查找数据
        String value = largeMap.get(random.nextInt(1000000));

        return ResponseEntity.ok(value);
    }
}
