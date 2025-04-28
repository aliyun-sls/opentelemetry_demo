package com.example.problempattern;

import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class ServiceDowntimeExecutor {

    private final Logger logger = LogManager.getLogger(ServiceDowntimeExecutor.class);

    public static final String EXIT_SERVICE_DOWNTIME_FLAG = "serviceDowntimeFlag";

    @Autowired
    private OpenFeatureAPI openFeatureAPI;


    public void downtimeExecute() {
        final Client client = openFeatureAPI.getClient();
        final boolean downtimeFlag = client.getBooleanValue(EXIT_SERVICE_DOWNTIME_FLAG, Boolean.FALSE);

        if (!downtimeFlag) {
            logger.info("Exit executor is disabled.");
            return;
        }
        // service downtime
        System.exit(0);
    }
}
