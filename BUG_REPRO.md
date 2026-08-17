# Bug 是什么
列表分页没有统一兜底，`page` 或 `size` 为 0、负数时会返回错误切片范围或错误结果。

# 如何触发
运行分页边界测试：

```bash
go test ./internal/service -run TestListMethodsClampPaginationBounds -count=1
```

# 错误信息
```text
page 0 size 0 returned len=0 total=3
```
