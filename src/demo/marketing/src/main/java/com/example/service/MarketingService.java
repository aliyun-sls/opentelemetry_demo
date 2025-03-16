package com.example.service;

import com.example.entity.MarketingEntity;

import java.util.List;

public interface MarketingService {
    List<MarketingEntity> listAds();

    MarketingEntity getAdById(Long id);
}
