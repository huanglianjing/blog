# 1. 介绍

Kafka是LinkedIn采用Scala开发的一个多分区、多副本、基于ZooKeeper协调的分布式消息系统，已被捐献给Apache基金会。Kafka定位为一个分布式流式处理平台，包含高吞吐、可持久化、可水平扩展、支持流数据处理等特性。Kafka有三大角色：

- 消息系统：Kafka具备系统解耦、冗余存储、流量削峰、缓冲、异步通信、扩展性、可恢复性等功能。
- 存储系统：Kafka把消息持久化到磁盘，并且有多副本机制，相比内存存储的系统降低了数据丢失的风险。
- 流式处理平台：Kafka为流式处理框架提供了可靠的数据来源。



# 2. 基本概念

## 2.1 体系架构

![kafka_architecture](image/kafka_architecture.png)

上图为Kafka的体系架构。一个Kafka体系架构包含若干Producer、若干Broker、若干Consumer，以及一个ZooKeeper集群。

- ZooKeeper：负责集群元数据的管理、控制器的选举。
- 生产者Producer：将消息发送到Broker。
- 消费者Consumer：从Broker订阅主题并消费消息。
- 服务代理节点Broker：将收到的消息存储到磁盘。Broker可以看作一个Kafka服务节点或Kafka服务实例，可以将多个Broker运行在不同的服务器上，也可以运行在同一个服务器但是配置不同的端口。

## 2.2 主题与分区

#### 主题Topic

Kafka的消息以主题进行归类，这是一个逻辑上的概念。生产者将消息发送到特定主题，而消费则订阅主题并进行消费。

#### 分区Partition

一个主题可以分为多个分区，而一个分区只属于某个主题。同一个主题下，不同分区包含的消息是不同的，分区在存储层面可以看做一个可追加的日志文件，消息追加到分区日志文件后会分配一个特定的偏移量offset。offset是消息在分区的唯一标识，offset可以保证消息在分区内的顺序性，每个分区有分别的offset，因此Kafka保证的是分区有序而不是主题有序。

一个主题的各个分区可以分布在不同的broker上，每条消息被发送到broker前，会根据分区规则选择存储到哪个具体的分区。因此在创建主题的时候通过参数指定分区的个数，或者创建完成后修改分区数量，以实现分区水平扩展，突破机器IO的性能瓶颈。

#### 副本Replica

分区的多副本机制可以增加副本数量以提升容灾能力，统一分区的不同副本保存相同的消息，副本之间为一主多从的关系，其中leader副本负责处理读写请求，follower副本之负责与leader副本的消息同步。副本处于不同的broker重，当leader副本出现故障时，便从follower副本重新选举新的leader副本对外提供服务。

下图展示了在一个有4个broker的集群中，一个主题配置了3个分区P1、P2、P3。副本因子为3，即每个分区包含有3个副本，一个是leader副本两个是follower副本，这3个副本存储在不同的broker中。

![kafka_partition_replica](image/kafka_partition_replica.png)

分区的所有副本统称为AR（Assigned Replicas），leader副本以及与leader副本保持一定程度同步的副本组成ISR（In-Sync Replicas），与leader副本同步滞后过多的副本组成OSR（Out-of-Sync Replicas），因此AR=ISR+OSR。leader副本维护和跟踪ISR集合与自己的滞后状态，滞后太多或失效的副本会被从ISR剔除，OSR中的副本追上leader副本则会被转移至ISR。默认配置下，leader副本发生故障后只会在ISR集合中选举新leader。

#### 偏移量offset

LogStartOffset为0，是日志文件的起始处，也就是第一条消息。

HW（High Watermark）俗称高水位，标识了一个offset，消费者只能拉取到HW之前的消息。

LEO（Log End Offset）标识当前日志文件下一条待写入消息的offset，相当于当前日志分区最后一条消息的offset加1。

分区的ISR集合每个副本都会维护自身的LEO，而ISR集合中最小的LEO就是分区的HW。

下图展示了一个日志文件，其中HW为6，LEO为9，因此消费者只能拉取偏移量0到5的消息，而下一条写入的消息偏移量将会是9。

![kafka_offset](image/kafka_offset.png)



# 3. 生产者



# 4. 消费者



# 参考

- [《深入理解Kafka：核心设计与实践原理》](https://book.douban.com/subject/30437872/)

