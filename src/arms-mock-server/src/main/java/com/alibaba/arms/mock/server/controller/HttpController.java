package com.alibaba.arms.mock.server.controller;

import com.alibaba.arms.mock.api.Invocation;
import com.alibaba.arms.mock.server.service.ComponentService;
import com.alibaba.arms.mock.server.service.HttpService;
import com.alibaba.arms.mock.server.service.MysqlService;
import com.alibaba.arms.mock.server.service.RedisService;
import io.swagger.annotations.Api;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/components/api/v1")
@Api("http")
@Slf4j
public class HttpController {

    @Autowired
    private ComponentService componentService;

    @Autowired
    private HttpService httpService;

    @Autowired
    private RedisService redisService;

    @Autowired
    private MysqlService mysqlService;

    @RequestMapping(value = "/http/success", method = RequestMethod.POST)
    public String success(@RequestBody(required = false) List<Invocation> children,
                          @RequestParam(required = false, value = "level") Integer level) {

        httpService.execute();
        // redis
        redisService.execute();
        // mysql
        mysqlService.execute();
        // httpService.execute();
        // return componentService.execute(children, level);
        return "success";
    }
}


