package com.alibaba.arms.mock.server.controller;

import com.alibaba.arms.mock.api.Invocation;
import com.alibaba.arms.mock.server.service.ComponentService;
import com.alibaba.arms.mock.server.service.MysqlService;
import com.alibaba.arms.mock.server.service.RedisService;
import com.alibaba.arms.mock.server.service.kafka.KafkaService;
import com.alibaba.arms.mock.server.service.trouble.*;
import io.swagger.annotations.Api;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/components/api/v1")
@Api("ads")
public class AdController {

    private static final Logger log = LoggerFactory.getLogger(AdController.class);
    @Autowired
    private ComponentService componentService;

    @Autowired
    private MysqlService mysqlService;

    @Autowired
    private RedisService redisService;

    @Autowired
    private KafkaService kafkaService;

    @Autowired
    private TroubleManager troubleManager;

    @RequestMapping(value = "/ads/data", method = RequestMethod.POST)
    public String getData(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        return componentService.execute(children, level);
    }

    @RequestMapping(value = "/ads/cache", method = RequestMethod.POST)
    public String getCache(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        redisService.execute();
        return componentService.execute(children, level);
    }

    @RequestMapping(value = "/ads/db", method = RequestMethod.POST)
    public String getDb(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        mysqlService.execute();

        GCTrouble trouble = (GCTrouble) troubleManager.getTroubleMaker(TroubleEnum.GC.getCode());
        if (trouble.isWorking()) {
            GarbageCollectionTrigger garbageCollectionTrigger = new GarbageCollectionTrigger();

            garbageCollectionTrigger.doExecute();
        }

        return componentService.execute(children, level);
    }

    @RequestMapping(value = "/ads/message", method = RequestMethod.POST)
    public String sendMessage(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        kafkaService.execute();
        return componentService.execute(children, level);
    }

}


