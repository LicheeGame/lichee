# 异常
defer延迟调用中recover捕获panic抛出的异常
```
func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    fmt.Println("Calling g.")
    panic(i)
    fmt.Println("Returned normally from g.")
}
```