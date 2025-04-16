   package com.example.entity;

   import lombok.Getter;
   import lombok.Setter;
   import javax.persistence.Entity;
   import javax.persistence.Id;
   import javax.persistence.Table;

   @Entity
   @Getter
   @Setter
   @Table(name = "marketing")
   public class MarketingEntity {
       @Id
       private Long id;
       
       private String redirectUrl;
       private String text;
   }
   