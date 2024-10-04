# MiniLSM

LSM-Tree(Log Structured Merge Tree) 存储引擎实现，实现一个轻量级的 kv 存储引擎，面向写多读少的场景。

## 需求

一个大概完整的 LSM-Tree 包含以下部分：
- [ ] 基于内存的 memtable：采用跳表实现。
- [ ] 读写分离的 memtable：分为可读可写的 active memtable，和只读的 read-only memtable，当 active 达到阈值的时候，将其切换为只读模式，并写入磁盘中。
- [ ] WAL（Write Ahead Log）：预写日志，避免基于内存的 memtable 数据丢失。
- [ ] 分层的 Level：Level0 由内存的 memtable 溢写得到，level1~k 由上一层和本层的文件合并得到。
- [ ] SSTable（Sorted String Table）：存储数据的磁盘文件称为 SSTable，以 Block 为单位进行数据拆分并建立索引。
- [ ] 布隆过滤器：针对 Block 为单位建立布隆过滤器，辅助数据校验。
- [ ] kv 记录的读/写
- [ ] 其他有待补充。。。

原理部分参考：[xiaoxuxiansheng/golsm](https://github.com/xiaoxuxiansheng/golsm)

## 未来计划

- [ ] 实现 Redis 协议
- [ ] 实现 Raft 算法