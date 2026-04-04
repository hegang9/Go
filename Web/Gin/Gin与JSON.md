# Gin 与 JSON 交互详解

在如今的 RESTful API 开发中，**JSON**（JavaScript Object Notation） 是最主要的数据交换格式。Gin 框架对此提供了极其精简和强大的内置支持，无需手动调用标准库 `encoding/json` 处理大量的字节流转换。

在 Gin 中与 JSON 交互，永远围绕两个核心方向展开：**下发（序列化为响应）** 和 **接收（反序列化从请求中解析）**。

---

## 📤 1. 生成并响应 JSON (下行输出)

当我们需要将服务端的数据打包成 JSON 格式返回给前端时，主要通过结构体（Struct）或 `map` 结合 `c.JSON()` 系列方法来实现。

### 1.1 使用 `gin.H` 输出动态 JSON
`gin.H` 是 Gin 框架内置的一种快捷别名，它的底层真实类型是 `map[string]any`。这是返回非固定层级、结构简单的 JSON 最常用的方法。

* **功能**：快捷拼装一个哈希表结构并下发。
* **代码示例**：
  ```go
  r.GET("/ping", func(c *gin.Context) {
      c.JSON(200, gin.H{
          "code": 200,
          "msg":  "PONG",
          "data": gin.H{ // 支持嵌套
              "version": "1.0.0",
          },
      })
  })
  ```
* **输出结果**：`{"code": 200, "data": {"version": "1.0.0"}, "msg": "PONG"}` (由于 Map 的无序性，键的顺序可能是随机的)。

### 1.2 使用 规范的 Struct 输出 JSON（推荐）
在真实的工程中，我们推荐定义**标准化的响应结构体**，以保证 API 输出的严谨性和统一结构。

* **🚨极易踩坑点**：如果希望字段能被 JSON 序列化并对外暴露，该字段的首字母 **必须大写（即 Public/Exported）**！如果是小写，底层的 `encoding/json` 库是不可见的，会被直接丢弃。

* **结构体定义与代码示例**：
  ```go
  // 定义标准响应体
  type BaseResponse struct {
      Code int    `json:"code"` // 注意末尾的 json tag
      Msg  string `json:"msg"`
      Data any    `json:"data"` // 任何类型
  }

  type UserInfo struct {
      UID     int    `json:"uid"`
      Account string `json:"account"`
      // Age 首字母大写，但希望转成 json 后键名小写，所以打 tag
      Age     int    `json:"age"` 
  }

  r.GET("/user/info", func(c *gin.Context) {
      user := UserInfo{UID: 1001, Account: "admin", Age: 18}
      
      c.JSON(200, BaseResponse{
          Code: 0,
          Msg:  "success",
          Data: user,
      })
  })
  ```

### 1.3 `c.JSON()` 的特殊变体 API
有时候仅仅输出普通 JSON 不够，如果遇到安全或特殊兼容场景，可以使用这几个 API：
1. **`c.AsciiJSON(code, obj)`**：如果 JSON 含有特殊中文字符等，它会将其自动转义成 `\uXXXX` 格式防止不可见字符引发前端解析错误。
2. **`c.PureJSON(code, obj)`**：默认 `c.JSON` 为了防范 XSS，会将 `<`，`>`，`&` 转义为类似 `\u003c` 的字串。如果想要返回原汁原味的 HTML 尖括号，就用此方法。
3. **`c.SecureJSON(code, obj)`**：防止 JSON 劫持，如果发送的数据是个数组（`[]` 开头），它会在数据头部强制前缀一段死循环脚本（如 `while(1);`）或其它指定字串，打断恶意挂载 `<script>` 行为。
4. **`c.JSONP(code, obj)`**：为古老的 JSONP 跨域请求技术定制提供支持。

---

## 📥 2. 接收并解析 JSON (上行请求)

当前端通过 `POST/PUT` 等请求并在 Header 携带 `Content-Type: application/json` 请求体时，Gin 会尝试提取这个 Body，将其反序列化到你的目标对象上。

### 2.1 绑定到结构体：`ShouldBindJSON`
这是目前处理入参最安全、最推荐的方法。不仅能自动映射 JSON，还能配合 validator 实现强大的参数校验。

* **结构体定义**：
  ```go
  type RegisterReq struct {
      Username string `json:"username" binding:"required,min=4"`
      Password string `json:"password" binding:"required,min=6"`
      Email    string `json:"email"    binding:"omitempty,email"` 
  }
  ```
