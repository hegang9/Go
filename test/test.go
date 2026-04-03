package main // 必须声明当前 go 文件属于哪个包，入口文件必须声明为 main 包

import (
	"net/http" // 导入 net/http 包，用于处理 HTTP 请求和响应
)


func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // 注册一个处理函数，当访问根路径时调用
		w.Write([]byte("Hello, World!"))
	})

	http.ListenAndServe("0.0.0.0:8080", nil)

}
