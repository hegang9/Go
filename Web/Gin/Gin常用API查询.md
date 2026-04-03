# Gin 常用 API 参考手册

本手册汇集了 Gin 框架最常用的核心 API，涵盖项目初始化、路由控制、请求参数处理、响应生成与中间件机制等。

---

## 🚀 1. 引擎与服务启动 (Engine & Serving)

### `gin.Default()`
* **概念与功能**：初始化一个带有默认中间件（Logger 和 Recovery）的 Gin 引擎实例。
* **返回值**：`*gin.Engine` (拥有路由器功能)
* **用法**：
  ```go
  r := gin.Default()
  ```

### `gin.New()`
* **概念与功能**：初始化一个完全空白的 Gin 引擎，不包含任何内置中间件。常用于按需挂载自定义日志和崩溃恢复中间件。
* **返回值**：`*gin.Engine`
* **用法**：
  ```go
  r := gin.New()
  // 手动挂载中间件
  r.Use(gin.Logger(), gin.Recovery())
  ```

### `Engine.Run(addr ...string)`
* **概念与功能**：以 HTTP 形式启动服务，监听指定端口。底层是对 `http.ListenAndServe` 的封装。
* **可选参数**：
  * **`addr`**：监听地址列表，如 `":8080"`, `"127.0.0.1:9090"`。如果不传，默认为 `:8080`（读取 `PORT` 环境变量，若无则使用 `:8080`）。
* **用法**：
  ```go
  r.Run(":8080")
  ```

### `Engine.RunTLS(addr, certFile, keyFile string)`
* **功能**：以 HTTPS 模式启动服务。
* **参数意义**：
  * **`addr`**：监听地址端口，如 `":443"`。
  * **`certFile`**：指定域名的 TLS 证书文件物理路径（如 `./server.pem`）。
  * **`keyFile`**：相对应的私钥文件物理路径（如 `./server.key`）。
* **用法**：
  ```go
  r.RunTLS(":443", "./server.pem", "./server.key")
  ```

### `gin.SetMode(value string)`
* **功能**：设置 Gin 的运行模式，降低或提高框架日志打印级别以提升性能。**必须在初始化 Engine 前或尽早调用**。
* **参数意义**：`value` 可选值为 `gin.DebugMode` (默认开启多彩路由日志)、`gin.ReleaseMode` (生产环境，将关闭默认控制台打印)、`gin.TestMode`。
* **用法**：
  ```go
  gin.SetMode(gin.ReleaseMode)
  r := gin.Default()
  ```

### `Engine.Static(relativePath, root string)` / `Engine.StaticFS(...)`
* **功能**：基于 URL 前缀提供本地静态文件（图片、CSS、JS等）的托管服务。
* **参数意义**：
  * **`relativePath`**：客户端访问的 URL 路由前缀（需形如 `/assets`）。
  * **`root`**：对应暴露给外部的本地物理目录相对或绝对路径（如 `./public`）。
* **用法**：
  ```go
  // 请求 /public/logo.png 将隐式返回服务端 ./assets_dir/logo.png 文件
  r.Static("/public", "./assets_dir")  
  ```

### `Engine.LoadHTMLGlob(pattern string)` / `Engine.LoadHTMLFiles(files ...string)`
* **功能**：从指定目录规则或文件列表中预加载并解析 HTML 模板文件，放入内存供后续 `c.HTML()` 调用直接渲染使用。
* **参数意义**：`pattern` 是类 unix glob 的匹配符。`files` 是不定长的指定模板文件列表。
* **用法**：
  ```go
  r.LoadHTMLGlob("templates/**/*")
  // r.LoadHTMLFiles("templates/index.html", "templates/login.html")
  ```

---

## 🚏 2. 路由 (Routing)

### `Engine.GET(relativePath string, handlers ...HandlerFunc)`
* **概念与功能**：注册处理 HTTP GET 请求的路由。同理还有 `POST()`, `PUT()`, `DELETE()`, `PATCH()`, `OPTIONS()`, `HEAD()`。
* **参数意义**：
  * **`relativePath`**：路由路径。可以是具体路径 `/users`，也可以带参数 `/users/:id`，或是模糊匹配 `/users/*action`。
  * **`handlers`**：变长参数。可以传一个或多个中间件及处理函数。
