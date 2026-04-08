# Gin 与请求参数交互详解

在 Web API 开发中，后端和前端（或客户端）的“交互核心”就是**请求参数**。  
Gin 已经把常见参数来源都统一到了 `*gin.Context` 上，让你可以用一致的方式处理：

- URL 路径参数（Path Param）
- 查询参数（Query Param）
- 表单参数（Form / Multipart）
- JSON / XML / YAML Body
- Header / Cookie

这篇文档按“从简单到工程化”的顺序，完整介绍 Gin 与请求参数的交互方式、绑定校验、错误处理和最佳实践。

---

## 1. 参数来源全景图

一个 HTTP 请求里，参数常见来源如下：

1. **路径参数**：`/user/:id`
2. **查询参数**：`/search?keyword=go&page=1`
3. **请求体参数**：`POST` 的 JSON / 表单
4. **请求头参数**：`Authorization: Bearer xxx`
5. **Cookie 参数**：会话、追踪等

Gin 里常见对应 API：

- `c.Param("id")`：读路径参数
- `c.Query("key")` / `c.DefaultQuery(...)` / `c.GetQuery(...)`：读查询参数
- `c.PostForm("name")` / `c.DefaultPostForm(...)` / `c.GetPostForm(...)`：读表单字段
- `c.ShouldBind(...)` / `c.ShouldBindJSON(...)`：将参数绑定到结构体
- `c.GetHeader("Authorization")`：读请求头
- `c.Cookie("token")`：读 Cookie

---

## 2. 路径参数（Path Param）

### 2.1 基础读取

```go
r.GET("/user/:id", func(c *gin.Context) {
	id := c.Param("id") // string
	c.JSON(200, gin.H{"user_id": id})
})
```

### 2.2 通配符路由

```go
r.GET("/files/*path", func(c *gin.Context) {
	path := c.Param("path") // 例如 /images/a.png
	c.String(200, path)
})
```

### 2.3 类型转换与校验

`c.Param()` 返回字符串，业务里通常需要转成整数：

```go
r.GET("/order/:id", func(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		c.JSON(400, gin.H{"code": 40001, "msg": "非法订单ID"})
		return
	}
	c.JSON(200, gin.H{"id": id})
})
```

---

## 3. 查询参数（Query Param）

### 3.1 常用方法

- `c.Query("k")`：没有就返回空字符串
- `c.DefaultQuery("k", "默认值")`：没有就返回默认值
- `c.GetQuery("k")`：返回 `(value, exists)`，可区分“没传”和“传了空值”

```go
r.GET("/search", func(c *gin.Context) {
	keyword := c.Query("keyword")
	page := c.DefaultQuery("page", "1")
	sort, exists := c.GetQuery("sort")

	c.JSON(200, gin.H{
		"keyword": keyword,
		"page":    page,
		"sort":    sort,
		"hasSort": exists,
	})
})
```

### 3.2 数组参数

前端常传：`?tag=go&tag=gin&tag=api`

```go
r.GET("/articles", func(c *gin.Context) {
	tags := c.QueryArray("tag")
	c.JSON(200, gin.H{"tags": tags})
})
```

也支持 map：

```go
r.GET("/filters", func(c *gin.Context) {
	filters := c.QueryMap("f")
	// 对应 URL: /filters?f[name]=tom&f[role]=admin
	c.JSON(200, gin.H{"filters": filters})
})
```

---

## 4. 表单参数（x-www-form-urlencoded / multipart/form-data）

### 4.1 读取普通表单

```go
r.POST("/login", func(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	remember := c.DefaultPostForm("remember", "false")

	c.JSON(200, gin.H{
		"username": username,
		"password": password,
		"remember": remember,
	})
})
```

### 4.2 读取上传文件

```go
r.POST("/upload", func(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"msg": "缺少文件"})
		return
	}

	dst := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(500, gin.H{"msg": "保存失败"})
		return
	}
	c.JSON(200, gin.H{"msg": "上传成功", "path": dst})
})
```

> 注意：生产环境应额外校验文件大小、后缀、MIME，并避免直接信任 `file.Filename`。

---

## 5. JSON 请求体（最常见）

### 5.1 推荐：`ShouldBindJSON`

```go
type CreateUserReq struct {
	Name  string `json:"name" binding:"required,min=2,max=20"`
	Age   int    `json:"age" binding:"required,gte=1,lte=120"`
	Email string `json:"email" binding:"omitempty,email"`
}

r.POST("/users", func(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"code": 40002, "msg": err.Error()})
		return
	}

	c.JSON(200, gin.H{"msg": "创建成功", "data": req})
})
```

### 5.2 `BindJSON` 与 `ShouldBindJSON` 区别

- `ShouldBindJSON`：只返回错误，不自动写响应（推荐，便于统一错误格式）。
- `BindJSON`：出错时会自动 `AbortWithError(400, err)`，灵活性较低。

### 5.3 同一个 Body 读取多次

HTTP Body 默认只能消费一次。如果你要先记录原始请求，再绑定结构体，可用：

