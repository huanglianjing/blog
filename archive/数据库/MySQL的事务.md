# 1. 事务

事务（transaction）是一组数据库操作的集合，把所有的操作作为一个整体一起向系统提交或撤销操作请求，要么全部成功，要么全部失败。

MySQL 中的事务是由存储引擎实现的，InnoDB 支持事务，而 MyISAM 和 Memory 不支持事务。

事务有 ACID 四大特性：

* 原子性（Atomicity）：事务的所有操作要么全部成功，要么全部失败，通过 undo log 实现；
* 一致性（Consistency）：事务执行前后必须保持一致状态，所有数据必须满足预定义的规则（约束、触发器等），通过主键、外键、唯一行、检查约束等保证；
* 隔离性（Isolation）：多个事务的执行互相隔离，不会互相干扰，通过锁机制或多版本并发控制（MVCC）实现；
* 持久性（Durability）：事务一旦提交，对数据库的修改就会永久生效，通过 redo log 实现；

## 1.1 事务启动

MySQL 支持显式地启动和提交事务。通过 begin 或 start transaction 启动事务，但并非立即启动事务，直到执行第一个操作语句时才是真正的启动事务。通过 start transaction with consistent snapshot 立即启动事务并立即创建一个一致性视图，通过 commit 提交事务，通过 rollback 回滚事务。

通过 set autocommit = 0，可以关闭该线程的自动提交，只要执行一个 select 语句，事务就会启动，且不会自动提交，事务持续到主动执行 commit 或 rollback 语句或断开连接为止。在这种设置下，长连接将会导致长事务。

建议总是使用 set autocommit = 1，通过显式语句来启动事务。在这种模式下，未声明事务的每个语句都会自动创建事务并执行，然后提交事务。commit 会提交事务，而 commit work and chain 则会提交事务并自动启动下一个事务。

## 1.2 长事务

长事务会在系统中保存很老的事务视图，造成回滚日志占用大量空间而迟迟无法释放，占用大量存储空间。

长事务还会占用锁资源，可能会拖垮整个库。

以下语句可以查询持续时间超过 60s 的事务：

```mysql
select * from information_schema.innodb_trx where TIME_TO_SEC(timediff(now(),trx_started)) > 60;
```

长事务应当尽量拆分为多个小事务，每个小事务执行一部份操作，减少单个事务对数据库的影响和锁的持有时间。

# 2. 隔离级别

当数据库有多个事务同时执行时，可能会出现脏读（dirty read）、不可重复读（non-repeatable read）、幻读（phantom read）问题。

对于事务的隔离性，有几种隔离级别：读未提交（read uncommitted）、读提交（read committed）、可重复读（repeatable read）和串行化（serializable ）。隔离级别越高，多个事务间执行的影响越小，从而执行效率越低。

* 读未提交：一个事务还没提交，它的变更就能被其他事务看到；
* 读提交：一个事务提交之后，它的变更才会被其他事务看到；
* 可重复读：一个事务执行过程中看到的数据，和事务启动时看到的是一致的，其他事务的提交不会对当前事务产生影响；
* 串行化：对于记录会加读写锁，出现读写锁冲突时，一个事务必须等待另一个事务执行完成释放锁，才能继续执行；

不同数据库的默认隔离级别不同，MySQL 的默认隔离级别是可重复读，Oracle 的默认隔离级别是读已提交，根据使用场景，可以通过修改配置参数来调整隔离级别。

数据库事务中可能出现的数据一致性问题有：

* 脏写（Dirty Write）：一个事务覆盖了另一个事务为提交的事务，如果被覆盖事务回滚，则会导致数据不一致。通过锁机制能防止事务同时修改同一数据；
* 脏读（Dirty Read）：一个事务读取了另一个事务未提交的数据，如果未提交的事务回滚，读取的数据将无效。读已提交隔离级别能够解决该问题；
* 不可重复读（Non-Repeatable Read）：在一个事务内多次读取同一数据，由于其他事务的修改，导致读取结果不一致。可重复读隔离级别能够解决该问题；
* 幻读（Phantom Read）：一个事务内多次查询同一条件数据集，由于其他事务的插入或删除操作，导致查询结果集发生变化。串行化隔离级别或间隙锁能够解决该问题，MySQL 的可重复读隔离级别通过 MVCC 和间隙锁也可以解决幻读问题；

