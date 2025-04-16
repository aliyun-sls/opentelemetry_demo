package com.example.problempattern;

import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.concurrent.ThreadLocalRandom;

@Component
public class HighCPUExecutor {

    public static final String HIGH_CPU_COST_LOOP_COUNT_FLAG = "marketingHighCpuCostLoopCountFlag";

    @Autowired
    private OpenFeatureAPI openFeatureAPI;

    private final Logger logger = LogManager.getLogger(HighCPUExecutor.class);

    public void execute() {
        final Client client = openFeatureAPI.getClient();
        Integer loopCount = client.getIntegerValue(HIGH_CPU_COST_LOOP_COUNT_FLAG, 0);

        if (loopCount == 0) {
            logger.info("HighCPUExecutor is disabled.");
            return;
        }

        double log = ThreadLocalRandom.current().nextDouble();
        for (int i = 0; i < loopCount; i++) {
            double sqrt = Math.sqrt(ThreadLocalRandom.current().nextDouble());
            double atan = Math.atan(Math.acos(sqrt));
            double pow = Math.pow(atan, sqrt);
            log = Math.log(pow);
        }

        logger.info("HighCPUExecutor is running. Loop count: {}", log);
    }

}
