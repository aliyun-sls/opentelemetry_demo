package com.example.repository;

import com.example.entity.MarketingEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface MarketingRepository extends JpaRepository<MarketingEntity, String> {
}
