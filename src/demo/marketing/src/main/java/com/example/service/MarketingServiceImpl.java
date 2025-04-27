package com.example.service;

import com.example.entity.MarketingEntity;
import com.example.problempattern.FrequentGCExecutor;
import com.example.problempattern.HighCPUExecutor;
import com.example.problempattern.SlowSQLExecutor;
import com.example.problempattern.ThreadPoolDepletionExecutor;
import com.example.repository.MarketingRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Collections;
import java.util.List;

@Service
public class MarketingServiceImpl implements MarketingService {
    @Autowired
    private MarketingRepository marketingRepository;

    @Autowired
    private HighCPUExecutor highCPUExecutor;

    @Autowired
    private ThreadPoolDepletionExecutor threadPoolDepletionExecutor;

    @Autowired
    private SlowSQLExecutor slowSQLExecutor;

    @Autowired
    private FrequentGCExecutor frequentGCExecutor;

    @Override
    public List<MarketingEntity> listMarketingEntity() {
        List<MarketingEntity> all = marketingRepository.findAll();
        //high cpu executor
        highCPUExecutor.execute();
        //thread pool depletion executor
        for(int i = 0; i < 4; i++) {
            threadPoolDepletionExecutor.execute();
        }
        //slowly sql
        slowSQLExecutor.execute();
        //frequent gc
        frequentGCExecutor.youngGcExecute();
        return all;
    }

    @Override
    public MarketingEntity getAdById(Long id) {
        return marketingRepository.getOne(String.valueOf(id));
    }
}
