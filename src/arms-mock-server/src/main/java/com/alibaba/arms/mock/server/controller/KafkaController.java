package com.alibaba.arms.mock.server.controller;

import com.alibaba.arms.mock.api.Invocation;
import com.alibaba.arms.mock.server.service.ComponentService;
import com.alibaba.arms.mock.server.service.kafka.KafkaService;
import io.swagger.annotations.Api;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/components/api/v1")
@Api("kafka")
public class KafkaController {


    @Autowired
    private ComponentService componentService;

    @Autowired
    private KafkaService kafkaService;

    @RequestMapping(value = "/kafka/success", method = RequestMethod.POST)
    public String success(@RequestBody(required = false) List<Invocation> children, @RequestParam(required = false, value = "level") Integer level) {
        kafkaService.execute();
        return componentService.execute(children, level);
    }
}
