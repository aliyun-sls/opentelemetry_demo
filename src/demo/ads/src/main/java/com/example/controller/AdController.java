package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.Scenario;
import com.example.service.AdService;
import com.example.service.ScenarioConfig;
import com.example.service.ScenarioHandler;
import com.example.service.ScenarioSelector;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import java.util.List;
import java.util.Random;
import java.util.concurrent.TimeUnit;

@Controller
public class AdController {

    @Autowired
    private AdService adService;

    @Autowired
    private ScenarioConfig scenarioConfig;

    @Autowired
    private ScenarioSelector scenarioSelector;

    @Autowired
    private ScenarioHandler scenarioHandler;

    @GetMapping("/listAds")
    public ResponseEntity<List> listAds() {
        List ads;

        String selectedScenario = scenarioSelector.selectScenario();

        try {
            scenarioHandler.handleScenario(selectedScenario);
        } catch (Exception e) {
            return ResponseEntity.status(500).body(List.of("Error: " + e.getMessage()));
        }

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

    // 通过 REST API 更新场景配置
    @PostMapping("/update-scenarios")
    public ResponseEntity<?> updateScenarios(@RequestBody List<Scenario> newScenarios) {
        scenarioConfig.setScenarios(newScenarios);
        scenarioSelector = new ScenarioSelector(newScenarios);
        return ResponseEntity.ok().build();
    }

}