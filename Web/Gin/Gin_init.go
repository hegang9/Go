package main

import "github.com/gin-gonic/gin"

/*
	Gin程序一般包含以下几个步骤：
	1. 创建一个 Gin 引擎实例。（初始化）
	2. 定义路由和处理函数。
	3. 启动服务器，监听指定的端口。
*/

// 返回json数据，可直接使用gin.H{}，它是一个 map[string]interface{} 的快捷方式，方便我们构建 JSON 响应。
func ResponseJSON(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello,Gin!",
		"code":    0,
	})
}

// 返回HTML数据
func ResponseHTML(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

// 文件响应
func ResponseFile(c *gin.Context) {
	// 浏览器直接请求接口唤醒文件下载，特点是要设置响应头，指定内容类型和文件名，且只能被GET请求访问
	c.Header("Content-Type", "application/octet-stream")                // 设置响应头，指定内容类型为二进制流
	c.Header("Content-Disposition", "attachment; filename=Gin_init.go") // 设置响应头，指定文件名和下载方式
	c.File("./Gin_init.go")
}

func main() {
	// 初始化
	gin.SetMode(gin.ReleaseMode) // 设置 Gin 的运行模式为 Release 模式，无 Debug 日志输出，适用于生产环境
	r := gin.Default()

	// 挂载路由 ResponseJSON)
	r.GET("/json", ResponseJSON)
	r.GET("/html", ResponseHTML)
	r.GET("/file", ResponseFile)
	// 绑定端口
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
