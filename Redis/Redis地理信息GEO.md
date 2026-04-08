# Redis 地理信息 (GEO) 简易指南 

## 1. 什么是 GEO？（大白话解释）
想象一下你现在打开了美团外卖，系统怎么知道**离你最近的奶茶店**在哪里？又或者你打开微信**附近的人**，系统是怎么算出谁在和你相距 500 米的？

这就需要用到一种专门存经纬度（地球上的坐标）的数据结构，Redis GEO 就是干这个的。它能帮你把现实世界的位置记下来，还能瞬间算出两个地方有多远，或者搜出你方圆几公里以内的东西。

## 2. 小白必须知道的剧透

Redis 为了实现 GEO，**并没有去发明一个新的底层结构**！
GEO 在 Redis 里的本质，其实就是一个**有序集合（ZSet）**！

- Redis 会把一个二维的地理坐标（经度 + 纬度），通过一种叫 GeoHash 的算法，压缩成一个一维的数字。
- 然后把这个数字当作分数（Score），把你存的名字（比如星巴克）当作成员（Member），统统塞进一个普通的 ZSet 里。
- **重点**：既然它本质是 ZSet，你也**可以用 ZSet 的命令去删改 GEO 的数据**！比如用 ZREM 去删除一个地点。

## 3. 核心大招（手把手实战）

 **避坑指南**：我们平时习惯说纬度和经度。但是在 Redis GEO 里添加坐标时，**永远是经度在前，纬度在后 （Longitude Latitude）**！不要搞反了，不然你会跑到海里去。

### 第一招：添加位置 (GEOADD)
把北京和上海的坐标加进名为 china:cities 的集合里：
```redis
# 格式：GEOADD 集合名 经度 纬度 名字
GEOADD china:cities 116.40 39.90 "Beijing"
GEOADD china:cities 121.47 31.23 "Shanghai"
```

### 第二招：算算有多远 (GEODIST)
美团外卖怎么算骑手离你有多远的？就靠这招：
```redis
# 算北京到上海的直线距离，km 表示千米，m 表示米
GEODIST china:cities Beijing Shanghai km
# 返回值：大概是 1067.59 千米
```

### 第三招：获取地点的经纬度 (GEOPOS)
如果你忘了某个已经加进去的地点的坐标：
```redis
# 这个命令会返回北京的经度、纬度
GEOPOS china:cities Beijing
```

### 第四招：附近的人 (GEOSEARCH / GEORADIUS)
*(注：Redis 6.2 以后推荐用 GEOSEARCH，以前老版本用 GEORADIUS，这里给你看最好理解的找法：)*

我要找出距离北京 **1500 公里内** 的所有城市：
```redis
# 以 Beijing 为圆心，找 1500 千米内的所有成员
GEORADIUSBYMEMBER china:cities Beijing 1500 km
# 会返回 Beijing 和 Shanghai
```

如果我想知道这些城市离北京到底有多少米？还可以加上 WITHDIST：
```redis
GEORADIUSBYMEMBER china:cities Beijing 1500 km WITHDIST
# 返回：
# 1) "Beijing"  -> 0.0000 km
# 2) "Shanghai" -> 1067.5980 km
```

## 4. 常见应用场景
学会了这几个简单的命令，你其实就能做很多很牛的功能了：
-  **滴滴打车**：查找距离自己 3 公里内的空闲司机。
-  **社交软件**：微信附近的人、探探同城匹配。
-  **共享单车**：找找离我最近的哈啰单车在哪里。
-  **本地生活**：美团、饿了么里的距您 800m。

---

## 5. 常用命令速查表

掌握这几个常用命令，你就基本可以在实际项目中把 GEO 玩转起来了！

| 命令 | 作用说明 | 常见应用场景 |
| :--- | :--- | :--- |
| **`GEOADD`** | `GEOADD key longitude latitude member` <br> 向指定的 key 中添加一或多个地理位置（经度 纬度 成员名）。 | “骑手上线”：记录外卖小哥当前的经纬度。<br>“商铺入驻”：保存新开星巴克的门店坐标。 |
| **`GEODIST`** | `GEODIST key member1 member2 [unit]` <br> 计算并返回两个已存位置之间的距离。单位可选 m、km、ft、mi。 | “距离估算”：查看用户距离想去的餐厅有多远。<br>“运费计算”：按两点间的直线距离收取跑腿费。 |
| **`GEOPOS`** | `GEOPOS key member [member ...]` <br> 从键里把某几个成员存的经纬度给查出来。 | “位置回显”：用户在地图上点击自己的历史常去地，读取它的准确坐标以在地图打点。 |
| **`GEOHASH`** | `GEOHASH key member [member ...]` <br> 返回一个或多个成员的 11位 Geohash 字符串表示形式。 | “近似位置共享”：想告诉别人“我在这一带”但不暴露具体坐标时，分享hash值即可。 |
| **`GEOSEARCH`** | `GEOSEARCH key FROMMEMBER member BYRADIUS radius unit` <br> (Redis 6.2+) 在指定范围内搜索返回匹配的位置成员。支持按成员、经纬度及矩形/圆形范围查。 | “附近的人/车”：以用户当前位置为圆心，找出方圆 3 公里内的所有共享单车。 |
| **`GEORADIUS`** | `GEORADIUS key longitude latitude radius unit` <br> (已废弃/老版本用) 找出指定经纬度半径范围内的所有元素，可带距离 `WITHDIST` 或坐标 `WITHCOORD`。 | “商圈搜索”：搜索特定商场坐标（经纬度）附近 500 米内的可用优惠商家。 |
| **`ZREM`** | `ZREM key member [member ...]` <br> GEO 底层就是 Zset，删除位置直接用 Zset 的删除命令即可。 | “骑手下线”：外卖小哥收工了，就把他的坐标从在线地图里删掉。 |