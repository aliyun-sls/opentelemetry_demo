package com.example.repository;

import com.example.entity.RecommendationEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface RecommendationRepository extends JpaRepository<RecommendationEntity, String> {
}
