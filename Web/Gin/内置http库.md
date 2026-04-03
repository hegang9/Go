## 内置http库

Go语言的标准库 `net/http` 提供了非常强大且易用的 HTTP 客户端和服务端实现。在Go语言中，甚至不需要借助第三方框架（如 Gin、Echo），仅仅使用标准库的 `net/http` 就足以构建高性能的 Web 服务。

### 核心运行机制

使用 `net/http` 构建服务端应用的核心机制围绕以下几个基本概念：

1. **Handler（处理器）**：所有处理 HTTP 请求的逻辑最终都必须实现 `http.Handler` 接口，该接口只有一个方法 `ServeHTTP(ResponseWriter, *Request)`。
2. **ServeMux（多路复用器/路由）**：它是 HTTP 请求的路由器，负责接收请求并将请求的 URL与注册的路由模式进行匹配，然后将请求转发给对应的 Handler 执行。
3. **Server（服务器）**：负责监听指定的网络端口，接收 HTTP 请求，并将请求交给 ServeMux 进行处理。

### 基本使用示例

以下是一个使用 `net/http` 搭建最简单 Web 服务器的例子：

```go
package main

import (
	"fmt"
	"net/http"
)

// 定义一个处理函数，签名为 func(http.ResponseWriter, *http.Request)
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 解析参数、表单等
	r.ParseForm()
	
	// 向客户端写入响应内容
	fmt.Fprintf(w, "Hello, Go Web! 请求的方法是: %s, 路径是: %s", r.Method, r.URL.Path)
}

func main() {
	// 1. 注册路由和处理函数 (默认将其绑定到全局的 DefaultServeMux)
	http.HandleFunc("/", helloHandler)

	// 2. 启动HTTP服务
	fmt.Println("服务器即将启动，监听地址: http://localhost:8080")
	// 第一个参数是监听的地址和端口，第二个参数是 handler，传 nil 则表示使用默认的多路复用器
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}
```

---

### 常用API概念、功能与用法

在原生的 Go Web 开发中，以下是 `net/http` 在开发过程中最为核心、常用的 API 列表：

#### 1. 路由与服务器启动接口
*   **`http.ListenAndServe(addr string, handler Handler) error`**
    *   **概念与功能**：配置并启动一个 HTTP 服务器。监听指定的 TCP 地址（如 `":8080"`），并拦截传入 HTTP 协议连接交由 handler 处理。
    *   **用法**：`http.ListenAndServe(":8080", nil)`。如果 handler 传 `nil`，则默认使用全局唯一的路由器 `DefaultServeMux`。
*   **`http.HandleFunc(pattern string, handler func(ResponseWriter, *Request))`**
    *   **概念与功能**：将指定的路由规则（`pattern`）和一个匹配的处理函数（`handler`）注册到默认的路由器上。
    *   **用法**：`http.HandleFunc("/api/user", userHandler)`。

#### 2. 自定义路由器实例
*   **`http.NewServeMux()`**
    *   **概念与功能**：创建一个新的、独立的路由器（ServeMux）实例。开发复杂项目时，为了避免路由污染，一般不会使用默认的全局路由，而是自己实例化并管理它。
    *   **用法**：
    ```go
    mux := http.NewServeMux()
    mux.HandleFunc("/api/login", loginHandler)
    http.ListenAndServe(":8080", mux) // 替换掉原来的 nil，注册自定义的 mux
    ```

#### 3. 响应构建：`http.ResponseWriter`
这是一个接口类型的参数，通过它可以构造和操控返回给客户端的 HTTP 响应体。
*   **`w.Write([]byte) (int, error)`**
    *   **功能**：向响应体中写入二进制数据或字符串。
    *   **用法**：`w.Write([]byte("响应成功"))`。
*   **`w.WriteHeader(statusCode int)`**
    *   **功能**：发送自定义的 HTTP 协议状态码（如 `http.StatusOK` 200，`http.StatusNotFound` 404）。
    *   **注意**：此方法必须在 `w.Write()` 之前调用生效。
*   **`w.Header().Set(key, value string)`**
    *   **功能**：设置响应头（Headers）数据。例如设置跨域、内容格式（返回JSON数据）等。
    *   **用法**：`w.Header().Set("Content-Type", "application/json")`。同样需要在写入响应内容之前设定。

#### 4. 请求解析：`*http.Request`
它是一个结构体类型的指针参数，内部包含了当前客户端发来的 HTTP 请求所包含的所有上下文信息。
*   **`r.Method`**
    *   **功能**：获取本次请求的方法类型（诸如 `"GET"`, `"POST"`, `"PUT"`, `"DELETE"` 等）。
*   **`r.URL.Path`**
    *   **功能**：获取本次请求的相对路由地址部分（例 `/api/profile`）。
*   **`r.URL.Query().Get(key string)`**
    *   **功能**：获取 URL 路径中携带的 Query 查询参数（例 `?id=1024&name=admin` 获取 `id` 或 `name`）。
*   **`r.Body`**
    *   **功能**：获取请求体。一般类型为 `io.ReadCloser`，常被用于读取 POST 请求传输的 JSON 或者原始流文件内容。需要使用 `io.ReadAll(r.Body)` 提取字节数据。
*   **`r.FormValue(key string)`**
    *   **功能**：快捷解析并返回表单中指定 `key` 的值。此函数会自动调用 `ParseForm`，并将 URL Query 参数和 HTTP Body 的表单数据进行统一合并提供查询。
    *   **用法**：`password := r.FormValue("password")`。
