package com.example.ad.service;

import com.example.ad.entity.AdEntity;
import com.example.ad.repository.AdRepository;
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
