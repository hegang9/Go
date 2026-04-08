# strings 与 strconv 标准库详细指南

`strings` 和 `strconv` 是 Go 语言中处理字符串最核心的两个包。

*   **`strings`** 包主要用于字符串的查找、替换、切割、拼接等操作。
*   **`strconv`** 包主要用于字符串与基本数据类型（如整型、布尔型、浮点型）之间的相互转换。

---

## `strings` 常用 API 说明

### 1. `strings.Contains` / `strings.HasPrefix` / `strings.HasSuffix`
*   **功能**：判断字符串中是否包含子串，或者是否以特定子串开头/结尾。
*   **参数说明**：
    *   `s string`：原字符串。
    *   `substr / prefix / suffix string`：要查找的子串。
*   **返回值**：`bool`（是否包含 / 匹配）。
*   **适用场景**：路由前缀判断（如是否以 `/api/` 开头）、简单的模糊搜索。

### 2. `strings.Split` / `strings.Join`
*   **功能**：字符串的切割与拼接。
*   **参数说明**：
    *   `Split(s, sep string) []string`：按 `sep` 分隔符将字符串 `s` 切割成切片。如果 `sep` 为空字符串 `""`，则会将 `s` 切割成每个 Unicode 字符（UTF-8 字符）。
    *   `Join(elems []string, sep string) string`：将字符串切片 `elems` 用 `sep` 连接成一个完整的字符串。
*   **注意事项**：`Join` 的性能通常优于在循环中使用 `+` 拼接字符串。
*   **适用场景**：解析 CSV/逗号分隔的数据、将切片中的 ID 拼成 `1,2,3` 格式传给 SQL 的 `IN` 语句。

### 3. `strings.TrimSpace` / `strings.Trim`
*   **功能**：去除字符串两端的特定字符。
*   **参数说明**：
    *   `TrimSpace(s string) string`：剥离字符串两端的所有空白字符（包括空格、制表符 `\t`、换行符 `\n`、回车符 `\r` 等）。
    *   `Trim(s, cutset string) string`：剥离字符串两端包含在 `cutset` 中的任意字符。
*   **适用场景**：处理用户通过表单提交的输入，防止用户不小心多敲了空格导致校验失败或数据库存入脏数据。

### 4. `strings.Replace` / `strings.ReplaceAll`
*   **功能**：内容替换。
*   **参数说明**：
    *   `Replace(s, old, new string, n int) string`：将 `s` 中的 `old` 替换为 `new`，只替换前 `n` 次。如果 `n < 0`，则无限制替换。
    *   `ReplaceAll(s, old, new string) string`：等价于 `Replace` 中 `n = -1`，替换所有匹配项。
*   **适用场景**：脏话过滤、模板占位符的简单替换。

### 5. `strings.Builder` (结构体)
*   **功能**：高效地进行字符串拼接。
*   **用法示例**：
    ```go
    var builder strings.Builder
    builder.WriteString("hello")
    builder.WriteString(" ")
    builder.WriteString("world")
    result := builder.String()
    ```
*   **注意事项**：在进行大量字符串拼接时（如在 `for` 循环中），使用 `strings.Builder` 的性能远高于使用 `+` 或 `fmt.Sprintf`，因为它在底层直接操作字节切片，减少了内存分配和拷贝。

---

## `strconv` 常用 API 说明

### 1. `strconv.Atoi` / `strconv.Itoa`
*   **功能**：这俩是使用最高频的快捷函数。实现**字符串与 int 型整数**的相互转换。
*   **用法示例**：
    *   `i, err := strconv.Atoi("100")` (ASCII to Integer)
    *   `s := strconv.Itoa(100)` (Integer to ASCII)
*   **注意事项**：`Atoi` 返回的整型是跟特定平台位宽相关的 `int` 类型。如果强转失败（如传了 `"abc"`），会返回错误。
*   **适用场景**：将 URL 路径参数（通常解析出来是 string，如 `/user/123`）转为 int 存入数据库，或者反向操作。

### 2. `strconv.ParseInt` / `strconv.ParseFloat` / `strconv.ParseBool`
*   **功能**：将字符串强转为指定的底层数值类型（可指定进制和位大小）。
*   **参数说明（以 `ParseInt` 为例）**：
    *   `s string`：原字符串。
    *   `base int`：进制数（2 到 36）。传 0 表示根据前缀自动判断（`0x` 为 16 进制，`0` 为 8 进制等）。
    *   `bitSize int`：指定结果必须能装入的整数类型（0, 8, 16, 32, 64 分别对应 `int`, `int8`, `int16`, `int32`, `int64`）。
*   **返回值**：统一返回 `int64` / `float64` 和 `error`，如果需要小类型的，自己在外部再做显式强转。
*   **适用场景**：处理严谨的数据格式转换，或者专门解析浮点数、布尔值（如把 `"true"`, `"T"`, `"1"` 转为 `true`）。

### 3. `strconv.FormatInt` / `strconv.FormatFloat` / `strconv.FormatBool`
*   **功能**：将各类数值类型转换为字符串，是 `Parse` 系列的逆过程。
*   **参数说明（以 `FormatInt` 为例）**：
    *   `i int64`：要转换的整数值。
    *   `base int`：转换后的进制表现形式（如 10 就是十进制字符串，16 就是十六进制字符串）。
*   **适用场景**：需要将非 `int` 基础类型（例如明确就是 `int64` 类型的雪花算法 ID）转成字符串传给前端处理时的场景（由于 JS 会丢失大数精度，长整型 ID 经常先被转为 string）。
