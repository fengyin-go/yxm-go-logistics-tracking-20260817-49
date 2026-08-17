# 物流追踪系统（logistics-tracking）

一个纯 Go 标准库（`net/http`，零第三方依赖）实现的物流追踪后端服务。采用 `cmd/internal/pkg` 标准工程分层，可编译、可测试、可运行。

## 功能特性

- 网点管理：物流网点 CRUD，被运单引用时禁止删除。
- 运单管理：运单号唯一，状态机 `created → picked → in_transit → delivering → delivered`，支持异常状态 `exception` 并可恢复。
- 包裹信息：一个运单关联一个包裹，记录物品描述与重量。
- 轨迹追踪：每次状态流转自动追加轨迹节点，支持按运单查询完整轨迹（时间升序）。

## 目录结构

```
origin/
├── cmd/server/main.go
├── internal/
│   ├── app/          # 依赖装配
│   ├── config/       # 环境变量配置
│   ├── model/        # 领域模型 + 校验 + 状态机
│   ├── store/        # 数据访问接口 + 内存实现
│   ├── service/      # 业务逻辑（状态流转/轨迹）
│   └── handler/      # HTTP 路由 + 处理器
└── pkg/
    ├── httpx/        # 统一响应/分页/JSON
    ├── idgen/        # ID 与短码生成
    └── logger/       # 分级日志
```

## 运行

```bash
cd origin
go run ./cmd/server
```

环境变量：`PORT`（默认 8080）、`ADDR`、`MAX_PAGE_SIZE`（默认 100）、`LOG_LEVEL`（默认 info）。

## API 一览

| 方法 | 路径 | 说明 |
|------|------|------|
| POST/GET/PUT/DELETE | `/api/stations` | 网点 CRUD |
| POST | `/api/waybills` | 创建运单 |
| GET | `/api/waybills` | 运单列表（`status`/`keyword`） |
| GET | `/api/waybills/{id}` | 运单详情 |
| POST | `/api/waybills/{id}/transition` | 状态流转（`status`/`station_id`/`description`） |
| GET | `/api/waybills/{id}/track` | 查询轨迹 |
| POST/GET/DELETE | `/api/parcels` | 包裹管理 |
| GET | `/api/track-events` | 轨迹列表（`waybill_id`/`status`） |

## 业务闭环示例

```bash
# 1. 建网点
curl -s -X POST localhost:8080/api/stations -d '{"name":"北京分拨","address":"北京"}'
curl -s -X POST localhost:8080/api/stations -d '{"name":"上海分拨","address":"上海"}'

# 2. 创建运单
curl -s -X POST localhost:8080/api/waybills -d '{"sender":"张三","receiver":"李四","origin_station_id":"<from>","dest_station_id":"<to>"}'

# 3. 逐节点流转
curl -s -X POST localhost:8080/api/waybills/<id>/transition -d '{"status":"picked","description":"已揽收"}'
curl -s -X POST localhost:8080/api/waybills/<id>/transition -d '{"status":"in_transit","description":"发往上海"}'
curl -s -X POST localhost:8080/api/waybills/<id>/transition -d '{"status":"delivering","station_id":"<to>","description":"到达派送"}'
curl -s -X POST localhost:8080/api/waybills/<id>/transition -d '{"status":"delivered","station_id":"<to>","description":"已签收"}'

# 4. 查看轨迹
curl -s localhost:8080/api/waybills/<id>/track
```

## 测试

```bash
go test ./...
```
