package com.example;

import com.example.config.ScenarioConfig;
import dev.openfeature.contrib.providers.flagd.FlagdOptions;
import dev.openfeature.contrib.providers.flagd.FlagdProvider;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@EntityScan(basePackages = "com.example.entity")
@ComponentScan("com.example")
@EnableConfigurationProperties(ScenarioConfig.class)
public class AdsApplication {

	public static void main(String[] args) {
		FlagdProvider flagd = new FlagdProvider(FlagdOptions.builder().host("flagd").port(8013).build());
		OpenFeatureAPI.getInstance().setProvider(flagd);
		SpringApplication.run(AdsApplication.class, args);
	}
}
