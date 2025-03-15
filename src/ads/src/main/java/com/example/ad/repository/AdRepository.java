package com.example.ad.repository;

import com.example.ad.entity.Ads;
import org.springframework.data.jpa.repository.JpaRepository;

public interface AdRepository extends JpaRepository<Ads, String> {
}
