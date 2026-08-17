# Bug 是什么
运单搜索没有正确归一化带空格的状态和关键词，刚揽收的运单会被过滤掉。

# 如何触发
运行运单搜索归一化测试：

```bash
go test ./internal/service -run TestWaybillSearchNormalizesKeywordAndStatus -count=1
```

# 错误信息
```text
normalized search len=0 total=0
```
