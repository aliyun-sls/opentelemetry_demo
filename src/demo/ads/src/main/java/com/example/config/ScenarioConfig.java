package com.example.config;

import com.example.entity.Scenario;
import lombok.Data;
import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.context.annotation.Configuration;

import java.util.List;

@Configuration
@ConfigurationProperties(prefix = "ads-scenarios")
@Data
public class ScenarioConfig {
    private List<Scenario> scenarios;

    // Getters and Setters
    public List<Scenario> getScenarios() {
        return scenarios;
    }

    public void setScenarios(List<Scenario> scenarios) {
        System.out.println(scenarios);
        this.scenarios = scenarios;
    }
}
