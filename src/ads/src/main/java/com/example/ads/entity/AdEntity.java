package com.example.ads.entity;

import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.Entity;
import org.springframework.data.annotation.Id;

import javax.persistence.Table;

@Entity
@Getter
@Setter
@Table(name = "ads")
public class AdEntity {
    @Id
    private String id;
    
    private String redirectUrl;
    private String text;
}
