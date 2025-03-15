package com.example.ad.entity;

import lombok.Getter;
import lombok.Setter;
import org.hibernate.annotations.Entity;
import org.springframework.data.annotation.Id;

import javax.persistence.Table;

@Entity
@Getter
@Setter
@Table(name = "ads")
public class Ads {
    @Id
    private String id;
    
    private String redirectUrl;
    private String text;
}
