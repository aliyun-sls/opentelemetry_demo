package com.example.repository;

import com.example.entity.MarketingEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;


public interface MarketingRepository extends JpaRepository<MarketingEntity, String> {
    @Query(value = "SELECT *, SLEEP(5) AS delay FROM marketing LIMIT 4", nativeQuery = true)
    List<MarketingEntity> findWithDelay();
}
