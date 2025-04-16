package com.example.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Service;

@Service
public class ScenarioHandler {
    @Autowired
    private JdbcTemplate jdbcTemplate; // 假设使用 Spring JDBC

    @Autowired
    private AdService adService;

    // 正常处理
    public void handleNormal() {
        adService.listAds();
    }

    // 程序延迟
    public void handleDelay() {
        try {
            Thread.sleep(2000);
            adService.listAds();
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    // SQL 延迟
    public void handleSQLDelay() {
        adService.findAdWithSleep();
    }

    // SQL 复杂查询
    public void handleSQLComplex() {
        adService.excuteComplexSQL();
    }

    // 大量插入操作
    public void handleMassInsert() {
        adService.excuteMassInsert();
    }

    public void handleTableNotExist() {
        adService.findFromInexistentTable();
    }

    // 根据场景名称调用对应方法
    public void handleScenario(String scenarioName) {
        switch (scenarioName) {
            case "normal":
                handleNormal();
                break;
            case "delay":
                handleDelay();
                break;
            case "sql_complex":
                handleSQLComplex();
                break;
            case "mass_insert":
                handleMassInsert();
                break;
            case "table_not_exist":
                handleTableNotExist();
                break;
            default:
                throw new IllegalArgumentException("未知场景: " + scenarioName);
        }
    }
}
