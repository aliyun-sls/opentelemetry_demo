package com.example.service;

import com.example.entity.AdEntity;
import com.example.repository.AdRepository;
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
