package com.example.problempattern;

import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.HashMap;

@Component
public class FrequentGCExecutor {

    private final Logger logger = LogManager.getLogger(FrequentGCExecutor.class);

    private static final String FREQUENT_YOUNG_GC_LOOP_COUNT_FLAG = "youngGcCountLoopCountFlag";

    private static final String FREQUENT_OLD_GC_LOOP_COUNT_FLAG = "oldGcLoopCountFlag";

    private static final String FREQUENT_JVM_OOM_FLAG = "jvmOOMFlag";

    private final HashMap<String, byte[]> bigObjectCache  = new HashMap<>();

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

    public void oldGcExecute() {
        final Integer count = openFeatureAPI.getClient().getIntegerValue(FREQUENT_OLD_GC_LOOP_COUNT_FLAG, 0);
        if (count < 1) {
            logger.info("Old GC executor is disabled.");
            return;
        }
        final int allocSize = 1024 * 1024;

        for (int i = 0; i < count; i++) {
            byte[] data = new byte[allocSize];
            bigObjectCache.put(String.valueOf(System.currentTimeMillis()), data);
            if (i++ % 100 == 0) {
                Runtime runtime = Runtime.getRuntime();
                logger.info("Memory Usage(Every 100): Size={},  Used={}MB, Max={}MB",
                        bigObjectCache.size(),
                        (runtime.totalMemory() - runtime.freeMemory()) / (1024 * 1024),
                        runtime.maxMemory() / (1024 * 1024));
            }
        }

        logger.info("Old GC Executed.");
    }

    public void oomExecute() {
        final boolean oomFlag = openFeatureAPI.getClient().getBooleanValue(FREQUENT_JVM_OOM_FLAG, Boolean.FALSE);
        if (!oomFlag) {
            logger.info("OOM executor is disabled.");
            return;
        }

        final int allocSize = 1024 * 1024;
        int i = 0;

        while (true) {
            byte[] data = new byte[allocSize];
            bigObjectCache.put(String.valueOf(System.currentTimeMillis()), data);

            if (i++ % 512 == 0) {
                Runtime runtime = Runtime.getRuntime();
                logger.info("Memory Usage(Every 512):CacheSize={},  Used={}MB, Max={}MB",
                        bigObjectCache.size(),
                        (runtime.totalMemory() - runtime.freeMemory()) / (1024 * 1024),
                        runtime.maxMemory() / (1024 * 1024));
            }
        }
    }

}
