package com.example.service;

import com.example.entity.Scenario;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Random;

@Service
public class ScenarioSelector {

    @Autowired
    ScenarioConfig scenarioConfig;
    private List<Scenario> scenarios;
    private final Random random = new Random();

    public void setScenarioSelector(List<Scenario> scenarios) {
        if (scenarios.isEmpty()) {
            this.scenarios = scenarioConfig.getScenarios();
        }
        // 验证概率总和是否为 1.0
        double total = scenarios.stream().mapToDouble(Scenario::getWeight).sum();
        if (Math.abs(total - 1.0) > 0.0001) {
            throw new IllegalArgumentException("概率总和必须为 1.0");
        }
    }

    public String selectScenario() {
        double randomValue = random.nextDouble();
        double accumulated = 0.0;
        for (Scenario scenario : scenarios) {
            accumulated += scenario.getWeight();
            if (randomValue <= accumulated) {
                return scenario.getName();
            }
        }
        return "normal"; // 默认返回正常场景
    }
}
