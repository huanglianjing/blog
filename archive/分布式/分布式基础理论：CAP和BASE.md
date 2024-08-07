# 1. 概述

当单一的本地程序由于硬件资源无法满足业务需求时，将其分散到可扩容资源的多个节点来解决问题，就是分布式要干的事情。它研究将一个需要非常大计算能力解决的问题分成很多小部分，分配给多个计算机进行处理，再把结果综合起来得到最终的结果。

分布式中很重要的一个概念就是事务，事务提供一种机制将一个活动涉及的所有操作纳入到一个不可分割的执行单元，组成事务的所有操作只有在所有操作均能正常执行的情况下方能提交，只要其中任一操作执行失败，都将导致整个事务的回滚。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_transaction.jpg)

事务基于数据进行操作，数据库事务支持 ACID 特性，它们分别为原子性（Atomicity）、一致性（Consistency）、隔离性（Isolation）、持久性（Durability）。

* 原子性指整个事务中的所有操作，要么全部完成，要么全部不完成，不会出现部分完成的情况
* 一致性指事物开始前和结束后，数据库数据的一致性约束没有被破坏，如两个账户金额之和是 100 元，发生转账之后还是 100 元
* 隔离性指事务进行期间的中间状态对其它事务是不可见的，一个事务正在修改某个数据，只要还未提交事务，其它事务对其读取就无法看到改动
* 持久性指事务对数据的修改是永久的，即使系统故障也不会丢失

而对于分布式事务来说，是无法实现 ACID 的，因为它受 CAP 理论的约束。

# 2. CAP理论

在设计一个大规模可扩展的网络服务时，会考虑到三个特性：一致性（Consistency）、可用性（Availability）、分区容错性（Partition Tolerance）。

CAP 理论指在一给分布式计算机系统中，一致性、可用性和分区容错性这三者无法同时得到满足，最多只能同时满足两个。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap.jpg)

下面通过一个场景来分析 CAP 中的每个特性。客户端发送请求，业务层处理请求，完成业务逻辑处理后存储数据，存储层内部分为主存储和从存储且存在数据同步，业务层从存储层去除数据，然后返回给客户端。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_case.png)

## 2.1 一致性

一致性指的是一旦数据完成更新并成功返回，那么分布式系统中的所有节点在同一时间数据完全一致。当数据分布在多个节点上时，写操作成功后，从任意节点读取到的数据都是最新的状态。

分布式系统中的一致性，由于存在数据同步的过程，写操作的响应会有一定的延迟。为了保证数据的一致性，会对资源暂时锁定，数据同步完成后再释放资源。

一致性分为三种级别：

* 强一致性：系统在某个节点中写入或修改了数据，之后在任意节点读取到的数据都必须是最新的数据，常用于银行交易系统
* 弱一致性：不一定能读取到最新的值，也不保证在一定时间后读取到最新值，只会尽量在某个时刻达到数据一致的状态
* 最终一致性：相较于弱一致性，保证在一定的时间内达到数据的最终一致性，常用于互联网应用服务

一致性的目标是将数据写入主存储，然后从从存储读取也成功，如果写入失败则从从存储读取也失败。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_c_1.png)

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_c_2.png)

为了实现一致性，在写入主存储后，需要锁定数据库以向从存储同步数据，待同步完成后再释放锁，避免写入新数据后读取到旧的数据。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_c_3.png)

## 2.2 可用性

可用性指服务一支可用，且可以在正常时间内响应。

可用性根据衡量标准分为不同等级，越是高的可用性一年可容忍的停机时间越少。

| 可用性                       | 可用水平 | 一年可容忍停机时间 |
| ---------------------------- | -------- | ------------------ |
| 商品可用性                   | 99%      | < 87.6 h           |
| 高可用性                     | 99.9%    | < 8.8 h            |
| 具有故障自动恢复能力的可用性 | 99.99%   | < 53 min           |
| 极高可用性                   | 99.999%  | < 5 min            |
| 容错可用性                   | 99.9999% | < 32 s             |

可用性的目标是，当主存储被更新是，如果从存储收到数据查询请求，应当立即能够响应数据查询结果，不允许出现响应超时或响应错误。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_a_1.png)

为了实现可用性，不可以将从存储中的资源锁定，即使数据还没有完成同步，从存储也要返回查询的数据而非返货错误或超时，即使它是旧数据。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_a_2.png)

## 2.3 分区容错性

分区容错性是指在分布式系统中，即使部分节点出现了故障或消息丢失，系统还可以继续运行。分区容错性是一个分布式系统的应当具备的基本能力。

