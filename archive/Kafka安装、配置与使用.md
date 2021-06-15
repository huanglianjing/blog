> 本文介绍Kafka在Linux与Mac下的安装步骤，配置介绍，以及自带脚本工具的使用。



# 1. 安装

Kafka运行环境需要先安装好Java环境。

从官网http://kafka.apache.org/downloads找到Kafka链接并下载

```bash
$ wget https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
```

解压安装包

```bash
$ tar zxf kafka_2.13-2.8.0.tgz -C /opt/
$ ln -s /opt/kafka_2.13-2.8.0/ /opt/kafka
$ cd /opt/kafka
```

启动ZooKeeper

```bash
$ ./bin/zookeeper-server-start.sh config/zookeeper.properties
```

启动Kafka

```bash
$ ./bin/kafka-server-start.sh config/server.properties
```



# 2. 参考