* **代码示例**：
  * **注意**：绑定时必须要传入结构体的**指针**（即前面加 `&`），因为函数内部需要修改它的值！
  ```go
  r.POST("/register", func(c *gin.Context) {
      var req RegisterReq
      
      // 注意传 &req
      if err := c.ShouldBindJSON(&req); err != nil {
          // 参数非法或缺失将走到这里
          c.JSON(400, gin.H{"error": err.Error()})
          return
      }

      // 这里 req 已经被安全地填充成了请求发来的数据
      c.JSON(200, gin.H{"status": "已注册", "user": req.Username})
  })
  ```

### 2.2 `BindJSON` vs `ShouldBindJSON`
Gin 提供了两个极其相似的 API，一定要分清：
* **`ShouldBindJSON()`**：解析失败仅返回 `error`。接下来怎么应对，发不发错误状态码 400，**完完全全由开发者自己写代码决定。**（推荐）
* **`BindJSON()`**：解析失败时，不仅返回 `error`，它**底层还会帮你自动执行 `c.AbortWithError(400, err)` 并发送 400 失败的 HTTP HTTP 头部响应出去。** 对于不需要定制错误格式的极简原型有用，但不推荐企业应用。

### 2.3 接收未知动态 JSON 到 Map
如果接口接收的 JSON 格式是千变万化不固定的（没有办法预先写明 struct），可以用万能容器 `map` 或 `any` 接收。

* **代码示例**：
  ```go
  r.POST("/dynamic/data", func(c *gin.Context) {
      // 准备一个空 map
      var data map[string]any
      if err := c.ShouldBindJSON(&data); err != nil {
          c.JSON(400, gin.H{"msg": "invalid json"})
          return
      }

      // 获取里面的任意节点，需要进行类型断言（因为是 any：任何类型）
      if nameParam, ok := data["name"]; ok {
          name := nameParam.(string) 
          c.JSON(200, gin.H{"提取到的动态name为": name})
      }
  })
  ```

---

## 🏷️ 3. Go 语言中的 JSON Tag 语法大全

这部分属于 Go 标准库 `net/http` 和 `encoding/json` 的规定，但它是 Gin 能否顺利干活的根基。在定义上面提到的那些接收或产出结构体时，使用不同的 `json:"xxx"` 属性，能引发不同控制效果：

| Tag 写法 | 作用解释 | 示例场景 |
| :--- | :--- | :--- |
| **`json:"xxx"`** | 将结构体现有的字段重命名，以 `xxx` 作为 JSON 中的 `key`。 | `Age int` `json:"age"` 导出 JSON 为 `{"age": 18}` |
| **`json:"-"`** | **绝对忽略此字段！** 该字段不论有值没值，在转为 JSON 响应时都不会出现在前端；反之解析 JSON 时也当它不存在。 | 诸如给前端屏蔽用户记录查出来的： `Password string` `json:"-"` |
| **`json:"xxx,omitempty"`** | **空值/零值自动隐藏。** 如果当前字段等于该类型的零值（如 数字0、布尔false、空字符串 `""`、空指针 nil），将**不会生出此 JSON Key。** | 若 `Age` 为 0，将不附带 `{"age": 0}` ，适合减少响应报文体积。 |
| **`json:",string"`** | 将 JSON 中的字符串形式和结构体中的数字/布尔互相转换对接兼容（极少使用）。 | 配合老旧接口传过来的数字全用 `""` 包裹的破破局。 |

---

## 📌 4. 总结与日常最佳实践

1. **结构体先行**：无论是入参请求，还是出参响应，尽可能约束在各自的 Struct 结构上。尽量使用统一的 Response 结构（如 `BaseResponse{Code, Msg, Data}`）发牌 `c.JSON()`，从而方便前端人员写固定的拦截器对接，并随时结合 `swagger` 自动生成强类型的文档。
2. **校验靠 Binding**：在 `ShouldBindJSON` 里配合 Gin 自带的校验标签 `binding:"required,min=xx"` 食用，远胜于接收完数据再写一堆 `if req.Age < 18` 这种啰嗦枯燥的判断条件代码。
3. **小心大小写**：再讲一遍，不暴露（小写开头）的字段无法被转化为 JSON！
4. **小心指针解引用冲突**：如果你在参数声明里写了 `var res *User` ，那 `ShouldBindJSON(&res)` 后很可能会面对指针 nil dereference，请永远直接声明普通的 **结构体变量实例** `var req RegisterReq` 然后取它的物理地址。