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

@Controller
public class PromotionController {

    @Autowired
    private PromotionService promotionService;
    @Autowired
    private AdsService adsService;

    @GetMapping("/listPromotion")
    public ResponseEntity<List> listPromotion() {
//        List<PromotionEntity> ads = promotionService.listAds();
        List<AdEntity> ads = adsService.listAds().block();
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/promotion/{id}")
    public ResponseEntity<PromotionEntity> getPromotionById(@PathVariable Long id) {
        PromotionEntity ad = promotionService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}