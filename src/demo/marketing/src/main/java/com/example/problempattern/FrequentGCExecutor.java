package com.example.problempattern;

import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

@Component
public class FrequentGCExecutor {

    private final Logger logger = LogManager.getLogger(FrequentGCExecutor.class);

    public static final String FREQUENT_YOUNG_GC_LOOP_COUNT_FLAG = "youngGcCountLoopCountFlag";

    @Autowired
    private OpenFeatureAPI openFeatureAPI;

    public void youngGcExecute() {

        final Integer count = openFeatureAPI.getClient().getIntegerValue(FREQUENT_YOUNG_GC_LOOP_COUNT_FLAG, 0);
        if (count < 1) {
            logger.info("Young GC executor is disabled.");
            return ;
        }
        final long loopCountSum = 1024L * count;
        final int youngGcAllocSize = 128;

        for (int i = 0; i < loopCountSum; i++) {
            //short life cycle object.
            byte[] tempArray = new byte[youngGcAllocSize];
            if (i % loopCountSum == 0) {
                try {
                    Thread.sleep(500);
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                }
            }
        }

        logger.info("Young GC Executed.");
    }
}