* **用法**：
  ```go
  r.GET("/ping", func(c *gin.Context) {
      c.String(200, "pong")
  })
  ```

### `Engine.Group(relativePath string, handlers ...HandlerFunc)`
* **概念与功能**：创建一个路由组。用于对拥有相同前缀和中间件的路由进行分组管理（例如 API 版本控制 `/api/v1`）。
* **参数意义**：同上。
* **用法**：
  ```go
  v1 := r.Group("/v1", AuthMiddleware())
  {
      v1.GET("/users", getUsers)
      v1.POST("/users", createUser)
  }
  ```

### `Engine.Any(relativePath string, handlers ...HandlerFunc)`
* **功能**：注册一个万能通用路由。它将**同时匹配所有标准 HTTP 请求方法**（GET、POST、PUT、PATCH、DELETE、HEAD、OPTIONS 等），常被用来做转发代理或无需区分动词的同一业务 API（例如微信支付 Webhook）。
* **参数意义**：与 `GET()` 等方法完全相同。
* **用法**：
  ```go
  r.Any("/webhook", webhookHandler)
  ```

### `Engine.Handle(httpMethod, relativePath string, handlers ...HandlerFunc)`
* **功能**：底层路由注册方法。可被用来注册非标准或自定义长度的不常用 HTTP 请求方法（例如 REST 外的自定义指令）。
* **参数意义**：大写且合法的 HTTP 动词字符传（如 `"PROPFIND"`、`"TRACE"` 等）。其余参数同基础路由。
* **用法**：
  ```go
  r.Handle("SEARCH", "/resources", func(c *gin.Context){ /* ... */ })
  ```

### `Engine.NoRoute(handlers ...HandlerFunc)`
* **功能**：引擎级别的全局缺省注册，即：针对 **未被匹配命中的 404 路由异常** 进行接管处置。默认其输出是纯文本 `404 page not found`，通过它可以输出更优雅一致的 404 JSON 出局响应页面。
* **参数意义**：同正常路由，传入自定义的中转或最终 HandlerFunc 逻辑闭包链。
* **用法**：
  ```go
  r.NoRoute(func(c *gin.Context) {
      c.JSON(404, gin.H{"code": 4040, "message": "该功能/页面未开放或已移除"})
  })
  ```

### `Engine.NoMethod(handlers ...HandlerFunc)`
* **功能**：全局注册方法处理未被允许的方法情况（路由路径匹配到了，但用户请求的 `Method` 不存在）。对应的状态码属于 `405 Method Not Allowed`，用法与 `NoRoute` 如出一辙。

---

## 📥 3. 请求解析与参数获取 (Request Context - gin.Context)

在 Gin 中，请求的所有操作几乎都由 `*gin.Context` (简称 `c`) 完成。

### 📌 3.1 路径参数 (Path Parameters)
#### `c.Param(key string)`
* **功能**：获取路由定义中的动态路径参数（形如 `:id`）。
* **参数意义**：`key` 是注册路由时的名称标识。
* **用法**：
  ```go
  // 注册: r.GET("/user/:id", handler)
  id := c.Param("id") // 请求 /user/123 -> id="123"
  ```

### 📌 3.2 URL查询参数 (Query Parameters)
#### `c.Query(key string)` / `c.DefaultQuery(key, defaultValue string)`
* **功能**：获取 URL 中的 `?key=value` 数据。
* **参数意义**：
  * **`key`**：参数名。
  * **`defaultValue`**：如果在 URL 中未找到该键，返回此默认值。
* **用法**：
  ```go
  // 请求 /search?name=john
  name := c.Query("name")               // "john"
  age := c.DefaultQuery("age", "18")    // 获取不到 age，默认返回 "18"
  ```

