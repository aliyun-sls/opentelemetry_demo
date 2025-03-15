package com.example.service;

import com.example.entity.PromotionEntity;
import com.example.repository.PromotionRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class PromotionServiceImpl implements PromotionService {
    @Autowired
    private PromotionRepository promotionRepository;
    @Override
    public List<PromotionEntity> listAds() {
        return promotionRepository.findAll();
    }

    @Override
    public PromotionEntity getAdById(Long id) {
        return promotionRepository.getOne(String.valueOf(id));
    }
}
