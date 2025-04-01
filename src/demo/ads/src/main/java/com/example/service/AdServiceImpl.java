package com.example.service;

import com.example.entity.AdEntity;
import com.example.repository.AdRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.persistence.EntityManager;
import javax.persistence.Query;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

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
        String sql = "SELECT * FROM ads WHERE SLEEP(10)";
        Query query = entityManager.createNativeQuery(sql, AdEntity.class);
        return query.getResultList();
    }

    @Override
    public List findFromInexistentTable() {
        String sql = "SELECT * FROM adEntity";
        Query query = entityManager.createNativeQuery(sql, AdEntity.class);
        return query.getResultList();
    }

    @Override
    public void excuteMassInsert() {
        List<AdEntity> adEntities = new ArrayList<>();
        for (int i = 0; i < 10000; i++) {
            AdEntity adEntity = new AdEntity();
            String randomSuffix = UUID.randomUUID().toString();
            adEntity.setText("testText-" + i);
            adEntity.setRedirectUrl("testUrl-" + i + randomSuffix);
            adEntities.add(adEntity);
        }
        List<AdEntity> res = adRepository.saveAll(adEntities);
        adRepository.deleteAll(res);
    }

    @Override
    public void excuteComplexSQL() {
        String sql = "SELECT EXP(LOG(POW(SQRT(ROUND(LENGTH(a1.text) + LENGTH(a2.text) + RAND() * 1000, 8)), ROUND(POW(RAND(), 3), 3)))) * COS(RADIANS(360 * RAND())) + SIN(RADIANS(180 * RAND())) + a3.id * a4.id / (RAND() + 0.1) AS complex_calculation FROM ads a1 CROSS JOIN ads a2 CROSS JOIN ads a3 CROSS JOIN ads a4 CROSS JOIN ads a5 WHERE a1.id < 35 AND a2.id < 35 AND a3.id < 35 AND a4.id < 35 AND a5.id < 35 LIMIT 5000000";
        Query query = entityManager.createNativeQuery(sql, String.class);
        query.getResultList();
    }
}
