> Kafka自2.8开始，移除了之前用于集群的元数据管理、控制器选举等的ZooKeeper的依赖，转而使用Kraft代替，本文来聊聊这一改动的差异和影响。



[TOC]



# 1. 概述

我们用过Kafka都知道，在安装Kafka之前，需要先安装Java和ZooKeeper。需要Java是因为ZooKeeper和Kafka都是用Java编写的，运行需要Java环境。而需要ZooKeeper则是因为，Kafka是使用ZooKeeper来保存集群的元数据信息和消费者信息。

# 2. Kafka与ZooKeeper的关系

Kafka系统架构

![kafka_architecture](image/kafka_architecture.png)

# 3. ZooKeeper在Kafka中的作用

# 4. Kraft模式

# 5. 改变的原因

# 6. 改变的影响



# 7. 参考

- [《从Paxos到Zookeeper》](https://book.douban.com/subject/26292004/)
- [《Kafka权威指南》](https://book.douban.com/subject/27665114/)

