package com.example.controller;

import com.example.entity.AdEntity;
import com.example.entity.NotificationEntity;
import com.example.service.AdsService;
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
    @Autowired
    private AdsService adsService;

    @GetMapping("/listNotification")
    public ResponseEntity<List> listNotification() {
//        List<NotificationEntity> notificationEntityList = notificationService.listNotification();
        List<AdEntity> ads = adsService.listAds().block();
        return ResponseEntity.ok(ads);
    }

    // 新增: 根据广告ID查询广告详情
    @GetMapping("/notification/{id}")
    public ResponseEntity<NotificationEntity> getNotificationById(@PathVariable Long id) {
        NotificationEntity notification = notificationService.getNotificationById(id);
        return notification != null ? ResponseEntity.ok(notification) : ResponseEntity.notFound().build();
    }
}