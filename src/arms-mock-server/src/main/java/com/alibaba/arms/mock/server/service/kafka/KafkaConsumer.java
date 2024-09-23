package com.alibaba.arms.mock.server.service.kafka;

import lombok.extern.slf4j.Slf4j;
import org.springframework.boot.autoconfigure.condition.ConditionalOnProperty;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

/**
 * @author zmh405877@alibaba-inc.com
 * @date 2023/11/17
 */
@Component
@ConditionalOnProperty(name = "service.name", havingValue = "insights-kafka-consumer")
@Slf4j
public class KafkaConsumer {

    @KafkaListener(topics = "${kafka.topics}", groupId = "spring.kafka.consumer.group-id")
    public void consume(String message) {
        System.out.println("Receive message: " + message);
        log.info("Receive message: " + message);
    }
}