### 📌 3.3 请求体与表单参数 (Form / Body)
#### `c.PostForm(key string)` / `c.DefaultPostForm(key, defaultValue string)`
* **功能**：获取存在 HTTP Request Body 中的 `x-www-form-urlencoded` 或 `multipart/form-data` 格式数据。常用于表单提交。
* **用法**：
  ```go
  username := c.PostForm("username")
  page := c.DefaultPostForm("page", "1")
  ```

### 📌 3.4 模型绑定 (Model Binding)
为了避免手动获取海量参数，Gin 提供了强大的自动解析到结构体的 API。

#### `c.ShouldBind(obj any)`
* **功能**：根据请求标头 `Content-Type` 自动推断格式，将请求数据绑定到结构体指针上。若校验失败返回 error。
* **参数意义**：
  * **`obj`**：目标结构体实例的**指针**。需要配合结构体 Tag（如 `json`, `form`, `binding`）使用。
* **用法**：
  ```go
  var loginReq struct {
      Username string `form:"username" binding:"required"`
      Password string `form:"password"`
  }
  if err := c.ShouldBind(&loginReq); err != nil {
      c.JSON(400, gin.H{"error": err.Error()})
      return
  }
  ```

#### `c.ShouldBindJSON(obj any)`
* **功能**：强行仅按 `application/json` 解析请求体。
* **用法**：同上。

#### `c.Bind(obj any)` 
* **功能**：功能与 `ShouldBind` 原理极其相似解析绑定结构体，**差别在于处理校验败的策略方式**。
* **差异说明**：如果使用 `Bind()`，它因内部实现会在发生错误时直接帮你强制写入 Header 并将 HTTP `400 状态响应`下发生效！如果你想自己包装一个特定的 `JSON 400 Msg` 给到你移动端/前端，你**千万不能**使用该方法，而必须去使用更温和灵活的带有 `Should` 前缀的 `ShouldBind()` 手动接管。
* **用法**：
  ```go
  // 若发生校验错会隐式自动向客户端刷入 err 信息（不再容许后续自定义任何 Header 和 状态码）
  c.Bind(&loginReq) 
  ```

#### `c.ShouldBindUri(obj any)` / `c.ShouldBindHeader(obj any)`
* **功能**：快捷包装方法，专为路由 URL 后边跟的参数（如 `/:uid/:pid` 系列）以及 HTTP 头中传输的一定批量附带属性（例如 `Token`、`Version`、`AppId`），进行强类型校验及结构映射提取。
* **参数意义**：必须提前把你要绑定的 `struct` 的 field 用 `uri`、`header` 标签（Tag）分别对应起来，否则抓不到包数据。
* **用法**：
  ```go
  var param struct {
      ID string `uri:"uid" binding:"required,uuid"`
  }
  // URL: GET /user/bbcd1000
  if err := c.ShouldBindUri(&param); err == nil {
      fmt.Println(param.ID) 
  }
  ```

### 📌 3.5 进阶请求与底层数据解析
#### `c.GetHeader(key string)`
* **功能**：便捷安全地获取 HTTP 原生的请求标头值。底层对 `c.Request.Header.Get(key)` 做了极简包装并作安全保障。
* **参数意义**：标准的或自定义的头关键字如 `"Authorization"`, `"User-Agent"`。
* **用法**：
  ```go
  token := c.GetHeader("Authorization")
  ```

#### `c.ClientIP()`
* **功能**：通过框架强大的逻辑猜测引擎获得真实的客户端公网 IP 地址。这在存在类似 Nginx、LB 的代理网络集群里极其实用。
* **参数意义**：如果启动了 `ForwardedByClientIP` 功能设定，引擎在读取本机 Socket 前会自动先从反向代理传输的最前端 Header (`X-Real-IP`, `X-Forwarded-For`) 中识别取原始IP进行欺骗穿透保护提取。
* **用法**：
  ```go
  ip := c.ClientIP() // 返回诸如 "192.168.1.10" 的字串
  ```

#### `c.Cookie(name string)`
* **功能**：通过给定的特定名尝试去捕获该次 HTTP 请求所对应的 Cookie 字串（比如登录后附带 Session，校验用户状态的极简场景）。
* **参数意义**：必须是匹配完全的名字，否则抛出未匹配异常 `ErrNoCookie`。
* **用法**：
  ```go
  sessionStr, err := c.Cookie("gin_session_id")
  if err != nil { /* 没登录 */ }
  ```

