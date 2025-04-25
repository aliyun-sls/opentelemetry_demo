package com.example.problempattern;

import com.example.entity.MarketingEntity;
import com.example.repository.MarketingRepository;
import dev.openfeature.sdk.Client;
import dev.openfeature.sdk.OpenFeatureAPI;
import org.apache.logging.log4j.LogManager;
import org.apache.logging.log4j.Logger;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;

import java.util.List;

@Component
public class SlowSQLExecutor {

    private final Logger logger = LogManager.getLogger(SlowSQLExecutor.class);

    public static final String SLOW_SQL_FLAG = "slowSQLFlag";

    @Autowired
    private OpenFeatureAPI openFeatureAPI;

    @Autowired
    private MarketingRepository marketingRepository;


    public void execute() {
        final Client client = openFeatureAPI.getClient();
        final boolean switchOnSlowSQL = client.getBooleanValue(SLOW_SQL_FLAG, Boolean.FALSE);

        if (!switchOnSlowSQL) {
            logger.info("SlowSQL executor is disabled.");
            return;
        }
        final List<MarketingEntity> data = marketingRepository.findWithDelay();
        logger.info("SlowSQL executor task executed.");
    }
}
