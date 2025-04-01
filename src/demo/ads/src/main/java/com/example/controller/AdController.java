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

import java.util.ArrayList;
import java.util.List;

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
        List ads = new ArrayList();
        String selectedScenario = scenarioSelector.selectScenario();
        try {
            scenarioHandler.handleScenario(selectedScenario);
        } catch (Exception e) {
            return ResponseEntity.status(500).build();
        }
        return ResponseEntity.ok(ads);
    }

    @GetMapping("/ads/listAds")
    public ResponseEntity<List> listAdsAll() {
        return listAds();
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ads/{id}")
    public ResponseEntity<AdEntity> getAdById(@PathVariable Long id) {
        AdEntity ad = adService.getAdById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }

    // 通过 REST API 更新场景配置
    @PostMapping("/ads/update-scenarios")
    public ResponseEntity<?> updateScenarios(@RequestBody List<Scenario> newScenarios) {
        scenarioConfig.setScenarios(newScenarios);
        scenarioSelector = new ScenarioSelector(newScenarios);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/ads/scenarios")
    public ResponseEntity<List<Scenario>> scenarios() {
        List<Scenario> s = scenarioConfig.getScenarios();
        return ResponseEntity.ok(s);
    }

}