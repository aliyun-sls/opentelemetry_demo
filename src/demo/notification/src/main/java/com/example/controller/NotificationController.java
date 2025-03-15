package com.example.controller;

import com.example.entity.NotificationEntity;
import com.example.service.NotificationService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;

import java.util.List;

@Controller
public class NotificationController {

    @Autowired
    private NotificationService notificationService;

    @GetMapping("/listAds")
    public ResponseEntity<List<NotificationEntity>> listAds() {
        List<NotificationEntity> ads = notificationService.listNotification();
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/ad/{id}")
    public ResponseEntity<NotificationEntity> getAdById(@PathVariable Long id) {
        NotificationEntity ad = notificationService.getNotificationById(id);
        return ad != null ? ResponseEntity.ok(ad) : ResponseEntity.notFound().build();
    }
}