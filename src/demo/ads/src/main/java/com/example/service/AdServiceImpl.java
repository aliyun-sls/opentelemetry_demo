package com.example.service;

import com.example.ads.entity.AdEntity;
import com.example.ads.repository.AdRepository;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.List;

public class AdServiceImpl implements AdService{
    @Autowired
    private AdRepository adRepository;
    @Override
    public List<AdEntity> listAds() {
        return adRepository.findAll();
    }
}
