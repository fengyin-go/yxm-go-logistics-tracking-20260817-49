# Bug 是什么
包裹列表返回了内存存储中的原始指针，调用方修改列表项会污染后续读取到的包裹数据。

# 如何触发
运行包裹列表隔离测试：

```bash
go test ./internal/store -run TestParcelListDoesNotExposeStoredPointers -count=1
```

# 错误信息
```text
stored parcel was mutated through list result: "被外部污染"
```
