   package com.example.entity;

   import lombok.Getter;
   import lombok.Setter;

   import javax.persistence.Entity;
   import javax.persistence.Id;
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
   