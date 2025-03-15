package com.example.service;

import com.example.entity.RecommendationEntity;

import java.util.List;

public interface RecommendationService {
    List<RecommendationEntity> listAds();

    RecommendationEntity getAdById(Long id);
}
