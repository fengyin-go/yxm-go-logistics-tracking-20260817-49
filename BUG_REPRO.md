# Bug 是什么
异常件进入 `exception` 后无法恢复到 `delivering`，并且异常轨迹节点没有完整保留下来。

# 如何触发
运行异常件恢复派送的服务层测试：

```bash
go test ./internal/service -run TestExceptionWaybillCanResumeDeliveryWithCompleteTrack -count=1
```

# 错误信息
```text
resume delivering after exception: 记录已存在或状态冲突
```
