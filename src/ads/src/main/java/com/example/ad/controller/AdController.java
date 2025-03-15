package com.example.ad.controller;

import com.example.ad.entity.AdEntity;
import com.example.ad.service.AdService;
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