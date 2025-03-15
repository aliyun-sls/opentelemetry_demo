package com.example.controller;

import com.example.entity.RecommendationEntity;
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

    @GetMapping("/listAds")
    public ResponseEntity<List<RecommendationEntity>> listAds() {
        List<RecommendationEntity> ads = recommendationService.listAds();
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ad/{id}")
    public ResponseEntity<RecommendationEntity> getAdById(@PathVariable Long id) {
        RecommendationEntity ad = recommendationService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}