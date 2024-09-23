## 如何部署


### 编译

mvn clean package

### 依赖的DB和表结构

1. 实例申请(略)

2. 创建用户

```
create user arms_mock identified by 'arms_mock!@#';
grant all on *.* to arms_mock;
flush privileges;
```

3. 创建DB和表

```
create database arms_mock;


CREATE TABLE `dummy_record` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `content` varchar(200) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1642 DEFAULT CHARSET=utf8mb4

```




### 搭建redis


## 如何使用

### 部署服务


#### 1. 创建配置文件

基于**mock-server/src/main/resources/application.properties**复制一份配置文件,然后进行修改.


1. 修改服务列表

![](resources/2021-06-29-11-41-20.png)

2. 修改DB配置

![](resources/2021-06-29-11-40-07.png)

3. 修改redis配置

![](resources/2021-06-29-11-40-50.png)


#### 2. 启动服务

nohup java -javaagent:xxx/arms-bootstrap-1.7.0-SNAPSHOT.jar -Darms.appName=Sanmu-Prom-Demo-0  -jar xxx/mock-server-2.0.5.RELEASE.jar --spring.config.location=file:xxx/sanmu-prom-demo-0-9190.properties 1>${PRGDIR}/logs/sanmu-prom-demo-0.log 2>&1 &


重复上述步骤,直至需要的服务都启动完毕.

### 注入流量

nohup java -jar  mock-client/target/mock-client-1.0-SNAPSHOT-jar-with-dependencies.jar & 



