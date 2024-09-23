package com.alibaba.arms.mock.server.service.kafka;

import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.CommonClientConfigs;
import org.apache.kafka.clients.producer.KafkaProducer;
import org.apache.kafka.clients.producer.ProducerConfig;
import org.apache.kafka.common.config.SaslConfigs;
import org.apache.kafka.common.config.SslConfigs;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.Profile;

import java.io.File;
import java.io.FileOutputStream;
import java.io.IOException;
import java.util.Properties;

/**
 * @author zmh405877@alibaba-inc.com
 * @date 2023/11/16
 */
//@Configuration
@Profile("arms-dev")
@Slf4j
public class ArmsDevKafkaConfiguration {

    @Value("${spring.kafka.bootstrap-servers}")
    private String bootstrapServers;

    @Value("${spring.kafka.consumer.properties.sasl.mechanism}")
    private String saslMechanismProp;

    @Value("${kafka.sasl.username}")
    private String saslUserName;

    @Value("${kafka.sasl.password}")
    private String saslPassword;

    @Value("${spring.kafka.ssl.trust-store-location}")
    private String sslTruststoreLocation;

    @Value("${spring.kafka.ssl.trust-store-password}")
    private String sslTruststorePassword;

    @Value("${spring.kafka.consumer.security.protocol}")
    private String securityProtocol;

    private static final String FILE_NAME = "/home/admin/sanmu-aiops-demo-server/kafka-conf/kafka_client_jaas.conf";

    @Bean(destroyMethod = "close")
    public KafkaProducer getKafkaProducer() {
        createJaasFile();
        Properties props = new Properties();
        //设置接入点，请通过控制台获取对应Topic的接入点
        props.put(ProducerConfig.BOOTSTRAP_SERVERS_CONFIG, bootstrapServers);
        String location = sslTruststoreLocation.contains("file:") ?
                sslTruststoreLocation.substring("file:".length()) : sslTruststoreLocation;
        //与sasl路径类似，该文件也不能被打包到jar中
        props.put(SslConfigs.SSL_TRUSTSTORE_LOCATION_CONFIG, location);
        //根证书store的密码，保持不变
        props.put(SslConfigs.SSL_TRUSTSTORE_PASSWORD_CONFIG, sslTruststorePassword);
        //接入协议，目前支持使用SASL_SSL协议接入
        props.put(CommonClientConfigs.SECURITY_PROTOCOL_CONFIG, securityProtocol);

        // 设置SASL账号
        String saslMechanism = saslMechanismProp;
        String username = saslUserName;
        String password = saslPassword;
        if (!isEmpty(username) && !isEmpty(password)) {
            String prefix = "org.apache.kafka.common.security.scram.ScramLoginModule";
            if ("PLAIN".equalsIgnoreCase(saslMechanism)) {
                prefix = "org.apache.kafka.common.security.plain.PlainLoginModule";
            }
            String jaasConfig = String.format("%s required username=\"%s\" password=\"%s\";", prefix, username, password);
            props.put(SaslConfigs.SASL_JAAS_CONFIG, jaasConfig);
        }

        //SASL鉴权方式，保持不变
        props.put(SaslConfigs.SASL_MECHANISM, saslMechanism);
        //Kafka消息的序列化方式
        props.put(ProducerConfig.KEY_SERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.StringSerializer");
        props.put(ProducerConfig.VALUE_SERIALIZER_CLASS_CONFIG, "org.apache.kafka.common.serialization.StringSerializer");
        //请求的最长等待时间
        props.put(ProducerConfig.MAX_BLOCK_MS_CONFIG, 30 * 1000);
        //设置客户端内部重试次数
        props.put(ProducerConfig.RETRIES_CONFIG, 5);
        //设置客户端内部重试间隔
        props.put(ProducerConfig.RECONNECT_BACKOFF_MS_CONFIG, 3000);

        //hostname校验改成空
        props.put(SslConfigs.SSL_ENDPOINT_IDENTIFICATION_ALGORITHM_CONFIG, "");

        // 构造Producer对象，注意，该对象是线程安全的，一般来说，一个进程内一个Producer对象即可；
        // 如果想提高性能，可以多构造几个对象，但不要太多，最好不要超过5个
        KafkaProducer<String, String> producer = new KafkaProducer<String, String>(props);
        return producer;
    }

    private void createJaasFile() {
        String content = "KafkaClient {\n" +
                "  org.apache.kafka.common.security.plain.PlainLoginModule required\n" +
                "  username=\"" + saslUserName + "\"\n" +
                "  password=\"" + saslPassword + "\";\n" +
                "};\n";
        try {
            File file = new File(FILE_NAME);
            File parentDir = file.getParentFile();
            if (!parentDir.exists()) {
                boolean created = parentDir.mkdirs();
                if (!created) {
                    throw new RuntimeException("Cannot create kafka conf dir.");
                }
                log.info("Create kafka conf dir success.");
            }

            boolean created = file.createNewFile();
            if (!created) {
                throw new RuntimeException("Cannot create kafka conf file.");
            }
            FileOutputStream fos = new FileOutputStream(file);
            fos.write(content.getBytes());
            fos.close();
            log.info("Create kafka conf file success.");
        } catch (IOException e) {
            log.error("Exception was threw wile create kafka jaas file", e);
        }
    }

    private boolean isEmpty(String str) {
        if (null == str) {
            return true;
        }
        if (0 == str.trim().length()) {
            return true;
        }

        return false;
    }
}
