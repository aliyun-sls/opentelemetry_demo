   package com.example.entity;

   import lombok.Getter;
   import lombok.Setter;

   import javax.persistence.*;

   @Entity
   @Getter
   @Setter
   @Table(name = "ads")
   public class AdEntity {
       @Id
       @GeneratedValue(strategy = GenerationType.IDENTITY)
       private Long id;
       
       private String redirectUrl;
       private String text;
   }
   