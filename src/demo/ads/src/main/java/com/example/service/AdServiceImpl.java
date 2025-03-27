package com.example.service;

import com.example.entity.AdEntity;
import com.example.repository.AdRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.persistence.EntityManager;
import javax.persistence.Query;
import java.util.List;

@Service
public class AdServiceImpl implements AdService{
    @Autowired
    private AdRepository adRepository;


    @Autowired
    private EntityManager entityManager;
    @Override
    public List<AdEntity> listAds() {
        return adRepository.findAll();
    }

    @Override
    public AdEntity getAdById(Long id) {
        return adRepository.getOne(String.valueOf(id));
    }

    @Override
    public List findAdWithSleep() {
        String sql = "SELECT * FROM AdEntity WHERE SLEEP(10)";
        Query query = entityManager.createNativeQuery(sql, AdEntity.class);
        return query.getResultList();
    }
}
