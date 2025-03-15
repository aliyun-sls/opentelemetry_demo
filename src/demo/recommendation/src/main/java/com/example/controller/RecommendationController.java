package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.RecommendationEntity;
import com.example.service.AdsService;
import com.example.service.RecommendationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@Controller
public class RecommendationController {

    @Autowired
    private RecommendationService recommendationService;
    @Autowired
    private AdsService adsService;

    @GetMapping("/listRecommendation")
    public ResponseEntity<List<AdEntity>> listRecommendation() {
//        List<RecommendationEntity> ads = recommendationService.listAds();
        List<AdEntity> ads = adsService.listAds().block();
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/recommendation/{id}")
    public ResponseEntity<RecommendationEntity> getRecommendationById(@PathVariable Long id) {
        RecommendationEntity ad = recommendationService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}