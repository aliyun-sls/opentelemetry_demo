package com.example.service;

import com.example.entity.NotificationEntity;

import java.util.List;

public interface NotificationService {
    List<NotificationEntity> listNotification();

    NotificationEntity getNotificationById(Long id);
}