一个普通的 select 语句在不同隔离级别的逻辑，以及解决的数据一致性问题：

* 读未提交：不加锁，直接读取记录的最新版本，可能发生脏读、不可重复读、幻读问题；
* 读提交：不加锁，每次执行 select 语句时都会生成一个 read-view，解决了脏读问题，可能发生不可重复读、幻读问题；
* 可重复读：不加锁，只在第一次执行 select 语句时生成一个 read-view，解决了脏读、不可重复读、幻读问题；
* 串行化：autocommit = 0 时，select 语句会被转为 select ... lock in share mode 语句，会加 S 锁；autocommit = 1 时，不加锁，会生成一个 read-view，解决了脏读、不可重复读、幻读问题；

通过命令来设置事务的隔离级别：

```mysql
SET [GLOBAL|SESSION] TRANSACTION ISOLATION LEVEL <level>;
```

## 2.1 MVCC

事务隔离是通过数据库的多版本并发控制（MVCC, Multiversion Concurrency Control）实现的。

在 MySQL 中，每条记录在更新的时候都会同时记录一条回滚操作，记录的最新值可以通过回滚操作得到之前状态的值。

假设一个值从 1 被按顺序改成了 2、3、4，在回滚日志里面就会有类似下面的记录。在查询这条记录的时候，不同时刻启动的事务会有不同的 read-view，只需要执行不同的回滚操作就可以得出值，看到的值也就不同。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/database/mysql_read_view.jpg)

旧的回滚日志没有被任何 read-view 用到时，也就不需要再被用到，会被系统判断到并删除。因此，建议尽量不要使用长事务，因为这会在系统中保留很多很老的事务视图，占用大量的存储空间。

InnoDB 中会给每个事务分配一个唯一的事务 ID 即 transaction id，在事务开始时获取，按申请顺序严格递增。表中的一行数据每次更新数据时，都会生成一个新的数据版本，并给这个数据版本记录 row trx_id = transaction id，也就是每一行数据会有多个版本（row）存在。

如下一行数据，存在 V1、V2、V3、V4 四个版本，每个版本记录了当前值和对应的 trx_id：

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/database/mysql_row_version.jpg)

行的更新将会生成 undo log（回滚日志），上图中的 V1、V2、V3 不再存在于表中，而是保存在了 undo log 中，当需要得到这一行之前版本的值时，根据当前版本和 undo log 依次执行 U3、U2、U1 算出来。

在可重复读的隔离级别中，一个事务启动后可以看到所有已提交的事务的结果，而看不见后面的其他事务的更新。InnoDB 为每个事务构造了一个数组，保存该事务启动瞬间，已启动未提交的活跃中的事务 ID。数组中事务 ID 最小值为低水位，系统已创建的最大事务 ID 加 1 为高水位，这个数组和高水位组成了当前事务的一致性视图（read-view）。

![](https://blog-1304941664.cos.ap-guangzhou.myqcloud.com/article_material/database/mysql_watermark.jpg)

在当前事务的一致性视图中，对于一行数据的某个数据版本 row trx_id，可能有这几种情况：

* 属于绿色部份，这个版本是已提交事务或当前事务生成的，可见；
* 属于黄色部份，如果 trx_id 在事务 ID 数组中，则是由还没提交的事务生成的，不可见，如果不在数组中，则是由已经提交的事务生成的，可见；
* 属于红色部份，这个版本是之后启动的事务生成的，不可见；

更新数据必须基于最新数据版本，而不能基于历史数据版本，否则就会丢失别的事务的更新。因此所有更新数据都是先读后写，读只能读当前的最新值，这称为当前读（current read），select 语句如果加锁也是当前读。

```mysql
select k from t where id=1 lock in share mode; # 读锁
select k from t where id=1 for update; # 写锁
```

通过以上机制，可以非常快速地为整个数据库创建快照版本，且能够准确地判断每个数据版本的事务是在该事务之前还是之后提交，从而实现了隔离级别中的可重复读。

# 3. 参考

* [03 | 事务隔离：为什么你改了我还看不见？-MySQL 实战 45 讲-极客时间](https://time.geekbang.org/column/article/68963)
* [08 | 事务到底是隔离的还是不隔离的？-MySQL 实战 45 讲-极客时间](https://time.geekbang.org/column/article/70562)

