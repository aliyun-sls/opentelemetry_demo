package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.Scenario;
import com.example.service.AdService;
import com.example.config.ScenarioConfig;
import com.example.service.ScenarioHandler;
import com.example.service.ScenarioSelector;
import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;

import javax.servlet.http.HttpServletRequest;
import java.util.ArrayList;
import java.util.Enumeration;
import java.util.List;

@Controller
public class AdController {

    private static final Logger logger = LoggerFactory.getLogger(AdController.class);
    private static final String ADS_FEATURE_FLAG = "adsFlag";
    @Autowired
    private AdService adService;

    @Autowired
    private ScenarioConfig scenarioConfig;

    @Autowired
    private ScenarioSelector scenarioSelector;

    @Autowired
    private ScenarioHandler scenarioHandler;

    private final OpenFeatureAPI openFeatureAPI;

    @Autowired
    public AdController(OpenFeatureAPI OFApi) {
        this.openFeatureAPI = OFApi;
    }

    @GetMapping("/listAds")
    public ResponseEntity<?> listAds(HttpServletRequest request) {
        List ads = new ArrayList();
//        String selectedScenario = scenarioSelector.selectScenario();
        final Client client = openFeatureAPI.getClient();
        String selectedScenario = client.getStringValue(ADS_FEATURE_FLAG, "normal");
        try {
            scenarioHandler.handleScenario(selectedScenario);
        } catch (Exception e) {
            logger.error("Error handling scenario: {}, Headers: {}, Exception: ", selectedScenario, logRequestHeaders(request), e);
            return ResponseEntity.status(500).body(e.getMessage());
        }
        logger.info("Selected scenario: {}, Headers: {}", selectedScenario, logRequestHeaders(request));
        return ResponseEntity.ok(ads);
    }

    private String logRequestHeaders(HttpServletRequest request) {
        StringBuilder headers = new StringBuilder();
        Enumeration<String> headerNames = request.getHeaderNames();
        if (headerNames != null) {
            while (headerNames.hasMoreElements()) {
                String headerName = headerNames.nextElement();
                String headerValue = request.getHeader(headerName);
                headers.append(headerName).append(" = ").append(headerValue).append(", ");
            }
        }
        return headers.toString();
    }

    @GetMapping("/ads/listAds")
    public ResponseEntity<?> listAdsAll(HttpServletRequest request) {
        return listAds(request);
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
        scenarioSelector = new ScenarioSelector(scenarioConfig);
        return ResponseEntity.ok().build();
    }

    @GetMapping("/ads/scenarios")
    public ResponseEntity<List<Scenario>> scenarios() {
        List<Scenario> s = scenarioConfig.getScenarios();
        return ResponseEntity.ok(s);
    }

}