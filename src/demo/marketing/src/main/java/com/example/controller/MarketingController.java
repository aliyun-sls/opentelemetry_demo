package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.MarketingEntity;
import com.example.service.AdsService;
import com.example.service.MarketingService;
import com.google.gson.Gson;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.Collections;
import java.util.List;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Controller
public class MarketingController {

    @Autowired
    private MarketingService marketingService;
    @Autowired
    private AdsService adsService;

    @GetMapping("/listMarketing")
    public ResponseEntity<List<AdEntity>> listMarketing() {
        Random random = new Random();
        if (random.nextDouble() < 0.1) {
            try {
                TimeUnit.SECONDS.sleep(10);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
            }
        }

        List<MarketingEntity> marketingEntities = marketingService.listMarketingEntity();
        List<AdEntity> ads = adsService.listAds().block();
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/marketing/{id}")
    public ResponseEntity<MarketingEntity> getMarketingById(@PathVariable Long id) {
        MarketingEntity ad = marketingService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}