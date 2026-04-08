# encoding/json 标准库详细指南

`encoding/json` 包实现了 JSON 对象的编码（也就是将 Go 语言的数据结构转换成 JSON 字符串，序列化）和解码（将 JSON 字符串转换成 Go 数据结，反序列化）。

## 核心接口说明 (marshaler 和 unmarshaler)

由于 Go 是静态类型语言，不像 JavaScript 可以任意增删属性，Go 必须借助结构体（`struct`）和标签（`tag`）将 JSON 绑定到特定的对象类型上，这就涉及了 Go 的反射机制。

`json` 中的大部分类型操作基于四个核心方法展开。

## 常用 API 说明

### 1. `json.Marshal` / `json.MarshalIndent`
* **功能**：将传入的 Go 数据结构（包括基本类型、切片、字典、特别是**结构体**）序列化为 JSON 格式的字节数组。
* **参数说明**：
  * `v any`：各种非指针、包含接口的复合类型的变量实例。
  * `prefix, indent string` (仅 `MarshalIndent`)：为了美观化格式输出 JSON，增加的换行符和缩进格式。例如 `""`, `  `（两个空格）。
* **返回值**：`([]byte, error)`，代表 JSON 数据的切片（可转为 string）。
* **注意事项（极其重要）**：
  * 该方法**仅仅导出首字母大写的所有公开字段**。如果你在结构体里的字段是小写首字母，那么在转换时它们将会默默地消失掉。
  * 配合**结构体标签 (Tag)** 控制重命名：\n```go\n  type User struct {\n      Name int    `json:"name,omitempty"` // omitempty：当由于类型默认产生的零值时（比如 int 为 0），在 JSON 结果中直接剔除此字段。\n      Password string `json:"-"` // 指定 JSON 序列化时直接忽略掉密码，不输出。\n  }\n```
* **适用场景**：所有组装好数据并打成字符串传递给前端或者其他微服务的场景。

### 2. `json.Unmarshal`
* **功能**：解析并挂载从外部取得的 JSON 字符串、切片数据，赋值到自己结构体中。
* **参数说明**：
  * `data []byte`：外部收到的 JSON 数据。
  * `v any`：用于接收的目标。**必须是指针（传地址）类型**，且其底层类型与 JSON 中的结构尽量对应匹配，若不匹配有些将报 `error`。如果目标是 `map[string]any`，则会按任意类型接收（如对象对应 `map[string]any`，数组对应 `[]any`，数字默认都是 `float64`）。
* **用法示例**：
```go
  var u User
  err := json.Unmarshal([]byte(`{"name":...}`), &u)
  // 或用一个不知道格式的表接住
  var res map[string]any
  _ = json.Unmarshal([]byte(`{"任意": ...}`), &res)
```
* **注意事项**：
  * 接收字段首字母需大写！若 JSON 数据有该字段但结构体没有，Go 会直接忽视并跳过它（并不会报错）。
  * 结构体字段若存在（比如数字 0），JSON 中没有传，则该对象此字段的解析结果就是默认值 0 而已（它无法分辨是用户原本就传了 0 还是没有传，除非你把它设为指针 `*int`）。
* **适用场景**：处理客户端 POST 或者 PUT 请求提交的 JSON 请求体 payload；接收下游接口返回的结果。

### 3. `json.NewEncoder` / `json.NewDecoder`
* **功能**：与前面两者的功能基本一样，但它们的出入口对接的是**流系统 (IO)**。它们被设计成直接与 `io.Reader` 或 `io.Writer` 对接，比如 `http.ResponseWriter` 或 `os.File`，可以流式逐行解析/输出。
* **参数说明**：`w io.Writer` 或 `r io.Reader`：将编码和解码逻辑直接挂载在流通道上。
* **返回值**：编码器或解码器的结构体实例。调用其 `.Encode(v)` / `.Decode(&v)` 即可。
* **用法示例**：
```go
func HandleUser(w http.ResponseWriter, r *http.Request) {
    // 省去了 ioutil.ReadAll(r.Body) 读到底层再 Unmarshal 的两次拷贝开销
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil { ... }
    
    // ...
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(response) // 直接写入响应流
}
```
* **注意事项**：它比 `json.Marshal` 高效，因为它不需要先把完整切片加载到内存中再写入到 HTTP 网络发送通道一次，而是直接在流中转换发送，极大地节省了海量 JSON 字符导致的内存暴涨。
* **适用场景**：读取大型 JSON 日志文件、API 的 HTTP 流（Streaming）发送和接收响应。
