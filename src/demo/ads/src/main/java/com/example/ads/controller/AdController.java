package com.example.ads.controller;

import com.example.ads.entity.AdEntity;
import com.example.ads.service.AdService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

import java.util.List;

@Controller
public class AdController {

    @Autowired
    AdService adService;
    @GetMapping("/listAds")
    public ResponseEntity<List<AdEntity>> listAds() {
        List<AdEntity> ads = adService.listAds();
        return ResponseEntity.ok(ads);
    }
}