#### `c.FormFile(name string)` / `c.MultipartForm()`
* **功能**：处理带有文件上传场景的 `multipart/form-data` 类型，可以无感知直接截取大表单或返回附带多种文件的集合数据体（即可以处理单图片、多附件上传的逻辑）。
* **参数意义**：与 `PostForm` 类似，指定 HTML `input type="file" name="uploadfile"` 的 `name` 属性值。
* **用法**：
  ```go
  // "avatar" 字段对应了 1 张上传图片
  fileHeader, err := c.FormFile("avatar") 
  ```

#### `c.SaveUploadedFile(file *multipart.FileHeader, dst string)`
* **功能**：非常强大的辅助轮子 API。允许一行代码直接将之前利用 `FormFile` 截取到的客户端在内存/临时区中的分块文件，落盘转移保存至服务端真正的目的硬盘绝对相对地址中。
* **参数意义**：传入已得到解析通过的 `*multipart.FileHeader` 实例和目的磁盘地址 `dst`。
* **用法**：
  ```go
  file, _ := c.FormFile("file")
  // `dst` 如为 "./uploads/" + file.Filename，必须确保 uploads 文件夹结构存在
  _ = c.SaveUploadedFile(file, "/data/nginx/images/"+file.Filename)
  ```

#### `c.GetRawData()`
* **功能**：无脑获取请求体 Body 中的原始二进制 `[]byte` 数据。如果想自己造轮子进行特殊的包体解密（如接收经过 RSA 的密文流），这种属于最能读到底层数据流进行处理的方法。**读完该 API 后流生命会消费耗尽无法对该请求做普通 JSON 等解析。**
* **用法**：
  ```go
  bodyBytes, _ := c.GetRawData()
  // 利用自制方法转换 bodyBytes...
  ```

---

## 📤 4. 响应输出 (Response Rendering)

### `c.JSON(code int, obj any)`
* **功能**：向客户端发送一个 JSON 响应并附加指定的 HTTP 状态码。
* **参数意义**：
  * **`code`**：HTTP 状态码。官方库 `net/http` 定义了常量如 `http.StatusOK` (200)。
  * **`obj`**：可以是 struct，也可以是 `gin.H` (即 `map[string]any`)。如果内部含特殊 HTML 字符，会自动转义。
* **用法**：
  ```go
  c.JSON(200, gin.H{
      "code": 0,
      "msg":  "success",
      "data": []string{"apple", "banana"},
  })
  ```

### `c.String(code int, format string, values ...any)`
* **功能**：返回带有具体格式的纯文本。
* **参数意义**：支持类似 `fmt.Printf` 的模版与变参替换。
* **用法**：
  ```go
  c.String(http.StatusOK, "Hello %s!", "World")
  ```

### `c.HTML(code int, name string, obj any)`
* **功能**：渲染注册好的 HTML 模板。
* **参数意义**：
  * **`name`**：模板文件的名称串。
  * **`obj`**：传入模板中用于渲染的变量数据模型。
* **用法**：
  ```go
  c.HTML(200, "index.tmpl", gin.H{"title": "HomePage"})
  ```

### `c.Redirect(code int, location string)`
* **功能**：让客户端执行网页重定向。
* **参数意义**：
  * **`code`**：常使用 `http.StatusMovedPermanently` (301) 或 `http.StatusFound` (302)。
  * **`location`**：目标重定向的绝对或相对 URL。

### `c.File(filepath string)` / `c.FileAttachment(filepath, filename string)`
* **功能**：用于提供文件下载或本地文件读取分发功能。
* **参数意义**：
  * **`filepath`**：文件在服务端物理硬盘真实路径。
  * **`filename`**：下载时客户端看见的保存文件名。

