package com.alibaba.arms.mock.server.controller;

import com.alibaba.arms.mock.api.Invocation;
import com.alibaba.arms.mock.server.service.ComponentService;
import com.alibaba.arms.mock.server.service.MysqlService;
import com.alibaba.arms.mock.server.service.RedisService;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.util.Arrays;

/**
 * @author changan.zca@alibaba-inc.com
 * @date 2021/06/10
 */
@RestController
@RequestMapping("/case/api/v1")
@Slf4j
public class CaseController {
    @Autowired
    private ComponentService componentService;

    @Autowired
    private MysqlService mysqlService;

    @Autowired
    private RedisService redisService;

    @RequestMapping(value = "/gateway/execute", method = RequestMethod.POST)
    public String gatewayExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }
    @RequestMapping(value = "/user_info/execute", method = RequestMethod.POST)
    public String userInfoExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        mysqlService.execute();
        return "success";
        // return componentService.execute(Arrays.asList(invocation), level);
    }
    @RequestMapping(value = "/user_cart/execute", method = RequestMethod.POST)
    public String userCartExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        redisService.execute();
        return "success";
        // return componentService.execute(Arrays.asList(invocation), level);
    }
    @RequestMapping(value = "/product/execute", method = RequestMethod.POST)
    public String productExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
        // return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/http/execute", method = RequestMethod.POST)
    public String httpExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/mysql/execute", method = RequestMethod.POST)
    public String mysqlExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/redis/execute", method = RequestMethod.POST)
    public String redisExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/local/execute", method = RequestMethod.POST)
    public String localExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/bad_sql/execute", method = RequestMethod.POST)
    public String badSqlExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        return componentService.execute(Arrays.asList(invocation), level);
    }

    @RequestMapping(value = "/mall_mysql/execute", method = RequestMethod.POST)
    public String mallMysqlExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        mysqlService.execute();
        return "success";
    }
    @RequestMapping(value = "/mall_redis/execute", method = RequestMethod.POST)
    public String mallRedisExecute(HttpServletRequest httpServletRequest, @RequestBody Invocation invocation, @RequestParam("level") int level) {
        log.info("execute {}", httpServletRequest.getRequestURL());
        redisService.execute();
        return "success";
    }
}
