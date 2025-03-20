package com.example.controller;

import com.example.entity.AdEntity;
import com.example.service.AdService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@Controller
public class AdController {

    @Autowired
    private AdService adService;

    @GetMapping("/listAds")
    public ResponseEntity<List<AdEntity>> listAds() {
        List<AdEntity> ads = adService.listAds();
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/ads/listAds")
    public ResponseEntity<List<AdEntity>> listAdsAll() {
        List<AdEntity> ads = adService.listAds();
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ad/{id}")
    public ResponseEntity<AdEntity> getAdById(@PathVariable Long id) {
        AdEntity ad = adService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}