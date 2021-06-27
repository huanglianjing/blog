> 本文介绍Kafka在Linux与Mac下的安装步骤，配置介绍，以及自带脚本工具的使用。



# 1. 安装

Kafka运行环境需要先安装好Java环境。

进入官网http://kafka.apache.org/downloads，选择相应的版本的Kafka链接并下载：

```bash
$ wget https://downloads.apache.org/kafka/2.8.0/kafka_2.13-2.8.0.tgz
```

解压安装包

```bash
$ tar zxf kafka_2.13-2.8.0.tgz -C /usr/local/
$ ln -s /usr/local/kafka_2.13-2.8.0/ /usr/local/kafka
$ cd /usr/local/kafka
```

启动ZooKeeper

```bash
$ ./bin/zookeeper-server-start.sh config/zookeeper.properties

# 后台运行
$ nohup ./bin/zookeeper-server-start.sh config/zookeeper.properties >> zookeeper.log 2>&1 &
```

启动Kafka

```bash
$ ./bin/kafka-server-start.sh config/server.properties

# 后台运行
$ nohup ./bin/kafka-server-start.sh config/server.properties >> kafka.log 2>&1 &
```



# 2. 配置系统服务单元

这一步是可选的，配置了之后通过systemctl命令启动和停止，也可以直接执行脚本来启动停止。

## 2.1 Zookeeper

创建系统服务单元

```bash
$ cd /etc/systemd/system
$ vi zookeeper.service
```

贴上以下内容

```
[Unit]
Description=Apache Zookeeper server
Documentation=http://zookeeper.apache.org
Requires=network.target remote-fs.target
After=network.target remote-fs.target

[Service]
Type=simple
ExecStart=/usr/local/kafka/bin/zookeeper-server-start.sh /usr/local/kafka/config/zookeeper.properties
ExecStop=/usr/local/kafka/bin/zookeeper-server-stop.sh
Restart=on-abnormal
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

操作命令

```bash
# 启动ZooKeeper
$ systemctl start zookeeper

# 查看ZooKeeper状态
$ systemctl status zookeeper

# 关闭ZooKeeper
$ systemctl stop zookeeper
```

## 2.2 Kafka

创建系统服务单元

```bash
$ cd /etc/systemd/system
$ vi kafka.service
```

贴上以下内容

```
[Unit]
Description=Apache Kafka Server
Documentation=http://kafka.apache.org/documentation.html
Requires=zookeeper.service

[Service]
Type=simple
ExecStart=/usr/local/kafka/bin/kafka-server-start.sh /usr/local/kafka/config/server.properties
ExecStop=/usr/local/kafka/bin/kafka-server-stop.sh
Restart=on-abnormal

[Install]
WantedBy=multi-user.target
```

操作命令

```bash
# 启动Kafka
$ systemctl start kafka

# 查看Kafka状态
$ systemctl status kafka

# 关闭Kafka
$ systemctl stop kafka
```



# 3. 目录结构

下面进入Kafka的目录，也就是`/usr/local/kafka`，看一下目录的结构。

```
|-- bin                           // Kafka和ZooKeeper的脚本工具
    |-- kafka-console-consumer.sh
    |-- kafka-console-producer.sh
    |-- kafka-server-start.sh
    |-- kafka-server-stop.sh
    |-- kafka-topics.sh
    |-- windows                   // windows下的bat脚本
    |-- zookeeper-server-start.sh
    |-- zookeeper-server-stop.sh
    `-- ...
|-- config                        // Kafka和ZooKeeper的配置文件
    |-- kraft                     // Kafka2.8开始移除ZooKeeper依赖的新启动配置，本文暂不介绍
    |-- server.properties
    |-- zookeeper.properties
    `-- ...
|-- libs                          // 一些依赖的jar包
|-- LICENSE
|-- licenses
|-- logs                          // 日志
|-- NOTICE
`-- site-docs                     // 文档
```



# 4. 脚本工具

Kafka提供了很多脚本工具，可以用来进行主题创建和查看、生产者、消费者等操作。

以下脚本执行需要先进入Kafka目录进行操作，脚本工具都在bin目录下。

## 4.1 kafka-topics.sh

与主题相关的脚本，用于查看主题、创建主题。

```bash
# 查看已创建的主题
$ ./bin/kafka-topics.sh --list --zookeeper localhost:2181

# 创建主题
$ ./bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
```

## 4.2 kafka-console-producer.sh

生产者脚本。

```bash
# 通过生产者发送消息，在终端输入然后回车发送消息
$ ./bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test
```

## 4.3 kafka-console-consumer.sh

消费者脚本。

```bash
# 通过消费者接收消息
$ ./bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test

# 通过消费者接收消息，从头开始
$ ./bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning
```



# 5. 配置



# 6. 参考

- [《深入理解Kafka：核心设计与实践原理》](https://book.douban.com/subject/30437872/)

