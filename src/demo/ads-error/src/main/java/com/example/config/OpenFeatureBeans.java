package com.example.config;

import dev.openfeature.contrib.providers.flagd.FlagdOptions;
import dev.openfeature.contrib.providers.flagd.FlagdProvider;
import dev.openfeature.sdk.exceptions.OpenFeatureError;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenFeatureBeans {

    @Bean
    public OpenFeatureAPI OpenFeatureAPI() {
        final OpenFeatureAPI openFeatureAPI = OpenFeatureAPI.getInstance();

        // Use flagd as the OpenFeature provider and use default configurations
        try {
            openFeatureAPI.setProviderAndWait(new FlagdProvider(FlagdOptions.builder().host("flagd").port(8013).deadline(5000).build()));
        } catch (OpenFeatureError e) {
            throw new RuntimeException("Failed to set OpenFeature provider", e);
        }

        return openFeatureAPI;
    }
}