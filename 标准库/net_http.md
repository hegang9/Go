# net/http 标准库详细指南

`net/http` 包实现了 HTTP 客户端（发出 HTTP 请求）和服务端（处理 HTTP 请求）。它是 Go 构建网络应用框架的基石。

## 常用 API 说明

### 1. `http.Get` / `http.Post`
* **功能**：基础的 HTTP 客户端请求方法。
* **参数说明**：
  * `url string`：请求地址。
  * `bodyType string` (仅 `Post`)：请求的 Content-Type，如 `"application/json"`。
  * `body io.Reader` (仅 `Post`)：请求体数据的读接口，通常使用 `bytes.NewReader(buf)`。
* **返回值**：`*http.Response`（响应结构体）和 `error`。
* **注意事项**：**必须手动关闭 Body**。在使用完响应内容后，务必调用 `defer resp.Body.Close()`，否则会造成协程泄漏和连接描述符泄漏。
* **适用场景**：简单、快速地拉取一个第三方接口时。

### 2. `http.Client`
* **功能**：灵活度更高的自定义 HTTP 客户端结构体。
* **核心字段**：`Timeout`（请求超时控制）、`Transport`（连接池及底层连接的配置）、`CheckRedirect`（重定向策略）。
* **使用方式**：
  * 创建请求：`req, err := http.NewRequest(method, url, body)`。
  * 发送请求：`resp, err := client.Do(req)`。
* **注意事项**：默认的 `http.DefaultClient` 是没有超时时间（Timeout = 0）的，直接发送容易永久阻塞。生成环境强烈建议新建带 Timeout 的 `Client`。
* **适用场景**：需要自定义 Header 发送请求，或者控制超时和连接配置的时候。

### 3. `http.HandleFunc` / `http.Handle`
* **功能**：用于向 HTTP 路由器注入路由规则和对应的处理逻辑（Handler 函数）。
* **参数说明**：
  * `pattern string`：路由路径，例如 `"/api/v1/user"` 或 `/`（匹配所有未匹配的）。
  * `handler func(w http.ResponseWriter, r *http.Request)`：用于处理业务请求的回调函数。
* **参数深度解析**：
  * `w http.ResponseWriter`：服务端的回信通道，我们可以往里面写 Header 信息（`w.Header().Set()`）、设置状态码（`w.WriteHeader(200)`）、写入 `Body` 数据（`w.Write(content)`）。
  * `r *http.Request`：包含了客户端发来的一切数据（URL 参数、请求体、Header、Method 等）。
* **注意事项**：Go 1.22 增强了原生的路由匹配能力，原生可以支持 `POST /books/{id}` 的匹配模式。
* **适用场景**：轻量级服务端接口搭建。对于复杂项目，依然建议使用 Gin 等基于此封装的高级框架。

### 4. `http.ListenAndServe`
* **功能**：启动并监听 TCP 网络地址，处理进来的 HTTP 连接。
* **参数说明**：
  * `addr string`：网络地址（如 `":8080"`, `"127.0.0.1:80"`）。
  * `handler http.Handler`：路由器实例，若传 `nil`，则默认使用 `http.DefaultServeMux`（对应了上面使用 `http.HandleFunc` 注册的路由器）。
* **注意事项**：
  * 这是一个阻塞方法，当程序正常运行时它会一直接管当前的协程。
  * 它不能很好地实现“优雅启停”或自定义超时控制。
* **适用场景**：启动基础和测试性质的一般 Web 接口服务。

### 5. `http.Server`
* **功能**：自定义的 HTTP 服务端。
* **参数说明**：配置诸如 `Addr`, `Handler`, `ReadTimeout`, `WriteTimeout` 等。
* **用法与注意事项**：使用 `server.ListenAndServe()` 启动。并可以使用 `server.Shutdown(ctx)` 进行优雅关闭（等待当前处理完的所有请求关闭）。
* **适用场景**：生产环境核心 HTTP 服务器部署。
