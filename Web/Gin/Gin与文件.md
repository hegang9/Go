# Gin 与文件交互详解

本文系统介绍 Gin 在 Web 开发中与“文件”的交互逻辑，覆盖：

1. 文件下载（内联预览、附件下载）
2. 文件上传（单文件、多文件）
3. 静态文件托管（CSS/JS/图片）
4. 结构体设计（请求/响应/文件元数据）
5. 常用 API 参数意义、用法与常见坑

---

## 1. 交互逻辑总览

### 1.1 下载流程

1. 客户端请求某个文件 URL。
2. Gin 路由命中处理函数。
3. 服务端校验权限、文件存在性、路径合法性。
4. 通过 `c.File(...)` 或 `c.FileAttachment(...)` 返回文件流。
5. 浏览器根据响应头决定“预览”或“下载”。

### 1.2 上传流程

1. 客户端以 `multipart/form-data` 提交文件。
2. Gin 解析表单，获取 `*multipart.FileHeader`。
3. 服务端做校验（大小、类型、后缀、路径安全）。
4. `c.SaveUploadedFile(...)` 保存到目标目录。
5. 返回 JSON 结果（路径、文件名、大小等）。

### 1.3 你最关心的区别：静态文件 vs 普通文件

很多同学会把这两个概念混在一起。它们都“返回文件”，但在 Web 架构中的定位完全不同：

- **静态文件（Static Assets）**：框架几乎不做业务处理，按路径直接把固定目录中的资源原样返回。
- **普通文件（Business File）**：必须经过业务 Handler（鉴权、日志、权限、审计、重命名、加水印等）后再返回。

对比表（核心差异）：

| 维度 | 静态文件（`r.Static`） | 普通文件（`c.File` / `c.FileAttachment`） |
| :--- | :--- | :--- |
| 入口 | 直接由静态路由处理 | 先进入你写的业务路由 Handler |
| 业务逻辑 | 几乎没有 | 可做任意业务逻辑 |
| 权限控制 | 不方便精细控制 | 非常适合做用户级/租户级鉴权 |
| 典型内容 | CSS、JS、图片、字体、前端构建产物 | 用户上传附件、合同、报表、私有文件 |
| URL 规则 | 通常固定前缀（如 `/static/...`） | 通常业务 URL（如 `/api/file/:id`） |
| 缓存策略 | 常配合强缓存、CDN | 常配合短缓存或不缓存（取决于业务） |
| 安全风险 | 主要是目录暴露范围过大 | 主要是越权下载、路径穿越 |

一句话判断：

- **公共资源，人人可访问** → 用 `r.Static`。
- **与用户身份/权限绑定** → 用 `c.File` 或 `c.FileAttachment`。

### 1.4 请求链路差异（为什么你会觉得“看起来都一样”）

静态文件链路：

`浏览器 -> /static/logo.png -> Gin静态路由 -> 直接读磁盘 -> 返回`

普通文件链路：

`浏览器 -> /api/files/123 -> 业务Handler -> 查DB确认归属 + 鉴权 -> 定位真实路径 -> c.File/c.FileAttachment`

所以它们不是“API 名字不同”，而是**系统职责不同**。

---

## 2. 结构体设计（推荐）

文件交互里，结构体通常用于“统一响应格式”和“文件元信息”。

## 2.1 文件元信息结构体

```go
type FileMeta struct {
	OriginalName string `json:"originalName"`
	StoredName   string `json:"storedName"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mimeType"`
	URL          string `json:"url"`
}
```

字段说明：

- `OriginalName`：客户端上传时原始文件名。
- `StoredName`：服务器落盘后的文件名（建议重命名，防冲突）。
- `Size`：文件大小（字节）。
- `MimeType`：文件 MIME 类型，如 `image/png`。
- `URL`：返回给前端用于访问该文件的地址。

## 2.2 统一响应结构体

```go
type APIResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
```

好处：前端解析统一，错误处理统一。

## 2.3 上传请求结构体（带额外字段场景）

当上传接口除了文件还要带业务字段（如目录类型、用户 ID）时：

```go
type UploadForm struct {
	BizType string `form:"bizType" binding:"required"`
	UserID  int64  `form:"userId" binding:"required"`
}
```

说明：

- `form` 标签用于解析表单字段。
- 文件本体仍通过 `c.FormFile("file")` 获取。

---

## 3. 核心 API 详解

## 3.1 文件下载类

### `c.File(filepath string)`

功能：直接返回指定路径文件，通常适合“浏览器可内联预览”的内容（如图片、PDF）。

- 参数 `filepath`：服务器本地文件路径（绝对或相对）。

用法：

```go
r.GET("/preview/:name", func(c *gin.Context) {
	name := c.Param("name")
	path := "./uploads/" + name
	c.File(path)
})
```

### `c.FileAttachment(filepath, filename string)`

功能：强制下载附件，自动设置 `Content-Disposition: attachment`。

- 参数 `filepath`：服务器文件路径。
- 参数 `filename`：客户端“另存为”看到的文件名。

用法：

```go
r.GET("/download/:name", func(c *gin.Context) {
	name := c.Param("name")
	c.FileAttachment("./uploads/"+name, name)
})
```

### `c.Header(key, value string)` + `c.Data(...)`

功能：自定义响应头和文件流（适合做更细粒度控制）。

示例：

```go
c.Header("Content-Disposition", "attachment; filename=report.txt")
c.Data(200, "text/plain; charset=utf-8", []byte("hello"))
```

---

## 3.2 文件上传类

### `c.FormFile(name string) (*multipart.FileHeader, error)`

功能：从 `multipart/form-data` 请求中读取单个文件。

- 参数 `name`：前端 `<input type="file" name="...">` 的 `name` 值。

用法：

```go
file, err := c.FormFile("file")
if err != nil {
	c.JSON(400, gin.H{"msg": "未找到上传文件"})
	return
}
```

### `c.SaveUploadedFile(file *multipart.FileHeader, dst string) error`

功能：把上传文件保存到本地磁盘。

- 参数 `file`：来自 `FormFile` 的文件头。
- 参数 `dst`：目标保存路径（包含文件名）。

用法：

```go
dst := "./uploads/" + file.Filename
if err := c.SaveUploadedFile(file, dst); err != nil {
	c.JSON(500, gin.H{"msg": "保存失败"})
	return
}
```

### `c.MultipartForm() (*multipart.Form, error)`

功能：读取整个 multipart 表单（含多文件、多字段）。

用法（多文件上传）：

```go
form, err := c.MultipartForm()
if err != nil {
	c.JSON(400, gin.H{"msg": "表单解析失败"})
	return
}

