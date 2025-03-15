package com.example.service;

import com.example.entity.PromotionEntity;

import java.util.List;

public interface PromotionService {
    List<PromotionEntity> listAds();

    PromotionEntity getAdById(Long id);
}
