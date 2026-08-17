package service

// defaultPageSize 是 page/size 被传入 0 或负数时的回退页大小。
const defaultPageSize = 20

// pageBounds 根据匹配到的总数 total 与分页参数 page/size，计算安全的切片区间
// [start, end)，调用方据此执行 items[start:end] 切片。
//
// 该函数对非法入参做兜底，保证 0 <= start <= end <= total 恒成立，因此切片
// 永不 panic，也避免前端传入 page=0 / 负数 / 越界大页时返回错误的空结果：
//   - page <= 0：回退为第 1 页；
//   - size <= 0：回退为默认页大小 defaultPageSize；
//   - (page-1)*size 整数溢出为负（page 极大）：按越界处理，返回空区间；
//   - start 已越过 total：返回空区间。
func pageBounds(total, page, size int) (start, end int) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = defaultPageSize
	}
	start = (page - 1) * size
	// start < 0 兜底 (page-1)*size 的整数溢出；start >= total 表示页号越界。
	if start < 0 || start >= total {
		return total, total // [total:total) 即空切片
	}
	end = start + size
	if end > total || end < start { // end < start 兜底加法溢出
		end = total
	}
	return start, end
}