### `c.XML(code int, obj any)` / `c.YAML(...)` / `c.ProtoBuf(...)`
* **功能**：与 JSON 一致，都是直接输出客户端所需的同等对应流的接口渲染。例如返回针对 RSS 阅读器的 `c.XML`；返回高性能序列化格式通讯的 `c.ProtoBuf` 格式响应给后端。
* **参数意义**：与 `c.JSON` 逻辑完全相同，需要配合在相应的 `struct` 里对应写上诸如 `xml`、`yaml` 的解析 Tag 标签。
* **用法**：
  ```go
  // 返回 XML 头 `<?xml version="1.0" encoding="UTF-8"?>` 和后续数据结构串
  c.XML(200, gin.H{"message": "hey", "status": 200})
  ```

### `c.AsciiJSON` / `c.PureJSON` / `c.SecureJSON` / `c.JSONP`
* **功能**：这些属于特定安全场景、或跨域的更**高级、更安全的 JSON 序列化变体定制。**
* **含义**：
  * **`AsciiJSON`**：如果有 Unicode 或非 ASCII 字符的数据体时会将它们转换为不发包异常的如 `\uXXXX` 格式响应。
  * **`PureJSON`**：在 JSON 输出时强行告诉原生的 JSON 序列化类库：“不要对包含如 `<` `>` 等进行可能导致 XSS 的安全转义！”，直接原汁原味地发出去包。（例如有些富文本 API 需要）
  * **`SecureJSON`**：为了防止严重的恶意通过 `Array` 的包含导致 JSON 黑客钓鱼漏洞。框架会在向客户端输出包含中括号包裹的 JSON 数据最前方强制添加防御的指定 `prefix` 魔术字符（默认是 `while(1);`，需要在 Engine 中提前配置 `SecureJSONPrefix` 属性更改）。
  * **`JSONP`**：前端使用 `script src` 原古跨域请求 API 时使用该方法返回对应的数据外层包裹一个前端给的回调函数包响应使用。

### `c.Data(code int, contentType string, data []byte)`
* **功能**：将已经编码/或者未知的纯生自定义类型格式写入 Body 的低级 API，并由 `contentType` 在标头上强制指明（如 `application/pdf`）。
* **参数意义**：`data` 只能吃自己算好的一串 `bytes` 流。
* **用法**：
  ```go
  var protoBufferBytes []byte = 已经加密的自定串();
  c.Data(200, "application/octet-stream", protoBufferBytes)
  ```

### `c.SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool)`
* **功能**：极其强力控制且标准的构建并下发向 HTTP 返回头内追加的 `Set-Cookie` 指令值，强迫让下端的如浏览器记住 Session 等信息。
* **参数意义**：
  * **`maxAge`**：此 Cookie 存在和过期时间（秒为单位）。
  * **`Domain/Path`**：允许生效回传的对应 Host 与 子目录地址范围。
  * **`Secure`**：强行标记该令牌只允许被客户端在开启了 HTTPS 下才可加密携带传送。
  * **`HttpOnly`**：禁止这块小甜饼被目标网页中的任何脚本 JavaScript 读取，是防 XSS 的绝好办法之一。
* **用法**：
  ```go
  // "LoginStatus", 值 "true", 存活 3600 秒, 全站 "/" 生效, 本地域名生效, HttpOnly 设为 true 无法被 JS xss 窃取
  c.SetCookie("LoginStatus", "true", 3600, "/", "localhost", false, true)
  ```

### `c.Header(key, value string)`
* **功能**：直接配置特定的响应 HTTP 标头（例如附加自己的网关、CDN节点、唯一跟踪追踪系统 `TracerID` 版本标记、或是设置允许 `Access-Control-Allow-Origin` 等 CORS 跨站安全参数使用的头）。
* **提醒与用法**：由于写协议的关系，`Header` 发送一定要比正文包 (`JSON`, `HTML`，`String`) **抢先注册写入！** 否则数据包出去的头早定型就没生效了。
  ```go
  // 在 API 返回前增加 X 标与跨域
  c.Header("X-Trace-Id", "1a2b3c4d")
  c.Header("Cache-Control", "no-cache") 
  ```

---

## ⛓️ 5. 中间件与流程控制 (Middlewares & Flow Control)

