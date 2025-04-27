package com.example.problempattern;

import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.stereotype.Component;

import java.util.concurrent.ExecutorService;

@Component
public class ThreadPoolDepletionExecutor {

    private final Logger logger = LogManager.getLogger(ThreadPoolDepletionExecutor.class);

    public static final String THREAD_POOL_DEPLETION_FLAG = "threadPoolDepletionFlag";

    @Autowired
    private OpenFeatureAPI openFeatureAPI;

    @Autowired
    @Qualifier("customThreadPool")
    private ExecutorService executorService;

    public void execute() {
        final Client client = openFeatureAPI.getClient();
        final boolean switchOnDepletion = client.getBooleanValue(THREAD_POOL_DEPLETION_FLAG, Boolean.FALSE);

        if (!switchOnDepletion) {
            logger.info("ThreadPoolDepletion executor is disabled.");
            return;
        }

        final Runnable blockingTask = () -> {
            try {
                logger.info(Thread.currentThread().getName() + " is running and will be blocked.");
                Thread.sleep(10000); // sleep 10s.
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                logger.error(Thread.currentThread().getName() + " was interrupted.");
            }
        };

        executorService.submit(blockingTask);
        logger.info("Task submitted successfully.");
    }
}
