package main

import "github.com/gin-gonic/gin"

/*
	Gin程序一般包含以下几个步骤：
	1. 创建一个 Gin 引擎实例。（初始化）
	2. 定义路由和处理函数。
	3. 启动服务器，监听指定的端口。
*/

func main() {
	// 初始化
	gin.SetMode(gin.ReleaseMode) // 设置 Gin 的运行模式为 Release 模式，无 Debug 日志输出，适用于生产环境
	r := gin.Default()

	// 挂载路由
	r.GET("/index", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	// 绑定端口
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}