```go
err := c.ShouldBindBodyWith(&req, binding.JSON)
```

它会缓存 Body，允许后续再次绑定（会有额外内存开销）。

---

## 6. 统一绑定：`ShouldBind` 自动按 Content-Type 推断

如果你希望同一接口兼容 JSON / Form，可用：

```go
type LoginReq struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

r.POST("/session", func(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(200, gin.H{"msg": "ok", "user": req.Username})
})
```

Gin 会根据 `Content-Type` 选择对应 binder（JSON、Form、XML 等）。

---

## 7. Header 与 Cookie 参数

### 7.1 Header

```go
r.GET("/profile", func(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" {
		c.JSON(401, gin.H{"msg": "缺少授权头"})
		return
	}
	c.JSON(200, gin.H{"token": token})
})
```

### 7.2 Cookie

```go
r.GET("/me", func(c *gin.Context) {
	sessionID, err := c.Cookie("session_id")
	if err != nil {
		c.JSON(401, gin.H{"msg": "未登录"})
		return
	}
	c.JSON(200, gin.H{"session_id": sessionID})
})
```

写入 Cookie：

```go
c.SetCookie("session_id", "abc123", 3600, "/", "localhost", false, true)
```

---

## 8. 绑定标签与验证规则

Gin 默认使用 `go-playground/validator`。常见 `binding` 标签：

- `required`：必填
- `min` / `max`：字符串长度最小/最大
- `gte` / `lte`：数值范围（大于等于/小于等于）
- `email`：邮箱格式
- `oneof=a b c`：枚举值
- `omitempty`：字段为空时跳过后续校验

示例：

```go
type PageReq struct {
	Page     int    `form:"page" binding:"omitempty,gte=1"`
	PageSize int    `form:"page_size" binding:"omitempty,gte=1,lte=100"`
	OrderBy  string `form:"order_by" binding:"omitempty,oneof=id created_at"`
}
```

---

## 9. 自定义验证器（进阶）

当内置标签不够时，可以注册自定义校验规则。

```go
type RegisterReq struct {
	Password string `json:"password" binding:"required,strongpwd"`
}

func registerValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("strongpwd", func(fl validator.FieldLevel) bool {
			s := fl.Field().String()
			if len(s) < 8 {
				return false
			}
			hasLetter := regexp.MustCompile(`[A-Za-z]`).MatchString(s)
			hasNumber := regexp.MustCompile(`[0-9]`).MatchString(s)
			return hasLetter && hasNumber
		})
	}
}
```

> 一般在程序启动时注册一次，例如 `main()` 初始化阶段。

---

## 10. 错误处理建议（工程实践）

### 10.1 统一错误响应结构

建议定义固定格式：

```go
type APIResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
```

参数错误时统一返回：

- 业务码：如 `40001`、`40002`
- HTTP 码：通常 `400 Bad Request`
- 可读消息：如“参数校验失败”

### 10.2 不要直接暴露内部细节

`err.Error()` 在开发阶段有帮助，但线上建议：

- 记录详细日志到服务端
- 返回对用户友好的错误文案

---

## 11. 常见坑位清单

1. **结构体字段未导出**：小写字段无法绑定/序列化。
2. **忘记传指针**：`ShouldBindJSON(req)` 应为 `ShouldBindJSON(&req)`。
3. **Tag 不一致**：`json:"name"` 与前端字段名不匹配。
4. **Body 重复读取失败**：应考虑 `ShouldBindBodyWith`。
5. **混用 Bind/ShouldBind**：导致错误响应风格不统一。
6. **缺少参数边界校验**：分页、时间范围、ID 等必须限制。

---

## 12. 一段可运行的综合示例

```go
package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateOrderReq struct {
	ProductID int `json:"product_id" binding:"required,gte=1"`
	Count     int `json:"count" binding:"required,gte=1,lte=99"`
}

func main() {
	r := gin.Default()

	r.POST("/users/:uid/orders", func(c *gin.Context) {
		uidStr := c.Param("uid")
		uid, err := strconv.Atoi(uidStr)
		if err != nil || uid <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40001, "msg": "uid 不合法"})
			return
		}

		traceID := c.GetHeader("X-Trace-Id")
		source := c.DefaultQuery("source", "web")

		var req CreateOrderReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 40002, "msg": "参数校验失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "ok",
			"data": gin.H{
				"uid":        uid,
				"source":     source,
				"trace_id":   traceID,
				"product_id": req.ProductID,
				"count":      req.Count,
			},
		})
	})

	_ = r.Run(":8080")
}
```

---

## 13. 总结

把 Gin 请求参数处理做好，本质上是三件事：

1. **拿对参数**：清楚来源（Path / Query / Body / Header / Cookie）
2. **绑对结构**：优先 `ShouldBind` + 结构体 Tag
3. **校验与报错统一**：保证接口稳定、可维护、可排查

当你把“参数绑定 + 参数校验 + 统一错误响应”这条链路固定下来后，后续接口开发速度和质量都会明显提升。

