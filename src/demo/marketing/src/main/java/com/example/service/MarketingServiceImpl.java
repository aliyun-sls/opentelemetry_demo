package com.example.service;

import com.example.entity.MarketingEntity;
import com.example.repository.MarketingRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class MarketingServiceImpl implements MarketingService {
    @Autowired
    private MarketingRepository marketingRepository;
    @Override
    public List<MarketingEntity> listAds() {
        return marketingRepository.findAll();
    }

    @Override
    public MarketingEntity getAdById(Long id) {
        return marketingRepository.getOne(String.valueOf(id));
    }
}
