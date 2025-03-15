package com.example.service;

import com.example.entity.RecommendationEntity;
import com.example.repository.RecommendationRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class RecommendationServiceImpl implements RecommendationService {
    @Autowired
    private RecommendationRepository recommendationRepository;
    @Override
    public List<RecommendationEntity> listAds() {
        return recommendationRepository.findAll();
    }

    @Override
    public RecommendationEntity getAdById(Long id) {
        return recommendationRepository.getOne(String.valueOf(id));
    }
}
