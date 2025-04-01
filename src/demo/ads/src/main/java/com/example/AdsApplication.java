package com.example;

import com.example.service.ScenarioConfig;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@EntityScan(basePackages = "com.example.entity")
@ComponentScan("com.example")
@EnableConfigurationProperties(ScenarioConfig.class)
public class AdsApplication {

	public static void main(String[] args) {
		SpringApplication.run(AdsApplication.class, args);
	}

	@Bean
	public CommandLineRunner commandLineRunner(ScenarioConfig scenarioConfig) {
		return args -> {
			System.out.println("Scenarios: " + scenarioConfig.getScenarios());
		};
	}
}
