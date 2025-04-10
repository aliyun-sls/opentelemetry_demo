package com.example.service;

import com.example.entity.MarketingEntity;
import com.example.problempattern.HighCPUExecutor;
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

    @Override
    public List<MarketingEntity> listMarketingEntity() {
        List<MarketingEntity> all = marketingRepository.findAll();
        highCPUExecutor.execute();
        return all;
    }

    @Override
    public MarketingEntity getAdById(Long id) {
        return marketingRepository.getOne(String.valueOf(id));
    }
}
