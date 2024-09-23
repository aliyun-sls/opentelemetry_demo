package com.alibaba.arms.mock.server.service.kafka;

import com.alibaba.arms.mock.server.service.AbstractComponent;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerRecord;
import org.apache.kafka.clients.producer.RecordMetadata;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

/**
 * @author zmh405877@alibaba-inc.com
 * @date 2023/11/16
 */
@Service
@Slf4j
public class KafkaService extends AbstractComponent {

    @Autowired
    private KafkaProducer kafkaProducer;

    @Value("${kafka.topics}")
    private String topics;

    private static final ExecutorService es = Executors.newFixedThreadPool(1000);

    @Override
    public String getComponentName() {
        return "kafka";
    }

    @Override
    public void execute() {
        //解析kafka topics
        String[] topicArray = topics.split(",");
        if (topicArray == null || topicArray.length == 0) {
            log.warn("nothing to do because of empty topic");
            return;
        }
        //构造一个Kafka消息
        String msgFormat = "[%s] [%s] [kafka message] this is a default order"; //消息的内容模板
        try {
            //批量获取 futures 可以加快速度, 但注意，批量不要太大
            List<Future<RecordMetadata>> futures = new ArrayList<Future<RecordMetadata>>(topicArray.length);
            for (String topic : topicArray) {
                //发送消息，并获得一个Future对象
                String msg = String.format(msgFormat, System.currentTimeMillis(), topic);
                ProducerRecord<String, String> kafkaMessage =  new ProducerRecord<String, String>(topic, msg);
                Future<RecordMetadata> metadataFuture = kafkaProducer.send(kafkaMessage);
                futures.add(metadataFuture);
            }
            kafkaProducer.flush();
            for (Future<RecordMetadata> future: futures) {
                //同步获得Future对象的结果
                try {
                    RecordMetadata recordMetadata = future.get();
                    log.info("Produce ok:" + recordMetadata.toString());
                } catch (Throwable t) {
                    log.error("kafka get future result error", t);
                    System.err.println("kafka get future result error: " + t.getMessage());
                    t.printStackTrace();
                }
            }
        } catch (Exception e) {
            //客户端内部重试之后，仍然发送失败，业务要应对此类错误
            //参考常见报错: https://help.aliyun.com/document_detail/68168.html?spm=a2c4g.11186623.6.567.2OMgCB
            log.error("kafka send error", e);
            System.err.println("kafka send error: " + e.getMessage());
            e.printStackTrace();
        }
    }

    @Override
    public Class getImplClass() {
            return KafkaService.class;
        }
}
