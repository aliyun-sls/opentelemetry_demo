package com.example.controller;

import com.example.entity.AdEntity;
import com.example.service.AdService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Controller
public class AdController {

    @Autowired
    private AdService adService;

    @GetMapping("/listAds")
    public ResponseEntity<List> listAds() {
        List ads;
        Random random = new Random();
        double n = random.nextDouble();
        if (n < 0.1) {
            ads = adService.findAdWithSleep();
        }else {
            if (n<0.4){
                try {
                    TimeUnit.SECONDS.sleep(10);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
            ads = adService.listAds();
        }
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/ads/listAds")
    public ResponseEntity<List<AdEntity>> listAdsAll() {
        List ads = null;
        Random random = new Random();
        double n = random.nextDouble();
        if (n < 0.1) {
            ads = adService.findAdWithSleep();
        }else {
            if (n<0.4){
                try {
                    TimeUnit.SECONDS.sleep(10);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
            ads = adService.listAds();
        }
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ad/{id}")
    public ResponseEntity<AdEntity> getAdById(@PathVariable Long id) {
        AdEntity ad = adService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}