package com.example.repository;

import com.example.entity.AdEntity;
import org.springframework.data.jpa.repository.JpaRepository;

public interface AdRepository extends JpaRepository<AdEntity, String> {
}
