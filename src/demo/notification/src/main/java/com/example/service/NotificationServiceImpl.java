package com.example.service;

import com.example.entity.NotificationEntity;
import com.example.repository.NotificationRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class NotificationServiceImpl implements NotificationService {
    @Autowired
    private NotificationRepository notificationRepository;
    @Override
    public List<NotificationEntity> listNotification() {
        return notificationRepository.findAll();
    }

    @Override
    public NotificationEntity getNotificationById(Long id) {
        return notificationRepository.getOne(String.valueOf(id));
    }
}
