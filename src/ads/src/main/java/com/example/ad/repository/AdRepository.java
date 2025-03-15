package com.example.ad.repository;

import com.example.ad.entity.AdEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;

import java.util.List;
import java.util.Set;

public interface AdRepository extends JpaRepository<AdEntity, String> {
}
