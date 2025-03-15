package com.example.ad.entity;

import jakarta.persistence.*;
import lombok.Data;

@Data
@Entity
@Table(name = "ads")
public class AdEntity {
    @Id
    private String id;
    
    private String redirectUrl;
    private String text;
}
