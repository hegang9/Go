# time 标准库详细指南

`time` 包提供了测量和显示时间的功能，不仅包含简单的耗费时间和休眠操作，还集成了极其重要的定时器触发机制和非常特殊的“Go 式时间格式化法则”。

## 结构体说明
所有的具体时间都体现为 `time.Time` 结构。它包含了经过时区（`*time.Location`）映射的秒、纳秒值。\n

## 常用 API 说明

### 1. `time.Now`
* **功能**：获取系统的当前时间。
* **参数说明**：无。
* **返回值**：`t time.Time` (带有本地时区信息的当前时间)。
* **适用场景**：一切需要获取当前时刻的地方，如生成日志时间戳、计算运行耗时等。

### 2. `time.Unix` / `time.UnixMilli` / `t.Unix()`
* **功能**：将时间转换为或从时间戳（自 Unix 纪元 1970年1月1日 00:00:00 UTC 的秒数）来初始化。
* **用法解析**：
  * `t.Unix()`：获取 `t time.Time` 对象的**秒级**时间戳。
  * `t.UnixMilli()` (Go 1.17+)：获取**毫秒级**时间戳。
  * `time.U  nix(sec int64, nsec int64)`：将时间戳转换回 `time.Time` 对象。
* **适用场景**：记录数据库的创建时间和更新时间、将时间传递给前端。

### 3. `t.Format` / `time.Parse`
* **功能**：Go 特色的时间**格式化与解析**。
* **注意事项（极其重要）**：
  * **在 Go 中，格式化的基准时间必须是 `2006-01-02 15:04:05`**！（可以用助记词：`1 月 2 日 下午 3 点 4 分 5 秒，2006 年` 对应 `123456`）。
  * 决不能像别的语言那样用 `%Y-%m-%d %H:%M:%S`。
* **用法示例**：
  * `now.Format("2006/01/02 15:04:05.000")` (带毫秒的格式化)。
  * `t, err := time.Parse("2006-01-02", "2023-10-01")` (按此模板解析字符串时间)。
* **注意事项**：`time.Parse` 会默认解析为 **UTC 时区**。在中国，应该使用 `time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)` 来解析字符串时间。

### 4. `time.Sleep` / `time.After`
* **功能**：
  * `Sleep` 阻塞当前 Goroutine 一段指定的时间间隔。
  * `After` 会在另一条不阻塞的协程里经过指间间隔后向一个通道 (`<-chan Time`) 发送当时的时间。
* **参数说明**：`d time.Duration` (时间段常量，例如 `time.Second * 5`)。
* **用法示例**：
```go
  time.Sleep(1 * time.Second) // 睡1秒
  select {
  case <-time.After(3 * time.Second): // 配合 select 实现异步超时监控
      fmt.Println("超时了!")
  }
```
* **适用场景**：节流、死循环中的停顿、非 context 的暴力超时控制。

### 5. `time.Since` / `time.Until`
* **功能**：快捷计算相对时间。
* **参数说明**：`t time.Time`。
* **返回值**：`time.Duration`，是基于纳秒精度的时差。
* **用法示例**：
```go
  start := time.Now()
  // ... 各种耗时操作 ...
  cost := time.Since(start)
```
* **适用场景**：性能监控，在中间件中统计 HTTP 请求处理耗时。

### 6. `time.Ticker` / `time.Timer`
* **功能**：
  * `Ticker`：提供了一个按给定时间段**重复**地触发周期的定时器（通过其 `C` channel 送达）。
  * `Timer`：提供了一个在此刻到达后**仅触发一次**的通道定时器。
* **用法说明**：
  * `ticker := time.NewTicker(time.Second)`。每次被读取 `<-ticker.C` 都能拿到跳动这一秒时刻，用于周期任务。
  * `timer := time.NewTimer(time.Second)`。当读出 `<-timer.C` 后生命周期结束。它可以用 `timer.Stop()` 或者 `timer.Reset()` 重置。
* **注意事项**：使用完毕后务必执行 `defer ticker.Stop()` / `timer.Stop()`，否则底层的定时器机制会永远维持它们泄漏内存资源。
* **适用场景**：长连接心跳包维护、缓存定期刷新、任务调度。
