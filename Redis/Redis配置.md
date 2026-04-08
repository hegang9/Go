## 配置文件
Redis的配置文件是`redis.conf`，它包含了Redis服务器的各种配置选项。

## 查看 Redis 配置
```bash
CONFIG GET (CONFIG_SETTING_NAME
```

例如，要查看Redis的端口配置，可以使用以下命令：
```bash
CONFIG GET port
```

## 修改 Redis 配置
```bash
CONFIG SET (CONFIG_SETTING_NAME) (CONFIG_SETTING_VALUE)
```

## Redis 常用配置参数说明

以下按**使用频率从高到低**排序（通用场景，单机与常见生产部署）。

| 排序 | 参数 | 作用 | 常见取值示例 |
| --- | --- | --- | --- |
| 1 | `port` | Redis 服务监听端口。 | `port 6379` |
| 2 | `bind` | 绑定监听地址，控制可访问网卡。 | `bind 127.0.0.1` / `bind 0.0.0.0` |
| 3 | `protected-mode` | 未显式配置安全策略时提供默认保护。 | `protected-mode yes` |
| 4 | `requirepass` | 设置访问密码（Redis 6+ 推荐配合 ACL）。 | `requirepass your_password` |
| 5 | `maxmemory` | 限制实例可使用的最大内存。 | `maxmemory 1gb` |
| 6 | `maxmemory-policy` | 内存达到上限后的淘汰策略。 | `maxmemory-policy allkeys-lru` |
| 7 | `appendonly` | 是否开启 AOF 持久化。 | `appendonly yes` |
| 8 | `appendfsync` | AOF 刷盘策略（性能/安全权衡）。 | `appendfsync everysec` |
| 9 | `save` | RDB 快照触发条件。 | `save 900 1`、`save 300 10` |
| 10 | `timeout` | 客户端空闲超时（秒），`0` 为不超时。 | `timeout 0` |
| 11 | `tcp-keepalive` | TCP 保活时间（秒），帮助清理无效连接。 | `tcp-keepalive 300` |
| 12 | `databases` | 逻辑数据库数量。 | `databases 16` |
| 13 | `loglevel` | 日志级别。 | `loglevel notice` |
| 14 | `logfile` | 日志文件路径。 | `logfile /var/log/redis/redis.log` |
| 15 | `dir` | RDB/AOF 等持久化文件目录。 | `dir /data/redis` |
| 16 | `dbfilename` | RDB 文件名。 | `dbfilename dump.rdb` |
| 17 | `appendfilename` | AOF 文件名。 | `appendfilename appendonly.aof` |
| 18 | `daemonize` | 是否以守护进程方式运行。 | `daemonize yes` |
| 19 | `tcp-backlog` | TCP 监听队列长度。 | `tcp-backlog 511` |
| 20 | `hz` | 服务器周期任务频率。 | `hz 10` |
| 21 | `slowlog-log-slower-than` | 慢查询阈值（微秒）。 | `slowlog-log-slower-than 10000` |
| 22 | `slowlog-max-len` | 慢查询日志最大条数。 | `slowlog-max-len 128` |
| 23 | `rename-command` | 重命名高风险命令，提升安全性。 | `rename-command FLUSHALL ""` |
| 24 | `aclfile` | ACL 用户规则文件路径（Redis 6+）。 | `aclfile /etc/redis/users.acl` |
| 25 | `stop-writes-on-bgsave-error` | RDB 持久化失败时是否停止写入。 | `stop-writes-on-bgsave-error yes` |
| 26 | `rdbcompression` | 是否压缩 RDB 文件。 | `rdbcompression yes` |
| 27 | `rdbchecksum` | 是否为 RDB 启用校验和。 | `rdbchecksum yes` |
| 28 | `auto-aof-rewrite-percentage` | AOF 自动重写触发增长百分比。 | `auto-aof-rewrite-percentage 100` |
| 29 | `auto-aof-rewrite-min-size` | AOF 自动重写最小文件大小。 | `auto-aof-rewrite-min-size 64mb` |
| 30 | `replicaof` | 配置当前实例为某主节点从库。 | `replicaof 10.0.0.10 6379` |
| 31 | `replica-read-only` | 从库是否只读。 | `replica-read-only yes` |
| 32 | `repl-backlog-size` | 主从复制积压缓冲区大小。 | `repl-backlog-size 64mb` |
| 33 | `repl-diskless-sync` | 全量复制是否使用无盘传输。 | `repl-diskless-sync yes` |
| 34 | `client-output-buffer-limit` | 限制不同客户端输出缓冲区。 | `client-output-buffer-limit pubsub 32mb 8mb 60` |
| 35 | `lazyfree-lazy-eviction` | 淘汰 key 时是否异步释放内存。 | `lazyfree-lazy-eviction yes` |
| 36 | `activedefrag` | 是否启用主动内存碎片整理。 | `activedefrag yes` |
| 37 | `io-threads` | I/O 线程数（Redis 6+，高并发场景）。 | `io-threads 4` |