### `Engine.Use(middleware ...HandlerFunc)`
* **功能**：为引擎或特定路由组挂载全局的中间件处理函数链。常用于安全拦截、限流、日志审计等。
* **用法**：
  ```go
  r.Use(MyCustomMiddleware())
  ```

### `c.Next()`
* **功能**：通常置于中间件内调用。意为挂起当前中间件流程，**先去执行本方法后面的所有后续中间件/路由逻辑**。等它们全部执行完毕返回后，再继续执行当前 `Next()` 后面的代码。常用于统计请求总耗时。
* **用法**：
  ```go
  func Logger() gin.HandlerFunc {
      return func(c *gin.Context) {
          start := time.Now()
          // 挂起，把控制权交给下一个层级
          c.Next()
          // 等到下面全部完成了再执行计时
          latency := time.Since(start)
          log.Println("请求耗时:", latency)
      }
  }
  ```

### `c.Abort()`
* **功能**：立刻掐断后续所有挂载的中间件和目标路由 Handler 的执行操作，但会正常执行完**当前**的这个函数体剩下的逻辑。
* **常见场景**：权限验证失败。

### `c.AbortWithStatusJSON(code int, jsonObj any)`
* **功能**：等同于调用 `c.Abort()` 后再调用 `c.JSON()` 返回错误数据，一种封装的快捷方式。
* **用法**：
  ```go
  if token != "123" {
      c.AbortWithStatusJSON(401, gin.H{"msg": "Unauthorized"})
      return
  }
  ```

### `c.IsAborted() bool`
* **功能**：返回判断这轮当前处理的 `*gin.Context` 生命周期内有没有执行过中断方法 `Abort()`，如果为 true 代表中间是被拦截杀掉进程的，它一般经常用于在整个中间件调用链末尾记录日志里做成监控图的丢弃监控或进行统计打标。
* **返回值含义**：`true` 表示是被框架强行中止了剩下的函数调用序列执行。
* **用法**：
  ```go
  latency := time.Since(start)
  if c.IsAborted() { // 在洋葱型外层全局 Log 收集如果发现
      log.Printf("被拒绝访问的来源耗时 %v", latency) 
  }
  ```

### `c.Error(err error)`
* **功能**：用来统一抛出这一个请求内的相关服务端系统内部、中间层产生的、业务不可抗力导致的 Error 给它收集并暂存在上下文栈中的 `c.Errors` （是一个错误列表串）中。后续配合一个最外层的中间错误接管函数一趟过统一写到后方的 Sentry 等上报分析等错误捕捉报警平台上。
* **参数意义**：传入标准的 `error` 信息实例。（**非常优雅，不在每个函数使用各种零散不一的 fmt.Println。**）
* **用法**：
  ```go
  func CheckDB(c *gin.Context) {
      if err := db.Ping(); err != nil {
          _ = c.Error(err)   // 只保存这个底层系统暴露错误（外部不泄密），让专属 Error中间层上报处理它。
          c.AbortWithStatusJSON(500, gin.H{"msg":"System Err!"})
          return
      }
  }
  ```

---

## 💾 6. 上下文传值 (Context Value Passing)

因为一个 HTTP 请求中的多个中间件和最终 Handler 共用一个 `*gin.Context`，Gin 提供了内置方法方便你在它们之间传递数据。

### `c.Set(key string, value any)`
* **功能**：在当前生命周期的 Context 内写入键值对数据。
* **用法**：
  ```go
  c.Set("UserID", 1024)
  ```

### `c.Get(key string) (value any, exists bool)`
* **功能**：安全地读取 Context 中指定 key 的数据内容。
* **用法**：
  ```go
  if v, exists := c.Get("UserID"); exists {
      log.Println("取得用户标别:", v) // 注意取出来是 any，如需使用具体类型需要断言 v.(int)
  }
  ```

### `c.MustGet(key string) any`
* **功能**：强行获取对应的数据。如果 key 不存在会导致 panic 崩溃。仅在你有绝对把握数据存在时使用。

---

## 🧬 7. 高级操作与原生结构 (Advanced Elements)

由于 `*gin.Context` 本质是 Go 自带框架的包装。有时候仍必须要操作底层实现并发设计或调用只吃 `net/http` 甚至 `context.Context` 接口的各种框架扩展接口，Gin 给它们留了后门和桥接口：

