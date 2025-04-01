package com.example.service;

import com.example.entity.AdEntity;

import java.util.List;

public interface AdService {
    List<AdEntity> listAds();

    AdEntity getAdById(Long id);

    List findAdWithSleep();

    List findFromInexistentTable();

    void excuteMassInsert();

    void excuteComplexSQL();
}
