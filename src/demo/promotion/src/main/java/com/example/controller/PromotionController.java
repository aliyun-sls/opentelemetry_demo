package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.PromotionEntity;
import com.example.service.AdsService;
import com.example.service.PromotionService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Controller
public class PromotionController {

    private final PromotionService promotionService;
    private final AdsService adsService;

    @Autowired
    public PromotionController(PromotionService promotionService, AdsService adsService) {
        this.promotionService = promotionService;
        this.adsService = adsService;
    }

    @GetMapping("/listPromotion")
    public ResponseEntity<List> listPromotion() {
//        Random random = new Random();
//        if (random.nextDouble() < 0.1) {
//            try {
//                TimeUnit.SECONDS.sleep(10);
//            } catch (InterruptedException e) {
//                Thread.currentThread().interrupt();
//            }
//        }
        List ads = adsService.listAds().block();
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/promotion/{id}")
    public ResponseEntity<PromotionEntity> getPromotionById(@PathVariable Long id) {
        PromotionEntity ad = promotionService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}