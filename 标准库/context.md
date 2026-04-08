# context 标准库详细指南

`context` 是 Go 并发编程最基本、最重要、几乎出现在所有开源库源码里的包。它定义了 `Context` 接口，它不仅能包含**截止时间、取消信号**，还能携带**请求作用域内的值**，在多个 Goroutine 之间安全传递。

## 核心接口说明 (Context)
* `Deadline() (deadline time.Time, ok bool)`：返回取消的时间，如果没有则 ok 为 false。
* `Done() <-chan struct{}`：当上下文被取消或到达超时时间，该通道会被关闭（通常与 `select` 配合使用）。
* `Err() error`：返回取消的原因（`context.Canceled` 或 `context.DeadlineExceeded`）。
* `Value(key any) any`：返回绑定在该 context 上的特定值。

## 常用 API 说明

### 1. `context.Background` / `context.TODO`
* **功能**：这两个函数返回非 nil 的空 `Context`。它们不能被取消，没有值，也没有截止期限。
* **参数说明**：无。
* **返回值**：`context.Context`（顶层基础上下文）。
* **注意事项**：不要往函数参数里传 `nil` 作为 context，不确定用什么时候就用 `context.TODO()`，在 `main` 函数或处理顶级请求时用 `context.Background()`。
* **适用场景**：这是所有上下文请求链路的根（Root），其他具有特定属性的上下文（例如设置了超时）都要从它们派生而来。

### 2. `context.WithCancel`
* **功能**：基于父 Context 派生出一个新的带取消功能的 Context。
* **参数说明**：`parent context.Context`：父节点上下文。
* **返回值**：`(ctx context.Context, cancel context.CancelFunc)`：新的派生上下文，和一个**用于触发取消的函数**。
* **注意事项**：调用了 `cancel()` 函数会导致 `ctx.Done()` 通道关闭。父节点的取消会自动扩散给所有的子孙节点，但子孙节点的取消不会影响父节点。
* **适用场景**：有一个 Goroutine 监控某个任务状态（例如接收到了退出信号或发起了重试），需要主动终止所有关联的子 Goroutine 以释放资源时。

### 3. `context.WithTimeout`
* **功能**：相当于基于父节点派生一个自带计时器的 Context。
* **参数说明**：
  * `parent context.Context`：父节点上下文。
  * `timeout time.Duration`：相对时间差，比如 `time.Second * 5`。
* **返回值**：`(ctx context.Context, cancel context.CancelFunc)`。
* **注意事项**：即使到了超时时间 context 自动取消了，依然在逻辑结束后尽早手动执行 `defer cancel()`，可以提前释放底层的定时器资源。如果不提前取消，会引发协程泄露和微小的内存积压。
* **适用场景**：HTTP 请求、数据库查询、微服务 RPC 时设置操作耗时上限，防止发生死锁或耗尽系统资源。

### 4. `context.WithDeadline`
* **功能**：类似于 `WithTimeout`，但它接收的不是一段时间跨度，而是一个**绝对的截止时间点**。
* **参数说明**：
  * `parent context.Context`：父节点上下文。
  * `d time.Time`：绝对时间，如明天中午 12 点 `time.Now().Add(...)`。
* **注意事项**：如果在 `d` 到来前，操作就完成了，记得尽早调用返回的 `cancel()`。如果 `d` 早于当前时间，它会立刻被判定为超时。
* **适用场景**：处理定时调度的脚本任务，在已知必须要结束的确切时刻之前。

### 5. `context.WithValue`
* **功能**：将全局链路级别的数据附带在 Context 节点上。
* **参数说明**：
  * `parent context.Context`：父节点上下文。
  * `key any`：键。
  * `val any`：值。
* **注意事项**：
  * `key` 的类型不该是内置的类型（如 `string` 或 `int`），这容易引发不同包之间的 key 冲突。官方推荐的做法是用自定义非导出的类型作为键：`type contextKey string; const myKey contextKey = "xxx"`。
  * 不应该用来传递函数的普通参数和业务数据，仅仅用来放一些**请求级生命周期的元数据**。
* **适用场景**：用来传递 API 网关打上的 Request ID (链路追踪 TraceID) 、请求携带的用户身份信息 (User Auth Token) 或者数据库连接池 Session。