files := form.File["files"]
for _, file := range files {
	_ = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
}
```

### `r.MaxMultipartMemory`

功能：控制 multipart 解析时的内存阈值，超过部分会落临时文件。

用法：

```go
r := gin.Default()
r.MaxMultipartMemory = 8 << 20 // 8 MiB
```

---

## 3.3 静态资源托管

### `r.Static(relativePath, root string)`

功能：将 URL 前缀映射到本地目录。

- `relativePath`：URL 前缀，如 `/static`。
- `root`：本地目录，如 `./static`。

注意：`root` 目录下的文件会被“公开可访问”。如果你把上传目录也映射成静态目录，往往会带来越权风险（例如别人猜文件名即可下载）。

```go
r.Static("/static", "./static")
```

### `r.StaticFS(relativePath string, fs http.FileSystem)`

功能：使用自定义文件系统托管静态资源（如嵌入式文件、只读 FS）。

### `r.StaticFile(relativePath, filepath string)`

功能：将“单个固定路由”绑定到“单个固定文件”。

```go
r.StaticFile("/favicon.ico", "./static/favicon.ico")
```

### 什么时候不要用静态路由

以下场景不建议用 `r.Static`：

1. 文件需要登录后访问。
2. 文件要按用户/部门/租户做权限隔离。
3. 下载前需要计费、审计、限流、签名校验。
4. 文件路径不能暴露真实磁盘结构。

这种情况应改用：`GET /api/file/:id` + 业务校验 + `c.FileAttachment(...)`。

---

## 4. 可运行示例（上传 + 下载）

```go
package main

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileMeta struct {
	OriginalName string `json:"originalName"`
	StoredName   string `json:"storedName"`
	Size         int64  `json:"size"`
	URL          string `json:"url"`
}

func main() {
	r := gin.Default()
	r.MaxMultipartMemory = 8 << 20

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "请选择文件"})
			return
		}

		dst := filepath.Join("uploads", file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"msg": "保存失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "上传成功",
			"data": FileMeta{
				OriginalName: file.Filename,
				StoredName:   file.Filename,
				Size:         file.Size,
				URL:          "/file/" + file.Filename,
			},
		})
	})

	r.GET("/file/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.File("uploads/" + name)
	})

	r.GET("/download/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.FileAttachment("uploads/"+name, name)
	})

	_ = r.Run(":8080")
}
```

---

## 5. 常见坑与安全建议

1. 不要直接信任 `file.Filename`，建议重命名（UUID/时间戳）防冲突与注入风险。
2. 严格限制大小（`MaxMultipartMemory` + 业务层大小校验）。
3. 校验后缀/MIME（防上传脚本木马）。
4. 禁止路径穿越（如 `../`），对用户输入做清洗。
5. 下载前做权限校验（避免越权访问）。
6. 不要把私有上传目录直接 `r.Static` 暴露出去。

---

## 5.1 一个直观例子（同一文件，两种访问方式）

```go
// 方式A：静态路由（公开）
r.Static("/public", "./uploads")
// 任何人只要知道 URL 都能访问 /public/abc.pdf

// 方式B：业务路由（可控）
r.GET("/api/download/:name", func(c *gin.Context) {
	name := c.Param("name")
	// 这里可以做登录态检查、归属校验、次数限制、操作日志
	if !hasPermission(c, name) {
		c.JSON(403, gin.H{"msg": "无权限"})
		return
	}
	c.FileAttachment("./uploads/"+name, name)
})
```

结论：**静态路由是“公开分发工具”，业务路由是“受控文件服务”。**

---

## 6. API 速查表

- `c.File(path)`：返回文件（偏预览）。
- `c.FileAttachment(path, name)`：返回附件（强制下载）。
- `c.FormFile("file")`：读取单文件。
- `c.MultipartForm()`：读取整表单/多文件。
- `c.SaveUploadedFile(file, dst)`：保存上传文件。
- `r.Static("/static", "./static")`：托管静态目录。
- `r.StaticFile("/a", "./a.txt")`：托管单文件。
- `r.MaxMultipartMemory = n`：设置 multipart 内存阈值。

掌握上面这些 API，你就能完成 Gin 中绝大多数“文件上传 + 下载 + 资源托管”场景。
