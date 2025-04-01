package com.example.entity;

import lombok.Data;

@Data
public class Scenario {
    private String name;
    private double weight;

    @Override
    public String toString() {
        return "Scenario{name='" + name + "', weight=" + weight + '}';
    }
}