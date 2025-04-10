package com.example.flagd;
import dev.openfeature.contrib.providers.flagd.FlagdOptions;
import dev.openfeature.contrib.providers.flagd.FlagdProvider;
import dev.openfeature.sdk.OpenFeatureAPI;
import dev.openfeature.sdk.exceptions.OpenFeatureError;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenFeatureBeans {

    @Value("${FLAGD_HOST:localhost}")
    private String flagdHost;

    @Value("${FLAGD_PORT:8013}")
    private int flagdPort;

    @Bean
    public OpenFeatureAPI openFeatureAPI() {
        final OpenFeatureAPI openFeatureAPI = OpenFeatureAPI.getInstance();

        FlagdOptions options =
                FlagdOptions.builder()
                        .host(flagdHost)
                        .port(flagdPort)
                        .withGlobalTelemetry(true)
                        .build();
        try {
            openFeatureAPI.setProviderAndWait(new FlagdProvider(options));
        } catch (OpenFeatureError e) {
            throw new RuntimeException("Failed to set OpenFeature provider", e);
        }
        return openFeatureAPI;
    }
}