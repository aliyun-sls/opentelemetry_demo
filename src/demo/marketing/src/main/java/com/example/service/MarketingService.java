package com.example.service;

import com.example.entity.MarketingEntity;

import java.util.List;

public interface MarketingService {
    List<MarketingEntity> listMarketingEntity();

    MarketingEntity getAdById(Long id);
}
