# Gin 与 HTML 交互详解

本文系统说明 Gin 在服务端渲染 HTML（SSR）场景下的完整交互逻辑：从请求进入、数据准备、模板渲染到响应输出；同时给出常用结构体设计方式与核心 API 用法。

---

## 1. 交互逻辑总览

Gin 与 HTML 交互通常分为 5 步：

1. 浏览器发起 HTTP 请求（如 `GET /profile`）。
2. Gin 路由命中对应 Handler。
3. Handler 组织页面数据（常见来源：数据库、配置、会话、URL 参数）。
4. 通过 `c.HTML(status, templateName, data)` 渲染模板。
5. Gin 返回 `text/html` 响应，浏览器解析并展示。

一句话：**Gin 负责控制流和数据注入，模板负责页面结构展示**。

---

## 2. 结构体设计（页面数据模型）

在 HTML 渲染场景中，建议不要直接把数据库模型原样传到模板，而是定义“页面视图模型（ViewModel）”。

### 2.1 基础结构体示例

```go
type UserCard struct {
	ID       int
	Username string
	Email    string
	IsVIP    bool
}

type ProfilePageData struct {
	Title      string
	CurrentNav string
	User       UserCard
	Tags       []string
	Now        string
}
```

说明：

- 字段首字母要大写（导出字段），模板引擎才能读取。
- 模板中可通过 `.Title`、`.User.Username`、`index .Tags 0` 访问。

### 2.2 通用响应结构体（可选）

如果你希望所有页面都带统一元信息，可封装：

```go
type PageMeta struct {
	Title       string
	Description string
	Keywords    []string
}

type HTMLResponse[T any] struct {
	Meta PageMeta
	Data T
}
```

适合中大型项目：页面 SEO 信息和业务数据分离，更清晰。

---

## 3. 核心 API 说明（重点）

## 3.1 模板加载

### `r.LoadHTMLGlob(pattern string)`

功能：按 glob 规则加载模板文件。

- 参数 `pattern`：例如 `"templates/**/*"`、`"templates/*.tmpl"`。
- 常见用途：模板较多、按目录组织（layout/、user/、admin/）。

示例：

```go
r.LoadHTMLGlob("templates/**/*")
```

### `r.LoadHTMLFiles(files ...string)`

功能：显式指定模板文件列表。

- 参数 `files`：可变参数，逐个给文件路径。
- 常见用途：模板数量少或希望精准控制加载文件。

示例：

```go
r.LoadHTMLFiles("templates/index.tmpl", "templates/profile.tmpl")
```

### `r.SetHTMLTemplate(t *template.Template)`

功能：注入自定义 `html/template` 实例（常用于预解析模板、挂载自定义函数、复杂模板组合）。

- 参数 `t`：标准库 `*template.Template`。

---

## 3.2 页面渲染

### `c.HTML(code int, name string, obj any)`

功能：渲染指定模板并返回 HTML。

- 参数 `code`：HTTP 状态码，页面通常用 `200`、错误页可用 `404/500`。
- 参数 `name`：模板名（通常是文件名，如 `"profile.tmpl"`）。
- 参数 `obj`：模板数据对象（结构体、`gin.H`、`map[string]any` 均可）。

示例：

```go
c.HTML(200, "profile.tmpl", ProfilePageData{
	Title:      "用户主页",
	CurrentNav: "profile",
	User:       UserCard{ID: 1, Username: "alice", Email: "a@demo.com", IsVIP: true},
	Tags:       []string{"Go", "Gin", "HTML"},
	Now:        "2026-04-04 17:30:00",
})
```

### `c.Status(code int)`

功能：仅设置状态码，不输出正文；可用于某些特殊响应流程。

---

## 3.3 模板函数与静态资源

### `r.SetFuncMap(funcMap template.FuncMap)`

功能：向模板注册自定义函数（在 `LoadHTMLGlob/LoadHTMLFiles` 之前调用）。

- 参数 `funcMap`：函数映射表，键是模板中调用名，值是 Go 函数。

示例：

```go
r.SetFuncMap(template.FuncMap{
	"toUpper": strings.ToUpper,
})
```

模板内可用：`{{ toUpper .User.Username }}`。

### `r.Static(relativePath, root string)`

功能：映射静态资源目录（CSS/JS/图片）。

- 参数 `relativePath`：URL 前缀，如 `"/static"`。
- 参数 `root`：本地目录，如 `"./static"`。

示例：

```go
r.Static("/static", "./static")
```

页面里即可引用：`<link rel="stylesheet" href="/static/app.css">`。

---

## 4. 完整示例（可运行）

```go
package main

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserCard struct {
	ID       int
	Username string
	Email    string
	IsVIP    bool
}

type HomePageData struct {
	Title string
	User  UserCard
	News  []string
}

func main() {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"toUpper": strings.ToUpper,
	})
	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		data := HomePageData{
			Title: "首页",
			User:  UserCard{ID: 1, Username: "alice", Email: "a@demo.com", IsVIP: true},
			News:  []string{"Gin 入门", "模板渲染", "静态资源配置"},
		}
		c.HTML(http.StatusOK, "home.tmpl", data)
	})

	_ = r.Run(":8080")
}
```

示例模板 `templates/home.tmpl`：

```html
<!doctype html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="/static/app.css">
</head>
<body>
  <h1>{{ .Title }}</h1>
  <p>用户：{{ toUpper .User.Username }}</p>
  <ul>
    {{ range .News }}
      <li>{{ . }}</li>
    {{ end }}
  </ul>
</body>
</html>
```

---

## 5. 与 JSON 交互的区别

- `c.JSON(...)`：返回结构化数据给前端 JS（前后端分离常用）。
- `c.HTML(...)`：返回完整页面给浏览器（服务端渲染常用）。

选择建议：

- 后台管理系统、SEO 页面：优先 HTML 渲染。
- 开放 API、前后端分离：优先 JSON。
- 混合场景：同一项目可以同时使用两种方式。

---

## 6. 常见问题与排查

### 6.1 报错 `html/template: "xxx" is undefined`

原因：模板名不匹配或未加载到。

排查：

- `LoadHTMLGlob` 路径是否正确。
- `c.HTML(..., "模板名", ...)` 与实际模板名是否一致。

### 6.2 模板拿不到字段

原因：结构体字段未导出（小写开头）。

解决：改成大写字段名。

### 6.3 自定义函数不生效

原因：`SetFuncMap` 调用顺序错误。

解决：先 `SetFuncMap`，再 `LoadHTMLGlob/LoadHTMLFiles`。

### 6.4 样式/脚本 404

原因：静态目录映射错误。

解决：检查 `r.Static("/static", "./static")` 与页面引用路径是否一致。

---

## 7. 最佳实践

1. 使用 ViewModel（页面结构体）与数据库模型解耦。
2. 统一布局模板（header/footer/layout）提升复用性。
3. 自定义模板函数放在单独文件集中维护。
4. 页面渲染时不做复杂业务计算，Handler 先准备好数据。
5. 对错误页（404/500）提供独立模板并统一渲染逻辑。

以上就是 Gin 与 HTML 交互的完整知识骨架。掌握 `LoadHTML* + c.HTML + ViewModel + Static + FuncMap`，基本就能覆盖 90% 的服务端渲染场景。