### `c.Copy()`
* **功能**：由于底层的 `gin.Context` 生命周期为了极限压缩内存采用 `sync.Pool` 在复用对象，克隆出一个这个生命周期的 Context 和其数据等内容的**浅值拷贝安全替身副本** (`*gin.Context`)。专门为需要在请求内部起 Goroutine 去并发做其他数据库异步工作而不用等待或堵塞客户端等使用。
* **致命注意场景**：如果开启新协程（如 `go asyncLog(...)`），绝不可在它内部再用原生的原入参引用 `c` （它将被客户端结束的响应给立刻回收覆盖并出现严重指针竞争锁宕机空指针读写错乱风险）。所以必须是向它传入 `c.Copy()` 得到的副本来读参数值传递协程安全机制！
* **用法**：
  ```go
  r.POST("/task", func(c *gin.Context) {
      cCp := c.Copy()   // 一定要有这个副本文去使用！
      // 耗时的业务或者调用另外一个内网消息发送：我们不再这个地方一直等待网络完成（防止客户端死等）
      go func(c2 *gin.Context) {
          time.Sleep(3 * time.Second) 
          log.Println("Done: 原来它的路径是：", c2.Request.URL.Path) // 原来的请求指针 c 早被清空或者给别的客户请求复用啦！
      }(cCp)

      c.JSON(200, gin.H{"msg": "OK我已经加急去通知你"})  // 这是立刻响应成功出去的内容
  })
  ```

### `c.Request` (*http.Request)
* **功能**：此属于向标准库妥协交互的数据：它是 Gin 从最初标准库传入时的原汁原味 Go 的核心请求的源指针（非常重要的 `*http.Request` 的对象类型）。
* **使用原因**：如果要使用兼容其它古董标准库的库件、获取诸如 `r.Host`、非常底层的 `Body io.ReadCloser` 数据或是获取标准库里面的某些信息等情况直接能拿这些字段调用和替换进去交互。
* **用法**：
  ```go
  host := c.Request.Host  // 得到 www.baidu.com
  url := c.Request.URL    // 一个非常庞大的 net/url.URL 标准结构库类
  ```

### `c.Writer` (gin.ResponseWriter)
* **功能**：同样也是兼容到标准的实现了 `http.ResponseWriter` 数据下行的接口的一个被 Gin 增加了统计拦截大小，获取返回状态监控属性功能而再次加强过写出包裹实现指针。
* **用法扩展**：当在使用 WebSockets 的 `Upgrader.Upgrade(c.Writer, c.Request, nil)` 接口把原有的简单返回给强行劫持（升级成长连接全双工网络通讯协议机制实现）流等底层大操作行为的时候我们一定会向其它那些老接口方法上用。还可以调用它的 `c.Writer.Size()` 属性查你最终输出到底输出总占用字节空间大小的计数用来对某些黑客流攻击封包熔断审计拦截监控。

### 作为一个 Go 标准的 `context.Context` 接口规范使用
* **功能**：你可以在大部分要求提供一个安全下传上下文数据要求（如 GRPC 的 `client.Call()` 接口的方法要求首入参就是这个接口类型时，或者用 `gorm` 处理携带这个含有带有本次取消信链事务处理 SQL 查询：`db.WithContext(c).Find(&users)`，如果客户还没等到 DB 出结果他等急了把页面给关了：这时候请求超时断联退出机制触发：那么底层数据库那端发过去的正在阻塞的语句同样也会神奇地触发了超时退出查询语句而不产生资源的白白性能浪费操作！！）
* **用法**：
  ```go
  // c 就是这个对象实例它遵守这个要求，它随时都可以带下去监控它的整个链接的超时跟链路数据的状态（如果用过 gRPC 或者中间件就会很了解它）
  // timeoutctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second) 等包装使用
  var result User
  // 给 gorm 的查询加个防客户端跑路的控制查询锁防占用，因为传了上下文c进去让它去联动查询！
  db.WithContext(c).Where("id=?", 1024).First(&result)
  ```