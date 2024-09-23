package com.alibaba.arms.mock.server.controller;

import com.alibaba.arms.mock.api.Invocation;
import com.alibaba.arms.mock.server.service.ComponentService;
import com.alibaba.arms.mock.server.service.MysqlService;
import com.alibaba.arms.mock.server.service.RedisService;
import com.alibaba.arms.mock.server.service.kafka.KafkaService;
import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/components/api/v1")
@Api("mall")
public class MallController {

    @Autowired
    private ComponentService componentService;

    @Autowired
    private MysqlService mysqlService;

    @Autowired
    private RedisService redisService;

    @Autowired
    private KafkaService kafkaService;

    @RequestMapping(value = "/mall/product", method = RequestMethod.POST)
    public String productEntry(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        return componentService.execute(children, level);
    }
    @RequestMapping(value = "/mall/user_info", method = RequestMethod.POST)
    public String userInfo(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        mysqlService.execute();
        return componentService.execute(children, level);
    }
    @RequestMapping(value = "/mall/user_cart", method = RequestMethod.POST)
    public String userCart(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        redisService.execute();
        return componentService.execute(children, level);
    }

    @RequestMapping(value = "/mall/sku_info", method = RequestMethod.POST)
    public String spuInfo(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        mysqlService.execute();
        return componentService.execute(children, level);
    }
    @RequestMapping(value = "/mall/spu_info", method = RequestMethod.POST)
    public String skuInfo(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        redisService.execute();
        return componentService.execute(children, level);
    }
    @RequestMapping(value = "/mall/save_order_info", method = RequestMethod.POST)
    public String saveOrderInfo(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        kafkaService.execute();
        return componentService.execute(children, level);
    }
}