分区容错性的目标是主存储向从存储同步数据失败时，不影响读写操作，其中部分节点出现故障不影响其他节点对外提供服务。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_p_1.png)

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_p_2.png)

为了实现分区容错性，可以使用异步同步的方式，将数据从主存储同步到从存储，从而在节点之间实现松耦合。添加备份的从存储节点，在一个节点挂掉后让备份的节点提供服务。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_p_3.png)

## 2.4 矛盾与取舍

对于没有发生分区的单一程序系统，就不存在分区容错性的必要了，这时候可以兼顾一致性和可用性。而对于一个分布式系统来说，分区容错性是必须要满足的基本特性，这时候一致性和可用性二者无法同时满足，必须选择其中一个。

以以下场景为例，一个分布式系统系统中有两个节点 Host1 和 Host2，它们之间通过网络联通，Host1 运行程序 Process1 和数据库 Data，Host2 运行程序 Process2 和数据库 Data。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_conflict_1.png)

满足一致性要求 Host1 和 Host2 中的 Data 的数据相同，满足可用性要求不管用户请求 Host1 还是 Host2 都会立即响应结果，满足分区容错性要求 Host1 或者 Host2 任一方脱离系统都不会影响分布式系统的正常运作。

以下是一个正常的运行流程。首先向 Host1 请求更新数据，将 Data 中的数据从 0 更新为 1。通过数据同步，Host1 中的 Data 将会同步到 Host2 中的 Data，使其更新为 1。然后向 Host2 请求读取数据，获取最新的数据为 1。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_conflict_2.png)

假如这时候发生了网络异常，Host1 和 Host2 之间的网络断开了。首先向 Host1 请求更新数据，将 Data 中的数据从 0 更新为 1。由于网络问题，Host2 中的数据无法同步至最新。用户向 Host2 发送读取请求是，无法立即将最新的数据返回给用户。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/distribute/distribute_theory_cap_conflict_3.png)

为了满足一致性，则需要阻塞等待，直到网络恢复连接然后数据同步完成，再将最新数据返回给用户，这牺牲了可用性。为了满足可用性，立即将当前的旧数据返回给用户，这又牺牲了一致性。因此这种情况下一致性和可用性无法同时得到满足。

既然一致性、可用性、分区容错性无法同时满足，就需要在它们之间作出取舍。

放弃分区容错性：不进行分区，从而不用考虑网络不通或节点故障的问题。常见的关系性数据库就是满足 CA 而放弃 P。对于分布式系统，不可能存在不满足 P，P 是基本要求，并且要想尽办法提升 P。

放弃可用性：不要求强的可用性，允许系统停机或者长时间无响应。常见的有 Redis、HBase、ZooKeeper 都是选择有限保证 CP 而放弃 A。典型的应用场景有金融行业如跨行转账，一次转账请求需要等待双方银行系统完成才算整个事务完成。

放弃一致性：接受查询到的数据在一定时间内不是最新的。很多分布式系统设计都会选择这个，放弃强一致性，退而求其次保证最终一致性。典型的应用场景有淘宝订单退款成功后一定时间才会到账，用户接受到账时间稍晚于退款时间，还有在 12306 买票或网上抢演唱会门票时，显示还有余票但是有时候点击进去却是无票了，就是数据的短暂不一致造成的。

# 3. BASE理论

对于分布式系统而言，CAP 理论不可能同时满足，而分区容错性优势不可或缺的，为了权衡一致性和可用性，于是出现了 BASE 理论，它包含基本可用（Basically Available）、软状态（Soft State）、最终一致性（Eventually Consistent）。

BASE 理论是对大规模互联网系统分布式实践的总结。其核心思想是，即使无法做到强一致性，每个应用都可以根据自身的业务特点，采用适当的方法来使系统达到最终一致性，通过牺牲强一致性，从而获得高可用性。

基本可用是对响应时间的妥协和对功能损失的妥协。一个请求在正常情况下需要在 500ms 内返回响应结果，由于出现故障，实际响应时间可能增加到 1-2s。一个电子商务网站在节日大促时，浏览和购买行为激增，为了保护系统的稳定性，部分请求可能会被引导到一个降级页面或失败页面。

软状态允许系统中的数据存在中间状态，即允许系统的多个不同节点的数据副本存在数据延迟。相对应地，要多多个节点的数据副本保持一致是一种硬状态。

最终一致性保证在一个时间期限过后，所有副本将会达到数据的一致性，对所有副本数据的访问都将或渠道最新值。这个时间期限取决于网络时延、系统负载、数据复制方案等因素。

# 4. 参考

* [分布式从ACID、CAP、BASE的理论推进](https://www.yuque.com/aceld/golang/ycp0nb)

