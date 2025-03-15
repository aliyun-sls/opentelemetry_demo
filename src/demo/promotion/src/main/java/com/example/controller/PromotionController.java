package com.example.controller;

import com.example.entity.PromotionEntity;
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

    @GetMapping("/listAds")
    public ResponseEntity<List<PromotionEntity>> listAds() {
        List<PromotionEntity> ads = promotionService.listAds();
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ad/{id}")
    public ResponseEntity<PromotionEntity> getAdById(@PathVariable Long id) {
        PromotionEntity ad = promotionService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}