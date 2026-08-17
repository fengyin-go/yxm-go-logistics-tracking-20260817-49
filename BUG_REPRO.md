# Bug 是什么
部分按关联键查询的未命中错误没有保持 `ErrNotFound` 语义，导致调用方把不存在的数据识别成冲突。

# 如何触发
运行未命中错误分类测试：

```bash
go test ./internal/store -run TestLookupMissesPreserveNotFoundSentinel -count=1
```

# 错误信息
```text
track event miss err = 记录已存在或状态冲突
```
