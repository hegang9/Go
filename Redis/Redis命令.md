# Redis 命令手册 (带语法格式说明)

本文档基于 [Redis命令手册](https://redis.com.cn/commands.html) 进行了补充，扩充了**使用格式 (Syntax)** 一列，方便开发过程中快速查阅具体的传参方式。

## 1. 键（Key）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **DEL** | DEL key [key ...] | 用于删除指定的 key |
| **DUMP** | DUMP key | 序列化给定 key，并返回被序列化的值 |
| **EXISTS** | EXISTS key [key ...] | 检查给定 key 是否存在 |
| **EXPIRE** | EXPIRE key seconds | 为给定 key 设置过期时间（秒） |
| **EXPIREAT** | EXPIREAT key timestamp | 用于为 key 设置过期时间，接受 UNIX 时间戳 |
| **PEXPIRE** | PEXPIRE key milliseconds | 设置 key 的过期时间，以毫秒计 |
| **PEXPIREAT** | PEXPIREAT key milliseconds-timestamp | 设置 key 过期时间的时间戳，以毫秒计 |
| **KEYS** | KEYS pattern | 查找所有符合给定模式的 key（如 KEYS *） |
| **MOVE** | MOVE key db | 将当前数据库的 key 移动到给定的数据库中 |
| **PERSIST** | PERSIST key | 移除 key 的过期时间，key 将持久保持 |
| **PTTL** | PTTL key | 以毫秒为单位返回 key 的剩余过期时间 |
| **TTL** | TTL key | 以秒为单位返回给定 key 的剩余生存时间 |
| **RANDOMKEY** | RANDOMKEY | 从当前数据库中随机返回一个 key |
| **RENAME** | RENAME key newkey | 修改 key 的名称，若 newkey 已存在则覆盖 |
| **RENAMENX** | RENAMENX key newkey | 仅当 newkey 不存在时，将 key 改名为 newkey |
| **TYPE** | TYPE key | 返回 key 所储存的值的数据类型 |

## 2. 字符串（String）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **SET** | SET key value [EX seconds｜PX ms] [NX｜XX] | 设置指定 key 的值及过期规则、覆盖规则 |
| **GET** | GET key | 获取指定 key 的值 |
| **GETRANGE** | GETRANGE key start end | 返回 key 中字符串值的子字符串（按偏移量） |
| **GETSET** | GETSET key value | 将给定 key 的值设为 value，并返回 key 的旧值 |
| **GETBIT** | GETBIT key offset | 获取指定偏移量上的位 (bit) |
| **MGET** | MGET key [key ...] | 批量获取一个或多个给定 key 的值 |
| **SETBIT** | SETBIT key offset value | 设置或清除指定偏移量上的位 (0或1) |
| **SETEX** | SETEX key seconds value | 设置 key 的值并同时将过期时间设为 seconds |
| **SETNX** | SETNX key value | 只有在 key 不存在时才设置 key 的值（锁常用） |
| **SETRANGE** | SETRANGE key offset value | 从 offset 开始用 value 覆写给定 key 的字符串值 |
| **STRLEN** | STRLEN key | 返回 key 所储存的字符串值的长度 |
| **MSET** | MSET key value [key value ...] | 同时设置一个或多个 key-value 对 |
| **MSETNX** | MSETNX key value [key value ...] | 所有 key 都不存在时，才同时设置多个 key-value 对 |
| **PSETEX** | PSETEX key milliseconds value | 以毫秒为单位设置 key 的生存时间及值 |
| **INCR** | INCR key | 将 key 中储存的数字值增一 |
| **INCRBY** | INCRBY key increment | 将 key 所储存的数字值增加给定的增量值 |
| **INCRBYFLOAT**| INCRBYFLOAT key increment | 将 key 所储存的数字值增加给定的浮点增量值 |
| **DECR** | DECR key | 将 key 中储存的数字值减一 |
| **DECRBY** | DECRBY key decrement | 将 key 所储存的数字值减去给定的减量值 |
| **APPEND** | APPEND key value | 将 value 追加到 key 原来的值的末尾 |

## 3. 哈希（Hash）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **HSET** | HSET key field value [field value ...] | 用于设置存储在 key 中的哈希表字段的值 |
| **HGET** | HGET key field | 获取存储在哈希表中指定字段的值 |
| **HGETALL** | HGETALL key | 获取在哈希表中指定 key 的所有字段和值 |
| **HMGET** | HMGET key field [field ...] | 获取所有给定字段的值 |
| **HDEL** | HDEL key field [field ...] | 用于删除哈希表中一个或多个字段 |
| **HEXISTS** | HEXISTS key field | 用于判断哈希表中指定的字段是否存在 |
| **HINCRBY** | HINCRBY key field increment | 为哈希表指定字段的整数值做增量运算 |
| **HKEYS** | HKEYS key | 获取哈希表中的所有字段（键）名 |
| **HLEN** | HLEN key | 获取哈希表中字段的数量 |
| **HVALS** | HVALS key | 用于获取哈希表中的所有值 |

## 4. 列表（List）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **LPUSH** | LPUSH key element [element ...] | 将一个或多个值插入到列表头部（左侧） |
| **RPUSH** | RPUSH key element [element ...] | 在列表中添加一个或多个值到尾部（右侧） |
| **LPOP** | LPOP key [count] | 移出并获取列表的第一个元素 |
| **RPOP** | RPOP key [count] | 移除并获取列表最后一个元素 |
| **BLPOP** | BLPOP key [key ...] timeout | 阻塞式移出并获取列表第一个元素 |
| **BRPOP** | BRPOP key [key ...] timeout | 阻塞式移出并获取列表最后一个元素 |
| **BRPOPLPUSH**| BRPOPLPUSH source destination timeout| 阻塞操作：从源列表弹出一个值插入目的列表并返回 |
| **LINDEX** | LINDEX key index | 通过索引获取列表中的元素 |
| **LINSERT** | LINSERT key BEFORE\|AFTER pivot element| 在列表特定的基准元素前或者后插入元素 |
| **LLEN** | LLEN key | 获取列表的长度 |
| **LPUSHX** | LPUSHX key element [element ...] | 将值插入到已存在的列表头部 |
| **RPUSHX** | RPUSHX key element [element ...] | 为已存在的列表的尾部添加值 |
| **LRANGE** | LRANGE key start stop | 获取列表指定范围内的元素 |
| **LREM** | LREM key count element | 移除列表中与 element 相等的 count 个元素 |
| **LSET** | LSET key index element | 通过索引去修改设置列表元素的值 |
| **LTRIM** | LTRIM key start stop | 对列表进行修剪，仅保留指定区间内的元素 |
| **RPOPLPUSH** | RPOPLPUSH source destination | 移除右侧最后一个元素，添加到目标列表左侧并返回 |

## 5. 集合（Set）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **SADD** | SADD key member [member ...] | 向集合添加一个或多个成员 |
| **SCARD** | SCARD key | 获取集合的成员总数 |
| **SISMEMBER** | SISMEMBER key member | 判断 member 元素是否是集合 key 的成员 |
| **SMEMBERS** | SMEMBERS key | 返回集合中的所有成员 |
| **SREM** | SREM key member [member ...] | 移除集合中一个或多个指定的成员 |
| **SDIFF** | SDIFF key [key ...] | 返回所有给定集合的差集 |
| **SDIFFSTORE**| SDIFFSTORE destination key [key ...] | 返回给定集合差集并存储在 destination 集合中 |
| **SINTER** | SINTER key [key ...] | 返回所有给定集合的交集 |
| **SINTERSTORE**|SINTERSTORE destination key [key ...]| 返回给定集合的交集并存储在 destination 集合中 |
| **SUNION** | SUNION key [key ...] | 返回所有给定集合的并集 |
| **SUNIONSTORE**|SUNIONSTORE destination key [key ...]| 所有给定集合的并集存储在 destination 集合中 |
| **SMOVE** | SMOVE source destination member | 将元素从原集合移动到目标集合 |
| **SPOP** | SPOP key [count] | 随机移除并返回集合中的一个或多个元素 |
| **SRANDMEMBER**| SRANDMEMBER key [count] | 返回集合中一个或多个随机元素（不移除） |
| **SSCAN** | SSCAN key cursor [MATCH pattern] | 基于游标迭代查找集合中的元素 |

## 6. 有序集合（Zset）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **ZADD** | ZADD key score member [score member...]| 向有序集合添加一个或多个成员及其分数值 |
| **ZCARD** | ZCARD key | 获取有序集合的成员数 |
| **ZCOUNT** | ZCOUNT key min max | 计算在指定区间分数的成员数 |
| **ZINCRBY** | ZINCRBY key increment member | 有序集合中对指定成员的分数加上增量 increment |
| **ZSCORE** | ZSCORE key member | 返回有序集中，对应成员的分数值 |
| **ZRANK** | ZRANK key member | 返回有序集合中指定成员的索引排名（分数从小到大）|
| **ZREVRANK** | ZREVRANK key member | 返回成员排名，分数递减（从大到小）排序 |
| **ZRANGE** | ZRANGE key min max [WITHSCORES] | 通过指定索引区间返回有序集合内的成员 |
| **ZREVRANGE** | ZREVRANGE key start stop | 通过索引区间返回成员，分数从高到底 |
| **ZRANGEBYSCORE**| ZRANGEBYSCORE key min max | 通过分数返回指定区间内的成员 |
| **ZREVRANGEBYSCORE**| ZREVRANGEBYSCORE key max min | 返回分数区间内的成员，倒序排列 |
| **ZRANGEBYLEX**| ZRANGEBYLEX key min max | 通过字典区间返回有序集合的成员 |
| **ZREM** | ZREM key member [member ...] | 移除有序集合中的一个或多个成员 |
| **ZREMRANGEBYRANK**| ZREMRANGEBYRANK key start stop| 移除给定排名（索引）区间的所有成员 |
| **ZREMRANGEBYSCORE**|ZREMRANGEBYSCORE key min max | 移除给定分数区间的所有成员 |
| **ZREMRANGEBYLEX**| ZREMRANGEBYLEX key min max | 移除给定的字典区间内的所有成员 |
| **ZINTERSTORE**| ZINTERSTORE dest numkeys key [...]| 计算多个有序集的交集并存储在新的 key 中 |
| **ZUNIONSTORE**| ZUNIONSTORE dest numkeys key [...]| 计算多个有序集的并集并存储在新的 key 中 |

## 7. HyperLogLog 命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **PFADD** | PFADD key element [element ...] | 添加指定元素到 HyperLogLog 中 |
| **PFCOUNT** | PFCOUNT key [key ...] | 返回给定基数估算值（UV数量估算） |
| **PFMERGE** | PFMERGE dest sourcekey [sourcekey...]| 将多个 HyperLogLog 结构合并为一个新的 |

## 8. 地理位置（Geo）命令
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **GEOADD** | GEOADD key longitude latitude member | 将指定的地理空间（经度、纬度、名称）添加到 key |
| **GEOPOS** | GEOPOS key member [member ...] | 获取给定位置元素的位置（经度和纬度） |
| **GEODIST** | GEODIST key member1 member2 [m\|km] | 返回两个给定位置之间的距离 |
| **GEOHASH** | GEOHASH key member [member ...] | 返回一个或多个位置元素的 Geohash 编码字符串 |
| **GEORADIUS** | GEORADIUS key lon lat radius m\|km | 以给定的经纬度为中心，根据半径找出范围内元素 |
| **GEORADIUSBYMEMBER**| GEORADIUSBYMEMBER key member r| 根据所给目标元素为中心，找出范围内的成员 |

## 9. 事务与发布订阅
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **MULTI** | MULTI | 标记一个事务块的开始 |
| **EXEC** | EXEC | 执行所有事务块内已缓存的命令队列 |
| **DISCARD** | DISCARD | 取消事务，放弃执行事务块内的所有命令 |
| **WATCH** | WATCH key [key ...] | 监视一个或多个 key，解决并发修改冲突 |
| **PUBLISH** | PUBLISH channel message | 将信息发送到指定的频道 |
| **SUBSCRIBE** | SUBSCRIBE channel [channel ...] | 订阅给定的一个或多个频道的信息 |
| **PSUBSCRIBE**| PSUBSCRIBE pattern [pattern ...] | 订阅匹配所给模式的频道（如 oo.*） |
| **PUBSUB** | PUBSUB subcommand [arg...] | 查看发布订阅操作的内部状态 |

## 10. 管理及服务类
| 命令 | 使用格式 (Syntax) | 功能说明 |
| --- | --- | --- |
| **PING** | PING [message] | 查看服务是否运行，或回显指定内容 |
| **SELECT** | SELECT index | 切换到指定数字索引的数据库 (默认0-15) |
| **AUTH** | AUTH [username] password | 服务提权，验证密码是否正确 |
| **INFO** | INFO [section] | 获取 Redis 服务器状态的各种信息和统计数值 |
| **CONFIG GET** | CONFIG GET parameter | 获取指定配置参数的值 |
| **CONFIG SET** | CONFIG SET parameter value | 修改 redis 配置参数，无需重启 |
| **BGSAVE** | BGSAVE | 在后台异步保存当前数据库的数据到磁盘RDB文件 |
| **BGREWRITEAOF**| BGREWRITEAOF | 异步执行一个 AOF 文件的压缩重写操作 |
| **FLUSHDB** | FLUSHDB [ASYNC] | 仅删除当前数据库的所有 key |
| **FLUSHALL** | FLUSHALL [ASYNC] | 删除实例中所有数据库的所有 key |
| **SLOWLOG** | SLOWLOG subcommand | 查看和管理 redis 的慢查询日志记录 |
| **EVAL** | EVAL script numkeys key... arg... | 将指定的 Lua 脚本内容交由服务器去解析并执行 |
