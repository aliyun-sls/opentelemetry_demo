package com.example.controller;

import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;


@org.springframework.stereotype.Controller
public class Controller {

    @GetMapping("/review")
    public ResponseEntity<?> review() {
        return ResponseEntity.ok().build();
    }
}